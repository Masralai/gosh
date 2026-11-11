package handlers

import "github.com/urfave/cli/v3"

// All returns every supported command, grouped by category.
func All() []*cli.Command {
	return []*cli.Command{
		// Shell utilities
		Cli(), Boom(), Echo(), Clear(), Help(), History(), Env(), Which(), Cd(), Pwd(), Exit(),
		// File operations
		Ls(), Mkdir(), Rm(), Touch(), Mv(), Cp(), Dir(), Cat(), Info(), Find(),
		// System monitoring
		Ps(), Ut(), Sys(), Mu(), Du(), Kill(),
		// Text processing
		Grep(), Head(), Tail(), Wc(), SortCmd(),
		// Networking
		Ping(),
		// Archive
		Zip(), Unzip(),
	}
}
