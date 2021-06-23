package main

import (
	"github.com/MarilynFranklin/waypoint-plugin-levant/platform"
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

func main() {
	// sdk.Main allows you to register the components which should
	// be included in your plugin
	// Main sets up all the go-plugin requirements

	sdk.Main(sdk.WithComponents(
		// Comment out any components which are not
		// required for your plugin
		&platform.Platform{},
	))
}
