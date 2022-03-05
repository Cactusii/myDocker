package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"myDocker/container"
	"os"
)

const usage = `Welcome!`

func main() {
	app := cli.NewApp()
	app.Name = "myDocker"
	app.Usage = usage

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

var runCommand = cli.Command{
	// run command
	Name: "run",
	Usage: `Create a container whit namespace and cgroups limit
			myDoker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container comand")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("it")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `Init container process run user's process in container. Do not call it outside`,

	Action: func(context *cli.Context) error {
		logrus.Infof("init come on")
		cmd := context.Args().Get(0)
		logrus.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
