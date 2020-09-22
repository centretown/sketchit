package main

import "github.com/centretown/sketchit/api"

var skills = []*api.Skill{
	{
		Task: api.Task_hello.String(),
		Alternates: []string{
			"ping", "hi",
		},
		Description: "Test server by saying hello.",
		Summary: &api.Summary{
			Usage:    "hello -f=<format> -d=<detail>",
			Examples: []string{"hello"},
		},
	},
	{
		Task: api.Task_help.String(),
		Alternates: []string{
			"info", "?",
		},
		Description: "Get information on usage.",
		Summary: &api.Summary{
			Usage:  "help <topic> <topic>... -f=<format> -d=<detail>",
			Syntax: "help <topic> where topic is a skill or collection",
			Examples: []string{
				"// outcome all topics brief",
				"help",
				"// outcome help for skill list presented in full detail",
				"help list",
				"// outcome help for collection devices in full detail",
				"help devices",
			},
		},
	},
	{
		Task:        api.Task_list.String(),
		Alternates:  []string{"ls"},
		Description: "List the contents of the destination.",
		Summary: &api.Summary{
			Usage:  "list <destination> -f=<format> -d=<detail>",
			Syntax: "list <destination>",
			Examples: []string{
				"list",
				"list devices",
				"list .devices.work",
			},
		},
	},
	{
		Task:        api.Task_goto.String(),
		Alternates:  []string{"cd", "go"},
		Description: "Go to a destination.",
		Summary: &api.Summary{
			Usage:  "goto <destination>",
			Syntax: "goto <destination> or go . up a level",
			Examples: []string{
				"goto .",
				"goto devices",
				"goto devices work",
				"goto .sketches.ESP32"},
		},
	},
	{
		Task:        api.Task_save.String(),
		Alternates:  []string{"store"},
		Description: "Save the current contents. Confirmation required.",
		Summary: &api.Summary{
			Usage:  "save",
			Syntax: "save",
			Examples: []string{
				"save",
			},
		},
	},
	{
		Task:        api.Task_remove.String(),
		Alternates:  []string{"rm", "del", "delete"},
		Description: "Remove the current contents. Confirmation required.",
		Summary: &api.Summary{
			Usage:  "remove",
			Syntax: "remove",
			Examples: []string{
				"remove",
			},
		},
	},
	{
		Task:        api.Task_exit.String(),
		Alternates:  []string{"x", "quit"},
		Description: "Exit the current session.",
		Summary: &api.Summary{
			Usage:  "exit",
			Syntax: "exit",
			Examples: []string{
				"exit",
			},
		},
	},
}
