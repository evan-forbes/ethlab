package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/pkg/errors"
)

/*
the most useful interface would be for the apis of that contract
*/

func All(path string) (map[string]contract, error) {
	if path == "" {
		p, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		path = p
	}
	filepaths, err := findAllFiles(path, ".sol", 0)
	if err != nil {
		return nil, err
	}
	output, err := solidity("solc", filepaths...)
	if err != nil {
		return nil, err
	}
	return output.Contracts, nil
}

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
				fmt.Println("failure to find files:", err)
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

type contract struct {
	BinRuntime                                  string `json:"bin-runtime"`
	SrcMapRuntime                               string `json:"srcmap-runtime"`
	Bin, SrcMap, Abi, Devdoc, Userdoc, Metadata string
	Hashes                                      map[string]string
}

// --combined-output format
type solcOutput struct {
	Contracts map[string]contract
	Version   string
}

// FORKED FROM GO ETHEREUM SEE GNU Lesser General Public License as published by
// the Free Software Foundation vvv

// solidity compiles all given solidity source files.
func solidity(solc string, sourcefiles ...string) (solcOutput, error) {
	if len(sourcefiles) == 0 {
		return solcOutput{}, errors.New("solc: no source files")
	}
	source, err := slurpFiles(sourcefiles)
	if err != nil {
		return solcOutput{}, err
	}
	s, err := compiler.SolidityVersion(solc)
	if err != nil {
		return solcOutput{}, err
	}
	args := []string{
		"--combined-json", "bin,bin-runtime,srcmap,srcmap-runtime,abi,userdoc,devdoc,metadata,hashes",
		"--optimize", // code optimizer switched on
		// "--allow-paths", "., ../", // default to support relative paths //
		"--",
	}
	cmd := exec.Command(s.Path, append(args, sourcefiles...)...)
	return run(cmd, source)
}

// run executes the command with the source string as
func run(cmd *exec.Cmd, source string) (out solcOutput, err error) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return out, fmt.Errorf("solc: %v\n%s", err, stderr.Bytes())
	}
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return out, err
	}
	return out, nil
}

func slurpFiles(files []string) (string, error) {
	var concat bytes.Buffer
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		concat.Write(content)
	}
	return concat.String(), nil
}

// // All compiles all given solidity files in a given path
// func All(path string) (map[string]*compiler.Contract, error) {
// 	// set default path to the current dir
// 	if path == "" {
// 		path = "."
// 	}
// 	sources, err := openAllFiles(path, 0)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failure to compile.All:")
// 	}
// 	return compiler.Solidity("solc", sources...)
// }

// // openAllFiles searches for and reads solidity source files up to four directories deep
// // to strings
// func openAllFiles(path string, recCount int) (out []string, err error) {
// 	files, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		return out, errors.Wrap(err, "failure to openAllFiles:")
// 	}
// 	for _, file := range files {
// 		// don't go over the recurse limit
// 		if recCount > 3 {
// 			fmt.Println("ignoring some files due to directory depth greater than 4")
// 			return out, nil
// 		}
// 		pathToFile := fmt.Sprintf("%s/%s", path, file.Name())
// 		// recursively search directories
// 		if file.IsDir() {
// 			nextSources, err := openAllFiles(pathToFile, recCount+1)
// 			if err != nil {
// 				fmt.Println(fmt.Println("failure to read files:", err))
// 				continue
// 			}
// 			out = append(out, nextSources...)
// 			continue
// 		}
// 		// check if the file is .sol
// 		if strings.Contains(file.Name(), ".sol") {
// 			// read the file

// 			source, err := ioutil.ReadFile(pathToFile)
// 			if err != nil {
// 				fmt.Println("failure to read file:", err)
// 				continue
// 			}
// 			out = append(out, string(source))
// 		}
// 	}
// 	return out, nil
// }
