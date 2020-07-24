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
	"github.com/evan-forbes/ethlab/module"
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

// ENS is a wrapper around bind.BoundContract, enforcing type checking and including
// QoL helper methods
type ENS struct {
	bind.BoundContract
}

// NewENS creates a new instance of ENS, bound to a specific deployed contract.
func NewENS(address common.Address, backend bind.ContractBackend) (*ENS, error) {
	a, err := abi.JSON(strings.NewReader(ENSABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, a, backend, backend, backend)
	return &ENS{*contract}, nil
}

//////////////////////////////////////////////////////
//		Deployment
////////////////////////////////////////////////////

// Deploy installs ENS to an ethereum node via the user provided
// by implementing module.Delpoyer
func Deploy(u *module.User) (addr common.Address, err error) {
	// ****************************************************
	////////  INSERT MODULE DEPLOYMENT CODE HERE   ////////
	// ***************************************************
	return addr, err
}

// DeployENS deploys a new Ethereum contract, binding an instance of ENS to it.
func DeployENS(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ENS, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ENSBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ENS{*contract}, nil
}

//////////////////////////////////////////////////////
//		Structs
////////////////////////////////////////////////////

//////////////////////////////////////////////////////
//		Data Calls
////////////////////////////////////////////////////

// Domains is a free data retrieval call binding the contract method 0xc722f177.
// - Solidity: function domains(bytes32 ) constant returns(address pointTo, address owner)
func (_ENS *ENS) Domains(opts *bind.CallOpts, arg0 [32]byte) (struct {
	PointTo common.Address
	Owner   common.Address
}, error) {
	ret := new(struct {
		PointTo common.Address
		Owner   common.Address
	})
	out := ret
	err := _ENS.Call(opts, out, "domains", arg0)
	return *ret, err
}

//////////////////////////////////////////////////////
//		Transactions
////////////////////////////////////////////////////

// Add is a paid mutator transaction binding the contract method 0x61641bdc.
// - Solidity: function add(bytes32 name, address addr) returns()
func (_ENS *ENS) Add(opts *bind.TransactOpts, name [32]byte, addr common.Address) (*types.Transaction, error) {
	return _ENS.Transact(opts, "add", name, addr)
}

// Change is a paid mutator transaction binding the contract method 0x33395e8f.
// - Solidity: function change(bytes32 name, address addr) returns()
func (_ENS *ENS) Change(opts *bind.TransactOpts, name [32]byte, addr common.Address) (*types.Transaction, error) {
	return _ENS.Transact(opts, "change", name, addr)
}

// LogTest is a paid mutator transaction binding the contract method 0x1361c394.
// - Solidity: function logTest() returns()
func (_ENS *ENS) LogTest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENS.Transact(opts, "logTest")
}

//////////////////////////////////////////////////////
//		Bin and ABI
////////////////////////////////////////////////////

// ENSBin is used to deploy the generated contract
const ENSBin = "0x608060405234801561001057600080fd5b50610328806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631361c3941461005157806333395e8f1461005b57806361641bdc14610087578063c722f177146100b3575b600080fd5b6100596100f6565b005b6100596004803603604081101561007157600080fd5b50803590602001356001600160a01b0316610123565b6100596004803603604081101561009d57600080fd5b50803590602001356001600160a01b03166101d4565b6100d0600480360360208110156100c957600080fd5b50356102a2565b604080516001600160a01b03938416815291909216602082015281519081900390910190f35b60405133907fd006e942fead40af3b8cfee5f041b69d1b14565ce6e18728850a676969bd896e90600090a2565b6000828152602081905260409020600101546001600160a01b0316331461017b5760405162461bcd60e51b815260040180806020018281038252602b8152602001806102c9602b913960400191505060405180910390fd5b60008281526020819052604080822080546001600160a01b0319166001600160a01b0385169081179091559051339285917f69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac9190a45050565b6000828152602081905260409020546001600160a01b031615610236576040805162461bcd60e51b8152602060048201526015602482015274646f6d61696e20616c72656164792065786973747360581b604482015290519081900360640190fd5b60008281526020819052604080822080546001600160a01b038086166001600160a01b031992831681178455600190930180549092163317918290559251921692909185917f32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf0991a45050565b600060208190529081526040902080546001909101546001600160a01b0391821691168256fe7573657220646f6573206e6f7420686176652072696768747320746f206465736972656420646f6d61696ea265627a7a72315820431e04355d1a8e2d572964ecf856fef24129023d7c88c1e3c950f28018fe424c64736f6c634300050f0032"

// ENSABI is used to communicate with the compiled solidity code of the generated contract
const ENSABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"AddDomain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ChangeDomain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"LogTest\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"add\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"change\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"domains\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pointTo\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logTest\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
