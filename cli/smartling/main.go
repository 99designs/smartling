package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/smartling"
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
	apiKey := c.GlobalString("apikey")
	projectId := c.GlobalString("projectid")
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
		if projectId == "" {
			projectId = ProjectConfig.ProjectId
		}
	}

	if apiKey == "" {
		log.Fatalln("ApiKey not specified in --apikey or", configFile)
	}
	if projectId == "" {
		log.Fatalln("ProjectId not specified in --projectid or", configFile)
	}

	var sc *smartling.Client
	if c.Bool("sandbox") {
		log.Println("Using sandbox")
		sc = smartling.NewSandboxClient(apiKey, projectId)
	} else {
		sc = smartling.NewClient(apiKey, projectId)
	}

	if timeout != 0 {
		sc.SetHttpTimeout(time.Duration(timeout) * time.Second)
	}

	client = &smartling.FaultTolerantClient{sc, 10}

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
		}, cli.BoolFlag{
			Name:   "sandbox",
			Usage:  "Use the sandbox",
			EnvVar: "SMARTLING_SANDBOX",
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
