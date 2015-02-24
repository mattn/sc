package main

import (
	"fmt"
	"github.com/mattn/sc"
)

func main() {
	sc.Cmds{
		{
			Name: "add",
			Flags: sc.Flags{
				{Name: "-n", Desc: "name", Type: sc.String},
				{Name: "-v", Desc: "verbose", Type: sc.Bool},
			},
			Desc: "untra add command",
			Run: func(c *sc.C, args []string) error {
				if len(args) == 0 {
					return sc.UsageError
				}
				// verbose
				if c.LookupFlag("-v").Bool() {
					for _, arg := range args {
						fmt.Println("added", arg)
					}
				} else {
					fmt.Println("added", args)
				}
				return nil
			},
		},
		{
			Name: "del",
			Flags: sc.Flags{
				{Name: "-n", Desc: "name", Type: sc.String},
				{Name: "-v", Desc: "verbose", Type: sc.Bool},
			},
			Desc: "untra del command",
			Run: func(c *sc.C, args []string) error {
				if len(args) == 0 {
					return sc.UsageError
				}
				if c.LookupFlag("-v").Bool() {
					for _, arg := range args {
						fmt.Println("deleted", arg)
					}
				} else {
					fmt.Println("deleted", args)
				}
				return nil
			},
		},
	}.Run(&sc.C{
		Default: "add",
		Desc:    "super cool command",
		Usage: func(c *sc.C) {
			fmt.Println(c.Name, "fooo")
			fmt.Println("    ", c.Desc)
			c.PrintCommands()
		},
	})
}
