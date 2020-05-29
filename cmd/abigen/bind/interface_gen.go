package bind

type tmplInterface struct {
	Package   string
	Contracts map[string]*tmplContract
}

const interfaceTemplate = `
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
{{range $contract := .Contracts}}
// {{.Type}} is a generated interface that describes
type {{.Type}} interface {
	{{range .Calls}}{{.Normalized.Name}}(opts *bind.CallOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} error)
	{{end}}
	{{range .Transacts}}{{.Normalized.Name}}(opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) (*types.Transaction, error)
	{{end}}
}
{{end}}

`
