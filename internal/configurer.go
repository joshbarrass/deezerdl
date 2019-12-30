package internal

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
)

func Configure(opts docopt.Opts, config *Configuration) {
	// set method
	if _, ok := opts["set"]; ok {
		if set, err := opts.Bool("set"); err != nil {
			logrus.Fatalf("failed to parse args: %s", err)
		} else if set {
			configureSet(opts, config)
			return
		}
	}
}

func configureSet(opts docopt.Opts, config *Configuration) {
	// set default format
	if _, ok := opts["DefaultFormat"]; ok {
		if selection, err := opts.Bool("DefaultFormat"); err != nil {
			logrus.Fatalf("failed to parse args: %s", err)
		} else if selection {
			var format string
			format, err = opts.String("<fmt>")
			if err != nil {
				logrus.Fatalf("failed to parse args: %s", err)
			}
			FormatStringToFormat(format)
			config.DefaultFormat = format
			config.SaveConfig()
			fmt.Printf("Set default format to %s\n", format)
			return
		}
	}
}
