-- v0.2.18: song 表新增 singer_id 和 upload_aws_status 列，修改唯一索引

-- 1. 删除旧索引
ALTER TABLE song DROP INDEX idx_album_name;

-- 2. 添加 singer_id 列（放在 album_id 前面）
ALTER TABLE song ADD COLUMN singer_id INT UNSIGNED DEFAULT NULL COMMENT '歌手ID（冗余列，主歌手）' AFTER id;

-- 3. 添加 upload_aws_status 列
ALTER TABLE song ADD COLUMN upload_aws_status TINYINT DEFAULT 0 COMMENT 'AWS上传状态：0-未上传 1-已上传 2-失败' AFTER is_hot;

-- 4. 创建新唯一索引（singer_id + album_id + name）
ALTER TABLE song ADD UNIQUE INDEX idx_singer_album_name (singer_id, album_id, name);

-- 5. 从专辑关联回填 singer_id（主歌手 = 专辑所属歌手）
UPDATE song s
INNER JOIN album a ON s.album_id = a.id
SET s.singer_id = a.singer_id;