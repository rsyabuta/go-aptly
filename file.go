package aptly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/smira/aptly/deb"
)

type FileService struct {
	Client *Client
}

type DebFileCollection struct {
	files []string
}

type DebDirCollection struct {
	dirs []string
}

func (service *FileService) UploadFile(file string) (*DebFileCollection, error) {
	dir, err := service.generateDir(file)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}
	w.Close()

	resp, err := service.Client.Post(fmt.Sprintf("files/%s", dir), w.FormDataContentType(), nil, &b)
	if err != nil {
		return nil, err
	}
	var collection DebFileCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.files)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (service *FileService) generateDir(file string) (string, error) {
	st, err := deb.GetControlFileFromDeb(file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", st["Package"], st["Version"]), nil
}

func (service *FileService) ListDirectories() (*DebDirCollection, error) {
	resp, err := service.Client.Get("files")
	if err != nil {
		return nil, err
	}
	var collection DebDirCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.dirs)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (service *FileService) ListFiles(dir string) (*DebFileCollection, error) {
	resp, err := service.Client.Get(fmt.Sprintf("files/%s", dir))
	if err != nil {
		return nil, err
	}

	var collection DebFileCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.files)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (service *FileService) Delete(item string) error {
	resp, err := service.Client.Delete(fmt.Sprintf("files/%s", item), nil)
	defer resp.Body.Close()
	return err
}
