package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99designs/smartling"
	smartlingNew "github.com/Smartling/api-sdk-go"
	"github.com/codegangsta/cli"
)

var client *smartling.FaultTolerantClient

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
	// Things needed to authenticate
	userID := c.GlobalString("userid")
	apiKey := c.GlobalString("apikey")
	projectID := c.GlobalString("projectid")
	configFile := c.GlobalString("configfile")

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
			projectID = ProjectConfig.ProjectId
		}
	}

	if apiKey == "" {
		log.Fatalln("ApiKey not specified in --apikey or", configFile)
	}
	if projectID == "" {
		log.Fatalln("ProjectId not specified in --projectid or", configFile)
	}

	sc := smartlingNew.NewClient(userID, apiKey)
	// FIXME: should projectID be passed to this faulttolerent client? probs not
	client = &smartling.FaultTolerantClient{sc, projectID, 10}

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
			Name:   "projectid, p",
			Usage:  "Smartling Project ID",
			EnvVar: "SMARTLING_PROJECTID",
		}, cli.StringFlag{
			Name:   "configfile,c",
			Usage:  "Project config file to use",
			EnvVar: "SMARTLING_CONFIGFILE",
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
