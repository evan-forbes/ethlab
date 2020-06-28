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

type Deployer func(u *User) (common.Address, error)

type User struct {
	Client *ethclient.Client
	priv   *ecdsa.PrivateKey
	from   common.Address
}

func (u *User) Deploy(deps ...Deployer) error {
	for _, dep := range deps {
		_, err := dep(u)
		if err != nil {
			return err
		}
	}
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
	txopts := out.NewTxOpts()
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
func (u *User) NewTxOpts() *bind.TransactOpts {
	out := bind.NewKeyedTransactor(u.priv)
	out.GasLimit = 3000000
	out.GasPrice = big.NewInt(10000000)
	return out
}
