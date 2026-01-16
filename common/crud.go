// Package common provides generic CRUD operations and utilities for building REST APIs.
//
// This package offers reusable DTOs (Data Transfer Objects) and pipe configurations
// for common database operations like Find, FindById, Search, Exists, and Delete.
//
// # Architecture
//
// The package uses a "Pipe" pattern where configuration objects (pipes) are passed through
// context to customize query behavior. Each DTO type has a corresponding Pipe type:
//
//   - ExistsDto + ExistsPipe: Check if a record exists by field value
//   - FindDto + FindPipe: Paginated list queries with sorting
//   - FindByIdDto + FindByIdPipe: Single record retrieval by ID
//   - SearchDto + SearchPipe: Lightweight search/autocomplete queries
//
// # Basic Usage Pattern
//
// All operations follow the same pattern:
//  1. Bind request parameters to a DTO
//  2. Create and configure a Pipe with desired options
//  3. Store the Pipe in context using SetPipe()
//  4. Call the DTO's query method
//
// # Security
//
// This package implements multiple layers of security:
//   - Sort column names are validated using whitelist (Sortable method)
//   - Exists column names are validated using whitelist (NewExistsPipe)
//   - ID formats are validated (UUID by default, customizable)
//   - All user inputs use parameterized queries to prevent SQL injection
package common

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// pipeKey is a custom type used as context key to avoid key collisions.
type pipeKey struct{}

// validColumnName validates column names to prevent SQL injection.
// Only allows letters, numbers, and underscores, starting with a letter or underscore.
var validColumnName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// validUUID validates UUID format (with or without hyphens).
var validUUID = regexp.MustCompile(`^[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{12}$`)

// IDValidator is a function type for custom ID validation.
type IDValidator func(id string) bool

// DefaultIDValidator validates that ID is a valid UUID.
var DefaultIDValidator IDValidator = func(id string) bool {
	return validUUID.MatchString(id)
}

// Controller defines the standard CRUD interface for API controllers.
type Controller interface {
	Create(ctx context.Context, c *app.RequestContext)
	Find(ctx context.Context, c *app.RequestContext)
	FindById(ctx context.Context, c *app.RequestContext)
	Update(ctx context.Context, c *app.RequestContext)
	Delete(ctx context.Context, c *app.RequestContext)
}

// SetPipe stores a pipe configuration in the context.
// The pipe will be retrieved later by DTO methods to customize query behavior.
func SetPipe(ctx context.Context, i any) context.Context {
	return context.WithValue(ctx, pipeKey{}, i)
}

// getPipe retrieves a typed pipe configuration from context.
// Returns the pipe and a boolean indicating whether the retrieval was successful.
func getPipe[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(pipeKey{}).(T)
	return v, ok
}

// ToOrderBy maps sort direction indicators to SQL ORDER BY suffixes.
// "1" maps to ascending (empty string), "-1" maps to "desc".
var ToOrderBy = map[string]string{
	"1":  "",
	"-1": "desc",
}

// ExistsDto is the data transfer object for checking record existence.
// Used to verify if a value already exists in a specific column.
//
// Query parameters:
//   - key: The column name to check
//   - q: The value to search for
type ExistsDto struct {
	Key string `query:"key,omitempty"`
	Q   string `query:"q,omitempty"`
}

// ExistsPipe configures which fields are allowed for existence checks.
// This provides a whitelist of valid column names for security.
type ExistsPipe struct {
	fields map[string]bool
}

// NewExistsPipe creates a new ExistsPipe with the specified allowed field names.
// Only fields in the whitelist can be used for existence checks.
//
// Example:
//
//	// Allow checking email, username, and phone fields
//	ctx = common.SetPipe(ctx, common.NewExistsPipe("email", "username", "phone"))
//
//	// Client request: GET /users/exists?key=email&q=test@example.com
//	// Returns: {"exists": true} if email exists
func NewExistsPipe(keys ...string) *ExistsPipe {
	fields := make(map[string]bool)
	for _, key := range keys {
		fields[key] = true
	}
	return &ExistsPipe{
		fields: fields,
	}
}

