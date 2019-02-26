package upload

import (
	"net/http"
	. "mime/multipart"
	"github.com/pkg/errors"
	"os"
	"io"
	"fmt"
	"io/ioutil"
	"github.com/oliamb/cutter"
	"image"
	"strconv"
	"image/jpeg"
	"time"
)

type Uploader struct {
	err error
}

type ImagePoint struct {
	X int
	Y int
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

	contentType := header.Header.Get("Content-Type")
	u.IsImageFile(contentType)

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

func (u *Uploader) DivideByFour(imageName string, folderPath string) map[int]string {
	if u.err != nil || imageName == "" || folderPath == "" {
		return nil
	}

	file, err := os.Open(folderPath + imageName)
	if err != nil {
		u.err = err
		return nil
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		u.err = err
		return nil
	}
	file.Close()

	imgSize  := img.Bounds().Size()
	newSizeX := imgSize.X / 2
	newSizeY := imgSize.Y / 2

	pointers := map[int]ImagePoint{
		1: {X: 0, Y: 0},
		2: {X: newSizeX, Y: 0},
		3: {X: newSizeX, Y: newSizeY},
		4: {X: 0, Y: newSizeY},
	}

	images := map[int]string{}

	for index, point := range pointers {
		croppedImg, err := cutter.Crop(img, cutter.Config{
			Width:  newSizeX,
			Height: newSizeY,
			Anchor: image.Point{point.X, point.Y},
		})

		imgName := strconv.FormatInt(time.Now().Unix(), 10) + "_" + strconv.Itoa(index) + ".jpg"
		out, err := os.Create(folderPath + imgName)
		if err != nil {
			u.err = err
		}
		defer out.Close()

		jpeg.Encode(out, croppedImg, nil)

		images[index] = imgName
	}

	return images
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