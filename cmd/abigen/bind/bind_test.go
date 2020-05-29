package bind

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

// findAllSolFiles searches for and returns the paths to solidity source files up to
// four directories deep
func findAllFiles(path, substr string, recCount int) (out []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return out, errors.Wrap(err, "failure to findAllSolFiles:")
	}
	for _, file := range files {
		// don't go over the recurse limit
		if recCount > 3 {
			fmt.Println("ignoring some files due to directory depth greater than 4")
			return out, nil
		}
		pathToFile := fmt.Sprintf("%s/%s", path, file.Name())
		// recursively search directories
		if file.IsDir() {
			next, err := findAllFiles(pathToFile, substr, recCount+1)
			if err != nil {
				fmt.Println(fmt.Println("failure to find files:", err))
				continue
			}
			out = append(out, next...)
			continue
		}
		// check if the file is .sol
		if strings.Contains(file.Name(), substr) {
			// add file path to output
			out = append(out, pathToFile)
		}
	}
	return out, nil
}

func TestInterfaceGen(t *testing.T) {
	files, err := findAllFiles("/home/evan/go/src/github.com/evan-forbes/bindings/uniswap/v2-core/uniswap", ".go", 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(files)
}
