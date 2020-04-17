package thereum

// I don't like this

// // API maps *Thereum methods to go-ethereum's bind.ContractBackend
// type API struct {
// 	*Thereum
// }

// func NewBackendAPI(config *Config, root *opts.Trans) (out API{}, err error) {
// 	if config == nil {
// 		config = &defaultConfig()
// 	}
// 	New(*config, )
// 	return API{}
// }

// // PendingCodeAt returns the code associated with an account in the pending state.
// func (t *Thereum) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
// 	t.mu.Lock()
// 	defer t.mu.Unlock()

// 	return t.latestState.GetCode(contract), nil
// }

// // PendingCallContract executes a contract call on the pending state.
// func (t *Thereum) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
// 	t.mu.Lock()
// 	defer t.mu.Unlock()
// 	defer t.latestState.RevertToSnapshot(t.latestState.Snapshot())

// 	rval, _, _, err := t.callContract(ctx, call, t.latestBlock, t.latestState)
// 	return rval, err
// }
