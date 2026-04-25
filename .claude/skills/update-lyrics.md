# Update Lyrics

Batch update lyrics in `song_rank` table by reading from .lrc files in folder.

## Usage
User invokes: `/update-lyrics`

## Parameters
- folderPath: Path containing mp3 and lrc files
- dbConfig: Database connection (from config.yaml)

## Steps
1. Scan folder for mp3/wav files, build file index
2. Query `song_rank` table for records with empty `lyric`
3. For each song:
   - Get singer name from `singer_rank` table or `full_name_singer` field
   - Match file by singer+songname key (handle multi-singer with `、`)
   - Read corresponding .lrc file
   - Update `lyric` field in database

## File Matching Logic
- Handle multi-singer: try first singer if full name fails
- Remove brackets from song name for alternative matching
- File format: `歌手-歌名.lrc` or `歌手-歌名.mp3`

## LRC Reading
```go
func readLrcFile(path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    if utf8.Valid(content) {
        result = string(content)
    } else {
        // GBK decode
        decoder := simplifiedchinese.GBK.NewDecoder()
        decoded, _, err := transform.Bytes(decoder, content)
    }
    return cleanLrcContent(result), nil
}
```

## Clean LRC Content
Remove non-standard tags: `[awlrc`, `[krc`, `[qlrc`
Remove UTF-8 BOM if present

## Example Output
- Found 100 songs with empty lyrics
- Updated 100 songs
- Failed: 0

## Related
- Run `/clean-music-files` first if files have spaces
- Run `/fix-nas-url-path` to fix database paths