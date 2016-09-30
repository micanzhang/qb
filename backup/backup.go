package backup

import "qiniupkg.com/api.v7/kodo"

type BackupProvider interface {
	// upload file
	Put(filepath string, key string) error
	// get file info by file's path or key
	Info(key string) (entry kodo.Entry, err error)
	// download file
	Get(key string, dirpath string) error
	// // remove file
	Remove(key string) error
	// // create one and remove the old one
	// Update()
	// // list file info
	List(prefix string, marker string, limit int) error
}

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrFileNotExists Err = "file not exists"
	ErrNotFound      Err = "not found"
	ErrDuplicated    Err = "duplicated"
)
