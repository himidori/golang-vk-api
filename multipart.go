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

func (client *VKClient) getPhotoMultipartReq(link string, files []string) (*http.Request, error) {
	size, err := GetFilesSizeMB(files)
	if err != nil {
		return nil, err
	}
	if size > 50 {
		return nil, errors.New("summary photos size can't be higher than 50mb")
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

	err = writer.Close()
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

func (client *VKClient) getPhotoByLinkMultipartReq(link string, urlFile string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	respGetFile, err := http.Get(urlFile)
	if err != nil {
		return nil, err
	}
	defer respGetFile.Body.Close()

	part, err := writer.CreateFormFile("file", "file.png")
	if err != nil {
		return nil, err
	}
	io.Copy(part, respGetFile.Body)

	err = writer.Close()
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

func (client *VKClient) getDocMultipartReq(link string, fileName string) (*http.Request, error) {
	size, err := GetFilesSizeMB([]string{fileName})
	if err != nil {
		return nil, err
	}
	if size > 200 {
		return nil, errors.New("document size can't be higher than 200mb")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	err = writer.Close()
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
