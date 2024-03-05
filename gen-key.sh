#!/bin/bash

# Set the desired key size (optional, defaults to 4096 bits)
KEY_SIZE=${1:-4096}

# Generate the key pair with a comment and passphrase
ssh-keygen -b $KEY_SIZE -t rsa -C "nguyenminhmannh2001@gmail.com" -f ./keys/id_rsa -m "pem" -N ""

# Check for errors during key generation
if [ $? -ne 0 ]; then
  echo "Error generating key pair!"
  exit 1
fi

echo "Your private key has been saved in: ./keys/id_rsa"
echo "Your public key has been saved in: ./keys/id_rsa.pub"