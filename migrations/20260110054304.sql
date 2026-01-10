-- Set comment to column: "code" on table: "permission"
COMMENT ON COLUMN "permission"."code" IS '授权标识';
-- Create index "idx_resource_action_resource_sort" to table: "resource_action"
CREATE UNIQUE INDEX "idx_resource_action_resource_sort" ON "resource_action" ("resource_id", "sort");
-- Set comment to column: "code" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."code" IS '路径';
-- Create index "idx_role_org_name" to table: "role"
CREATE UNIQUE INDEX "idx_role_org_name" ON "role" ("org_id", "name");
