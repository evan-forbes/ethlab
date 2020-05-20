package compile

import "github.com/urfave/cli/v2"

func Compile(c *cli.Context) error {
	// find and compile all .sol files in the provided directory
	path := "."
	contracts, err := All(path)
	if err != nil {
		return err
	}
	for _, con := range contracts {
		// use bind to generate go bindings
	}
	// unmarshal solc output into solcOutput
	// save the abi and bin seperately
	// use the abi and bin to generate go bindings

}
