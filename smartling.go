// Package smartling is a client implementation of the Smartling Translation API as documented at
// https://docs.smartling.com/display/docs/Smartling+Translation+API
package smartling

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	baseUrlProd    = "https://api.smartling.com/v1"
	baseUrlSandbox = "https://sandbox-api.smartling.com/v1"
)

type Client struct {
	BaseUrl    string
	ApiKey     string
	ProjectId  string
	httpClient *http.Client
}

func (c *Client) addCredentials(v *url.Values) {
	v.Set("apiKey", c.ApiKey)
	v.Set("projectId", c.ProjectId)
}

func (c *Client) doRequestAndUnmarshalData(url string, req interface{}, res interface{}) error {
	v, err := query.Values(req)
	if err != nil {
		return err
	}
	c.addCredentials(&v)

	httpResponse, err := c.httpClient.PostForm(c.BaseUrl+url, v)
	if err != nil {
		return err
	}

	return unmarshalResponseData(httpResponse, &res)
}

func NewClient(apiKey string, projectId string) *Client {
	return &Client{
		BaseUrl:    baseUrlProd,
		ApiKey:     apiKey,
		ProjectId:  projectId,
		httpClient: http.DefaultClient,
	}
}

func NewSandboxClient(apiKey string, projectId string) *Client {
	return &Client{
		BaseUrl:    baseUrlSandbox,
		ApiKey:     apiKey,
		ProjectId:  projectId,
		httpClient: http.DefaultClient,
	}
}

type FileType string

const (
	Android        FileType = "android"
	Ios            FileType = "ios"
	Gettext        FileType = "gettext"
	Html           FileType = "html"
	JavaProperties FileType = "javaProperties"
	Yaml           FileType = "yaml"
	Xliff          FileType = "xliff"
	Xml            FileType = "xml"
	Json           FileType = "json"
	Docx           FileType = "docx"
	Pptx           FileType = "pptx"
	Xlsx           FileType = "xlsx"
	Idml           FileType = "idml"
	Qt             FileType = "qt"
	Resx           FileType = "resx"
	Plaintext      FileType = "plaintext"
	Csv            FileType = "csv"
	Stringsdict    FileType = "stringsdict"
)

var extensionToType = map[string]FileType{
	".yml":         Yaml,
	".yaml":        Yaml,
	".html":        Html,
	".htm":         Html,
	".xlf":         Xliff,
	".xliff":       Xliff,
	".json":        Json,
	".docx":        Docx,
	".pptx":        Pptx,
	".xlsx":        Xlsx,
	".txt":         Plaintext,
	".ts":          Qt,
	".idml":        Idml,
	".resx":        Resx,
	".resw":        Resx,
	".csv":         Csv,
	".stringsdict": Stringsdict,
	".strings":     Ios,
	".po":          Gettext,
	".pot":         Gettext,
	".xml":         Xml,
}

// FileTypeByExtension returns the FileType associated with the file extension ext.
// The extension ext should begin with a leading dot, as in ".html".
// When ext has no associated type, FileTypeByExtension returns "".
func FileTypeByExtension(ext string) FileType {
	t, _ := extensionToType[ext]
	return t
}

type RetrievalType string

const (
	Pending                     RetrievalType = "pending"
	Published                   RetrievalType = "published"
	Pseudo                      RetrievalType = "pseudo"
	ContextMatchingInstrumented RetrievalType = "contextMatchingInstrumented"
)

type smartlingResponseWrapper struct {
	Response SmartlingResponse `json:"response"`
}

type SmartlingResponse struct {
	Code     string          `json:"code"`
	Messages []string        `json:"messages"`
	Data     json.RawMessage `json:"data"`
}

func (sr SmartlingResponse) IsError() bool {
	return sr.Code != "SUCCESS"
}

func (sr SmartlingResponse) IsResourceLockedError() bool {
	return sr.Code == "RESOURCE_LOCKED"
}

func (sr SmartlingResponse) IsValidationError() bool {
	return sr.Code == "VALIDATION_ERROR"
}

func (sr SmartlingResponse) Error() string {
	return fmt.Sprintf("Smartling error: %s: %v", sr.Code, sr.Messages)
}

func (sr SmartlingResponse) IsNotFoundError() bool {
	for _, m := range sr.Messages {
		return strings.Contains(m, "No row with the given identifier exists")
	}

	return false
}

func unmarshalResponseData(r *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	p := smartlingResponseWrapper{}
	err = json.Unmarshal(body, &p)

	if err == nil && p.Response.IsError() {
		return p.Response
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", r.Status)
	}
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	return json.Unmarshal(p.Response.Data, &data)
}
