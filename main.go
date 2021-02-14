package main

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/tn-platform/waypoint-plugin-levant/platform"
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
