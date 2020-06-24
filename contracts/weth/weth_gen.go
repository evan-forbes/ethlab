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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// WETH9 is a wrapper around bind.BoundContract, enforcing type checking and including
// QoL helper methods
type WETH9 struct {
	bind.BoundContract
}

// NewWETH9 creates a new instance of WETH9, bound to a specific deployed contract.
func NewWETH9(address common.Address, backend bind.ContractBackend) (*WETH9, error) {
	a, err := abi.JSON(strings.NewReader(WETH9ABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, a, backend, backend, backend)
	return &WETH9{*contract}, nil
}

//////////////////////////////////////////////////////
//		Deployment
////////////////////////////////////////////////////

// DeployWETH9 deploys a new Ethereum contract, binding an instance of WETH9 to it.
func DeployWETH9(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WETH9, error) {
	parsed, err := abi.JSON(strings.NewReader(WETH9ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WETH9Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WETH9{*contract}, nil
}

//////////////////////////////////////////////////////
//		Structs
////////////////////////////////////////////////////

//////////////////////////////////////////////////////
//		Data Calls
////////////////////////////////////////////////////

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
// - Solidity: function allowance(address , address ) constant returns(uint256)
func (_WETH9 *WETH9) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "allowance", arg0, arg1)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
// - Solidity: function balanceOf(address ) constant returns(uint256)
func (_WETH9 *WETH9) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "balanceOf", arg0)
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
// - Solidity: function decimals() constant returns(uint8)
func (_WETH9 *WETH9) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "decimals")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
// - Solidity: function name() constant returns(string)
func (_WETH9 *WETH9) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "name")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
// - Solidity: function symbol() constant returns(string)
func (_WETH9 *WETH9) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "symbol")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
// - Solidity: function totalSupply() constant returns(uint256)
func (_WETH9 *WETH9) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WETH9.Call(opts, out, "totalSupply")
	return *ret0, err
}

//////////////////////////////////////////////////////
//		Transactions
////////////////////////////////////////////////////

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
// - Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WETH9 *WETH9) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Transact(opts, "approve", guy, wad)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
// - Solidity: function deposit() returns()
func (_WETH9 *WETH9) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.Transact(opts, "deposit")
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
// - Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WETH9 *WETH9) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Transact(opts, "transfer", dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
// - Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WETH9 *WETH9) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Transact(opts, "transferFrom", src, dst, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
// - Solidity: function withdraw(uint256 wad) returns()
func (_WETH9 *WETH9) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Transact(opts, "withdraw", wad)
}

//////////////////////////////////////////////////////
//		Bin and ABI
////////////////////////////////////////////////////

// WETH9Bin is used to deploy the generated contract
const WETH9Bin = "0x60c0604052600d60808190527f577261707065642045746865720000000000000000000000000000000000000060a090815261003e91600091906100a3565b506040805180820190915260048082527f57455448000000000000000000000000000000000000000000000000000000006020909201918252610083916001916100a3565b506002805460ff1916601217905534801561009d57600080fd5b5061013e565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100e457805160ff1916838001178555610111565b82800160010185558215610111579182015b828111156101115782518255916020019190600101906100f6565b5061011d929150610121565b5090565b61013b91905b8082111561011d5760008155600101610127565b90565b6106728061014d6000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde0381146100b8578063095ea7b31461014257806318160ddd1461017a57806323b872dd146101a15780632e1a7d4d146101cb578063313ce567146101e357806370a082311461020e57806395d89b411461022f578063a9059cbb14610244578063d0e30db0146100ae578063dd62ed3e14610268575b6100b661028f565b005b3480156100c457600080fd5b506100cd6102de565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101075781810151838201526020016100ef565b50505050905090810190601f1680156101345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561014e57600080fd5b50610166600160a060020a036004351660243561036c565b604080519115158252519081900360200190f35b34801561018657600080fd5b5061018f6103d2565b60408051918252519081900360200190f35b3480156101ad57600080fd5b50610166600160a060020a03600435811690602435166044356103d7565b3480156101d757600080fd5b506100b660043561050b565b3480156101ef57600080fd5b506101f86105a0565b6040805160ff9092168252519081900360200190f35b34801561021a57600080fd5b5061018f600160a060020a03600435166105a9565b34801561023b57600080fd5b506100cd6105bb565b34801561025057600080fd5b50610166600160a060020a0360043516602435610615565b34801561027457600080fd5b5061018f600160a060020a0360043581169060243516610629565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103645780601f1061033957610100808354040283529160200191610364565b820191906000526020600020905b81548152906001019060200180831161034757829003601f168201915b505050505081565b336000818152600460209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b303190565b600160a060020a0383166000908152600360205260408120548211156103fc57600080fd5b600160a060020a038416331480159061043a5750600160a060020a038416600090815260046020908152604080832033845290915290205460001914155b1561049a57600160a060020a038416600090815260046020908152604080832033845290915290205482111561046f57600080fd5b600160a060020a03841660009081526004602090815260408083203384529091529020805483900390555b600160a060020a03808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561052757600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f19350505050158015610566573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103645780601f1061033957610100808354040283529160200191610364565b60006106223384846103d7565b9392505050565b6004602090815260009283526040808420909152908252902054815600a165627a7a723058208000072bb2b644a18d2878ef70f67fa7e4a13a8a77c5d35e5d260666f89d98e80029"

// WETH9ABI is used to communicate with the compiled solidity code of the generated contract
const WETH9ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"
