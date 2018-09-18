// Package smartling is a client implementation of the Smartling Translation
// API v2 as documented at https://help.smartling.com/v1.0/reference
package smartling

import (
	"net/http"
	"time"
)

type (
	// LogFunction represents abstract logger function interface which
	// can be used for setting up logging of library actions.
	LogFunction func(format string, args ...interface{})
)

var (
	// Version is a API SDK version, sent in User-Agent header.
	Version = "1.0"

	// DefaultBaseURL specifies base URL which will be used for calls unless
	// other is specified in the Client struct.
	DefaultBaseURL = "https://api.smartling.com"

	// DefaultHTTPClient specifies default HTTP client which will be used
	// for calls unless other is specified in the Client struct.
	DefaultHTTPClient = http.Client{Timeout: 60 * time.Second}

	// DefaultUserAgent is a string that will be sent in User-Agent header.
	DefaultUserAgent = "smartling-api-sdk-go"
)

// Client represents Smartling API client.
type Client struct {
	BaseURL     string
	Credentials *Credentials

	HTTP *http.Client

	Logger struct {
		Infof  LogFunction
		Debugf LogFunction
	}

	UserAgent string
}

// NewClient returns new Smartling API client with specified authentication
// data.
func NewClient(userID string, tokenSecret string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Credentials: &Credentials{
			UserID: userID,
			Secret: tokenSecret,
		},

		HTTP: &DefaultHTTPClient,

		Logger: struct {
			Infof  LogFunction
			Debugf LogFunction
		}{
			Infof:  func(string, ...interface{}) {},
			Debugf: func(string, ...interface{}) {},
		},

		UserAgent: DefaultUserAgent + "/" + Version,
	}
}

// SetInfoLogger sets logger function which will be called for logging
// informational messages like progress of file download and so on.
func (client *Client) SetInfoLogger(logger LogFunction) *Client {
	client.Logger.Infof = logger

	return client
}

// SetDebugLogger sets logger function which will be called for logging
// internal information like HTTP requests and their responses.
func (client *Client) SetDebugLogger(logger LogFunction) *Client {
	client.Logger.Debugf = logger

	return client
}
