package smartling

import "time"

func isResourceLockedError(err error) bool {
	if err != nil {
		if sErr, ok := err.(SmartlingResponse); ok {
			return sErr.IsResourceLockedError()
		}
	}
	return false
}

// FaultTolerantClient decorates a Client and retries
// requests when Smartling returns with a Resource Locked error
type FaultTolerantClient struct {
	*Client
	RetriesWhenResourceLocked int
}

func (c *FaultTolerantClient) execWithRetry(f func() error) {
	err := f()

	retries := c.RetriesWhenResourceLocked
	backoff := 1 * time.Second

	for isResourceLockedError(err) && retries > 0 {
		// log.Println("Resource locked, retrying")
		time.Sleep(backoff)
		f()
		retries--
		backoff = backoff * 2
	}
}

func (c *FaultTolerantClient) Upload(localFilePath string, req *UploadRequest) (r *UploadResponse, err error) {
	c.execWithRetry(func() error {
		r, err = c.Client.Upload(localFilePath, req)
		return err
	})
	return
}

func (c *FaultTolerantClient) Get(req *GetRequest) (b []byte, err error) {
	c.execWithRetry(func() error {
		b, err = c.Client.Get(req)
		return err
	})
	return
}

func (c *FaultTolerantClient) List(req ListRequest) (ff []File, err error) {
	c.execWithRetry(func() error {
		ff, err = c.Client.List(req)
		return err
	})
	return
}

func (c *FaultTolerantClient) Status(fileUri, locale string) (f File, err error) {
	c.execWithRetry(func() error {
		f, err = c.Client.Status(fileUri, locale)
		return err
	})
	return
}

func (c *FaultTolerantClient) Rename(oldFileUri, newFileUri string) (err error) {
	c.execWithRetry(func() error {
		return c.Client.Rename(oldFileUri, newFileUri)
	})
	return
}

func (c *FaultTolerantClient) Delete(fileUri string) (err error) {
	c.execWithRetry(func() error {
		return c.Client.Delete(fileUri)
	})
	return
}

func (c *FaultTolerantClient) LastModified(req LastModifiedRequest) (ii []LastModifiedItem, err error) {
	c.execWithRetry(func() error {
		ii, err = c.Client.LastModified(req)
		return err
	})
	return
}
