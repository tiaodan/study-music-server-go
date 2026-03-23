############# v0.1
-- 跑通一个产品(本地) v0.1 -> 原java项目，转成go语言, 上传, 备份
1. 把原项目: https://github.com/tiaodan/study-music-server-java.git


############# v0.2
-- 跑通一个产品(本地) v0.2 -> 把下载的歌曲, 存sql，归整到HDD(硬盘)

目的：把下载的歌曲，根据 歌手-专辑，把名字、歌词名字重新格式化，规整进数据库，并移动到HDD
仅后端实现相关接口
需求拆分：
1. 名字重新格式化。
	- 实现：api指定 歌手-专辑路径，自动遍历文件(mp3/wav/layic)名字,重新格式化(多作者-歌名.文件类型)
2. 移动到HDD
	- 实现：api指定路径，from - to ,根据路径，移动；移动后，查询是否移动成功
3. 规整进数据库
	- 实现：api指定路径，path:xx, 自动遍历文件(mp3/wav/layic)名字,
		-- 插入作者。如果有新作者，插入数据库
		-- 把歌曲名，插入song表
		-- 关联song 和 singer关系

重新设计数据库 （保证数据正确性），
sing/song表应该是有关联关系的，关联时删，不能删，修改时，同时改




############# sql相关
1. resource/sql/tp_music.sql_v0.1 原项目自带
2. resource/sql/tp_music.sql_v0.2 修改后，song 和 singer 增加关联关系