// Get retrieves the ExistsPipe from context.
// Returns an error if the pipe is not found.
func (x *ExistsDto) Get(ctx context.Context) (*ExistsPipe, error) {
	p, ok := getPipe[*ExistsPipe](ctx)
	if !ok {
		return nil, help.E(0, "ExistsPipe not found in context")
	}
	return p, nil
}

// ExistsResult represents the response for existence check queries.
type ExistsResult struct {
	Exists bool `json:"exists"`
}

// Exists checks if a record with the specified value exists in the given column.
// Returns ExistsResult with Exists=true if a matching record is found.
//
// Security: The column name is validated against the pipe's whitelist.
// Only columns explicitly allowed via NewExistsPipe can be queried.
//
// Example:
//
//	func (x *Controller) Exists(ctx context.Context, c *app.RequestContext) {
//	    var dto ExistsDto
//	    c.BindAndValidate(&dto)
//
//	    ctx = common.SetPipe(ctx, common.NewExistsPipe("email", "phone"))
//	    result, err := dto.Exists(ctx, db.Model(&User{}))
//	    if err != nil {
//	        c.Error(err)
//	        return
//	    }
//	    c.JSON(200, result)  // {"exists": true/false}
//	}
func (x *ExistsDto) Exists(ctx context.Context, do *gorm.DB) (result ExistsResult, err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	// Whitelist validation - only allowed columns can be queried
	if !p.fields[x.Key] {
		err = help.E(0, fmt.Sprintf(`field [%s] is not allowed for existence check`, x.Key))
		return
	}
	var count int64
	if err = do.
		Where(`? = ?`, clause.Column{Name: x.Key}, x.Q).
		Count(&count).Error; err != nil {
		return
	}
	result = ExistsResult{
		Exists: count != 0,
	}
	return
}

// FindDto is the data transfer object for paginated list queries.
//
// Headers:
//   - x-pagesize: Number of records per page (default: 1000, max: 1000)
//   - x-page: Page number (0-indexed)
//
// Query parameters:
//   - q: Search keyword for filtering
//   - sort: Sort rules in format "column:direction" (e.g., "name:1", "created_at:-1")
type FindDto struct {
	PageSize int64    `header:"x-pagesize" vd:"omitempty,min=0,max=1000"`
	Page     int64    `header:"x-page" vd:"omitempty,min=0"`
	Q        string   `query:"q,omitempty"`
	Sort     []string `query:"sort,omitempty" vd:"omitempty,dive,sort"`
}

// GetPageSize returns the page size, defaulting to 1000 if not specified.
func (x *FindDto) GetPageSize() int {
	if x.PageSize == 0 {
		x.PageSize = 1000
	}
	return int(x.PageSize)
}

// GetOffset calculates the offset for pagination based on page number and page size.
func (x *FindDto) GetOffset() int {
	return int(x.Page) * int(x.PageSize)
}

