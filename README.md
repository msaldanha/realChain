# RealChain

## Introduction

This is my pet project where I can exercise my learning of Go programming language and blockchain related technologies. 

RealChain is a transaction chain, similar to a blockchain but made of individual transactions. 
Inspired by Nano (https://nano.org), 
it this current version, it is simple and uses a simplistic consensus algorithm.

## Components

RealChain consists in one executable for the sake of deployment/usage simplicity. 
This executable has the following functionalities:

- **Wallet** – keeps the address keys, creates transactions and submit them to the network.
- **Node** – manages the transaction chain.

#### Wallet

The RealChain network must receive two transactions: one SEND and one RECEIVE. 
So, it is the wallet responsibility to create/get both transaction and send them to the network.

One scenario would be:

1. Wallet A sends a payment request for address X to Wallet B.
2. Wallet B checks account balance and creates and sends the SEND transaction to Wallet A.
3. Wallet A checks the received SEND transaction, creates a RECEIVE transaction linked to the SEND transaction and send both transactions to the network.
4. Wallet A waits for the network confirmation and informs the results to the user.

The flow described above will be used as basis to the inter-wallet communication protocol which, right now, is yet to be designed and implemented.

As we do not have this inter-wallet communication protocol yet, this version of the wallet only works for addresses which it has the keys.

So, for now, the wallet can only do transfers between address that it manages (has the private keys).

#### Node

As stated previously, the RealChain network needs two transactions to validate and confirm a transfer.

The node then does the following:

1. The node receives the two transaction and validates them.
2. The node starts a voting process, and asks other nodes to validate the transfer.
   * Each node validates the two transactions and votes to accept or not them.
3. The node collects all votes and, if all nodes had accepted the transfer, the node writes the transactions to the ledger.
4. The node sends the voting result to the other nodes.
   * Each node validates the voting and if all nodes had accepted the transfer, the node writes the transactions to the ledger.
5. The node returns success or failure (in case the transaction validation or the voting failed) to the client.

Note that by using this naive consensus algorithm, the network is not scalable as each node needs to know and contact 
all nodes in the network.

As RealChain is inspired by Nano (previously called RaiBlocks), a look at Nano's white paper will give more details about its protocol: https://www.raiblocks.net/media/RaiBlocks_Whitepaper__English.pdf

One difference between RealChain and Nano networks is that RealChain does not use representatives and must receive the SEND and RECEIVE transactions.

## Requirements

- Go 1.11+

## Build

```
$ go build
```

## Setup

### Setup overview

To run the network it is necessary at least 2 nodes. 

#### Nodes

The steps to setup the network are:

1. Setup the first node and the ledger (genesis transaction and address).
2. Copy the keys db of the genesis address to be used by the wallet.
3. Setup the second node and copy the ledger db to its data dir.
4. Configure each node as a peer of the other.
5. Run both nodes.

#### Wallet

The steps to setup the wallet are:

1. Copy the keys db containing the genesis address to wallets data dir.
2. Start using the wallet.

### Setup details (running nodes and wallet on the same machine)

1. Build the application executable
    ```
    go build
    ```
2. Create a working dir
    ```
    mkdir ~/realChain
    ```
3. Create a folder in the working dir for the first node and copy the executable
    ```
    mkdir ~/realChain/node1
    mkdir ~/realChain/node1/data
    cp realChain ~/realChain/node1
    cp config.yaml ~/realChain/node1
    ```
4. Create a folder in the working dir for the second node and copy the executable
    ```
    mkdir ~/realChain/node2
    mkdir ~/realChain/node2/data
    cp realChain ~/realChain/node2
    cp config.yaml ~/realChain/node2
    ```
5. Create a folder in the working dir for the wallet and copy the executable
    ```
    mkdir ~/realChain/wallet
    mkdir ~/realChain/wallet/data
    cp realChain ~/realChain/wallet
    cp config.yaml ~/realChain/wallet
    ```
6. Setup the first node. 
   * Set the configuration property `node.server` to `'127.0.0.1:4000'`
   * Include the ip/port of the second node (`'127.0.0.1:4002'`) in the configuration property `node.peers`.
   * Run node's init command. 
    ```
    cd ~/realChain/node1
    vi config.yaml  #edit properties here, save end exit
    ./realChain node init
    ```
7. Use the first node to create the ledger with a starting amount.
   * Copy the generated keys to the wallet's data dir.
   * Copy the ledger to the second node's data dir.
    ```
    cd ~/realChain/node1
    ./realChain ledger init 1000000000 waddresses.db
    cp data/waddresses.db ../wallet/data
    cp data/chain.db ../node2/data
    ```
8. Setup the second node. 
   * Set the configuration property `node.server` to `'127.0.0.1:4002'`
   * Include the ip/port of the first node (`'127.0.0.1:4000'`) in the configuration property `node.peers`.
   * Run node's init command. 
    ```
    cd ~/realChain/node2
    vi config.yaml  #edit properties here, save end exit
    ./realChain node init
    ```
9. Run first node on a separated terminal.
    ```
    cd ~/realChain/node1
    ./realChain node serve
    ```
10. Run second node on a separated terminal.
    ```
    cd ~/realChain/node2
    ./realChain node serve
    ```
11. Create a new address on the wallet and take note of the address.
    ```
    cd ~/realChain/wallet
    ./realChain wallet create
    ```
    
### Setup test

After the network is up and running and the wallet is setup, it is possible to test
by transferring values from one address to another. 

1. For example, transfer 1000 and 500 from the genesis address to the new created address (from wallet setup above).
    ```
    cd ~/realChain/wallet
    ./realChain wallet send <genesis address> <new created address> 1000
    ./realChain wallet send <genesis address> <new created address> 500
    ```
2. Get the transaction statement for the new created address.
    ```
    ./realChain wallet statement <new created address>
    ```
    
## Usage

The user interacts with the wallet command. Use `help` command option to get a list 
of available wallet commands:

```
./realChain wallet help
Wallet related commands

Usage:
  realChain wallet [command]

Available Commands:
  create      Creates an address
  list        Lists all managed addresses
  send        Sends [amount] from [FROM address] to [TO address]
  statement   Lists all transactions for [address]

Flags:
  -h, --help   help for wallet
```

To list all addresses managed by the wallet:
```
./realChain wallet list
```
To create a new address on the wallet:
```
./realChain wallet create
```
To send values from one address to another (in this version, both addresses must be managed by the wallet):
```
./realChain wallet send <from address> <to address> <value>
```
To get the transaction statement (all transactions of an address):
```
./realChain wallet statement <address>
```

## Next steps

I hope to add (as time permits):

- **Auto Seeding** – ability for a new node get the latest chain when joining the network.
- **Consesus algorithm** – replace the current naive implementation by a real one.
- **Secure wallet's keys db** – use cryptography to secure address keys stored by the wallet.
- **Add inter-wallet communication protocol** – add ability to wallets to exchange transactions. 

