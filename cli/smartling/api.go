package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/smartling"
	smartlingNew "github.com/Smartling/api-sdk-go"
	"github.com/codegangsta/cli"
)

func ListConditionSlice(cc []string) []smartling.ListCondition {
	ll := []smartling.ListCondition{}
	for _, c := range cc {
		ll = append(ll, smartling.ListCondition(c))
	}
	return ll
}

func removeEmptyStrings(ss []string) []string {
	newSs := []string{}
	for _, s := range ss {
		if s != "" {
			newSs = append(newSs, s)
		}
	}
	return newSs
}

func PrintList(uriMask string, olderThan time.Duration, long bool, conditions []string) {
	req := smartling.ListRequest{
		UriMask: uriMask,
	}
	conditions = removeEmptyStrings(conditions)
	if len(conditions) > 0 {
		req.Conditions = ListConditionSlice(conditions)
	}
	if olderThan > 0 {
		t := smartling.Iso8601Time(time.Now().Add(-olderThan))
		req.LastUploadedBefore = &t
	}

	files, err := client.List(req)
	logAndQuitIfError(err)

	// TODO: fix this "long"
	if long {
		fmt.Println("total", files.TotalCount)
		for _, f := range files.Items {
			// t := time.Time(f.LastUploaded).Format("2 Jan 15:04")
			fmt.Println(f.FileURI)

			// fmt.Printf("%3d strings  %s  %s\n", f.StringCount, t, f.FileUri)
		}
	} else {
		for _, f := range files.Items {
			fmt.Println(f.FileURI)
		}
	}
}

var LsCommand = cli.Command{
	Name:        "ls",
	Usage:       "list remote files",
	Description: "ls [<uriMask>]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "conditions",
		}, cli.StringFlag{
			Name: "older-than",
		}, cli.BoolFlag{
			Name: "long,l",
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

		conditions := strings.Split(c.String("conditions"), ",")

		PrintList(uriMask, d, c.Bool("long"), conditions)
	},
}

func PrintFileStatus(remotepath, locale string) {
	f, err := client.Status(remotepath, locale)
	logAndQuitIfError(err)

	fmt.Println("File                    ", f.FileURI)
	fmt.Println("String Count            ", f.TotalStringCount)
	fmt.Println("Word Count              ", f.TotalWordCount)
	fmt.Println("Authorized String Count ", f.AuthorizedStringCount)
	fmt.Println("Completed String Count  ", f.CompletedStringCount)
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
			b, err = client.DownloadTranslation(locale, smartlingNew.FileDownloadRequest{
				FileURIRequest: smartlingNew.FileURIRequest{FileURI: remotepath},
			})
		}
		logAndQuitIfError(err)

		fmt.Println(string(b))
	},
}

var PutCommand = cli.Command{
	Name:        "put",
	Usage:       "uploads a local file",
	Description: "put <local file>",
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
		if len(c.Args()) != 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: put <local file>")
		}

		localpath := c.Args().Get(0)

		ft := smartlingNew.FileType(c.String("filetype"))
		if ft == "" {
			ft = smartlingNew.GetFileTypeByExtension(filepath.Ext(localpath))
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
		if err != nil {
			logAndQuitIfError(err)
		}

		r, err := client.Upload(&smartlingNew.FileUploadRequest{
			File:           f,
			FileType:       ft,
			Authorize:      c.Bool("approve"),
			FileURIRequest: smartlingNew.FileURIRequest{FileURI: localpath},
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

		items, err := client.LastModified(smartling.LastModifiedRequest{
			FileUri: remotepath,
		})
		logAndQuitIfError(err)

		for _, i := range items {
			t := time.Time(i.LastModified).Format("2 Jan 3:04")
			fmt.Printf("%s %s\n", i.Locale, t)
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

		// Should this somehow be extracted out to a method?
		pd, err := client.Client.GetProjectDetails("09bd710ee")
		logAndQuitIfError(err)
		r := pd.TargetLocales
		// r, err := client.Locales()

		for _, l := range r {
			if l.Enabled {
				fmt.Printf("%-5s  %s\n", l.LocaleID, l.Description)
			}
		}
	},
}
