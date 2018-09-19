package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Smartling/api-sdk-go"
	"github.com/codegangsta/cli"
)

func PrintList(uriMask string, olderThan time.Duration) {
	req := smartling.FilesListRequest{
		URIMask: uriMask,
	}

	if olderThan > 0 {
		req.LastUploadedBefore = smartling.UTC{Time: time.Now().Add(-olderThan)}
	}

	files, err := client.List(req)
	logAndQuitIfError(err)

	for _, f := range files.Items {
		fmt.Println(f.FileURI)
	}
}

var LsCommand = cli.Command{
	Name:        "ls",
	Usage:       "list remote files",
	Description: "ls [<uriMask>]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "older-than",
		},
	},
	Before: cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: ls [<uriMask>]")
		}
		uriMask := c.Args().Get(0)

		var d time.Duration
		if len(c.String("older-than")) > 0 {
			var err error
			d, err = time.ParseDuration(c.String("older-than"))
			logAndQuitIfError(err)
		}

		PrintList(uriMask, d)
	},
}

func PrintFileStatus(remotepath, locale string) {
	f, err := client.Status(remotepath)
	logAndQuitIfError(err)
	fst, err := f.GetFileStatusTranslation(locale)
	logAndQuitIfError(err)

	fmt.Println("File                    ", f.FileURI)
	fmt.Println("String Count            ", f.TotalStringCount)
	fmt.Println("Authorized String Count ", fst.AuthorizedStringCount)
	fmt.Println("Completed String Count  ", fst.CompletedStringCount)
	fmt.Println("Excluded String Count   ", fst.ExcludedStringCount)
	fmt.Println("Last Uploaded           ", f.LastUploaded)
	fmt.Println("File Type               ", f.FileType)
}

var StatusCommand = cli.Command{
	Name:        "stat",
	Usage:       "display the translation status of a remote file",
	Description: "stat <remote file> <locale>",
	Before:      cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: stat <remote file> <locale>")
		}

		remotepath := c.Args().Get(0)
		locale := c.Args().Get(1)

		PrintFileStatus(remotepath, locale)
	},
}

var GetCommand = cli.Command{
	Name:        "get",
	Usage:       "downloads a remote file",
	Description: "get <remote file>",
	Before:      cmdBefore,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "locale",
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) != 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: get <remote file>")
		}

		remotepath := c.Args().Get(0)
		locale := c.String("locale")

		var (
			b   []byte
			err error
		)

		if locale == "" {
			b, err = client.Download(remotepath)
		} else {
			b, err = client.DownloadTranslation(locale, smartling.FileDownloadRequest{
				FileURIRequest: smartling.FileURIRequest{FileURI: remotepath},
			})
		}
		logAndQuitIfError(err)

		fmt.Println(string(b))
	},
}

var PutCommand = cli.Command{
	Name:        "put",
	Usage:       "uploads a local file",
	Description: "put <local file> <remote file>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "filetype",
		}, cli.StringFlag{
			Name: "parserconfig",
		}, cli.BoolFlag{
			Name: "approve",
		},
	},
	Before: cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: put <local file> <remote file>")
		}

		localpath := c.Args().Get(0)
		remotepath := c.Args().Get(1)

		ft := smartling.FileType(c.String("filetype"))
		if ft == "" {
			ft = smartling.GetFileTypeByExtension(filepath.Ext(localpath))
		}

		parserconfig := map[string]string{}
		if c.String("parserconfig") != "" {
			parts := strings.Split(c.String("parserconfig"), ",")
			if len(parts)%2 == 1 {
				log.Fatalln("parserconfig must be in the format --parserconfig=key1,value1,key2,value2")
			}
			for i := 0; i < len(parts); i += 2 {
				parserconfig[parts[i]] = parts[i+1]
			}
		}

		f, err := ioutil.ReadFile(localpath)
		logAndQuitIfError(err)

		r, err := client.Upload(&smartling.FileUploadRequest{
			File:           f,
			FileType:       ft,
			Authorize:      c.Bool("approve"),
			FileURIRequest: smartling.FileURIRequest{FileURI: remotepath},
		})

		logAndQuitIfError(err)

		fmt.Println("Overwritten: ", r.Overwritten)
		fmt.Println("String Count:", r.StringCount)
		fmt.Println("Word Count:  ", r.WordCount)
	},
}

var RenameCommand = cli.Command{
	Name:        "rename",
	Usage:       "renames a remote file",
	Description: "rename <remote file name> <new smartling file name>",
	Before:      cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: rename <remote file> <new smartling file name>")
		}

		remotepath := c.Args().Get(0)
		newremotepath := c.Args().Get(1)

		err := client.Rename(remotepath, newremotepath)

		logAndQuitIfError(err)
	},
}

var RmCommand = cli.Command{
	Name:        "rm",
	Usage:       "removes a remote file",
	Description: "rm <remote file>...",
	Before:      cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) < 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: rm <remote file>...")
		}

		for _, remotepath := range c.Args() {
			logAndQuitIfError(client.Delete(remotepath))
		}
	},
}

var LastmodifiedCommand = cli.Command{
	Name:        "lastmodified",
	Usage:       "shows when a remote file was modified last",
	Description: "lastmodified <remote file>",
	Before:      cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) != 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: lastmodified <remote file>")
		}

		remotepath := c.Args().Get(0)

		locales, err := client.LastModified(smartling.FileLastModifiedRequest{
			FileURIRequest: smartling.FileURIRequest{FileURI: remotepath},
		})
		logAndQuitIfError(err)

		for _, i := range locales.Items {
			t := time.Time(i.LastModified.Time).Format("2 Jan 3:04")
			fmt.Printf("%s %s\n", i.LocaleID, t)
		}
	},
}

var LocalesCommand = cli.Command{
	Name:   "locales",
	Usage:  "list the locales for the project",
	Before: cmdBefore,
	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: locales")
		}

		tl, err := client.Locales()
		logAndQuitIfError(err)

		for _, l := range tl {
			if l.Enabled {
				fmt.Printf("%-5s  %s\n", l.LocaleID, l.Description)
			}
		}
	},
}
