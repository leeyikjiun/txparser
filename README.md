# TX Parser

Tx Parser is a simple parser for the ethereum blockchain.

It allows anyone to listen to transactions sent or received to an address.

See main.go for example usage.

## Notes
- Not sure if `GetTransactions` only retrieves transactions only for subscribed addresses or all addresses
- Does not deal with block re-orgs
- A better way to test this is to run a local hardhat node so that we can subscribe to a wallet and spam transactions
- Skips some basic validation on address such as 40 characters
- Configs are hard-coded