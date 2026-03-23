-- v0.6 -> v0.7 迁移脚本
-- 功能：重命名字段

-- 步骤1：nas_url -> nas_url_path
ALTER TABLE song CHANGE COLUMN nas_url nas_url_path VARCHAR(255);

-- 步骤2：spider_url_path -> spider_url（改为存完整URL）
ALTER TABLE song CHANGE COLUMN spider_url_path spider_url VARCHAR(500);
