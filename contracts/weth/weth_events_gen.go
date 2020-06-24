package weth

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

//////// Withdrawal ////////

// WithdrawalID is the hex of the Topic Hash
const WithdrawalID = "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"

// Withdrawal represents a Withdrawal event raised by the WETH9 contract.
type Withdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// UnpackWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WETH9 *WETH9) UnpackWithdrawal(log types.Log) (*Withdrawal, error) {
	event := new(Withdrawal)
	if err := _WETH9.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// Approval ////////

// ApprovalID is the hex of the Topic Hash
const ApprovalID = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"

// Approval represents a Approval event raised by the WETH9 contract.
type Approval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// UnpackApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WETH9 *WETH9) UnpackApproval(log types.Log) (*Approval, error) {
	event := new(Approval)
	if err := _WETH9.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// Transfer ////////

// TransferID is the hex of the Topic Hash
const TransferID = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// Transfer represents a Transfer event raised by the WETH9 contract.
type Transfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// UnpackTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WETH9 *WETH9) UnpackTransfer(log types.Log) (*Transfer, error) {
	event := new(Transfer)
	if err := _WETH9.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

//////// Deposit ////////

// DepositID is the hex of the Topic Hash
const DepositID = "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c"

// Deposit represents a Deposit event raised by the WETH9 contract.
type Deposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// UnpackDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WETH9 *WETH9) UnpackDeposit(log types.Log) (*Deposit, error) {
	event := new(Deposit)
	if err := _WETH9.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	return event, nil
}
