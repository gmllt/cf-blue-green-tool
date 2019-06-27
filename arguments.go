package main

import (
	"flag"
)

// Arguments
type Arguments struct {
	Action        string
	DeleteOldApps bool
	Manifest      string
}

// NewArgs : Get arguments
func NewArgs(osArgs []string) Arguments {
	args := Arguments{}
	args.Action = extractAction(osArgs)

	// Only use FlagSet so that we can pass string slice to Parse
	f := flag.NewFlagSet("blue-green-tool", flag.ExitOnError)

	f.BoolVar(&args.DeleteOldApps, "delete-old-apps", false, "")
	f.StringVar(&args.Manifest, "f", "manifest.yml", "")

	f.Parse(extractArgs(osArgs))

	return args
}

// indexOfAction : Get index of param action
func indexOfAction(osArgs []string) int {
	index := 0
	for i, arg := range osArgs {
		if arg == "blue-green-tool" || arg == "bgt" {
			index = i + 1
			break
		}
	}
	if len(osArgs) > index {
		return index
	}
	return -1
}

// extractAction : extract param action
func extractAction(osArgs []string) string {
	// Assume an app name will be passed - issue #27
	index := indexOfAction(osArgs)
	if index >= 0 {
		return osArgs[index]
	}
	return ""
}

// extractArgs : extract arguments
func extractArgs(osArgs []string) []string {
	index := indexOfAction(osArgs)
	if index >= 0 && len(osArgs) > index+1 {
		return osArgs[index+1:]
	}

	return []string{}
}
