#!/bin/bash

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root"
  exit 1
fi

# Install Go
apt-get update
apt-get install -y golang

# Set up Go environment variables
echo 'export GOPATH=$HOME/go' >> ~/.profile
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile
source ~/.profile

# Build the Go application
go build -o bk-setup-02 main.go

# Move the binary to a known location
mv bk-setup-02 /usr/local/bin/

echo "Installation complete. Run the application using 'bk-setup-02'."

