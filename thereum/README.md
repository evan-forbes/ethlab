## Usage 

Thereum describes a bare minimum blockchain. A ton of code was forked from go-ethereum in order to be able to plug into existing interfaces

## Notes

To be fully honest, I am personally not satifsied with the quality of this (my own) code. Most is molded to fit into go-ethereum and the ethereum json-rpc, which was not designed to be interoperable or this flexible. Does it work with a high enough degree of confidence? yes. Is it easy to debug. NO. Unfortunately, if we want to use the tools that go-etherum uses to create an EVM based blockchain, this is what we have to do...