#!/bin/bash

echo "Installing IBSE application..."
mkdir ~/.ibse

echo "Creating config..."
touch ~/.ibse/config.yaml

echo "Copying encryptor library..."
cp rustgo/ib_self_encryption_rust.wasm ~/.ibse/ib_self_encryption_rust.wasm

echo "Installation done!"
