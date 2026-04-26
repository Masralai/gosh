package handlers

import "github.com/urfave/cli/v3"

// All returns every supported command, grouped by category.
func All() []*cli.Command {
	return []*cli.Command{
		// Shell utilities
		Cli(), Boom(), Echo(), Cd(), Pwd(), Exit(),
		// File operations
		Ls(), Mkdir(), Rm(), Touch(), Mv(), Cp(), Dir(), Cat(), Info(),
		// System monitoring
		Ps(), Ut(), Sys(), Mu(), Du(), Kill(),
		// Text processing
		Grep(), Head(), Tail(),
		// Networking
		Ping(),
		// Archive
		Zip(), Unzip(),
	}
}
