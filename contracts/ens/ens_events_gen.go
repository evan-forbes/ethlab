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

//////////////////////////////////////////////////////
//		Events
////////////////////////////////////////////////////

//////// AddDomain ////////

// AddDomainID is the hex of the Topic Hash
const AddDomainID = "0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09"

// AddDomain represents a AddDomain event raised by the ENS contract.
type AddDomain struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackAddDomain is a log parse operation binding the contract event 0x32aa30f216d1137a37fb469a1675afd5e738e2fb578bc6fa4c1f6d5a3275bf09
// Solidity: event AddDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_ENS *ENS) UnpackAddDomain(log types.Log) (*AddDomain, error) {
	event := new(AddDomain)
	if err := _ENS.UnpackLog(event, "AddDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// ChangeDomain ////////

// ChangeDomainID is the hex of the Topic Hash
const ChangeDomainID = "0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac"

// ChangeDomain represents a ChangeDomain event raised by the ENS contract.
type ChangeDomain struct {
	Name    [32]byte
	Pointer common.Address
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// UnpackChangeDomain is a log parse operation binding the contract event 0x69e450e01952bb1748a7ad9b0c8f92f2cdb0d7d141f070e1892a920f41da50ac
// Solidity: event ChangeDomain(bytes32 indexed name, address indexed pointer, address indexed owner)
func (_ENS *ENS) UnpackChangeDomain(log types.Log) (*ChangeDomain, error) {
	event := new(ChangeDomain)
	if err := _ENS.UnpackLog(event, "ChangeDomain", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// LogTest ////////

// LogTestID is the hex of the Topic Hash
const LogTestID = "0xd006e942fead40af3b8cfee5f041b69d1b14565ce6e18728850a676969bd896e"

// LogTest represents a LogTest event raised by the ENS contract.
type LogTest struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// UnpackLogTest is a log parse operation binding the contract event 0xd006e942fead40af3b8cfee5f041b69d1b14565ce6e18728850a676969bd896e
// Solidity: event LogTest(address indexed owner)
func (_ENS *ENS) UnpackLogTest(log types.Log) (*LogTest, error) {
	event := new(LogTest)
	if err := _ENS.UnpackLog(event, "LogTest", log); err != nil {
		return nil, err
	}
	return event, nil
}