// GetKeyword returns the search keyword wrapped with SQL LIKE wildcards.
// Example: "test" becomes "%test%"
func (x *FindDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

// FindPipe configures the behavior of Find queries.
type FindPipe struct {
	ts       bool            // Whether to handle timestamp fields (created_at, updated_at)
	sort     bool            // Whether to apply sorting
	page     bool            // Whether to apply pagination
	keys     []string        // Specific columns to select
	omit     []string        // Columns to exclude from results
	sortable map[string]bool // Whitelist of sortable column names
}

// Get retrieves the FindPipe from context.
// Returns an error if the pipe is not found.
func (x *FindDto) Get(ctx context.Context) (*FindPipe, error) {
	p, ok := getPipe[*FindPipe](ctx)
	if !ok {
		return nil, help.E(0, "FindPipe not found in context")
	}
	return p, nil
}

// NewFindPipe creates a new FindPipe with default settings.
// By default, timestamp handling, sorting, and pagination are all enabled.
//
// Default behavior:
//   - Omits created_at and updated_at columns
//   - Orders by created_at desc if no sort specified
//   - Applies pagination (default 1000 per page)
//
// Example:
//
//	// Basic usage with defaults
//	ctx = common.SetPipe(ctx, common.NewFindPipe())
//
//	// Custom configuration
//	ctx = common.SetPipe(ctx, common.NewFindPipe().
//	    SkipTs().                                    // Don't auto-omit timestamps
//	    Sortable("name", "email", "created_at").     // Whitelist sortable columns
//	    Omit("password", "secret").                  // Exclude sensitive fields
//	    SkipPage())                                  // Disable pagination
func NewFindPipe() *FindPipe {
	return &FindPipe{
		ts:   true,
		sort: true,
		page: true,
	}
}

// SkipTs disables automatic timestamp field handling.
// When disabled, created_at and updated_at won't be automatically omitted.
func (x *FindPipe) SkipTs() *FindPipe {
	x.ts = false
	return x
}

// SkipSort disables sorting. No ORDER BY clause will be applied.
func (x *FindPipe) SkipSort() *FindPipe {
	x.sort = false
	return x
}

// SkipPage disables pagination. All matching records will be returned.
func (x *FindPipe) SkipPage() *FindPipe {
	x.page = false
	return x
}

// Select specifies which columns to include in the query results.
// When set, only these columns will be returned.
func (x *FindPipe) Select(keys ...string) *FindPipe {
	x.keys = keys
	return x
}

// Omit specifies which columns to exclude from the query results.
// This is ignored if Select is used.
func (x *FindPipe) Omit(keys ...string) *FindPipe {
	x.omit = keys
	return x
}

// Sortable sets the whitelist of columns that can be used for sorting.
// If not set, any valid column name format will be accepted (less secure).
// For better security, always specify the allowed sortable columns.
//
// Example:
//
//	pipe := NewFindPipe().Sortable("name", "created_at", "email")
func (x *FindPipe) Sortable(keys ...string) *FindPipe {
	x.sortable = make(map[string]bool)
	for _, key := range keys {
		x.sortable[key] = true
	}
	return x
}

// Factory builds a GORM query with the configured options from FindPipe.
// Applies column selection/omission, sorting, and pagination.
//
// Security: Sort column names are validated against whitelist (if configured)
// or basic format validation to prevent SQL injection.
func (x *FindDto) Factory(ctx context.Context, do *gorm.DB) (*gorm.DB, error) {
	p, err := x.Get(ctx)
	if err != nil {
		return nil, err
	}
	if len(p.keys) != 0 {
		do = do.Select(p.keys)
	} else {
		if len(p.omit) == 0 && p.ts {
			do = do.Omit(`created_at`, `updated_at`)
		}
		if len(p.omit) != 0 {
			do = do.Omit(p.omit...)
		}
	}

	if p.sort {
		if len(x.Sort) == 0 && p.ts {
			do = do.Order("created_at desc")
		}
		for _, v := range x.Sort {
			rule := strings.Split(v, ":")
			if len(rule) != 2 {
				return nil, help.E(0, fmt.Sprintf(`invalid sort format: %s`, v))
			}
			columnName := rule[0]
			// Validate column name using whitelist if configured, otherwise use regex
			if len(p.sortable) > 0 {
				if !p.sortable[columnName] {
					return nil, help.E(0, fmt.Sprintf(`column [%s] is not sortable`, columnName))
				}
			} else {
				// Fallback to basic format validation
				if !validColumnName.MatchString(columnName) {
					return nil, help.E(0, fmt.Sprintf(`invalid column name in sort: %s`, columnName))
				}
			}
			order, ok := ToOrderBy[rule[1]]
			if !ok {
				return nil, help.E(0, fmt.Sprintf(`invalid sort direction: %s`, rule[1]))
			}
			do = do.Order(fmt.Sprintf(`%s %s`, columnName, order))
		}
	}

	if p.page {
		do = do.Limit(x.GetPageSize()).Offset(x.GetOffset())
	}
	return do, nil
}

// Find executes a paginated query and scans results into the provided slice.
//
// Example:
//
//	func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
//	    var dto FindDto
//	    c.BindAndValidate(&dto)
//
//	    ctx = common.SetPipe(ctx, common.NewFindPipe().
//	        Sortable("name", "created_at").
//	        Omit("password"))
//
//	    do := db.Model(&User{})
//	    if dto.Q != "" {
//	        do = do.Where("name LIKE ?", dto.GetKeyword())
//	    }
//
//	    var results []User
//	    if err := dto.Find(ctx, do, &results); err != nil {
//	        c.Error(err)
//	        return
//	    }
//	    c.JSON(200, results)
//	}
//
// Client request examples:
//
//	GET /users                           // Default pagination and sorting
//	GET /users?sort=name:1               // Sort by name ASC
//	GET /users?sort=created_at:-1        // Sort by created_at DESC
//	GET /users?q=john                    // Search with keyword
//	Headers: X-Page: 0, X-PageSize: 20   // Pagination control
func (x *FindDto) Find(ctx context.Context, do *gorm.DB, i any) (err error) {
	db, err := x.Factory(ctx, do)
	if err != nil {
		return
	}
	return db.Find(i).Error
}

// FindByIdDto is the data transfer object for single record retrieval.
//
// Path parameters:
//   - id: The record ID
//
// Query parameters:
//   - full: Set to 1 to retrieve all fields (full mode)
type FindByIdDto struct {
	ID   string `path:"id"`
	Full int    `query:"full,omitempty"`
}

// IsFull returns true if full mode is requested (all fields should be returned).
func (x *FindByIdDto) IsFull() bool {
	return x.Full == 1
}

// FindByIdPipe configures the behavior of FindById queries.
// Supports different column configurations for normal and full modes.
type FindByIdPipe struct {
	ts          bool        // Whether to handle timestamp fields
	keys        []string    // Columns to select in normal mode
	omit        []string    // Columns to omit in normal mode
	fKeys       []string    // Columns to select in full mode
	fOmit       []string    // Columns to omit in full mode
	idValidator IDValidator // Custom ID validator function
}

// Get retrieves the FindByIdPipe from context.
// Returns an error if the pipe is not found.
func (x *FindByIdDto) Get(ctx context.Context) (*FindByIdPipe, error) {
	p, ok := getPipe[*FindByIdPipe](ctx)
	if !ok {
		return nil, help.E(0, "FindByIdPipe not found in context")
	}
	return p, nil
}

// NewFindByIdPipe creates a new FindByIdPipe with default settings.
// By default, timestamp handling is enabled and UUID validation is used.
//
// Supports two modes:
//   - Normal mode: Returns limited fields (default)
//   - Full mode (?full=1): Returns all fields for editing
//
// Example:
//
//	// Basic usage
//	ctx = common.SetPipe(ctx, common.NewFindByIdPipe())
//
//	// Custom configuration with different fields for normal/full modes
//	ctx = common.SetPipe(ctx, common.NewFindByIdPipe().
//	    Omit("password", "secret").         // Normal mode: hide sensitive fields
//	    FullOmit("password").               // Full mode: only hide password
//	    SkipIDValidation())                 // For non-UUID IDs
//
//	// Numeric ID validation
//	ctx = common.SetPipe(ctx, common.NewFindByIdPipe().
//	    SetIDValidator(func(id string) bool {
//	        _, err := strconv.Atoi(id)
//	        return err == nil
//	    }))
func NewFindByIdPipe() *FindByIdPipe {
	return &FindByIdPipe{
		ts:          true,
		idValidator: DefaultIDValidator,
	}
}

// SkipTs disables automatic timestamp field handling.
func (x *FindByIdPipe) SkipTs() *FindByIdPipe {
	x.ts = false
	return x
}

// Select specifies which columns to include in normal mode.
func (x *FindByIdPipe) Select(keys ...string) *FindByIdPipe {
	x.keys = keys
	return x
}

// Omit specifies which columns to exclude in normal mode.
func (x *FindByIdPipe) Omit(keys ...string) *FindByIdPipe {
	x.omit = keys
	return x
}

// FullSelect specifies which columns to include in full mode.
func (x *FindByIdPipe) FullSelect(keys ...string) *FindByIdPipe {
	x.fKeys = keys
	return x
}

// FullOmit specifies which columns to exclude in full mode.
func (x *FindByIdPipe) FullOmit(keys ...string) *FindByIdPipe {
	x.fOmit = keys
	return x
}

// SetIDValidator sets a custom ID validator function.
// Use this to customize ID format validation (e.g., for non-UUID IDs).
//
// Example:
//
//	pipe := NewFindByIdPipe().SetIDValidator(func(id string) bool {
//	    _, err := strconv.Atoi(id)
//	    return err == nil
//	})
func (x *FindByIdPipe) SetIDValidator(v IDValidator) *FindByIdPipe {
	x.idValidator = v
	return x
}

// SkipIDValidation disables ID format validation.
// Use with caution - only when you trust the input source.
func (x *FindByIdPipe) SkipIDValidation() *FindByIdPipe {
	x.idValidator = nil
	return x
}

// Take retrieves a single record by ID with the configured column selection.
// Uses normal or full mode configuration based on the Full query parameter.
//
// Security: ID format is validated before query execution.
//
// Example:
//
//	func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
//	    var dto FindByIdDto
//	    c.BindAndValidate(&dto)
//
//	    ctx = common.SetPipe(ctx, common.NewFindByIdPipe().
//	        Omit("password").
//	        FullOmit("password"))
//
//	    var result User
//	    if err := dto.Take(ctx, db.Model(&User{}), &result); err != nil {
//	        c.Error(err)
//	        return
//	    }
//	    c.JSON(200, result)
//	}
//
// Client request examples:
//
//	GET /users/:id           // Normal mode, limited fields
//	GET /users/:id?full=1    // Full mode, all fields (for editing)
func (x *FindByIdDto) Take(ctx context.Context, do *gorm.DB, i any) (err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	// Validate ID format if validator is configured
	if p.idValidator != nil && !p.idValidator(x.ID) {
		return help.E(0, fmt.Sprintf(`invalid ID format: %s`, x.ID))
	}
	if !x.IsFull() {
		if len(p.keys) != 0 {
			do = do.Select(p.keys)
		} else {
			if len(p.omit) == 0 && p.ts {
				do = do.Omit(`created_at`, `updated_at`)
			}
			if len(p.omit) != 0 {
				do = do.Omit(p.omit...)
			}
		}
	} else {
		if len(p.fKeys) != 0 {
			do = do.Select(p.fKeys)
		} else {
			if len(p.fOmit) != 0 {
				do = do.Omit(p.fOmit...)
			}
		}
	}

	return do.Where(`id = ?`, x.ID).Take(i).Error
}

// SearchDto is the data transfer object for lightweight search/autocomplete queries.
// Designed for quick lookups with minimal data transfer.
//
// Query parameters:
//   - m: Search mode (optional, for custom filtering)
//   - q: Search keyword
//   - ids: Comma-separated list of IDs to prioritize in results
type SearchDto struct {
	M   string `query:"m,omitempty"`
	Q   string `query:"q,omitempty"`
	IDs string `query:"ids,omitempty"`
}

// GetKeyword returns the search keyword wrapped with SQL LIKE wildcards.
func (x *SearchDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

// SearchPipe configures the behavior of Search queries.
type SearchPipe struct {
	keys        []string    // Columns to return (default: id, name)
	async       bool        // Whether to limit results for async/autocomplete use
	idValidator IDValidator // Custom ID validator function for IDs parameter
}

// SkipAsync disables the result limit. By default, search returns max 50 results.
func (x *SearchPipe) SkipAsync() *SearchPipe {
	x.async = false
	return x
}

// NewSearchPipe creates a new SearchPipe with the specified columns to return.
// Defaults to ["id", "name"] if no columns are specified.
// Async mode (50 result limit) is enabled by default.
// UUID validation for IDs is enabled by default.
//
// Designed for autocomplete/dropdown data sources with minimal payload.
//
// Example:
//
//	// Basic usage - returns id and name, max 50 results
//	ctx = common.SetPipe(ctx, common.NewSearchPipe())
//
//	// Custom columns
//	ctx = common.SetPipe(ctx, common.NewSearchPipe("id", "name", "avatar"))
//
//	// Disable result limit for full search
//	ctx = common.SetPipe(ctx, common.NewSearchPipe().SkipAsync())
func NewSearchPipe(keys ...string) *SearchPipe {
	search := &SearchPipe{
		async:       true,
		idValidator: DefaultIDValidator,
	}
	if len(keys) == 0 {
		search.keys = []string{"id", "name"}
	} else {
		search.keys = keys
	}
	return search
}

// SetIDValidator sets a custom ID validator function for the IDs parameter.
func (x *SearchPipe) SetIDValidator(v IDValidator) *SearchPipe {
	x.idValidator = v
	return x
}

// SkipIDValidation disables ID format validation for the IDs parameter.
func (x *SearchPipe) SkipIDValidation() *SearchPipe {
	x.idValidator = nil
	return x
}

// Get retrieves the SearchPipe from context.
// Returns an error if the pipe is not found.
func (x *SearchDto) Get(ctx context.Context) (*SearchPipe, error) {
	p, ok := getPipe[*SearchPipe](ctx)
	if !ok {
		return nil, help.E(0, "SearchPipe not found in context")
	}
	return p, nil
}

// Factory builds a GORM query with the configured options from SearchPipe.
// Applies column selection and optional result limiting.
func (x *SearchDto) Factory(ctx context.Context, do *gorm.DB) (*gorm.DB, error) {
	p, err := x.Get(ctx)
	if err != nil {
		return nil, err
	}
	if p.async {
		do = do.Limit(50)
	}
	return do.Select(p.keys), nil
}

// Find executes a search query with optional ID prioritization.
// If IDs are provided, those records appear first in results (using UNION ALL).
//
// Security: IDs are validated before query execution if validator is configured.
//
// Example:
//
//	func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
//	    var dto SearchDto
//	    c.BindAndValidate(&dto)
//
//	    ctx = common.SetPipe(ctx, common.NewSearchPipe("id", "name"))
//
//	    do := db.Model(&User{})
//	    if dto.Q != "" {
//	        do = do.Where("name LIKE ?", dto.GetKeyword())
//	    }
//
//	    var results []SearchResult
//	    if err := dto.Find(ctx, do, &results); err != nil {
//	        c.Error(err)
//	        return
//	    }
//	    c.JSON(200, results)
//	}
//
// Client request examples:
//
//	GET /users/search?q=john                    // Search by keyword
//	GET /users/search?ids=uuid1,uuid2           // Get specific IDs first
//	GET /users/search?ids=uuid1,uuid2&q=john    // Prioritize IDs, then search
func (x *SearchDto) Find(ctx context.Context, do *gorm.DB, i any) (err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	if x.IDs != "" {
		ids := strings.Split(x.IDs, ",")
		// Validate each ID if validator is configured
		if p.idValidator != nil {
			for _, id := range ids {
				id = strings.TrimSpace(id)
				if id != "" && !p.idValidator(id) {
					return help.E(0, fmt.Sprintf(`invalid ID format in IDs: %s`, id))
				}
			}
		}
		factory, err := x.Factory(ctx, do.WithContext(ctx))
		if err != nil {
			return err
		}
		return do.Raw(`(?) union all (?)`,
			do.WithContext(ctx).Select(p.keys).Where(`id in (?)`, ids),
			factory.Where(`id not in (?)`, ids),
		).Find(i).Error
	}
	factory, err := x.Factory(ctx, do)
	if err != nil {
		return
	}
	return factory.Find(i).Error
}

// SearchResult is a standard response structure for search queries.
type SearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeleteDto is the data transfer object for batch delete operations.
type DeleteDto struct {
	IDs []string `json:"ids"`
}
