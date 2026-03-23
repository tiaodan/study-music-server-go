-- v0.4 -> v0.5 迁移脚本
-- 功能：重构 song 表 URL 相关字段

-- 步骤1：将现有 url 数据迁移到 nas_url
ALTER TABLE song ADD COLUMN nas_url VARCHAR(255) AFTER lyric;
UPDATE song SET nas_url = url;

-- 步骤2：删除旧的 url 字段
ALTER TABLE song DROP COLUMN url;

-- 步骤3：新增 5 个字段
ALTER TABLE song
ADD COLUMN spider_url VARCHAR(500) AFTER nas_url,
ADD COLUMN spider_url_https VARCHAR(500) AFTER spider_url,
ADD COLUMN aws_url VARCHAR(500) AFTER spider_url_https,
ADD COLUMN aws_url_temp VARCHAR(500) AFTER aws_url,
ADD COLUMN full_name VARCHAR(255) AFTER aws_url_temp;
