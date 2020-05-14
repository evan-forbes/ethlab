package compile

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func All(path string) {
	// set default path to the current dir
	if path == "" {
		path = "."
	}

}

// openAllFiles searches for and reads solidity source files up to four directories deep
// to strings
func openAllFiles(path string, recCount int) (out []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return out, err
	}
	for _, file := range files {
		// don't go over the recurse limit
		if recCount > 3 {
			return out, nil
		}
		pathToFile := fmt.Sprintf("%s/%s", path, file.Name())
		// recursively search directories
		if file.IsDir() {
			nextSources, err := openAllFiles(pathToFile, recCount+1)
			if err != nil {
				fmt.Println(fmt.Println("failure to read files:", err))
				continue
			}
			out = append(out, nextSources...)
			continue
		}
		// check if the file is .sol
		if strings.Contains(file.Name(), ".sol") {
			// read the file

			source, err := ioutil.ReadFile(pathToFile)
			if err != nil {
				fmt.Println("failure to read file:", err)
				continue
			}
			out = append(out, string(source))
		}
	}
	return out, nil
}
