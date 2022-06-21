# ID-based self-encryption application

## Overview
The prototype of a Hyperledger Fabric application for interracting with the ledger and IPFS to securely store ownership-preserved self-encrpyted data as part of the Delft University of Technology Bachelor's Thesis.

This source code is based on Hyperledger Fabric Samples, https://github.com/hyperledger/fabric-samples   

Self-Encryption library: https://github.com/ilyagrishkov/ib-self-encryption-rust  
Smart contract: https://github.com/ilyagrishkov/ib-self-encryption-smart-contract

## Installation
* Golang v1.18.3
* Docker v19.03.8
* Hyperledger test-network v2.4.3

## Before run
* Deploy the Hyperledger Fabric test-network with certificate authorities (`-ca` flag) based on the instructions: https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html and create a default channel (`mychannel`)
* Deploy the smart contract (https://github.com/ilyagrishkov/ib-self-encryption-smart-contract) to the test-network


## Run

#### Clone the repository
```shell
git clone https://github.com/ilyagrishkov/ib-self-encryption.git
```

#### Navigate to the application directory
```shell
cd ib-self-encryption
```

#### Install the application by running `install.sh`
```shell
sh install.sh
```

#### Run help command to get information about all available commands
```shell
./ibse help
```
