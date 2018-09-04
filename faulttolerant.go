package smartling

import (
	"log"
	"net"
	"strings"
	"time"

	smartlingNew "github.com/Smartling/api-sdk-go"
)

func isResourceLockedError(err error) bool {
	if err != nil {
		if sErr, ok := err.(SmartlingResponse); ok {
			return sErr.IsResourceLockedError()
		}
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

// FaultTolerantClient decorates a Client and retries
// requests when Smartling returns with an error
type FaultTolerantClient struct {
	*smartlingNew.Client
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

func (c *FaultTolerantClient) Upload(req *smartlingNew.FileUploadRequest) (r *smartlingNew.FileUploadResult, err error) {
	c.execWithRetry(func() error {
		r, err = c.Client.UploadFile(c.ProjectID, *req)

		return err
	})
	return
}

func (c *FaultTolerantClient) Get(req *GetRequest) (b []byte, err error) {
	c.execWithRetry(func() error {
		// c.Client.DownloadFile
		// b, err = c.Client.Get(req)
		// return err

		return nil
	})
	return
}

func (c *FaultTolerantClient) List(req ListRequest) (ff *smartlingNew.FilesList, err error) {
	c.execWithRetry(func() error {
		ff, err = c.Client.ListFiles(c.ProjectID, smartlingNew.FilesListRequest{})
		return err
	})
	return
}

func (c *FaultTolerantClient) Status(fileUri, locale string) (f FileStatus, err error) {
	c.execWithRetry(func() error {
		// f, err = c.Client.Status(fileUri, locale)
		// return err
		return nil

	})
	return
}

func (c *FaultTolerantClient) Rename(oldFileUri, newFileUri string) (err error) {
	c.execWithRetry(func() error {
		// return c.Client.Rename(oldFileUri, newFileUri)
		return nil

	})
	return
}

func (c *FaultTolerantClient) Delete(fileUri string) (err error) {
	c.execWithRetry(func() error {
		// return c.Client.Delete(fileUri)
		return nil

	})
	return
}

func (c *FaultTolerantClient) LastModified(req LastModifiedRequest) (ii []LastModifiedItem, err error) {
	c.execWithRetry(func() error {
		// ii, err = c.Client.LastModified(req)
		// return err
		return nil

	})
	return
}
