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
			Name:   "remote.server",
			Usage:  "ownCloud server",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "remote.folder",
			Usage:  "remote folder",
			EnvVar: "PLUGIN_REMOTE_FOLDER",
		},
		cli.StringSliceFlag{
			Name:   "local.files",
			Usage:  "local files",
			EnvVar: "PLUGIN_FILES",
		},
		cli.StringFlag{
			Name:   "name",
			Usage:  "suffix for compressed file",
			EnvVar: "PLUGIN_NAME",
		},
		cli.StringFlag{
			Name:   "auth.user",
			Usage:  "ownCloud username",
			EnvVar: "OWNCLOUD_USERNAME",
		},
		cli.StringFlag{
			Name:   "auth.pass",
			Usage:  "ownCloud password",
			EnvVar: "OWNCLOUD_PASSWORD",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "be verbose",
			EnvVar: "PLUGIN_VERBOSE",
		},
		cli.StringFlag{
			Name:   "commit.tag",
			Usage:  "commit tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.repo",
			Usage:  "repo name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.BoolTFlag{
			Name:   "parentdir",
			Usage:  "Include directory structure. Defualts to true",
			EnvVar: "PLUGIN_PARENTDIR",
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
		Files:     c.StringSlice("local.files"),
		Name:      c.String("name"),
		Parentdir: c.Bool("parentdir"),
		Auth: Auth{
			User: c.String("auth.user"),
			Pass: c.String("auth.pass"),
		},
		Commit: Commit{
			Tag:  c.String("commit.tag"),
			Sha:  c.String("commit.sha"),
			Repo: c.String("commit.repo"),
		},
		Verbose: c.Bool("verbose"),
	}

	return plugin.Exec()
}
