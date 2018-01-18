package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/ladydascalie/vex/scaffold"
	"github.com/spf13/cobra"
)

// VexFileName defines the name of our config file
const VexFileName = "VexFile"

var vexfile scaffold.VexFile
var boldGreen = color.New(color.FgHiGreen, color.Bold)

func init() {
	log.SetFlags(0)
}

func main() {
	rootCmd := &cobra.Command{
		Use: "vex",
		Run: runCommands,
	}

	initCmd := &cobra.Command{
		Use:        "init",
		Short:      "Create a default VexFile",
		Long:       `init creates a default staring vexfile`,
		Run:        initVexFile,
		SuggestFor: []string{"vex"},
	}

	listCmd := &cobra.Command{
		Use:        "list",
		Short:      "List available commands",
		Run:        listCommands,
		SuggestFor: []string{"vex"},
	}
	runCmd := &cobra.Command{
		Use:        "run",
		Short:      "Runs a command specified by name",
		Run:        runNamedCommands,
		SuggestFor: []string{"vex"},
	}
	addCmd := &cobra.Command{
		Use:        "add",
		Short:      "add a new command named into the VexFile",
		Run:        addNewCommand,
		SuggestFor: []string{"vex"},
	}

	rootCmd.AddCommand(initCmd, listCmd, runCmd, addCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
func decodeVexFile() {
	_, err := toml.DecodeFile(VexFileName, &vexfile)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Fatalln("VexFile not found.\nPlease ensure a VexFile is present in your directory, or add one with vex init")
		}

		log.Fatalf("%s\n\tPlease refer to the documentation, or generate an example VexFile with vex init", color.HiRedString("Invalid VexFile"))
	}
}

func addNewCommand(cmd *cobra.Command, args []string) {
	decodeVexFile()
	for _, a := range args {
		vexfile.Vex.Cmd = append(vexfile.Vex.Cmd, scaffold.FormattedCmd{
			Name: a,
			Runs: "",
			Desc: "",
		})
	}
	f, err := os.OpenFile("VexFile", os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot open vexfile: %v", err)
	}
	enc := toml.NewEncoder(f)
	enc.Indent = "\t"
	enc.Encode(vexfile)
}

// prints all commands
func listCommands(cmd *cobra.Command, args []string) {
	decodeVexFile()
	if len(vexfile.Vex.Pre) > 0 {
		for _, pre := range vexfile.Vex.Pre {
			pre.Print()
		}
	}
	if len(vexfile.Vex.Cmd) > 0 {
		for _, cmd := range vexfile.Vex.Cmd {
			cmd.Print()
		}
	}
	if len(vexfile.Vex.Post) > 0 {
		for _, post := range vexfile.Vex.Post {
			post.Print()
		}
	}
}

// run all commands
func runCommands(cmd *cobra.Command, args []string) {
	decodeVexFile()
	executePreCommands()
	defer executePostCommands()

	// if run order is specified, execute that
	if vexfile.Vex.All != nil {
		runNamedCommands(cmd, vexfile.Vex.All)
		return
	}

	if len(vexfile.Vex.Cmd) > 0 {
		for _, cmd := range vexfile.Vex.Cmd {
			color.HiGreen("running: %s", cmd.Name)
			execute(cmd.Runs)
		}
	}
}

// run a named command
func runNamedCommands(cmd *cobra.Command, args []string) {
	decodeVexFile()
	executePreCommands()
	defer executePostCommands()

	if len(vexfile.Vex.Cmd) > 0 {
		for _, cmd := range vexfile.Vex.Cmd {
			for _, a := range args {
				if cmd.Name == a {
					color.HiGreen("running: %s", cmd.Name)
					execute(cmd.Runs)
				}
			}
		}
	}
	return
}

func executePostCommands() {
	if len(vexfile.Vex.Post) > 0 {
		boldGreen.Println("Executing post commands")
		for _, post := range vexfile.Vex.Post {
			color.HiGreen("running: %s", post.Name)
			execute(post.Runs)
		}
	}
}

func executePreCommands() {
	if len(vexfile.Vex.Pre) > 0 {
		color.HiGreen("Executing pre commands")
		for _, pre := range vexfile.Vex.Pre {
			color.HiGreen("running: %s", pre.Name)
			execute(pre.Runs)
		}
	}
}

// initVexFile creates a new VexFile
func initVexFile(cmd *cobra.Command, args []string) {
	_, err := os.Stat(VexFileName)
	if err == nil {
		return
	}

	if err := ioutil.WriteFile(VexFileName, []byte(scaffold.VexTpl), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return
}

// vex the command
func execute(command string) {
	// invoke the naked shell
	cmd := exec.Command("/bin/sh", "-c", command)
	// only pipe the outputs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run and log any errors
	if err := cmd.Run(); err != nil {
		log.Fatalf("cannot execute command: %v", err)
	}
}
