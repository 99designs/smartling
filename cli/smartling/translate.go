package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/99designs/smartling"
)

var cachePath = findCachePath()

func findCachePath() string {
	var cachePath string
	usr, err := user.Current()
	if err == nil && usr.HomeDir != "" {
		cachePath = filepath.Join(usr.HomeDir, ".smartling", "cache")
	} else {
		log.Panicln("Can't locate a cache directory")
	}

	_ = os.MkdirAll(cachePath, 0755)

	return cachePath
}

func mustSha1File(filePath string) string {
	file, err := os.Open(filePath)
	panicIfErr(err)

	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	panicIfErr(err)

	b := []byte{}
	return hex.EncodeToString(hash.Sum(b))
}

type stringSlice []string

func (ss stringSlice) contains(s string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}

var uploadMutex = sync.Mutex{}
var tempFilesUploaded = stringSlice{}

func uploadAsTempFile(localpath string, filetype smartling.FileType, parserConfig map[string]string) (remotepath string, err error) {
	uploadMutex.Lock()
	defer uploadMutex.Unlock()

	tmppath := "/tmp/" + mustSha1File(localpath)
	if tempFilesUploaded.contains(tmppath) {
		return tmppath, nil
	}

	// upload
	_, err = client.Upload(localpath, &smartling.UploadRequest{
		FileUri:      tmppath,
		FileType:     filetype,
		ParserConfig: parserConfig,
	})
	if err != nil {
		return
	}

	tempFilesUploaded = append(tempFilesUploaded, tmppath)

	return tmppath, nil
}

func cacheHash(locale, localpath string, filetype smartling.FileType, parserConfig map[string]string) string {
	file, err := os.Open(localpath)
	panicIfErr(err)
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	panicIfErr(err)

	_, err = hash.Write([]byte(fmt.Sprintf("%#v%#v%#v", locale, filetype, parserConfig)))
	panicIfErr(err)

	b := []byte{}
	return hex.EncodeToString(hash.Sum(b))
}

func translateViaCache(locale, localpath string, filetype smartling.FileType, parserConfig map[string]string) (hit bool, b []byte, err error, ch string) {

	ch = cacheHash(locale, localpath, filetype, parserConfig)
	cacheFilePath := filepath.Join(cachePath, ch)

	// get cached file
	if cacheFile, err := os.Open(cacheFilePath); err == nil {
		if cfStat, err := cacheFile.Stat(); err == nil {
			if cfStat.ModTime().After(time.Now().Add(time.Duration(-cacheTtl))) {
				if b, err = ioutil.ReadFile(cacheFilePath); err == nil {
					return true, b, nil, ch // return the cached data
				}
			}
		}
	}

	// translate
	b, err = translate(locale, localpath, filetype, parserConfig)
	if err != nil {
		return
	}

	// write cache
	err = ioutil.WriteFile(cacheFilePath, b, 0644)
	if err != nil {
		return
	}

	return
}

func translate(locale, localpath string, filetype smartling.FileType, parserConfig map[string]string) (b []byte, err error) {
	tmppath, err := uploadAsTempFile(localpath, filetype, parserConfig)
	if err != nil {
		return
	}

	// download
	b, err = client.Get(&smartling.GetRequest{
		FileUri: tmppath,
		Locale:  locale,
	})
	if err != nil {
		return
	}

	return
}

func cleanupTempFiles() {
	var wg sync.WaitGroup
	for _, f := range tempFilesUploaded {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.Delete(f)
			if err != nil {
				log.Println(err.Error())
			}
		}()

	}
	wg.Wait()
}
