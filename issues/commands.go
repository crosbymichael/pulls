package main

import (
	"github.com/codegangsta/cli"
)

func loadCommands(app *cli.App) {
	app.Action = mainCmd

	app.Flags = []cli.Flag{
		cli.StringFlag{"assigned", "", "display issues assigned to <user>. Use '*' for all assigned, or 'none' for all unassigned."},
		cli.BoolFlag{"no-trunc", "do not truncate the issue name"},
		cli.IntFlag{"votes", -1, "display the number of votes '+1' filtered by the <number> specified."},
		cli.BoolFlag{"vote", "add '+1' to an specific issue."},
		cli.StringFlag{"labels", "", "add lables to the search"},
		cli.BoolFlag{"triage", "display only issues without labels"},
		cli.StringSliceFlag{"apply", &cli.StringSlice{}, "apply lables to an issue"},
	}

	app.Commands = []cli.Command{
		{
			Name:   "alru",
			Usage:  "Show the Age of the Least Recently Updated issue for this repo. Lower is better.",
			Action: alruCmd,
		},
		{
			Name:   "repo",
			Usage:  "List information about the current repository",
			Action: repositoryInfoCmd,
		},
		{
			Name:   "take",
			Usage:  "Assign an issue to your github account",
			Action: takeCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{"overwrite", "overwrites a taken issue"},
			},
		},
		{
			Name:   "auth",
			Usage:  "Add a github token for authentication",
			Action: authCmd,
			Flags: []cli.Flag{
				cli.StringFlag{"add", "", "add new token for authentication"},
			},
		},
		{
			Name:   "labels",
			Usage:  "Show all labels",
			Action: labelsCmd,
		},
	}
}
