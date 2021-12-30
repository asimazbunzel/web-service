package mesa

import (
   "fmt"
   "io"
   "os"
)


// function retrieve from this post:
// https://stackoverflow.com/questions/17863821/how-to-read-last-lines-from-a-big-file-with-go-every-10-secs
func GetLastLineWithSeek (filepath string) string {

   fileHandle, err := os.Open(filepath)

    if err != nil {
       return "file open error"
    }

    defer fileHandle.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := fileHandle.Stat()
    filesize := stat.Size()
    for {
        cursor -= 1
        fileHandle.Seek(cursor, io.SeekEnd)

        char := make([]byte, 1)
        fileHandle.Read(char)

        if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
            break
        }

        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

        if cursor == -filesize { // stop if we are at the begining
            break
        }
    }

    return line
}

