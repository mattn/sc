# sc

*THIS IS EXPERIMENTAL PACKAGE, AND WORK IN PROGRESS!*

Sub Command

## Usage

```go
sc.Cmds{
	{
		Name: "add",
		Flags: sc.Flags{
			{Name: "-n", Desc: "name", Type: sc.String, Usage: func() { fmt.Println("invalid name") }},
			{Name: "-v", Desc: "verbose", Type: sc.Bool, Usage: func() { fmt.Println("invalid option") }},
		},
		Desc: "untra add command",
		Run: func(c *sc.C, args []string) error {
			if len(args) == 0 {
				return sc.UsageError
			}
			fmt.Println(c.LookupFlag("-n").String())
			fmt.Println("added", args)
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
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
