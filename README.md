# go-aptos

Golang SDK for [Aptos Blockchain](https://aptos.dev)

See package page for documentation. All codes resides inside [aptos](./aptos) directory.

[![Go Reference](https://pkg.go.dev/badge/github.com/fardream/go-aptos.svg)](https://pkg.go.dev/github.com/fardream/go-aptos)

## Import

Add to your project

```sh
go get github.com/fardream/go-aptos@latest
```

## CLI

There is also a cli provide, install

```sh
go install github/fardream/go-aptos/aptos/cmd/aptos-aux@latest
```

## Quick Tour

The code stems from the [aux dex](https://aux.exchange), where we look for a lightweight golang library to interact with the chain. BCS support is provided in [go-bcs](https://pkg.go.dev/github.com/fardream/go-bcs/bcs)

Aptos right now only supports rest api, without a websocket or streaming api.

Aptos has multiple networks:

- Mainnet, chain id 1
- Testnet, chain id 2
- Devnet, chain id reset to previous one plus 1 every time chain resets. Devnet resets approximately every 2 weeks.
- Localnet, chain id 4.
- Customnet

### Account, Keys, and Address

Aptos uses ed25519 for signing, and support Bitcoin style mnemonic Words. Right now only [`LocalAccount`](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#LocalAccount) with single private key is supported. It can be generated from either [ed25519 private key](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#NewLocalAccountFromPrivateKey), or from [mnemonic words](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#NewLocalAccountFromMnemonic).

[Resource account address calculation](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#CalculateResourceAddress) is also supported.

### Rest Api

All interactions are handled through `aptos.Client`. A network must be specified - the default endpoint for that network will be used if the URL is left to be empty string.

For the supported methods, see the full [api](https://fullnode.mainnet.aptoslabs.com/v1/spec#/) documentation of aptos fullnode.

```go
client, err := new aptos.NewClient(aptos.Devnet, "..<url>..")
```

All methods that interact with the rest endpoint will take a `context.Context` parameter, and do not have special logics for cancelling or timing out.

```go
// create a new context with 5 minutes timeout
ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 5)
defer cancel()
// now the request will wait for 5 minutes before timing out or when cancel is called.
resp, err := client.GetAccountResources(ctx, &aptos.GetAccountResourcesRequest{Address: aptos.MustParseAddress("0x1")})
```

A detailed example can be found [here](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#example-Client.GetAccountResources)

### Submit/Simulate Transaction

The transactions can be encoded locally with [`EncodeTransaction`](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#EncodeTransaction) function. **Note**: techanically speaking, the transaction hash can also be backed out locally, however, right now the method doesn't match the one returned from the chain.

Right now only entry function payload (more on what is ["entry function"](https://aptos.dev/guides/system-integrators-guide/#types-of-transactions)) and single ed25519 signature is supported.

[`clientFillTransactionData`](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#Client.FillTransactionData) can fill in some of the missing data for a transaction (gas price, sequence number, chain id) by requesting the data from the chain.

```go
// create a new random account
localAccount := aptos.NewLocalAccountWithRandomKey()
// create some transaction
tx := &aptos.Transaction{<...>}
// filling missing data
client.FillTransactionData(ctx, tx, false)
// sign
sig, _ := localAccount.Sign(tx)
// submit
client.SubmitTransaction(ctx, &aptos.SubmitTransactionRequest{Transaction: tx, Signature: *sig})
```

A detailed example can be found here [here](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#example-Client.SubmitTransaction)

Signing, submitting transaction, and waiting for transaction completion can be done in one method call [SignSubmitTransactionWait](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#Client.SignSubmitTransactionWait).

### Aux Client

Aux exchange contract codes are located [here](https://github.com/aux-exchange/aux-exchange/tree/main/aptos/contract/aux). Use [`aptos.AuxClientConfig`](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#AuxClientConfig) to build transactions. All transaction builder methods are named ModuleName_MethodName (for example, clob_market::place_order is named [ClobMarket_PlaceOrder](https://pkg.go.dev/github.com/fardream/go-aptos@main/aptos#AuxClientConfig.ClobMarket_PlaceOrder))

```go
// get configuration
auxConfig, err := aptos.GetAuxClientConfig(aptos.Devnet)

// build a transaction
tx := info.ClobMarket_PlaceOrder(...)

client, err := client.SignSubmitTransactionAndWait(...)
```
