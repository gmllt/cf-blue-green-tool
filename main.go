package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"log"
	"os"
)

var PluginVersion string

type CfPlugin struct {
	Connection plugin.CliConnection
	Deploy     BlueGreenDeploy
}

func (p *CfPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	arguments := NewArgs(args)

	p.Connection = cliConnection

	p.Deploy.Setup(cliConnection)

	var inError = false
	if arguments.Action == "" {
		inError = true
	} else {
		inError = true
		if arguments.Action == "deploy" {
			p.RunDeploy(arguments)
			inError = false
		}
		if arguments.Action == "green" {
			p.RunGreen(arguments)
			inError = false
		}
		if arguments.Action == "rollback" {
			p.RunRollback(arguments)
			inError = false
		}
		if arguments.Action == "approve" {
			p.RunApprove(arguments)
			inError = false
		}
	}
	if inError {
		log.Fatal("Action must be provided and must be one of 'deploy' or 'green' or 'rollback' or 'approve'.")
	}
}

func (p *CfPlugin) GetMetadata() plugin.PluginMetadata {
	var major, minor, build int
	fmt.Sscanf(PluginVersion, "%d.%d.%d", &major, &minor, &build)

	return plugin.PluginMetadata{
		Name: "blue-green-tool",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: build,
		},
		Commands: []plugin.Command{
			{
				Name:     "blue-green-tool",
				Alias:    "bgt",
				HelpText: "BlueGreen deployment tool",
				UsageDetails: plugin.Usage{
					Usage: "blue-green-tool ACTION [-f MANIFEST_FILE] [--delete-old-apps]",
					Options: map[string]string{
						"f":               "Path to manifest",
						"delete-old-apps": "Delete old app instance(s)",
					},
				},
			},
		},
	}
}

func (p *CfPlugin) RunDeploy(arguments Arguments) () {
	p.RunGreen(arguments)
	p.RunApprove(arguments)
}

func (p *CfPlugin) RunGreen(arguments Arguments) {
	var manifestParse ManifestParse
	manifest, _ := manifestParse.parseFile(arguments.Manifest)
	newManifest, _ := manifestParse.parseFile(arguments.Manifest)

	// Push new app
	p.Deploy.pushNewApp(newManifest)

	// Map new app
	p.Deploy.MapNewApp(newManifest, manifest)
}

func (p *CfPlugin) RunRollback(arguments Arguments) {
	var manifestParse ManifestParse
	newManifest, _ := manifestParse.parseFile(arguments.Manifest)
	newManifest.SuffixApp("new")
	// Unmap old route from new app
	p.Deploy.UnMapNewRoute(newManifest)
	p.Deploy.RemoveOldApp(newManifest)

}

func (p *CfPlugin) RunApprove(arguments Arguments) {
	var manifestParse ManifestParse
	manifest, _ := manifestParse.parseFile(arguments.Manifest)
	newManifest, _ := manifestParse.parseFile(arguments.Manifest)
	oldManifest, _ := manifestParse.parseFile(arguments.Manifest)
	oldManifest.SuffixApp("old")
	newManifest.SuffixApp("new")

	// Unmap old route from new app
	p.Deploy.UnMapNewRoute(newManifest)

	// Rename old app
	p.Deploy.RenameApp(oldManifest, manifest)

	p.Deploy.MapOldApp(oldManifest)

	p.Deploy.UnMapOldApp(oldManifest, manifest)

	p.Deploy.RenameApp(manifest, newManifest)

	if arguments.DeleteOldApps {
		p.Deploy.RemoveOldRoute(oldManifest)
		p.Deploy.RemoveOldApp(oldManifest)
	}
}

func main() {

	log.SetFlags(0)

	p := CfPlugin{
		Deploy: BlueGreenDeploy{
			ErrorFunc: func(message string, err error) {
				log.Fatalf("%v - %v", message, err)
			},
			Out: os.Stdout,
		},
	}

	plugin.Start(&p)
}
