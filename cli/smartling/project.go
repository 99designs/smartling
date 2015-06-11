package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
	"text/template"

	"github.com/99designs/smartling"
	"github.com/codegangsta/cli"
)

var ProjectCommand = cli.Command{
	Name:  "project",
	Usage: "manage local project files",
	Before: func(c *cli.Context) (err error) {
		return loadProjectErr
	},
	After: func(c *cli.Context) error {
		cleanupTempFiles()
		return nil
	},
	Subcommands: []cli.Command{
		projectStatusCommand,
		projectPullCommand,
		projectPushCommand,
	},
}

func fetchRemoteFileList() stringSlice {
	files := stringSlice{}
	listFiles, err := client.List(smartling.ListRequest{})
	panicIfErr(err)

	for _, fs := range listFiles {
		files = append(files, fs.FileUri)
	}

	return files
}

func fetchLocales() []string {
	ll := []string{}
	locales, err := client.Locales()
	panicIfErr(err)
	for _, l := range locales {
		ll = append(ll, l.Locale)
	}

	return ll
}

var projectStatusCommand = cli.Command{
	Name:        "status",
	Usage:       "show the status of the project's remote files",
	Description: "status [<prefix>]",
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: status [<prefix>]")
		}

		prefix := prefixOrGitPrefix(c.Args().Get(0))

		locales, err := client.Locales()
		panicIfErr(err)

		var wg sync.WaitGroup
		statuses := make(map[string]map[string]smartling.FileStatus)

		remoteFiles := fetchRemoteFileList()

		for _, projectFilepath := range ProjectConfig.Files() {

			prefixedProjectFilepath := filepath.Clean("/" + prefix + "/" + projectFilepath)
			if !remoteFiles.contains(prefixedProjectFilepath) {
				prefixedProjectFilepath = filepath.Clean("/" + projectFilepath)
			}

			for _, l := range locales {
				wg.Add(1)
				go func(remotefile, locale string) {
					defer wg.Done()

					fs, err := client.Status(remotefile, locale)
					panicIfErr(err)

					_, ok := statuses[remotefile]
					if !ok {
						mm := make(map[string]smartling.FileStatus)
						statuses[remotefile] = mm
					}
					statuses[remotefile][locale] = fs
				}(prefixedProjectFilepath, l.Locale)
			}

		}
		wg.Wait()

		fmt.Print("\n")
		fmt.Println("Translation counts: Awaiting Authorization -> In Progress -> Completed")
		fmt.Print("\n")

		// Format in columns
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		fmt.Fprint(w, " ")
		for _, locale := range locales {
			fmt.Fprint(w, "\t", locale.Locale)
		}
		fmt.Fprint(w, "\n")

		for projectFilepath, _ := range statuses {
			fmt.Fprint(w, projectFilepath)
			for _, locale := range locales {
				status := statuses[projectFilepath][locale.Locale]
				fmt.Fprint(w, "\t", status.AwaitingAuthorizationStringCount(), "->", status.InProgressStringCount(), "->", status.CompletedStringCount)
			}
			fmt.Fprint(w, "\n")
		}
		w.Flush()
	},
}

var projectPullCommand = cli.Command{
	Name:  "pull",
	Usage: "translate local project files using Smartling as a translation memory",

	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: project pull")
		}

		locales, err := client.Locales()
		panicIfErr(err)

		var wg sync.WaitGroup
		for _, projectFilepath := range ProjectConfig.Files() {
			for _, l := range locales {
				wg.Add(1)
				go func(locale, projectFilepath string) {
					defer wg.Done()

					hit, b, err, _ := translateViaCache(
						locale,
						localRelativeFilePath(projectFilepath),
						filetypeForProjectFile(projectFilepath),
						ProjectConfig.FileConfig.ParserConfig,
					)
					panicIfErr(err)

					fp := localPullFilePath(projectFilepath, locale)
					cached := ""
					if hit {
						cached = "(using cache)"
					}
					fmt.Println(fp, cached)
					err = ioutil.WriteFile(fp, b, 0644)
					panicIfErr(err)
				}(l.Locale, projectFilepath)
			}
		}
		wg.Wait()
	},
}

