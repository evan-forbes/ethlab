package main

import (
	"log"
	"os"

	"github.com/evan-forbes/ethlab/cmd"
	"github.com/urfave/cli/v2"
)

// TODO:
/*
 - port the abigen command from buddy
	. alter to provide super easy muxing
 - txpool that can link transactions
 - pause and set break points for the chain
 - option to stop the chain when nothing is happening
 - txpool should reccomend a gas price? keep track of highest lowest median?
 - test txpool with a massive number of txs
*/

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "ethlab"

	// abiFlags are flags for the subcommand abigen
	abiFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "abi, a",
			Value: ".",
			Usage: "path to abi (.json or .abi)",
			// Destination: &abiPath,
		},
		&cli.StringFlag{
			Name:  "bin, b",
			Value: ".",
			Usage: "path to contract binary (usually a .bin)",
			// Destination: &binPath,
		},
		&cli.StringFlag{
			Name:  "type, t",
			Value: "",
			Usage: "specify the main type",
			// Destination: &tp,
		},
		&cli.StringFlag{
			Name:  "pkg, p",
			Value: "",
			Usage: "specify the package name",
			// Destination: &tp,
		},
		&cli.StringFlag{
			Name:  "out, o",
			Value: "",
			Usage: "specify the output file name (default = type_gen.go",
			// Destination: &tp,
		},
	}

	bootFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "address, a",
			Value: "",
			Usage: "specify address for server to use (can also enter in config file)",
		},
		&cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "path to config file (.json)",
		},
	}

	// subcommands
	app.Commands = []*cli.Command{
		{
			Name:   "abigen",
			Usage:  "generate interface and mock friendly go bindings",
			Action: cmd.ABIgen,
			Flags:  abiFlags,
		},
		{
			Name:   "boot",
			Usage:  "start serving access to a configured blockchain",
			Action: cmd.Boot,
			Flags:  bootFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
