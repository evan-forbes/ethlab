package ens

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

// Ens is a wrapper around BoundContract, enforcing type checking and including
// QoL helper methods
type Ens struct {
	bind.BoundContract
}

// NewEns creates a new instance of Ens, bound to a specific deployed contract.
func NewEns(address common.Address, backend bind.ContractBackend) (*Ens, error) {
	a, err := abi.JSON(strings.NewReader(EnsABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, a, backend, backend, backend)
	return &Ens{*contract}, nil
}

//////////////////////////////////////////////////////
//		Deployment
////////////////////////////////////////////////////

// DeployEns deploys a new Ethereum contract, binding an instance of Ens to it.
func DeployEns(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ens, error) {
	parsed, err := abi.JSON(strings.NewReader(EnsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EnsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ens{*contract}, nil
}

//////////////////////////////////////////////////////
//		Structs
////////////////////////////////////////////////////

//////////////////////////////////////////////////////
//		Data Calls
////////////////////////////////////////////////////

// Domains is a free data retrieval call binding the contract method 0xc722f177.
// - Solidity: function domains(bytes32 ) constant returns(address pointTo, address owner)
func (_Ens *Ens) Domains(opts *bind.CallOpts, arg0 [32]byte) (struct {
	PointTo common.Address
	Owner   common.Address
}, error) {
	ret := new(struct {
		PointTo common.Address
		Owner   common.Address
	})
	out := ret
	err := _Ens.Call(opts, out, "domains", arg0)
	return *ret, err
}

//////////////////////////////////////////////////////
//		Transactions
////////////////////////////////////////////////////

// Add is a paid mutator transaction binding the contract method 0x61641bdc.
// - Solidity: function add(bytes32 name, address addr) returns()
func (_Ens *Ens) Add(opts *bind.TransactOpts, name [32]byte, addr common.Address) (*types.Transaction, error) {
	return _Ens.Transact(opts, "add", name, addr)
}

// Change is a paid mutator transaction binding the contract method 0x33395e8f.
// - Solidity: function change(bytes32 name, address addr) returns()
func (_Ens *Ens) Change(opts *bind.TransactOpts, name [32]byte, addr common.Address) (*types.Transaction, error) {
	return _Ens.Transact(opts, "change", name, addr)
}

//////////////////////////////////////////////////////
//		Events
////////////////////////////////////////////////////

//////// AddDomain ////////

// AddDomainID is the hex of the Topic Hash
const AddDomainID = "0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09"

// AddDomainLog represents a AddDomain event raised by the Ens contract.
type AddDomainLog struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackAddDomainLog is a log parse operation binding the contract event 0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09
// Solidity: event AddDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_Ens *Ens) UnpackAddDomainLog(log types.Log) (*AddDomainLog, error) {
	event := new(AddDomainLog)
	if err := _Ens.UnpackLog(event, "AddDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// ChangeDomain ////////

// ChangeDomainID is the hex of the Topic Hash
const ChangeDomainID = "0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac"

// ChangeDomainLog represents a ChangeDomain event raised by the Ens contract.
type ChangeDomainLog struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackChangeDomainLog is a log parse operation binding the contract event 0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac
// Solidity: event ChangeDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_Ens *Ens) UnpackChangeDomainLog(log types.Log) (*ChangeDomainLog, error) {
	event := new(ChangeDomainLog)
	if err := _Ens.UnpackLog(event, "ChangeDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

/*
Mux can be copied and pasted to save ya a quick minute when distinguishing between log data
func Mux(c *Ens, log types.Log) error {
	switch log.Topics[0].Hex() {
	case ens.EnsAddDomainID:
		ulog, err := c.UnpackAddDomainLog(log)
		if err != nil {
			return err
		}
		// insert additional code here

	case ens.EnsChangeDomainID:
		ulog, err := c.UnpackChangeDomainLog(log)
		if err != nil {
			return err
		}
		// insert additional code here

	}
}
*/

//////////////////////////////////////////////////////
//		Bin and ABI
////////////////////////////////////////////////////

// EnsBin is used to deploy the generated contract
const EnsBin = "0x608060405234801561001057600080fd5b506105c6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806333395e8f1461004657806361641bdc14610094578063c722f177146100e2575b600080fd5b6100926004803603604081101561005c57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610183565b005b6100e0600480360360408110156100aa57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506102ef565b005b61010e600480360360208110156100f857600080fd5b8101908080359060200190929190505050610502565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390f35b60008083815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461023c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180610567602b913960400191505060405180910390fd5b8060008084815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac60405160405180910390a45050565b600073ffffffffffffffffffffffffffffffffffffffff1660008084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146103c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f646f6d61696e20616c726561647920657869737473000000000000000000000081525060200191505060405180910390fd5b8060008084815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503360008084815260200190815260200160002060010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060008083815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf0960405160405180910390a45050565b60006020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508256fe7573657220646f6573206e6f7420686176652072696768747320746f206465736972656420646f6d61696ea265627a7a72315820166e98dc1a2fc889b6e334165cdf69fef0eaea336661e6b6237bf46436b79a1664736f6c634300050f0032"

// EnsABI is used to communicate with the compiled solidity code of the generated contract
const EnsABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"AddDomain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ChangeDomain\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"add\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"change\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"domains\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pointTo\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
