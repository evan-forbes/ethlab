package main

import (
	"log"
	"os"

	"github.com/evan-forbes/ethlab/cmd/abigen"
	"github.com/evan-forbes/ethlab/cmd/boot"
	"github.com/evan-forbes/ethlab/cmd/compile"
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

// test subscribing to heads without serving and figure up what's blocking

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

	// bootFlags are the flags for boo
	bootFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "ip",
			Value: "",
			Usage: "specify ip address for server to use (can also enter in config file). default == 127.0.0.1:84384",
		},
		&cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "*optional* path to config file (.json)",
		},
	}

	// bootFlags are the flags for boo
	compileFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "pkg, p",
			Value: "",
			Usage: "specify the package name",
		},
		// &cli.StringFlag{
		// 	Name:  "config, c",
		// 	Value: "",
		// 	Usage: "*optional* path to config file (.json)",
		// },
	}

	// subcommands
	app.Commands = []*cli.Command{
		{
			Name:   "abigen",
			Usage:  "generate ethlab and go-ethereum go bindings for an abi or bin",
			Flags:  abiFlags,
			Action: abigen.Generate,
		},
		{
			Name:   "boot",
			Usage:  "start serving access to a configured Thereum blockchain",
			Flags:  bootFlags,
			Action: boot.Boot,
		},
		{
			Name:   "compile",
			Usage:  "combine solc and abigen with a simple naming scheme",
			Flags:  compileFlags,
			Action: compile.Compile,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
