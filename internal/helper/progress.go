package helper

import (
	"fmt"
	"math"
	"sync"
)

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	Store     *sync.Map
	Username  string
	Filename  string
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	pr.SaveToRedis()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("File upload in progress: %d - %d%%\n", pr.BytesRead, 30+int(math.Floor(float64(pr.BytesRead)*70.0/float64(pr.TotalSize))))
}

func (pr *Progress) SaveToRedis() {
	pr.Store.Store(pr.Username+"upload-profile-picture-progress", map[string]interface{}{
		"fileFormat":          "default",
		"fileName":            pr.Filename,
		"currentPercent":      30 + int(math.Floor(float64(pr.BytesRead)*70.0/float64(pr.TotalSize))),
		"currentUploadedSize": FormatFileSize(float64(pr.BytesRead), 1024.0),
		"totalSize":           FormatFileSize(float64(pr.TotalSize), 1024.0),
		"uploading":           true,
		"finish":              false,
	})
}
