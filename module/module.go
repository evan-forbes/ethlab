package module

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

// Deployer wraps around functions intending to deploy a set of smart contracts.
type Deployer func(u *User) (common.Address, error)

// User represents the data needed to use the ethereuem block chain
type User struct {
	Client *ethclient.Client
	priv   *ecdsa.PrivateKey
	From   common.Address
	nonce  *big.Int
}

// Deploy runs multiple deploy functions using user u's private key
func (u *User) Deploy(deps ...Deployer) error {
	for _, dep := range deps {
		_, err := dep(u)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewUser inits a new user
func NewUser() (*User, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("COULD NOT GENERATE PRIVATE KEY")
		return nil, err
	}
	out := &User{
		priv:  priv,
		nonce: big.NewInt(1),
	}
	txopts := out.NewTxOpts()
	out.From = txopts.From

	return out, nil
}

// StarterKit generates a new account and requests eth to it.
func StarterKit(host string) (*User, error) {
	user, err := NewUser()
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(fmt.Sprintf("http://%s", host))
	if err != nil {
		return nil, err
	}
	user.Client = client
	err = RequestETH(host, user.From.Hex(), big.NewInt(1000000000000000000))
	return user, err
}

func (u *User) Balance() (*big.Int, error) {
	return u.Client.BalanceAt(context.Background(), u.From, nil)
}

// NewTxOpts issues a new transact opt with sane defaults and signs using User
// u's private key
func (u *User) NewTxOpts() *bind.TransactOpts {
	out := bind.NewKeyedTransactor(u.priv)
	out.GasLimit = 30000
	out.GasPrice = big.NewInt(1000000000)
	out.Nonce = u.nonce
	return out
}

// NewSend creates a signed transaction to send ETH
func (u *User) NewSend(to common.Address, amount, price *big.Int, lim uint64) (*types.Transaction, error) {
	txopts := u.NewTxOpts()
	txopts.GasLimit = lim
	txopts.GasPrice = price
	txopts.Value = amount
	utx := types.NewTransaction(txopts.Nonce.Uint64(), to, amount, lim, price, nil)
	return u.Sign(utx)
}

// Sign uses user data to sign a transaction
func (u *User) Sign(tx *types.Transaction) (*types.Transaction, error) {
	// increment and update the nonce
	u.IncrNonce()
	return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), u.priv)
}

func (u *User) IncrNonce() {
	nonce := new(big.Int).Add(u.nonce, big.NewInt(1))
	u.nonce = nonce
}

// RequestETH asks the server to dish out some eth to an address
func RequestETH(host, address string, amount *big.Int) error {

	type faucetPay struct {
		Address string   `json:"address"`
		Amount  *big.Int `json:"amount"`
	}

	type faucetResp struct {
		Message string `json:"message"`
	}

	data := faucetPay{
		Address: address,
		Amount:  amount,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/requestETH", host), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var out faucetResp
	json.Unmarshal(rawResp, &out)
	if out.Message != "success" {
		return errors.Errorf("failure to send eth: %s", out.Message)
	}

	return nil
}

// ENSAddress asks the host for the hex address of the ens contract
func ENSAddress(host string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/ens", host), strings.NewReader("hiya"))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(rawResp) == 32 {
		return "", errors.Errorf("no viable address for ens at host: %s: recieved %s", host, string(rawResp))
	}
	if string(rawResp[0]) != "0" && string(rawResp[1]) == "x" {
		return "", errors.Errorf("failure to get ens address at host %s: unexpected format")
	}
	return string(rawResp), nil
}
