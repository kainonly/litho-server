package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 用户名
	Username string `bson:"username" json:"username"`

	// 密码
	Password string `bson:"password" json:"password"`

	// 权限组
	Roles []primitive.ObjectID `bson:"roles" json:"roles"`

	// 独立授权页面
	Pages []primitive.ObjectID `bson:"pages" json:"pages"`

	// 只读权限
	Readonly []primitive.ObjectID `bson:"readonly" json:"readonly"`

	// 显示名称
	Name string `bson:"name" json:"name"`

	// 电子邮件
	Email []string `bson:"email" json:"email"`

	// 头像
	Avatar string `bson:"avatar" json:"avatar"`

	// 标签
	Labels []Value `bson:"labels" json:"labels"`

	// 状态
	Status *bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewUser(username string, password string) *User {
	return &User{
		Username:   username,
		Password:   password,
		Roles:      []primitive.ObjectID{},
		Pages:      []primitive.ObjectID{},
		Readonly:   []primitive.ObjectID{},
		Email:      []string{},
		Labels:     []Value{},
		Status:     Bool(true),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *User) SetRoles(v []primitive.ObjectID) *User {
	x.Roles = v
	return x
}

func (x *User) SetPages(v []primitive.ObjectID) *User {
	x.Pages = v
	return x
}

func (x *User) SetReadonly(v []primitive.ObjectID) *User {
	x.Readonly = v
	return x
}

func (x *User) SetName(v string) *User {
	x.Name = v
	return x
}

func (x *User) SetEmail(v []string) *User {
	x.Email = v
	return x
}

func (x *User) SetAvatar(v string) *User {
	x.Avatar = v
	return x
}
