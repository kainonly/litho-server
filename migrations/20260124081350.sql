-- Modify "menu" table
ALTER TABLE "menu" ADD COLUMN "icon" text NOT NULL DEFAULT '';
-- Set comment to column: "icon" on table: "menu"
COMMENT ON COLUMN "menu"."icon" IS '图标';
-- Drop index "idx_resource_action_resource_sort" from table: "resource_action"
DROP INDEX "idx_resource_action_resource_sort";
-- Set comment to column: "code" on table: "resource_action"
COMMENT ON COLUMN "resource_action"."code" IS '编码';
