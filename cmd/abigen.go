package cmd

import (
	"context"
	"log"

	"github.com/evan-forbes/ethlab/server"
	"github.com/evan-forbes/ethlab/thereum"
	"github.com/urfave/cli/v2"
)

func ABIgen(c *cli.Context) error {

	return nil
}

func Boot(c *cli.Context) error {
	// TODO: read/load config

	mngr := NewManager(context.Background(), nil)
	go mngr.Listen()

	eth, err := thereum.New(thereum.DefaultConfig(), nil)
	if err != nil {

	}
	// start producing blocks on the fresh chain
	mngr.WG.Add(1)
	go eth.Run(mngr.Ctx, mngr.WG)

	// pass the backend and configs to create a new server
	// TODO: make address a flag or read from a configs
	srvr := server.NewServer("127.0.0.1:8000", eth)
	go func() {
		log.Println((srvr.ListenAndServe()))
	}()
	<-mngr.Done()
	return nil
}