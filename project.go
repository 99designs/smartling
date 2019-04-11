package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/99designs/smartling/gc"

	"github.com/99designs/api-sdk-go"
	"github.com/codegangsta/cli"
)

var ProjectCommand = cli.Command{
	Name:  "project",
	Usage: "manage local project files",
	Before: func(c *cli.Context) error {
		err := cmdBefore(c)
		if err != nil {
			return nil
		}

		return loadProjectErr
	},
	Subcommands: []cli.Command{
		projectFilesCommand,
		projectStatusCommand,
		projectPullCommand,
		projectPushCommand,
		projectGCCommand,
	},
}

var prefixFlag = cli.StringFlag{
	Name:   "prefix",
	Usage:  "Prefix to use for uploaded file names",
	EnvVar: "SMARTLING_PREFIX",
}

func fetchRemoteFileList() stringSlice {
	files := stringSlice{}
	listFiles, err := client.List(smartling.FilesListRequest{})
	logAndQuitIfError(err)

	for _, fs := range listFiles.Items {
		files = append(files, fs.FileURI)
	}

	return files
}

var remoteFileList = stringSlice{}
var remoteFileListFetched = false

func getRemoteFileList() stringSlice {
	if !remoteFileListFetched {
		remoteFileList = fetchRemoteFileList()
	}

	return remoteFileList
}

func fetchLocales() []string {
	ll := []string{}
	locales, err := client.Locales()
	logAndQuitIfError(err)
	for _, l := range locales {
		ll = append(ll, l.LocaleID)
	}

	return ll
}

var projectFilesCommand = cli.Command{
	Name:  "files",
	Usage: "lists the local files",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: files")
		}

		for _, projectFilepath := range ProjectConfig.Files() {
			fmt.Println(projectFilepath)
		}
	},
}

var projectStatusCommand = cli.Command{
	Name:  "status",
	Usage: "show the status of the project's remote files",
	Flags: []cli.Flag{
		prefixFlag,
		cli.BoolFlag{
			Name:  "awaiting-auth",
			Usage: "Output the number of strings Awaiting Authorization",
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) > 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: status")
		}

		prefix := prefixOrGitPrefix(c.String("prefix"))
		locales := fetchLocales()
		statuses := GetProjectStatus(prefix, locales)

		if c.Bool("awaiting-auth") {
			fmt.Println(statuses.AwaitingAuthorizationCount())
		} else {
			fmt.Print("\n")
			PrintProjectStatusTable(statuses, locales)
			fmt.Print("\n")
			fmt.Printf("Awaiting Authorization: %4d\n", statuses.AwaitingAuthorizationCount())
			fmt.Printf("Total:                  %4d\n", statuses.TotalStringsCount())
		}

	},
}

var projectPullCommand = cli.Command{
	Name:  "pull",
	Usage: "translate local project files using Smartling as a translation memory",
	Flags: []cli.Flag{
		prefixFlag,
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) > 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: pull")
		}

		prefix := prefixOrGitPrefix(c.String("prefix"))

		pullAllProjectFiles(prefix)
	},
}

func pullAllProjectFiles(prefix string) {
	locales, err := client.Locales()
	logAndQuitIfError(err)

	// do this first to cache result and prevent races in the goroutines
	_ = getRemoteFileList()

	var wg sync.WaitGroup
	for _, projectFilepath := range ProjectConfig.Files() {
		for _, l := range locales {
			wg.Add(1)
			go func(locale, projectFilepath string) {
				defer wg.Done()

				pullProjectFile(projectFilepath, locale, prefix)
			}(l.LocaleID, projectFilepath)
		}
	}
	wg.Wait()
}

func pullProjectFile(projectFilepath, locale, prefix string) {
	hit, b, err := translateProjectFile(projectFilepath, locale, prefix)
	logAndQuitIfError(err)

	fp := localPullFilePath(projectFilepath, locale)
	cached := ""
	if hit {
		cached = "(using cache)"
	}
	err = ioutil.WriteFile(fp, b, 0644)
	logAndQuitIfError(err)
	fmt.Println("Wrote", fp, cached)
}

func cleanPrefix(s string) string {
	s = path.Clean("/" + s)
	if s == "/" {
		return ""
	}
	return s
}

func prefixOrGitPrefix(prefix string) string {
	if prefix == "" {
		prefix = pushPrefix()
	}

	prefix = cleanPrefix(prefix)

	if prefix != "" {
		log.Println("Using prefix", prefix)
	}
	return prefix
}

