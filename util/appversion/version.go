// Package appversion processes version information for app.
package appversion

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

type versionInfo struct {
	AppVersion  string `json:"version"`
	GitRevision string `json:"commit"`
	BuildTime   string `json:"build_time"`
}

func (v versionInfo) print(w io.Writer) {
	json.NewEncoder(w).Encode(v)
}

// Print prints version information to stdout.
func Print() {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		var v versionInfo
		v.AppVersion = buildInfo.Main.Version // default to mod's version

		for _, s := range buildInfo.Settings {
			switch s.Key {
			case "-tags":
				if v.AppVersion != "" { // update with -tags if doesn't exist
					v.AppVersion = s.Value
				}
			case "vcs.revision":
				v.GitRevision = s.Value
			case "vcs.time":
				v.BuildTime = s.Value
			}
		}
		v.print(os.Stdout)
	} else {
		fmt.Println("No version information available.")
	}

	os.Exit(0)
}

// PrintExit prints version to stdout and os.Exit(0) if `-version` flag is true.
// Call it after `flag.Parse()`.
func PrintExit() {
	if VersionFlag() {
		Print()
		os.Exit(0)
	}
}
