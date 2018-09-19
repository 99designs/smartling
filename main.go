package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/api-sdk-go"
	"github.com/codegangsta/cli"
)

var client *FaultTolerantClient

var Version = "dev"

func init() {
	log.SetFlags(0)
}

func logAndQuitIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

var cmdBefore = func(c *cli.Context) error {
	userID := c.GlobalString("userid")
	apiKey := c.GlobalString("apikey")
	projectID := c.GlobalString("projectid")
	configFile := c.GlobalString("configfile")
	timeout := c.GlobalInt("timeout")

	if configFile == "" {
		configFile = "smartling.yml"
	}

	var err error
	ProjectConfig, err = loadConfig(configFile)
	if err != nil {
		loadProjectErr = fmt.Errorf("Error loading %s: %s", configFile, err.Error())
	}

	if ProjectConfig != nil {
		if apiKey == "" {
			apiKey = ProjectConfig.ApiKey
		}
		if projectID == "" {
			projectID = ProjectConfig.ProjectID
		}
		if userID == "" {
			userID = ProjectConfig.UserID
		}
	}

	if apiKey == "" {
		log.Fatalln("ApiKey not specified in --apikey or", configFile)
	}
	if projectID == "" {
		log.Fatalln("ProjectID not specified in --projectid or", configFile)
	}
	if userID == "" {
		log.Fatalln("UserID not specified in --userid or", configFile)
	}

	sc := smartling.NewClient(userID, apiKey)

	if timeout != 0 {
		sc.HTTP.Timeout = (time.Duration(timeout) * time.Second)
	}

	client = &FaultTolerantClient{sc, projectID, 10}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "smartling"
	app.Usage = "manage translation files using Smartling"
	app.Version = Version
	app.Before = func(c *cli.Context) error {
		if c.Bool("version") {
			cli.VersionPrinter(c)
			os.Exit(0)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apikey, k",
			Usage:  "Smartling ApiKey",
			EnvVar: "SMARTLING_APIKEY",
		}, cli.StringFlag{
			Name:   "userid, u",
			Usage:  "Smartling User Identifier",
			EnvVar: "SMARTLING_USERID",
		}, cli.StringFlag{
			Name:   "projectid, p",
			Usage:  "Smartling Project ID",
			EnvVar: "SMARTLING_PROJECTID",
		}, cli.StringFlag{
			Name:   "configfile,c",
			Usage:  "Project config file to use",
			EnvVar: "SMARTLING_CONFIGFILE",
		}, cli.IntFlag{
			Name:   "timeout,t",
			Value:  60,
			Usage:  "Maximum time in seconds for an API request to take",
			EnvVar: "SMARTLING_API_TIMEOUT",
		},

		cli.VersionFlag,
	}
	app.Commands = []cli.Command{
		LsCommand,
		StatusCommand,
		GetCommand,
		PutCommand,
		RenameCommand,
		RmCommand,
		LastmodifiedCommand,
		LocalesCommand,
		ProjectCommand,
	}

	err := app.Run(os.Args)
	logAndQuitIfError(err)
}
