package scaffold

import (
	"fmt"

	"github.com/fatih/color"
)

// VexFile defines the VexFile structure to follow
type VexFile struct {
	Vex struct {
		Env  []string       `json:"env" toml:"env"`
		All  []string       `json:"runs" toml:"runs"`
		Pre  []FormattedCmd `json:"pre" toml:"pre"`
		Cmd  []FormattedCmd `json:"cmd" toml:"cmd"`
		Post []FormattedCmd `json:"post" toml:"post"`
	} `json:"vex" toml:"vex"`
}

type FormattedCmd struct {
	Name string `json:"name" toml:"name"`
	Desc string `json:"desc" toml:"desc"`
	Runs string `json:"runs" toml:"runs"`
}

func (f *FormattedCmd) Print() {
	name := color.HiGreenString("command %s:", f.Name)
	desc := fmt.Sprintf("\n\t%s\n", f.Desc)
	fmt.Println(name, desc)
}

// vexTpl defines the basic VexFile template
const VexTpl = `
	[vex]
	all = [""]
	[[vex.cmd]]
	name = ""
	desc = "" # optional
	runs = ""
`
