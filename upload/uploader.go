package upload

import (
	"net/http"
	. "mime/multipart"
	"github.com/pkg/errors"
	"os"
	"io"
	"fmt"
	"io/ioutil"
)

type Uploader struct {
	err error
}

func (u *Uploader) GetFile(r *http.Request, fileName string) (File, *FileHeader) {
	if fileName == "" {
		u.err = errors.New("fileName cannot be empty")
		return nil, nil
	}

	file, header, err := r.FormFile(fileName)
	if err != nil {
		u.err = errors.Wrap(err, "unable to read file")
	}
	defer file.Close()

	if u.err != nil {
		return nil, nil
	}

	return file, header
}

func (u *Uploader) MkDir(uploadDir string) {
	if u.err != nil || uploadDir == "" {
		return
	}

	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		u.err = errors.Wrap(err, "unable to create folder")
	}
}

func (u *Uploader) Copy(dst io.Writer, src io.Reader) {
	if u.err != nil {
		return
	}

	if _, err := io.Copy(dst, src); err != nil {
		u.err = err
	}
}

func (u *Uploader) CreateFile(filePath string) *os.File {
	if u.err != nil || filePath == "" {
		return nil
	}

	resultFile, err := os.Create(filePath)
	if err != nil {
		u.err = err
	}

	return resultFile
}

func (u *Uploader) MoveFile(file File, filePath string) {
	if u.err != nil || filePath == "" {
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		u.err = err
	}

	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		u.err = err
	}
}

func (u *Uploader) IsImageFile(contentType string) bool {
	if u.err != nil || contentType == "" {
		return false
	}

	if !(contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif") {
		u.err = errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
		return false
	}

	return true
}

func (u *Uploader) GetError() error {
	return u.err
}