func cleanPrefix(s string) string {
	s = filepath.Clean("/" + s)
	if s == "/" {
		return ""
	}
	return s
}

func prefixOrGitPrefix(prefix string) string {
	if prefix == "" {
		prefix = pushPrefix()
	}
	if prefix == "master" {
		prefix = ""
	}

	prefix = cleanPrefix(prefix)

	if prefix != "" {
		fmt.Println("Using prefix", prefix)
	}
	return prefix
}

type RemoteFileStatus struct {
	RemoteFilePath string
	Statuses       map[string]*smartling.FileStatus
}

func (r *RemoteFileStatus) NotCompletedStringCount() int {
	c := 0
	for _, fs := range r.Statuses {
		c += fs.NotCompletedStringCount()
	}
	return c
}

func fetchStatusForLocales(remoteFilePath string, locales []string) RemoteFileStatus {
	ss := RemoteFileStatus{
		RemoteFilePath: remoteFilePath,
		Statuses:       map[string]*smartling.FileStatus{},
	}

	var wg sync.WaitGroup
	for _, locale := range locales {
		wg.Add(1)
		go func(f, l string) {
			defer wg.Done()

			s, err := client.Status(f, l)
			panicIfErr(err)
			ss.Statuses[l] = &s

		}(remoteFilePath, locale)
	}
	wg.Wait()

	return ss
}

var projectPushCommand = cli.Command{
	Name:  "push",
	Usage: "upload local project files with new strings, using the git branch or user name as a prefix",
	Description: `push [<prefix>]
Outputs the uploaded files for the given prefix
`,
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: push [<prefix>]")
		}

		prefix := prefixOrGitPrefix(c.Args().Get(0))
		locales := fetchLocales()

		var wg sync.WaitGroup
		for _, projectFilepath := range ProjectConfig.Files() {
			wg.Add(1)
			go func(prefix, projectFilepath string) {
				defer wg.Done()

				remoteFile := filepath.Clean(prefix + "/" + projectFilepath)

				_, err := client.Upload(projectFilepath, &smartling.UploadRequest{
					FileUri:      remoteFile,
					FileType:     filetypeForProjectFile(projectFilepath),
					ParserConfig: ProjectConfig.FileConfig.ParserConfig,
				})
				panicIfErr(err)

				remoteFileStatuses := fetchStatusForLocales(remoteFile, locales)

				// when using a prefix, we don't want to see files with
				// completely translated content
				if prefix != "" && remoteFileStatuses.NotCompletedStringCount() == 0 {
					err := client.Delete(remoteFile)
					panicIfErr(err)
				} else {
					fmt.Println(remoteFile)
				}
			}(prefix, projectFilepath)
		}
		wg.Wait()
	},
}

func filetypeForProjectFile(projectFilepath string) smartling.FileType {
	ft := smartling.FileTypeByExtension(filepath.Ext(projectFilepath))
	if ft == "" {
		ft = ProjectConfig.FileConfig.FileType
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
	fp, err := filepath.Rel(".", filepath.Join(ProjectConfig.path, remotepath))
	panicIfErr(err)
	return fp
}

func localPullFilePath(p, locale string) string {
	parts := FilenameParts{
		Path:   p,
		Dir:    filepath.Dir(p),
		Base:   filepath.Base(p),
		Ext:    filepath.Ext(p),
		Locale: locale,
	}

	dt := defaultPullDestination
	if dt != "" {
		dt = ProjectConfig.FileConfig.PullFilePath
	}

	out := bytes.NewBufferString("")
	tmpl := template.New("name")
	tmpl.Funcs(template.FuncMap{
		"TrimSuffix": strings.TrimSuffix,
	})
	_, err := tmpl.Parse(dt)
	panicIfErr(err)

	err = tmpl.Execute(out, parts)
	panicIfErr(err)

	return localRelativeFilePath(out.String())
}
