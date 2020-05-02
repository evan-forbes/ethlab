package cmd

// func ABIgen(c *cli.Context) error {

// 	return nil
// }

// func Boot(c *cli.Context) error {
// 	// TODO: read/load config
// 	// look for config
// 	config, err := thereum.LoadConfig(c.String("config"))
// 	if err != nil {
// 		return errors.Wrapf(err, "Could not load config from path: \"%s\"", c.String("config"))
// 	}
// 	// TODO: save the config? at least save the initial allocations?
// 	mngr := NewManager(context.Background(), nil)
// 	go mngr.Listen()

// 	eth, err := thereum.New(config, nil)
// 	if err != nil {

// 	}
// 	// start producing blocks on the fresh chain
// 	mngr.WG.Add(1)
// 	go eth.Run(mngr.Ctx, mngr.WG)

// 	// pass the backend and configs to create a new server
// 	// TODO: make address a flag or read from a configs
// 	srvr := server.NewServer("127.0.0.1:8000", eth)
// 	go func() {
// 		log.Println((srvr.ListenAndServe()))
// 	}()
// 	<-mngr.Done()
// 	return nil
// }
