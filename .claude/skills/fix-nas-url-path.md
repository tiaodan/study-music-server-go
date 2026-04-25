# Fix NAS URL Path

Clean up `nas_url_path` field in database by removing spaces. This ensures database paths match cleaned file names.

## Usage
User invokes: `/fix-nas-url-path`

## Target Tables
- `song_rank` - 排行榜歌曲表
- May extend to other tables with `nas_url_path` field

## Steps
1. Query database for records where `nas_url_path LIKE '% %'`
2. For each record:
   - Generate new path: `strings.ReplaceAll(nasUrlPath, " ", "")`
   - Update database: `UPDATE table SET nas_url_path = newPath WHERE id = ?`
3. Report: updated count, remaining count with spaces

## Example
Before: `rank/kugou/top100/白小白 - 人生路漫漫.mp3`
After: `rank/kugou/top100/白小白-人生路漫漫.mp3`

## Database Connection
```go
dsn := "root:password@tcp(127.0.0.1:3306)/music?charset=utf8mb4&parseTime=True&loc=Local"
```

## Related
- Run `/clean-music-files` first to rename actual files
- Then run this to update database paths