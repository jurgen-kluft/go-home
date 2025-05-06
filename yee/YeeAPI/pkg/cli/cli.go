// package cli contains code for the command line client interface
package cli

import "flag"

// CommonFlags are flags common to all subcommands.
type CommonFlags struct {
	Interface string
	// Verbose bool
}

func (c *CommonFlags) Register(fs *flag.FlagSet) {
	fs.StringVar(&c.Interface, "interface", "enp0s31f6", "name of network interface")
	// fs.BoolVar(&c.Verbose, "v", false, "verbose output")
}
