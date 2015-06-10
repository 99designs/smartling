package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/99designs/smartling"
	"github.com/codegangsta/cli"
)

func statusOfProjectFile(remotepath, locale string) string {
	r, err := client.Status(remotepath, locale)
	panicIfErr(err)

	var pctComplete float64
	if r.CompletedStringCount != 0 {
		pctComplete = float64(r.CompletedStringCount) * 100 / float64(r.StringCount)
	}

	return fmt.Sprintf("  %-5s %d%%\n", locale, int(math.Floor(pctComplete)))
}

func printStatusOfProjectFile(projectFilepath string, locales []smartling.Locale) {
	tmpfile, err := uploadAsTempFile(
		localRelativeFilePath(projectFilepath),
		filetypeForProjectFile(projectFilepath),
		ProjectConfig.FileConfig.ParserConfig,
	)
	panicIfErr(err)

	statuses := ""

	var wgLocales sync.WaitGroup
	for _, l := range locales {
		wgLocales.Add(1)
		go func(locale string) {
			defer wgLocales.Done()
			s := statusOfProjectFile(tmpfile, locale)
			statuses += s
		}(l.Locale)
	}
	wgLocales.Wait()

	fmt.Printf("%s:\n%s", projectFilepath, statuses)
}

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
		{
			Name:        "status",
			Usage:       "show the status of the project's local files",
			Description: "status",

			Action: func(c *cli.Context) {
				locales, err := client.Locales()
				panicIfErr(err)

				var wgFiles sync.WaitGroup
				for _, projectFilepath := range ProjectConfig.Files {
					wgFiles.Add(1)
					go func(projectFilepath string, locales []smartling.Locale) {
						defer wgFiles.Done()
						printStatusOfProjectFile(projectFilepath, locales)
					}(projectFilepath, locales)
				}
				wgFiles.Wait()
			},
		}, {
			Name:  "pull",
			Usage: "translate local project files using Smartling as a translation memory",

			Action: func(c *cli.Context) {
				locales, err := client.Locales()
				panicIfErr(err)

				var wg sync.WaitGroup
				for _, projectFilepath := range ProjectConfig.Files {
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
		}, {
			Name:        "push",
			Usage:       "upload local project files, using the git branch or user name as a prefix",
			Description: "push",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "Use the specified prefix instead of the default",
				},
			},
			Action: func(c *cli.Context) {
				prefix := c.String("prefix")
				if prefix == "" {
					prefix = pushPrefix()
				}
				prefix = filepath.Clean("/" + prefix)
				fmt.Println("Using prefix", prefix)

				var wg sync.WaitGroup
				for _, projectFilepath := range ProjectConfig.Files {
					wg.Add(1)
					go func(projectFilepath string) {
						defer wg.Done()

						f := filepath.Clean(prefix + "/" + projectFilepath)

						fmt.Println("Uploading", f)
						_, err := client.Upload(projectFilepath, &smartling.UploadRequest{
							FileUri:      f,
							FileType:     filetypeForProjectFile(projectFilepath),
							ParserConfig: ProjectConfig.FileConfig.ParserConfig,
						})
						panicIfErr(err)
					}(projectFilepath)
				}
				wg.Wait()
			},
		},
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
