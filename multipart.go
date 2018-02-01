package vkapi

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func (client *VKClient) getMultipartRequest(link string, fileType string, files []string) (*http.Request, error) {
	if len(files) > 6 {
		return nil, errors.New("you can't upload more than 6 files")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	idx := 1
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return nil, err
		}

		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile("file"+strconv.Itoa(idx), fi.Name())
		if err != nil {
			return nil, err
		}
		io.Copy(part, file)
		file.Close()
		idx++
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", link, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}
