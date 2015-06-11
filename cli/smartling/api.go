package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/smartling"
	"github.com/codegangsta/cli"
)

var LsCommand = cli.Command{
	Name:  "ls",
	Usage: "list remote files",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: ls")
		}

		files, err := client.List(smartling.ListRequest{})
		panicIfErr(err)

		fmt.Println("total", len(files))
		for _, f := range files {
			t := time.Time(f.LastUploaded).Format("2 Jan 3:04")
			fmt.Printf("%3d strings  %s  %s\n", f.StringCount, t, f.FileUri)
		}

	},
}

var StatusCommand = cli.Command{
	Name:        "stat",
	Usage:       "display the translation status of a remote file",
	Description: "stat <remote file> <locale>",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: stat <remote file> <locale>")
		}

		remotepath := c.Args().Get(0)
		locale := c.Args().Get(1)

		r, err := client.Status(remotepath, locale)
		panicIfErr(err)

		fmt.Println("File                  ", r.FileUri)
		fmt.Println("String Count          ", r.StringCount)
		fmt.Println("Word Count            ", r.WordCount)
		fmt.Println("Approved String Count ", r.ApprovedStringCount)
		fmt.Println("Completed String Count", r.CompletedStringCount)
		fmt.Println("Last Uploaded         ", r.LastUploaded)
		fmt.Println("File Type             ", r.FileType)
	},
}

var GetCommand = cli.Command{
	Name:        "get",
	Usage:       "downloads a remote file",
	Description: "get <remote file>",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: get <remote file>")
		}

		remotepath := c.Args().Get(0)

		b, err := client.Get(&smartling.GetRequest{
			FileUri: remotepath,
		})
		panicIfErr(err)

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
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: put <local file> <remote file>")
		}

		localpath := c.Args().Get(0)
		remotepath := c.Args().Get(1)

		ft := smartling.FileType(c.String("filetype"))
		if ft == "" {
			ft = smartling.FileTypeByExtension(filepath.Ext(localpath))
		}

		parserconfig := map[string]string{}
		for _, q := range strings.Split(c.String("parserconfig"), ",") {
			pc := strings.SplitN(q, "=", 2)
			parserconfig[pc[0]] = pc[1]
		}

		r, err := client.Upload(localpath, &smartling.UploadRequest{
			FileUri:      remotepath,
			FileType:     ft,
			ParserConfig: parserconfig,
		})
		panicIfErr(err)

		fmt.Println("Overwritten: ", r.OverWritten)
		fmt.Println("String Count:", r.StringCount)
		fmt.Println("Word Count:  ", r.WordCount)
	},
}

var RenameCommand = cli.Command{
	Name:        "rename",
	Usage:       "renames a remote file",
	Description: "rename <remote file> <new smartling file>",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: rename <remote file> <new smartling file>")
		}

		remotepath := c.Args().Get(0)
		newremotepath := c.Args().Get(0)

		err := client.Rename(remotepath, newremotepath)
		panicIfErr(err)
	},
}

var RmCommand = cli.Command{
	Name:        "rm",
	Usage:       "removes a remote file",
	Description: "rm <remote file>...",
	Action: func(c *cli.Context) {
		if len(c.Args()) < 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: rm <remote file>...")
		}

		for _, remotepath := range c.Args() {
			panicIfErr(client.Delete(remotepath))
		}
	},
}

var LastmodifiedCommand = cli.Command{
	Name:        "lastmodified",
	Usage:       "shows when a remote file was modified last",
	Description: "lastmodified <remote file>",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: lastmodified <remote file>")
		}

		remotepath := c.Args().Get(0)

		items, err := client.LastModified(smartling.LastModifiedRequest{
			FileUri: remotepath,
		})
		panicIfErr(err)

		for _, i := range items {
			t := time.Time(i.LastModified).Format("2 Jan 3:04")
			fmt.Printf("%s %s\n", i.Locale, t)
		}
	},
}

var LocalesCommand = cli.Command{
	Name:  "locales",
	Usage: "list the locales for the project",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: locales")
		}

		r, err := client.Locales()
		panicIfErr(err)

		for _, l := range r {
			fmt.Printf("%-5s  %-23s  %s\n", l.Locale, l.Name, l.Translated)
		}
	},
}
