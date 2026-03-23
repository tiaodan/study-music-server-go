-- v0.2 -> v0.3 迁移脚本
-- 功能：新增专辑表和歌曲歌手中间表，改用多对多关系

-- 步骤1：创建 album 表（专辑表）
CREATE TABLE IF NOT EXISTS `album` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `singer_id` int unsigned NOT NULL,
  `pic` varchar(255) DEFAULT NULL,
  `introduction` varchar(255) DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_singer_id` (`singer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- 步骤2：创建 song_singer 中间表（歌曲歌手多对多关系）
CREATE TABLE IF NOT EXISTS `song_singer` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `song_id` int unsigned NOT NULL,
  `singer_id` int unsigned NOT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_song_singer` (`song_id`, `singer_id`),
  KEY `idx_song_id` (`song_id`),
  KEY `idx_singer_id` (`singer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- 步骤3：删除 song 表原有的外键约束（先解除外键关联，保留 singer_id 字段用于数据迁移）
ALTER TABLE song DROP FOREIGN KEY fk_song_singer;

-- 步骤4：将现有 singer_id 数据迁移到 song_singer 中间表
INSERT INTO song_singer (song_id, singer_id, create_time)
SELECT id, singer_id, NOW()
FROM song
WHERE singer_id IS NOT NULL;

-- 步骤5：删除 song 表的 singer_id 字段
ALTER TABLE song DROP COLUMN singer_id;

-- 步骤6：为 song 表增加 album_id 字段
ALTER TABLE song ADD COLUMN album_id int unsigned DEFAULT NULL AFTER id;

-- 步骤7：为 song 表的 album_id 添加索引
ALTER TABLE song ADD KEY `idx_album_id` (`album_id`);

-- 步骤8：添加 album_id 外键约束
ALTER TABLE song
ADD CONSTRAINT fk_song_album
FOREIGN KEY (album_id)
REFERENCES album(id)
ON UPDATE CASCADE
ON DELETE SET NULL;
