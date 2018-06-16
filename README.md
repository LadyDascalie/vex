# Vex

*Yet another vexing Makefile replacement you didn't want nor need.*

Vex aims to replace error-prone Makefiles for teams that don't require the full power of `make`, but need the easy access build scripts.

Vex is `toml` based, and its syntax is as obvious as possible, so it's really easy to pickup.

Command line helpers come bundled to help you be productive quickly.

## Install:

`go get github.com/ladydascalie/vex`

## Usage:

```text
Usage:
  vex [flags]
  vex [command]

Available Commands:
  add         add a new command named into the VexFile
  help        Help about any command
  init        Create a default VexFile
  list        List available commands
  run         Runs a command specified by name

Flags:
  -h, --help   help for vex

Use "vex [command] --help" for more information about a command.
```

## Example VexFile: 

```toml
[vex]

	# This is the same as a Makefile all command.
	# Commands are comma separated, ex: all = ["cmd1", "cmd2"].	
	all = ["build"]
	# No need to specify the pre and post commands here, they will run regardless.


	# Pre commands run before anything else is ran.
	# You may use this block to perform any initialization you need.
	[[vex.pre]]
	name = "setup"
	desc = "Does x and y"
	runs = "..."

	# Main command block, this is where you should write most commands.
	# For convenience 'vex add command_name' will add new commands in here for you.
	[[vex.cmd]]
	name = "Build"
	desc = "Performs the build tasks"
	runs = '''
		echo "VexFiles support multi-line commands!"
		echo "Hello from the VexFile!"
	'''

	# Post commands will run after everything else has ran.
	# This is where you should write your cleanup scripts if needed.
	[[vex.post]]
	name = "Cleanup"
	desc = "Does x and y to cleanup after itself"
	runs = "..."

```