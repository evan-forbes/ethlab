package compile

import (
	"fmt"
	"io/ioutil"
	"os"
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
	pkg := c.String("pkg")
	if pkg == "" {
		wkdir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failure to compile, no package specified and could not read working directory")
		}
		pkg = isolatePkg(wkdir) // not the best idea because of of go's strict package naming
	}
	var (
		types []string
		abis  []string
		bins  []string
	)
	for path, con := range contracts {
		pathData := strings.Split(path, ":")
		typeName := pathData[len(pathData)-1]
		// don't bother generating code for interfaces and libraries
		if len(con.Bin) < 3 {
			// generate interface code instead
			continue
		}
		types = append(types, typeName)
		abis = append(abis, con.Abi)
		bins = append(bins, con.Bin)
	}
	code, eventCode, err := bind.Bind(
		types,
		abis,
		bins,
		pkg,
	)
	if err != nil {
		return errors.Wrap(err, "failure to generate bindings during compilation")
	}
	// write to file
	filename := fmt.Sprintf("%s_gen.go", pkg)
	err = ioutil.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return errors.Wrap(err, "failure to write file:")
	}
	filename = fmt.Sprintf("%s_events_gen.go", pkg)
	err = ioutil.WriteFile(filename, []byte(eventCode), 0644)
	if err != nil {
		return errors.Wrap(err, "failure to write file:")
	}
	return nil
}

func isolatePkg(path string) string {
	// todo: include windows support "\"
	splt := strings.Split(path, "/")
	return splt[len(splt)-1]
}
