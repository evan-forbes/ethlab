package compile

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/evan-forbes/ethlab/cmd/abigen/bind"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// Compile slurps all solidity files in the provided path, compiles them using
// the solc installed on the system, and generates go bindings. By default,
// Compile will overwrite .abi, .bin, and _gen.go files.
func Compile(c *cli.Context) error {
	// find and compile all .sol files in the provided directory
	path := "."
	contracts, err := All(path)
	if err != nil {
		return err
	}
	for path, con := range contracts {
		pathData := strings.Split(path, ":")
		typeName := pathData[len(pathData)-1]
		code, err := bind.Bind(
			[]string{typeName},
			[]string{con.Abi},
			[]string{con.Bin},
			"ens",
		)
		if err != nil {
			return errors.Wrap(err, "failure to generate bindings during compilation")
		}
		// write to file
		filename := fmt.Sprintf("%s_gen.go", typeName)
		err = ioutil.WriteFile(filename, []byte(code), 0644)
		if err != nil {
			return errors.Wrap(err, "failure to write file:")
		}
	}
	return nil
}
