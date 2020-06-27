package module

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evan-forbes/ethlab/server"
)

// Module descibes a deployment script that can be
// compiled
type Module interface {
	Deploy() error
}

type User struct {
	Client *ethclient.Client
	priv   *ecdsa.PrivateKey
	from   common.Address
}

// NewUser inits a new user
func NewUser() (*User, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("COULD NOT GENERATE PRIVATE KEY")
		return nil, err
	}
	out := &User{
		priv: priv,
	}
	txopts := out.NewTxOpts(300000000, big.NewInt(10000000))
	out.from = txopts.From
	return out, nil
}

// StarterKit generates a new account and requests
// eth to it.
func StarterKit(host string) (*User, error) {
	user, err := NewUser()
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(host)
	if err != nil {
		return nil, err
	}
	user.Client = client
	server.RequestETH(host, user.from.Hex())
	return user, nil
}

// NewTxOpts issues a new transact opt with sane defaults for
// user u
func (u *User) NewTxOpts(gasLim uint64, gasPrice *big.Int) *bind.TransactOpts {
	out := bind.NewKeyedTransactor(u.priv)
	out.GasLimit = gasLim
	out.GasPrice = gasPrice
	if gasLim == 0 {
		out.GasLimit = 3000000
	}
	if gasPrice == nil || gasPrice.Cmp(big.NewInt(0)) == 0 {
		out.GasPrice = big.NewInt(10000000)
	}
	return out
}
