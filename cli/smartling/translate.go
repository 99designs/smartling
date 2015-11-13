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

type stringSlice []string

func (ss stringSlice) contains(s string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}

func projectFileHash(projectFilepath string) string {
	return hash(
		localRelativeFilePath(projectFilepath),
		filetypeForProjectFile(projectFilepath),
		ProjectConfig.ParserConfig,
	)
}

func hash(localpath string, filetype smartling.FileType, parserConfig map[string]string) string {
	file, err := os.Open(localpath)
	logAndQuitIfError(err)
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	logAndQuitIfError(err)

	_, err = hash.Write([]byte(fmt.Sprintf("%#v%#v", filetype, parserConfig)))
	logAndQuitIfError(err)

	b := []byte{}
	return hex.EncodeToString(hash.Sum(b))
}

func translateProjectFile(projectFilepath, locale, prefix string) (hit bool, b []byte, err error, h string) {
	localpath := localRelativeFilePath(projectFilepath)
	filetype := filetypeForProjectFile(projectFilepath)

	h = hash(localpath, filetype, ProjectConfig.ParserConfig)
	cacheFilePath := filepath.Join(cachePath, locale, h)

	// check cache
	hit, b = getCachedTranslations(cacheFilePath)
	if hit {
		return hit, b, nil, h
	}

	// translate
	b, err = translateViaSmartling(projectFilepath, locale, prefix, filetype, ProjectConfig.ParserConfig)
	if err != nil {
		return
	}

	// write to cache
	err = ioutil.WriteFile(cacheFilePath, b, 0644)
	if err != nil {
		return
	}

	return
}

func getCachedTranslations(cacheFilePath string) (hit bool, b []byte) {
	if cacheFile, err := os.Open(cacheFilePath); err == nil {
		if cfStat, err := cacheFile.Stat(); err == nil {
			if time.Now().Sub(cfStat.ModTime()) < ProjectConfig.cacheMaxAge() {
				if b, err = ioutil.ReadFile(cacheFilePath); err == nil {
					return true, b // return the cached data
				}
			}
		}
	}

	return
}

func translateViaSmartling(projectFilepath, locale, prefix string, filetype smartling.FileType, parserConfig map[string]string) (b []byte, err error) {
	remotePath := projectFileRemoteName(projectFilepath, prefix)

	b, err = client.Get(&smartling.GetRequest{
		FileUri: remotePath,
		Locale:  locale,
	})

	// file might not exist, try pushing first
	if err != nil {

		remotePath := pushProjectFile(projectFilepath, prefix)

		b, err = client.Get(&smartling.GetRequest{
			FileUri: remotePath,
			Locale:  locale,
		})
	}

	return
}
