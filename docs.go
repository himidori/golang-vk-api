package vkapi

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

type DocAttachment struct {
	ID         int    `json:"id"`
	OwnerID    int    `json:"owner_id"`
	Title      string `json:"title"`
	Size       int    `json:"size"`
	Extenstion string `json:"ext"`
	URL        string `json:"url"`
	Date       int64  `json:"date"`
	Type       int    `json:"type"`
	IsLicensed int    `json:"is_licensed"`
}

type Docs struct {
	Count     int              `json:"count"`
	Documents []*DocAttachment `json:"items"`
}

func (client *VKClient) docGetWallUploadServer(groupID int) (string, error) {
	gidStr := strconv.Itoa(groupID)
	if gidStr[0] == '-' {
		gidStr = gidStr[1:]
	}

	params := url.Values{}
	params.Set("group_id", gidStr)
	resp, err := client.MakeRequest("docs.getWallUploadServer", params)
	if err != nil {
		return "", err
	}

	m := map[string]string{}
	json.Unmarshal(resp.Response, &m)

	return m["upload_url"], err
}

func (client *VKClient) docWallUpload(groupID int, fileName string) (string, error) {
	uploadURL, err := client.docGetWallUploadServer(groupID)
	if err != nil {
		return "", err
	}

	req, err := client.getDocMultipartReq(uploadURL, fileName)
	if err != nil {
		return "", err
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	m := map[string]string{}
	json.Unmarshal(body, &m)

	return m["file"], nil
}

func (client *VKClient) UploadGroupWallDoc(groupID int, fileName string) (*DocAttachment, error) {
	file, err := client.docWallUpload(groupID, fileName)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("file", file)

	resp, err := client.MakeRequest("docs.save", params)
	if err != nil {
		return nil, err
	}

	var docs []*DocAttachment
	json.Unmarshal(resp.Response, &docs)

	return docs[0], err
}

func (client *VKClient) GetDocsString(docs []*DocAttachment) string {
	s := []string{}

	for _, d := range docs {
		s = append(s, "doc"+strconv.Itoa(d.OwnerID)+"_"+strconv.Itoa(d.ID))
	}

	return strings.Join(s, ",")
}

func (client *VKClient) DocsSearch(query string, count int, params url.Values) (int, []*DocAttachment, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("q", query)
	params.Set("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("docs.search", params)
	if err != nil {
		return 0, nil, err
	}

	var docs *Docs
	json.Unmarshal(resp.Response, &docs)
	return docs.Count, docs.Documents, nil
}
