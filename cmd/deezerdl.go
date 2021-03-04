package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joshbarrass/deezerdl/internal"
	"github.com/sirupsen/logrus"
)

const VERSION = "0.2.3"

const doc = `deezerdl

Usage:
  deezerdl login <arl>
  deezerdl download track <ID> [-f <fmt> | --format=<fmt>]
  deezerdl download album <ID> [-f <fmt> | --format=<fmt>]
  deezerdl config set DefaultFormat <fmt>

Options:
  -f --format=<fmt>    Specifies the download format. Valid options are FLAC, MP3_320, MP3_256.
`

var config *internal.Configuration
var envConfig *internal.EnvConfig

func main() {
	var err error
	config, err = internal.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load config: %s", err)
	}

	envConfig, err = internal.GetEnvConfig()
	if err != nil {
		logrus.Warnf("failed to load environment variables: %s", err)
		envConfig = internal.NewEnvConfig()
	}

	argv := os.Args[1:]

	parser := &docopt.Parser{
		HelpHandler: docopt.PrintHelpOnly,
	}
	opts, err := parser.ParseArgs(doc, argv, VERSION)
	if envConfig.DebugMode {
		logrus.Info(opts)
	}
	if err != nil {
		// err is "" if no valid argument
		if err.Error() != "" {
			logrus.Fatalf("failed to create parser: %s", err)
		} else {
			return
		}
	}

	// login method
	if _, ok := opts["login"]; ok {
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
	}

	// download method
	if _, ok := opts["download"]; ok {
		if dl, err := opts.Bool("download"); err != nil {
			logrus.Fatalf("failed to parse args: %s", err)
		} else if dl {
			internal.Download(opts, config)
			return
		}
	}

	// config method
	if _, ok := opts["config"]; ok {
		if cfg, err := opts.Bool("config"); err != nil {
			logrus.Fatalf("failed to parse args: %s", err)
		} else if cfg {
			internal.Configure(opts, config)
			return
		}
	}
}
