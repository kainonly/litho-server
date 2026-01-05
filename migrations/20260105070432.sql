-- Create "menu" table
CREATE TABLE "menu" (
  "id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "sort" smallint NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_menu_created_at" to table: "menu"
CREATE INDEX "idx_menu_created_at" ON "menu" ("created_at");
-- Set comment to table: "menu"
COMMENT ON TABLE "menu" IS '导航表';
-- Set comment to column: "sort" on table: "menu"
COMMENT ON COLUMN "menu"."sort" IS '排序';
-- Set comment to column: "active" on table: "menu"
COMMENT ON COLUMN "menu"."active" IS '状态';
-- Set comment to column: "name" on table: "menu"
COMMENT ON COLUMN "menu"."name" IS '导航名称';
-- Create "org" table
CREATE TABLE "org" (
  "id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_org_created_at" to table: "org"
CREATE INDEX "idx_org_created_at" ON "org" ("created_at");
-- Set comment to table: "org"
COMMENT ON TABLE "org" IS '组织表';
-- Set comment to column: "active" on table: "org"
COMMENT ON COLUMN "org"."active" IS '状态';
-- Set comment to column: "name" on table: "org"
COMMENT ON COLUMN "org"."name" IS '组织名称';
-- Create "permission" table
CREATE TABLE "permission" (
  "id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "code" text NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_permission_code" to table: "permission"
CREATE UNIQUE INDEX "idx_permission_code" ON "permission" ("code");
-- Create index "idx_permission_created_at" to table: "permission"
CREATE INDEX "idx_permission_created_at" ON "permission" ("created_at");
-- Set comment to table: "permission"
COMMENT ON TABLE "permission" IS '特定授权表';
-- Set comment to column: "active" on table: "permission"
COMMENT ON COLUMN "permission"."active" IS '状态';
-- Set comment to column: "code" on table: "permission"
COMMENT ON COLUMN "permission"."code" IS '授权编码';
-- Set comment to column: "description" on table: "permission"
COMMENT ON COLUMN "permission"."description" IS '描述';
-- Create "resource" table
CREATE TABLE "resource" (
  "id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "sort" smallint NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true,
  "name" text NOT NULL,
  "code" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_resource_code" to table: "resource"
CREATE UNIQUE INDEX "idx_resource_code" ON "resource" ("code");
-- Create index "idx_resource_created_at" to table: "resource"
CREATE INDEX "idx_resource_created_at" ON "resource" ("created_at");
-- Set comment to table: "resource"
COMMENT ON TABLE "resource" IS '资源表';
-- Set comment to column: "sort" on table: "resource"
COMMENT ON COLUMN "resource"."sort" IS '排序';
-- Set comment to column: "active" on table: "resource"
COMMENT ON COLUMN "resource"."active" IS '状态';
-- Set comment to column: "name" on table: "resource"
COMMENT ON COLUMN "resource"."name" IS '资源名称';
-- Set comment to column: "code" on table: "resource"
COMMENT ON COLUMN "resource"."code" IS '路径';
-- Create "resource_action" table
CREATE TABLE "resource_action" (
  "id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "sort" smallint NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true,
  "name" text NOT NULL,
  "code" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_resource_action_code" to table: "resource_action"
CREATE UNIQUE INDEX "idx_resource_action_code" ON "resource_action" ("code");
-- Create index "idx_resource_action_created_at" to table: "resource_action"
CREATE INDEX "idx_resource_action_created_at" ON "resource_action" ("created_at");
-- Create index "idx_resource_action_resource_id" to table: "resource_action"
CREATE INDEX "idx_resource_action_resource_id" ON "resource_action" ("resource_id");
-- Set comment to table: "resource_action"
COMMENT ON TABLE "resource_action" IS '资源操作表';
-- Set comment to column: "resource_id" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."resource_id" IS '资源ID';
-- Set comment to column: "sort" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."sort" IS '排序';
-- Set comment to column: "active" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."active" IS '状态';
-- Set comment to column: "name" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."name" IS '操作名称';
-- Set comment to column: "code" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."code" IS '编码';
-- Create "role" table
CREATE TABLE "role" (
  "id" bigint NOT NULL,
  "org_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "sort" smallint NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true,
  "name" text NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_role_created_at" to table: "role"
CREATE INDEX "idx_role_created_at" ON "role" ("created_at");
-- Create index "idx_role_org_id" to table: "role"
CREATE INDEX "idx_role_org_id" ON "role" ("org_id");
-- Set comment to table: "role"
COMMENT ON TABLE "role" IS '权限表';
-- Set comment to column: "org_id" on table: "role"
COMMENT ON COLUMN "role"."org_id" IS '所属组织ID';
-- Set comment to column: "sort" on table: "role"
COMMENT ON COLUMN "role"."sort" IS '排序';
-- Set comment to column: "active" on table: "role"
COMMENT ON COLUMN "role"."active" IS '状态';
-- Set comment to column: "name" on table: "role"
COMMENT ON COLUMN "role"."name" IS '权限名称';
-- Set comment to column: "description" on table: "role"
COMMENT ON COLUMN "role"."description" IS '权限描述';
-- Create "role_menu" table
CREATE TABLE "role_menu" (
  "id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  "menu_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_role_menu" to table: "role_menu"
CREATE UNIQUE INDEX "idx_role_menu" ON "role_menu" ("role_id", "menu_id");
-- Set comment to table: "role_menu"
COMMENT ON TABLE "role_menu" IS '权限导航表';
-- Set comment to column: "role_id" on table: "role_menu"
COMMENT ON COLUMN "role_menu"."role_id" IS '权限ID';
-- Set comment to column: "menu_id" on table: "role_menu"
COMMENT ON COLUMN "role_menu"."menu_id" IS '导航ID';
-- Create "role_permission" table
CREATE TABLE "role_permission" (
  "id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  "permission_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_role_permission" to table: "role_permission"
CREATE UNIQUE INDEX "idx_role_permission" ON "role_permission" ("role_id", "permission_id");
-- Set comment to table: "role_permission"
COMMENT ON TABLE "role_permission" IS '权限特定授权表';
-- Set comment to column: "role_id" on table: "role_permission"
COMMENT ON COLUMN "role_permission"."role_id" IS '权限ID';
-- Set comment to column: "permission_id" on table: "role_permission"
COMMENT ON COLUMN "role_permission"."permission_id" IS '特定授权ID';
-- Create "role_route" table
CREATE TABLE "role_route" (
  "id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  "route_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_role_route" to table: "role_route"
CREATE UNIQUE INDEX "idx_role_route" ON "role_route" ("role_id", "route_id");
-- Set comment to table: "role_route"
COMMENT ON TABLE "role_route" IS '权限路由表';
-- Set comment to column: "role_id" on table: "role_route"
COMMENT ON COLUMN "role_route"."role_id" IS '权限ID';
-- Set comment to column: "route_id" on table: "role_route"
COMMENT ON COLUMN "role_route"."route_id" IS '路由ID';
-- Create "route" table
CREATE TABLE "route" (
  "id" bigint NOT NULL,
  "menu_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "sort" smallint NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true,
  "pid" bigint NOT NULL DEFAULT 0,
  "name" text NOT NULL,
  "type" smallint NOT NULL DEFAULT 1,
  "icon" text NOT NULL,
  "link" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_route_created_at" to table: "route"
CREATE INDEX "idx_route_created_at" ON "route" ("created_at");
-- Create index "idx_route_menu_id" to table: "route"
CREATE INDEX "idx_route_menu_id" ON "route" ("menu_id");
-- Set comment to table: "route"
COMMENT ON TABLE "route" IS '路由表';
-- Set comment to column: "menu_id" on table: "route"
COMMENT ON COLUMN "route"."menu_id" IS '导航ID';
-- Set comment to column: "sort" on table: "route"
COMMENT ON COLUMN "route"."sort" IS '排序';
-- Set comment to column: "active" on table: "route"
COMMENT ON COLUMN "route"."active" IS '状态';
-- Set comment to column: "pid" on table: "route"
COMMENT ON COLUMN "route"."pid" IS '父级ID';
-- Set comment to column: "name" on table: "route"
COMMENT ON COLUMN "route"."name" IS '路由名称';
-- Set comment to column: "type" on table: "route"
COMMENT ON COLUMN "route"."type" IS '路由类型';
-- Set comment to column: "icon" on table: "route"
COMMENT ON COLUMN "route"."icon" IS '字体图标';
-- Set comment to column: "link" on table: "route"
COMMENT ON COLUMN "route"."link" IS '链接';
-- Create "route_resource_action" table
CREATE TABLE "route_resource_action" (
  "id" bigint NOT NULL,
  "route_id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "action_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_route_resource_action" to table: "route_resource_action"
CREATE UNIQUE INDEX "idx_route_resource_action" ON "route_resource_action" ("route_id", "resource_id", "action_id");
-- Set comment to table: "route_resource_action"
COMMENT ON TABLE "route_resource_action" IS '路由资源表';
-- Set comment to column: "route_id" on table: "route_resource_action"
COMMENT ON COLUMN "route_resource_action"."route_id" IS '路由ID';
-- Set comment to column: "resource_id" on table: "route_resource_action"
COMMENT ON COLUMN "route_resource_action"."resource_id" IS '资源ID';
-- Set comment to column: "action_id" on table: "route_resource_action"
COMMENT ON COLUMN "route_resource_action"."action_id" IS '操作ID';
-- Create "user" table
CREATE TABLE "user" (
  "id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "email" text NOT NULL,
  "phone" text NOT NULL,
  "name" text NOT NULL,
  "password" text NOT NULL,
  "avatar" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_created_at" to table: "user"
CREATE INDEX "idx_user_created_at" ON "user" ("created_at");
-- Create index "idx_user_email" to table: "user"
CREATE UNIQUE INDEX "idx_user_email" ON "user" ("email");
-- Create index "idx_user_phone" to table: "user"
CREATE INDEX "idx_user_phone" ON "user" ("phone");
-- Set comment to table: "user"
COMMENT ON TABLE "user" IS '用户表';
-- Set comment to column: "active" on table: "user"
COMMENT ON COLUMN "user"."active" IS '状态';
-- Set comment to column: "email" on table: "user"
COMMENT ON COLUMN "user"."email" IS '电子邮件';
-- Set comment to column: "phone" on table: "user"
COMMENT ON COLUMN "user"."phone" IS '手机号';
-- Set comment to column: "name" on table: "user"
COMMENT ON COLUMN "user"."name" IS '姓名';
-- Set comment to column: "password" on table: "user"
COMMENT ON COLUMN "user"."password" IS '密码';
-- Set comment to column: "avatar" on table: "user"
COMMENT ON COLUMN "user"."avatar" IS '头像';
-- Create "user_org_role" table
CREATE TABLE "user_org_role" (
  "id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "org_id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_org_role" to table: "user_org_role"
CREATE UNIQUE INDEX "idx_user_org_role" ON "user_org_role" ("user_id", "org_id");
-- Create index "idx_user_org_role_role_id" to table: "user_org_role"
CREATE INDEX "idx_user_org_role_role_id" ON "user_org_role" ("role_id");
-- Set comment to table: "user_org_role"
COMMENT ON TABLE "user_org_role" IS '用户组织权限表';
-- Set comment to column: "user_id" on table: "user_org_role"
COMMENT ON COLUMN "user_org_role"."user_id" IS '用户ID';
-- Set comment to column: "org_id" on table: "user_org_role"
COMMENT ON COLUMN "user_org_role"."org_id" IS '组织ID';
-- Set comment to column: "role_id" on table: "user_org_role"
COMMENT ON COLUMN "user_org_role"."role_id" IS '权限ID';
