// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind

import "github.com/ethereum/go-ethereum/accounts/abi"

// tmplData is the data structure required to fill the binding template.
type tmplData struct {
	Package   string                   // Name of the package to place the generated file in
	Contracts map[string]*tmplContract // List of contracts to generate into this file
	Libraries map[string]string        // Map the bytecode's link pattern to the library name
	Structs   map[string]*tmplStruct   // Contract struct type definitions
	Events    map[string]*tmplEvent
}

// TODO: don't do anything for log ids that have already been generated
// do something about the shitload of files generated for more complex projects

// tmplContract contains the data needed to generate an individual contract binding.
type tmplContract struct {
	Type        string                 // Type name of the main contract binding
	InputABI    string                 // JSON ABI used as the input to generate the binding from
	InputBin    string                 // Optional EVM bytecode used to denetare deploy code from
	FuncSigs    map[string]string      // Optional map: string signature -> 4-byte signature
	Constructor abi.Method             // Contract constructor for deploy parametrization
	Calls       map[string]*tmplMethod // Contract calls that only read state data
	Transacts   map[string]*tmplMethod // Contract calls that write state data
	Events      map[string]*tmplEvent  // Contract events accessors
	Libraries   map[string]string      // Same as tmplData, but filtered to only keep what the contract needs
	Library     bool                   // Indicator whether the contract is a library
}

// tmplMethod is a wrapper around an abi.Method that contains a few preprocessed
// and cached data fields.
type tmplMethod struct {
	Original   abi.Method // Original method as parsed by the abi package
	Normalized abi.Method // Normalized version of the parsed method (capitalized names, non-anonymous args/returns)
	Structured bool       // Whether the returns should be accumulated into a struct
}

// tmplEvent is a wrapper around an a
type tmplEvent struct {
	Original   abi.Event // Original event as parsed by the abi package
	Normalized abi.Event // Normalized version of the parsed fields
	Topic      string
	Type       string
}

// tmplField is a wrapper around a struct field with binding language
// struct type definition and relative filed name.
type tmplField struct {
	Type    string   // Field type representation depends on target binding language
	Name    string   // Field name converted from the raw user-defined field name
	SolKind abi.Type // Raw abi type information
}

// tmplStruct is a wrapper around an abi.tuple contains a auto-generated
// struct name.
type tmplStruct struct {
	Name   string       // Auto-generated struct name(before solidity v0.5.11) or raw name.
	Fields []*tmplField // Struct fields definition depends on the binding language.
}

