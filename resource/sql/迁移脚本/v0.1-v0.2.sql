-- v0.1 -> v0.2 迁移脚本
-- 功能：为 song 表添加外键约束，建立 singer 与 song 的关联关系

-- 步骤1：修改 singer_id 类型为 int unsigned（与 singer.id 一致）
ALTER TABLE song MODIFY singer_id INT UNSIGNED DEFAULT NULL;

-- 步骤2：为 song 表添加外键约束
-- 关联表：song -> singer
-- 外键字段：song.singer_id -> singer.id
-- 删除时限制：删除 singer 时，如果有关联 song 记录，则拒绝删除（RESTRICT）
-- 更新时级联：更新 singer.id 时，关联的 song.singer_id 自动同步更新（CASCADE）
ALTER TABLE song
ADD CONSTRAINT fk_song_singer
FOREIGN KEY (singer_id)
REFERENCES singer(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;
