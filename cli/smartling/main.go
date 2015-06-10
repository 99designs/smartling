package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99designs/smartling"
	"github.com/codegangsta/cli"
)

var client *smartling.FaultTolerantClient

func init() {
	log.SetFlags(0)
}

func panicIfErr(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "smartling"
	app.Usage = "manage translation files using Smartling"
	app.Version = "0.1.0"
	app.Before = func(c *cli.Context) error {
		apiKey := c.String("apikey")
		projectId := c.String("projectid")
		configFile := c.String("configfile")
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
			if projectId == "" {
				projectId = ProjectConfig.ProjectId
			}
		}

		var sc *smartling.Client
		if c.Bool("sandbox") {
			log.Println("Using sandbox")
			sc = smartling.NewSandboxClient(apiKey, projectId)
		} else {
			sc = smartling.NewClient(apiKey, projectId)
		}

		client = &smartling.FaultTolerantClient{sc, 3}

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
		}, cli.BoolFlag{
			Name:   "sandbox",
			Usage:  "Use the sandbox",
			EnvVar: "SMARTLING_SANDBOX",
		}, cli.StringFlag{
			Name:   "configfile,c",
			Usage:  "Project config file to use",
			EnvVar: "SMARTLING_CONFIGFILE",
		},
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
	panicIfErr(err)
}