var projectPushCommand = cli.Command{
	Name:  "push",
	Usage: "upload local project files that contain untranslated strings",
	Flags: []cli.Flag{
		prefixFlag,
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) > 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: push")
		}

		prefix := prefixOrGitPrefix(c.String("prefix"))

		pushAllProjectFiles(prefix)
	},
}

// if prefix is empty, don't append the hash also
func projectFileRemoteName(projectFilepath, prefix string) string {
	remoteFile := projectFilepath
	if prefix != "" {
		remoteFile = fmt.Sprintf("%s/%s/%s", prefix, projectFileHash(projectFilepath), projectFilepath)
	}

	return path.Clean("/" + remoteFile)
}

func readFile(projectFilepath string) []byte {
	f, err := ioutil.ReadFile(projectFilepath)
	logAndQuitIfError(err)
	return f
}

func pushProjectFile(projectFilepath, prefix string) string {
	remoteFile := projectFileRemoteName(projectFilepath, prefix)

	req := &smartling.FileUploadRequest{
		FileURIRequest: smartling.FileURIRequest{FileURI: remoteFile},
		FileType:       filetypeForProjectFile(projectFilepath),
		File:           readFile(projectFilepath),
	}
	req.Smartling.Directives = ProjectConfig.ParserConfig
	_, err := client.Upload(req)
	logAndQuitIfError(err)

	fmt.Println("Uploaded", remoteFile)
	return remoteFile
}

func pushProjectFileIfNotExists(projectFilepath, prefix string) (string, bool) {
	remoteFiles := getRemoteFileList()
	remoteFileName := projectFileRemoteName(projectFilepath, prefix)

	if prefix != "" && remoteFiles.contains(remoteFileName) {
		return remoteFileName, false
	}

	return pushProjectFile(projectFilepath, prefix), true
}

func pushAllProjectFiles(prefix string) {
	pushedFilesCount := 0

	// do this first to cache result and prevent races in the goroutines
	_ = getRemoteFileList()

	var wg sync.WaitGroup
	for _, projectFilepath := range ProjectConfig.Files() {
		wg.Add(1)
		go func(projectFilepath string) {
			defer wg.Done()
			_, pushed := pushProjectFileIfNotExists(projectFilepath, prefix)
			if pushed {
				pushedFilesCount++
			}
		}(projectFilepath)
	}
	wg.Wait()

	if pushedFilesCount == 0 {
		fmt.Println("Nothing to do")
	}
}

func filetypeForProjectFile(projectFilepath string) smartling.FileType {
	ft := smartling.GetFileTypeByExtension(path.Ext(projectFilepath))
	if ft == "" {
		ft = ProjectConfig.FileType
	}
	if ft == "" {
		log.Panicln("Can't determine file type for " + projectFilepath)
	}

	return ft
}

type FilenameParts struct {
	Path           string
	Base           string
	Dir            string
	Ext            string
	PathWithoutExt string
	Locale         string
}

func localRelativeFilePath(remotepath string) string {
	fp, err := filepath.Rel(".", path.Join(ProjectConfig.path, remotepath))
	logAndQuitIfError(err)
	return fp
}

func localPullFilePath(p, locale string) string {
	parts := FilenameParts{
		Path:   p,
		Dir:    path.Dir(p),
		Base:   path.Base(p),
		Ext:    path.Ext(p),
		Locale: locale,
	}

	dt := defaultPullDestination
	if dt != "" {
		dt = ProjectConfig.PullFilePath
	}

	out := bytes.NewBufferString("")
	tmpl := template.New("name")
	tmpl.Funcs(template.FuncMap{
		"TrimSuffix": strings.TrimSuffix,
		"Truncate": func(s string, n int) string {
			return s[:n]
		},
	})
	_, err := tmpl.Parse(dt)
	logAndQuitIfError(err)

	err = tmpl.Execute(out, parts)
	logAndQuitIfError(err)

	return localRelativeFilePath(out.String())
}

var projectGCCommand = cli.Command{
	Name:  "gc",
	Usage: "Collect garbage strings",
	Flags: []cli.Flag{
		prefixFlag,
	},
	Subcommands: []cli.Command{
		{
			Name:  "branch",
			Usage: "Collect garbage strings for the current branch",
			Action: func(c *cli.Context) error {
				return gc.Branch()
			},
		},
	},
	Action: func(c *cli.Context) error {
		return gc.Project()
	},
}
