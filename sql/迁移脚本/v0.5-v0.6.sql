-- v0.5 -> v0.6 迁移脚本
-- 功能：重命名字段

-- 步骤1：spider_url -> spider_url_path
ALTER TABLE song CHANGE COLUMN spider_url spider_url_path VARCHAR(500);

-- 步骤2：full_name -> full_name_singer
ALTER TABLE song CHANGE COLUMN full_name full_name_singer VARCHAR(255);
