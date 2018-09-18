package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Smartling/api-sdk-go"
)

func isResourceLockedError(err error) bool {
	if err != nil {
		return strings.Contains(err.Error(), "RESOURCE_LOCKED")
	}
	return false
}

func isNetworkErrClosing(err error) bool {
	if err != nil {
		return strings.Contains(err.Error(), "use of closed network connection")
	}
	return false
}

func isTimeoutError(err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			return netErr.Timeout()
		}
	}
	return false
}

func isRetryableError(err error) bool {
	return isResourceLockedError(err) || isNetworkErrClosing(err) || isTimeoutError(err)
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

// FaultTolerantClient decorates a Client and retries
// requests when Smartling returns with an error
type FaultTolerantClient struct {
	*smartling.Client
	ProjectID      string
	RetriesOnError int
}

func (c *FaultTolerantClient) execWithRetry(f func() error) {
	retries := c.RetriesOnError
	backoff := 1 * time.Second

	err := f()

	for retries > 0 && isRetryableError(err) {
		log.Printf("%s, Retrying...\n", err.Error())

		time.Sleep(backoff)
		err = f()

		retries--
		backoff = backoff * 2
	}
}

func (c *FaultTolerantClient) Upload(req *smartling.FileUploadRequest) (r *smartling.FileUploadResult, err error) {
	c.execWithRetry(func() error {
		r, err = c.Client.UploadFile(c.ProjectID, *req)

		return err
	})
	return
}

func (c *FaultTolerantClient) Download(fileURI string) (b []byte, err error) {
	c.execWithRetry(func() error {
		var r io.Reader
		r, err = c.Client.DownloadFile(c.ProjectID, fileURI)
		b = streamToByte(r)
		return err
	})
	return
}

func (c *FaultTolerantClient) DownloadTranslation(locale string, req smartling.FileDownloadRequest) (b []byte, err error) {
	c.execWithRetry(func() error {
		var r io.Reader
		r, err = c.Client.DownloadTranslation(c.ProjectID, locale, req)
		b = streamToByte(r)
		return err
	})
	return
}

func (c *FaultTolerantClient) List(req smartling.FilesListRequest) (ff *smartling.FilesList, err error) {
	c.execWithRetry(func() error {
		ff, err = c.Client.ListFiles(c.ProjectID, req)
		return err
	})
	return
}

func (c *FaultTolerantClient) Status(fileUri, locale string) (f *smartling.FileStatusExtended, err error) {
	c.execWithRetry(func() error {
		f, err = c.Client.GetFileStatusExtended(c.ProjectID, fileUri, locale)
		return err

	})
	return
}

func (c *FaultTolerantClient) Rename(oldFileUri, newFileUri string) (err error) {
	c.execWithRetry(func() error {
		err = c.Client.RenameFile(c.ProjectID, oldFileUri, newFileUri)
		return err
	})
	return
}

func (c *FaultTolerantClient) Delete(fileUri string) (err error) {
	c.execWithRetry(func() error {
		err = c.Client.DeleteFile(c.ProjectID, fileUri)
		return err
	})
	return
}

func (c *FaultTolerantClient) LastModified(req smartling.FileLastModifiedRequest) (f *smartling.FileLastModifiedLocales, err error) {
	c.execWithRetry(func() error {
		f, err = c.Client.LastModified(c.ProjectID, req)

		return err
	})
	return
}

func (c *FaultTolerantClient) Locales() (tl []smartling.Locale, err error) {
	c.execWithRetry(func() error {
		var pd *smartling.ProjectDetails
		pd, err = c.Client.GetProjectDetails(c.ProjectID)

		if err != nil {
			return err
		}

		tl = pd.TargetLocales

		return nil
	})
	return
}
