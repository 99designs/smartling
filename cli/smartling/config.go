package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/smartling"
	"gopkg.in/yaml.v2"
)

var ProjectConfig *Config
var loadProjectErr error

var defaultPullDestination = "{{ TrimSuffix .Path .Ext }}.{{.Locale}}{{.Ext}}"

type Config struct {
	path         string
	ApiKey       string             `yaml:"ApiKey"`
	ProjectId    string             `yaml:"ProjectId"`
	CacheMaxAge  string             `yaml:"CacheMaxAge"`
	FileGlobs    []string           `yaml:"Files"`
	FileType     smartling.FileType `yaml:"FileType"`
	ParserConfig map[string]string  `yaml:"ParserConfig"`
	PullFilePath string             `yaml:"PullFilePath"`
	hasGlobbed   bool
	files        []string
}

var ErrConfigFileNotExist = errors.New("smartling.yml not found")

func (c *Config) Files() []string {
	if !c.hasGlobbed {
		for _, g := range c.FileGlobs {
			ff, err := filepath.Glob(g)
			panicIfErr(err)
			c.files = append(c.files, ff...)
		}
	}

	return c.files
}

func (c *Config) cacheMaxAge() time.Duration {
	if c.CacheMaxAge != "" {
		d, err := time.ParseDuration(c.CacheMaxAge)
		panicIfErr(err)
		return d
	}

	return time.Duration(4 * time.Hour)
}

func gitBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	return strings.TrimSpace(out.String())
}

func pushPrefix() string {
	// prefer branch
	b := gitBranch()
	if b == "master" {
		return "/"
	} else if b != "" {
		return "/branch/" + b
	}

	// fall back to username
	u, err := user.Current()
	panicIfErr(err)
	if u.Username == "" {
		log.Panicln("Can't find a prefix")
	}

	return "/user/" + u.Username
}

func loadConfig(configfilepath string) (*Config, error) {
	if _, err := os.Stat(configfilepath); err != nil {
		return nil, ErrConfigFileNotExist
	}

	b, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	c.path = filepath.Dir(configfilepath)

	return &c, nil
}
