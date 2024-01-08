package fileserver

import "time"

type File struct {
	Content      []byte    // File content
	FilePath     string    // Storage path of the file
	OriginalName string    // Original name of the file
	Size         int64     // Size of the file in bytes
	MIMEType     string    // MIME type of the file
	UploadedBy   string    // Identifier of the user who uploaded the file
	UploadedAt   time.Time // Timestamp of when the file was uploaded
	DataID       int       // Optional: Link to DataFiles ID
}

type DataFiles struct {
	DataID int
	Name   string
}

type Versions struct {
	DataID    int
	Version   int
	FilePath  string
	UpdatedAt time.Time
	VersionID string
}
