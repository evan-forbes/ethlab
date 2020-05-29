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

// ENS is a wrapper around BoundContract, enforcing type checking and including
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

//////////////////////////////////////////////////////
//		Events
////////////////////////////////////////////////////

//////// AddDomain ////////

// AddDomainID is the hex of the Topic Hash
const AddDomainID = "0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09"

// AddDomainLog represents a AddDomain event raised by the ENS contract.
type AddDomainLog struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackAddDomainLog is a log parse operation binding the contract event 0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09
// Solidity: event AddDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_ENS *ENS) UnpackAddDomainLog(log types.Log) (*AddDomainLog, error) {
	event := new(AddDomainLog)
	if err := _ENS.UnpackLog(event, "AddDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// ChangeDomain ////////

// ChangeDomainID is the hex of the Topic Hash
const ChangeDomainID = "0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac"

// ChangeDomainLog represents a ChangeDomain event raised by the ENS contract.
type ChangeDomainLog struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackChangeDomainLog is a log parse operation binding the contract event 0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac
// Solidity: event ChangeDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_ENS *ENS) UnpackChangeDomainLog(log types.Log) (*ChangeDomainLog, error) {
	event := new(ChangeDomainLog)
	if err := _ENS.UnpackLog(event, "ChangeDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

/*
Mux can be copied and pasted to save ya a quick minute when distinguishing between log data
func Mux(c *ENS, log types.Log) error {
	switch log.Topics[0].Hex() {
	case ens.ENSAddDomainID:
		ulog, err := c.UnpackAddDomainLog(log)
		if err != nil {
			return err
		}
		// insert additional code here

	case ens.ENSChangeDomainID:
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

// ENSBin is used to deploy the generated contract
const ENSBin = "0x608060405234801561001057600080fd5b506102e8806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806333395e8f1461004657806361641bdc14610074578063c722f177146100a0575b600080fd5b6100726004803603604081101561005c57600080fd5b50803590602001356001600160a01b03166100e3565b005b6100726004803603604081101561008a57600080fd5b50803590602001356001600160a01b0316610194565b6100bd600480360360208110156100b657600080fd5b5035610262565b604080516001600160a01b03938416815291909216602082015281519081900390910190f35b6000828152602081905260409020600101546001600160a01b0316331461013b5760405162461bcd60e51b815260040180806020018281038252602b815260200180610289602b913960400191505060405180910390fd5b60008281526020819052604080822080546001600160a01b0319166001600160a01b0385169081179091559051339285917f69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac9190a45050565b6000828152602081905260409020546001600160a01b0316156101f6576040805162461bcd60e51b8152602060048201526015602482015274646f6d61696e20616c72656164792065786973747360581b604482015290519081900360640190fd5b60008281526020819052604080822080546001600160a01b038086166001600160a01b031992831681178455600190930180549092163317918290559251921692909185917f32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf0991a45050565b600060208190529081526040902080546001909101546001600160a01b0391821691168256fe7573657220646f6573206e6f7420686176652072696768747320746f206465736972656420646f6d61696ea265627a7a72315820e822a806a325f020587a08056cd40b11ba23f635e7dac4f1d3f6995c15db04a764736f6c634300050f0032"

// ENSABI is used to communicate with the compiled solidity code of the generated contract
const ENSABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"AddDomain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pointer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ChangeDomain\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"add\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"change\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"domains\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pointTo\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