// tmplSourceGo is the Go source template use to generate the contract binding
// based on.
const tmplSourceGo = `
{{$pkg := .Package}}
package {{$pkg}}

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

{{$structs := .Structs}}
{{range $contract := .Contracts}}

// {{.Type}} is a wrapper around bind.BoundContract, enforcing type checking and including
// QoL helper methods
type {{.Type}} struct {
	bind.BoundContract
}

// New{{.Type}} creates a new instance of {{.Type}}, bound to a specific deployed contract.
func New{{.Type}}(address common.Address, backend bind.ContractBackend) (*{{.Type}}, error) {
	a, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, a, backend, backend, backend)
	return &{{.Type}}{*contract}, nil
}

{{if .InputBin}}
//////////////////////////////////////////////////////
//		Deployment
////////////////////////////////////////////////////

// Deploy installs {{.Type}} to an ethereum node via the user provided
// by implementing module.Delpoyer
func Deploy(u *module.User) (addr common.Address, err error) {
	// ****************************************************
	////////  INSERT MODULE DEPLOYMENT CODE HERE   ////////
	// ***************************************************
	return addr, err
}

// Deploy{{.Type}} deploys a new Ethereum contract, binding an instance of {{.Type}} to it.
func Deploy{{.Type}}(auth *bind.TransactOpts, backend bind.ContractBackend {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type $structs}}{{end}}) (common.Address, *types.Transaction, *{{.Type}}, error) {
  parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
  if err != nil {
	return common.Address{}, nil, nil, err
  }
  address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex({{.Type}}Bin), backend {{range .Constructor.Inputs}}, {{.Name}}{{end}})
  if err != nil {
	return common.Address{}, nil, nil, err
  }
  return address, tx, &{{.Type}}{*contract}, nil
}
{{end}}

//////////////////////////////////////////////////////
//		Structs
////////////////////////////////////////////////////
{{range $structs}}
// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
}
{{end}}

//////////////////////////////////////////////////////
//		Data Calls
////////////////////////////////////////////////////

{{range .Calls}}
// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
// - Solidity: {{formatmethod .Original $structs}}
func (_{{$contract.Type}} *{{$contract.Type}}) {{.Normalized.Name}}(opts *bind.CallOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} error) {
	{{if .Structured}}ret := new(struct{
		{{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}}
		{{end}}
	}){{else}}var (
		{{range $i, $_ := .Normalized.Outputs}}ret{{$i}} = new({{bindtype .Type $structs}})
		{{end}}
	){{end}}
	out := {{if .Structured}}ret{{else}}{{if eq (len .Normalized.Outputs) 1}}ret0{{else}}&[]interface{}{
		{{range $i, $_ := .Normalized.Outputs}}ret{{$i}},
		{{end}}
	}{{end}}{{end}}
	err := _{{$contract.Type}}.Call(opts, out, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
	return {{if .Structured}}*ret,{{else}}{{range $i, $_ := .Normalized.Outputs}}*ret{{$i}},{{end}}{{end}} err
}
{{end}}

//////////////////////////////////////////////////////
//		Transactions
////////////////////////////////////////////////////

{{range .Transacts}}
// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
// - Solidity: {{formatmethod .Original $structs}}
func (_{{$contract.Type}} *{{$contract.Type}}) {{.Normalized.Name}}(opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) (*types.Transaction, error) {
	return _{{$contract.Type}}.Transact(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}
{{end}}

//////////////////////////////////////////////////////
//		Bin and ABI
////////////////////////////////////////////////////

// {{.Type}}Bin is used to deploy the generated contract
const {{.Type}}Bin = "0x{{.InputBin}}"

// {{.Type}}ABI is used to communicate with the compiled solidity code of the generated contract
const {{.Type}}ABI = "{{.InputABI}}"
{{end}}
`

// eventsTmpl is the template for generating go bindings to events and logs in ethereuem.
const eventsTmpl = `
{{$pkg := .Package}}
package {{$pkg}}

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

{{$structs := .Structs}}
//////////////////////////////////////////////////////
//		Events
////////////////////////////////////////////////////

{{range .Events}}
//////// {{.Normalized.Name}} ////////

// {{.Normalized.Name}}ID is the hex of the Topic Hash
const {{.Normalized.Name}}ID = "{{.Topic}}"

// {{.Normalized.Name}} represents a {{.Normalized.Name}} event raised by the {{.Type}} contract.
type {{.Normalized.Name}} struct { {{range .Normalized.Inputs}}
	{{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}; {{end}}
	Raw types.Log // Blockchain specific contextual infos
}

// Unpack{{.Normalized.Name}} is a log parse operation binding the contract event {{.Topic}}
// Solidity: {{formatevent .Original $structs}}
func (_{{.Type}} *{{.Type}}) Unpack{{.Normalized.Name}}(log types.Log) (*{{.Normalized.Name}}, error) {
	event := new({{.Normalized.Name}})
	if err := _{{.Type}}.UnpackLog(event, "{{.Original.Name}}", log); err != nil {
		return nil, err
	}
	return event, nil
}

{{end}}
`

// // for event in events make events + unpacker + parser
// // {{.Type}}{{.Normalized.Name}} represents a {{.Normalized.Name}} event raised by the {{.Type}} contract.
// type {{.Type}}{{.Normalized.Name}} struct { {{range .Normalized.Inputs}}
// 	{{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}; {{end}}
// 	Raw types.Log // Blockchain specific contextual infos
// }
// func (_{{.Type}}) *{{.Type}} {{}}

// /*
// Mux can be copied and pasted to save ya a quick minute when distinguishing between log data
// func Mux(c *{{.Type}}, log types.Log) error {
// 	switch log.Topics[0].Hex() { {{range .Events}}
// 	case {{$pkg}}.{{.Type}}{{.Normalized.Name}}ID:
// 		ulog, err := c.Unpack{{.Normalized.Name}}Log(log)
// 		if err != nil {
// 			return err
// 		}
// 		// insert additional code here
// 	{{end}}
// 	}
// }
// */
