# ID-based self-encryption application

## Overview
The prototype of a Hyperledger Fabric application for interacting with the ledger and IPFS to securely store ownership-preserved self-encrypted data as part of the Delft University of Technology Bachelor's Thesis.

The original paper: http://resolver.tudelft.nl/uuid:77406422-688c-4158-93f1-a83ab97810b4

This source code is based on Hyperledger Fabric Samples, https://github.com/hyperledger/fabric-samples   

ID-based self-encryption library: https://github.com/ilyagrishkov/ib-self-encryption-rust  
Smart contract: https://github.com/ilyagrishkov/ib-self-encryption-smart-contract

## Installation
* Golang v1.18.3
* Docker v20.10.14
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
