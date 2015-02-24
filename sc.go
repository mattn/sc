package sc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type usageError int

func (u usageError) Error() string {
	return "usage error"
}

const (
	UsageError usageError = iota
	usageErrorFire
)

type FlagType int

const (
	String FlagType = iota
	Bool
	Int64
	Float64
)

type Flags []*F
type Cmds []*C

type F struct {
	Name  string
	Desc  string
	Type  FlagType
	value interface{}
}

func (f *F) IsSet() bool {
	return f.value != nil
}

func (f *F) Int64() int64 {
	if v, ok := f.value.(int64); ok {
		return v
	}
	v, err := strconv.ParseInt(f.String(), 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func (f *F) Float64() float64 {
	if v, ok := f.value.(float64); ok {
		return v
	}
	v, err := strconv.ParseFloat(f.String(), 64)
	if err != nil {
		return 0.0
	}
	return v
}

func (f *F) Bool() bool {
	if v, ok := f.value.(bool); ok {
		return v
	}
	v, err := strconv.ParseBool(f.String())
	if err != nil {
		return false
	}
	return v
}

func (f *F) String() string {
	if s, ok := f.value.(string); ok {
		return s
	}
	return fmt.Sprint(f.value)
}

type C struct {
	Name    string
	Flags   Flags
	Run     func(*C, []string) error
	Desc    string
	Usage   func(*C)
	Cmds    Cmds
	Default string
	main    string
}

func (t *C) Main() string {
	return t.main
}

func (t *C) PrintCommands() {
	for _, c := range t.Cmds {
		fmt.Println("    ", c.Name, ":", c.Desc)
	}
}

func (t *C) PrintFlags() {
	for _, c := range t.Flags {
		fmt.Println("    ", c.Name, ":", c.Desc)
	}
}

func (t *C) LookupFlag(arg string) *F {
	var f *F
flagLoop:
	for _, ff := range t.Flags {
		if ff.Name == arg {
			f = ff
			break flagLoop
		}
	}
	return f
}

func (t *C) run(args []string) error {
	if t == nil {
		return errors.New("sc: invalid argument")
	}
	i := 1
	for ; i < len(args); i++ {
		arg := args[i]
		if arg != "" && arg[0] == '-' {
			f := t.LookupFlag(arg)
			var hs bool
			var s string
		flagLoop:
			for _, ff := range t.Flags {
				if ff.Name == arg {
					f = ff
					break flagLoop
				}
				if strings.HasPrefix(arg, ff.Name+"=") {
					hs = true
					s = arg[len(ff.Name)+1:]
					f = ff
					break flagLoop
				}
			}
			if f == nil {
				return UsageError
			}
			if f.Type != Bool && !hs {
				if i < len(args)-1 {
					s = args[i+1]
					i++
				} else {
					return UsageError
				}
			}
			switch f.Type {
			case Bool:
				if hs {
					if v, err := strconv.ParseBool(s); err == nil {
						f.value = v
					}
				} else {
					f.value = true
				}
			case String:
				f.value = args[i+1]
			case Float64:
				if v, err := strconv.ParseFloat(s, 64); err == nil {
					f.value = v
				} else {
					return UsageError
				}
			case Int64:
				if v, err := strconv.ParseInt(args[i+1], 10, 64); err == nil {
					f.value = v
				} else {
					return UsageError
				}
			}
		} else {
			if t.Run != nil {
				return t.Run(t, args[i:])
			}
			var c *C
		cmdLoop:
			for _, cc := range t.Cmds {
				cc.main = t.main
				if cc.Name == arg {
					c = cc
					break cmdLoop
				}
			}
			arg = t.Default
		defLoop:
			for _, cc := range t.Cmds {
				if cc.Name == arg {
					c = cc
					break defLoop
				}
			}
			if c == nil {
				if t.Usage != nil {
					t.Usage(t)
					return errors.New("invalid flag")
				}
			}
			err := c.run(args[i:])
			if err == UsageError {
				if c.Usage != nil {
					c.Usage(c)
					return usageErrorFire
				}
				fmt.Println(c.main+": "+c.Name, ":", c.Desc)
				c.PrintFlags()
				return usageErrorFire
			}
			return err
		}
	}
	if t.Run != nil {
		return t.Run(t, args[i:])
	}
	return UsageError
}

func (cmds Cmds) Run(m *C) error {
	m.Cmds = cmds
	return cmds.RunWith(m, os.Args)
}

func (cmds Cmds) RunWith(m *C, args []string) error {
	if m.Name == "" {
		m.Name = filepath.Base(args[0])
		if runtime.GOOS == "windows" {
			if ext := filepath.Ext(m.Name); ext != "" {
				m.Name = m.Name[:len(m.Name)-len(ext)]
			}
		}
	}
	m.Cmds = cmds
	m.main = m.Name
	err := m.run(args)
	if err == UsageError {
		if m.Usage != nil {
			m.Usage(m)
			return UsageError
		}
		fmt.Println(m.Name, ":", m.Desc)
		m.PrintFlags()
		m.PrintCommands()
		return err
	}
	return err
}
