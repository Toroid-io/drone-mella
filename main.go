package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {

	app := cli.NewApp()
	app.Name = "mella plugin"
	app.Usage = "mella plugin"
	app.Action = run
	app.Version = fmt.Sprintf("0.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "remote.server",
			Usage: "ownCloud server",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name: "remote.folder",
			Usage: "remote folder",
			EnvVar: "PLUGIN_REMOTE_FOLDER",
		},
		cli.StringFlag{
			Name: "local.folder",
			Usage: "local folder",
			EnvVar: "PLUGIN_LOCAL_FOLDER",
		},
		cli.StringFlag{
			Name: "local.files",
			Usage: "local files",
			EnvVar: "PLUGIN_LOCAL_FILES",
		},
		cli.StringFlag{
			Name: "auth.user",
			Usage: "ownCloud username",
			EnvVar: "OWNCLOUD_USERNAME",
		},
		cli.StringFlag{
			Name: "auth.pass",
			Usage: "ownCloud password",
			EnvVar: "OWNCLOUD_PASSWORD",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {

	plugin := Plugin{
		Remote: Remote{
			Server: c.String("remote.server"),
			Folder: c.String("remote.folder"),
		},
		Local: Local{
			Folder: c.String("local.folder"),
			Files: c.String("local.files"),
		},
		Auth: Auth{
			User: c.String("auth.user"),
			Pass: c.String("auth.pass"),
		},
	}
	
	return plugin.Exec()
}