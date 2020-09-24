package main

import "github.com/centretown/sketchit/api"

var features = []*api.Feature{
	{
		Flag:        api.Feature_f,
		Label:       "format",
		Description: "The format feature can format outcomes as yaml, json or xml.",
		Summary: &api.Summary{
			Usage:  "<task> -f=<format> <destination>",
			Syntax: "-f=<format>",
			Examples: []string{
				"// outcomes formatted as YAML",
				"-f=yaml",
				"// outcomes formatted as JSON",
				"-f=json",
				"// outcomes formatted as XML",
				"-f=xml",
			},
		},
	},
	{
		Flag:        api.Feature_d,
		Label:       "detail",
		Description: "The detail feature can project outcomes as full, summary or brief. More than one projection may be specified.  Each projection corresponds to a level of detail in the data. The last projection is applied to all remaining levels.",
		Summary: &api.Summary{
			Usage:  "<task> -d=<detail> <destination>",
			Syntax: "-d=<detail1,detail2...>",
			Examples: []string{
				"// all levels are presented in full",
				"-d=full",
				"// all levels are summarized",
				"-d=summary",
				"// all levels are brief",
				"-d=brief",
				"// 1st level full, remaining are brief",
				"-d=full,brief",
				"// 1st level full, 2nd level summarized, remaining are brief",
				"-d=full,summary,brief",
			},
		},
	},
	{
		Flag:        api.Feature_auto,
		Label:       "auto",
		Description: "The auto feature can automatically reply to confirmation requests.",
		Summary: &api.Summary{
			Usage:  "<task> -auto=<reply> <destination>",
			Syntax: "-auto=<reply>",
			Examples: []string{
				"// auto reply off",
				"-auto=off",
				"// auto reply yes",
				"-auto=y",
				"// auto reply no",
				"-auto=n",
			},
		},
	},
}
