package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joshbarrass/deezerdl/internal"
	"github.com/sirupsen/logrus"
)

const VERSION = "0.0.1"

const doc = `deezerdl

Usage:
  deezerdl login <arl>
  deezerdl download track <ID>
`

var config *internal.Configuration

func main() {
	var err error
	config, err = internal.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load config: %s", err)
	}

	argv := os.Args[1:]

	parser := &docopt.Parser{
		HelpHandler:  docopt.PrintHelpOnly,
		OptionsFirst: true,
	}
	opts, err := parser.ParseArgs(doc, argv, VERSION)
	if err != nil {
		logrus.Info(argv)
		logrus.Fatalf("failed to create parser: %s", err)
	}
	// fmt.Println(opts)

	// login method
	if login, err := opts.Bool("login"); err != nil {
		logrus.Fatalf("failed to parse args: %s", err)
	} else if login {
		config.ARLCookie, err = opts.String("<arl>")
		if err != nil {
			logrus.Fatalf("failed to set arl cookie: %s", err)
		}
		config.SaveConfig()
		fmt.Println("Saved arl! You can now use the rest of the program.")
		return
	}

	// download method
	if dl, err := opts.Bool("download"); err != nil {
		logrus.Fatalf("failed to parse args: %s", err)
	} else if dl {
		internal.Download(&opts, config)
	}
}
