-- v0.3 -> v0.4 迁移脚本
-- 功能：为 song 表的 (album_id, name) 添加唯一索引

-- 添加唯一索引 (album_id, name)
-- 注意：MySQL 唯一索引允许多个 NULL 值
ALTER TABLE song ADD UNIQUE INDEX uk_album_name (album_id, name);
