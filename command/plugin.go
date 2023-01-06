package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"log"
	"os"
)

// SyslogDrainPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at "code.cloudfoundry.org/cli/plugin/plugin.go"
type SyslogDrainPlugin struct{}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second parameter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
func (c *SyslogDrainPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	l := log.New(os.Stderr, "", 0)
	switch args[0] {
	case "list-syslog-drains":
		err := ListSyslogDrains(cliConnection, l, os.Stdout)
		handleIfErr(err)
	case "list-org-syslog-drains":
		err := ListOrgSyslogDrains(cliConnection, l, os.Stdout)
		handleIfErr(err)
	case "list-space-syslog-drains":
		err := ListSpaceSyslogDrains(cliConnection, l, os.Stdout)
		handleIfErr(err)
	case "CLI-MESSAGE-UNINSTALL":
		os.Exit(0)
	default:
		fmt.Printf("unsupported command %s\n", args[0])
		os.Exit(1)
	}
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as separate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (c *SyslogDrainPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "SyslogDrainPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 7,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "list-syslog-drains",
				HelpText: "List all syslog drains across orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "list-syslog-drains",
				},
			},
			{
				Name:     "list-org-syslog-drains",
				HelpText: "List all syslog drains in the currently targeted org",
				UsageDetails: plugin.Usage{
					Usage: "list-org-syslog-drains",
				},
			},
			{
				Name:     "list-space-syslog-drains",
				HelpText: "List all syslog drains in the currently targeted space",
				UsageDetails: plugin.Usage{
					Usage: `list-space-syslog-drains`,
				},
			},
		},
	}
}

func handleIfErr(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
