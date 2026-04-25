# Clean Music Files

Remove spaces from music file names (mp3, wav, lrc) in a folder. This prevents playback issues when accessing files via URL.

## Usage
User invokes: `/clean-music-files`

## Steps
1. Ask user for folder path if not provided
2. Scan folder for files with extensions: `.mp3`, `.wav`, `.lrc`
3. For each file with spaces in name:
   - Generate new name by removing all spaces
   - Rename file using `os.Rename(oldPath, newPath)`
4. Report: renamed count, skipped count (no spaces), failed count

## Example
Input folder: `C:\Music\lx-download`
Output:
- 200 files renamed (100 mp3 + 100 lrc)
- 0 skipped
- 0 failed

## Code Pattern
```go
newName := strings.ReplaceAll(originalName, " ", "")
os.Rename(filepath.Join(folder, originalName), filepath.Join(folder, newName))
```

## Related
- Also update database `nas_url_path` field to remove spaces
- See `/fix-nas-url-path` skill for database cleanup