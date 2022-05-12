#!/usr/bin/env	sh

if [ $# -ne 2 ]
then
	echo "Usage: script.sh ChannelName ContractName"
	exit 1
fi

CHANNEL=$1
CONTRACT=$2

echo "Clear test network"
cd test-network


export FABRIC_CFG_PATH=${PWD}/../config/ 
./setOrgEnv.sh

./network.sh down
if [ $? -ne 0 ]
then
	echo "Can't clear previous test network"
fi

echo "Set up test network"
./network.sh up
if [ $? -ne ]
then
	echo "Can't create test network"
fi

echo "Create test channel"
./network.sh up createChannel -c $CHANNEL
if [ $? -ne 0 ]
then
	echo "Can't create test channel"
	exit 1
fi

echo "Building go-lang chaincode"
cd ../chaincode/go
go mod tidy && go mod vendor

if [ $? -ne 0 ]
then
	echo "Chaincode is not correct"
	exit 1
fi

echo "Load chaincode"
cd ../../test-network
./network.sh deployCC -ccn $CONTRACT -ccp ../chaincode/go -ccl go -c $CHANNEL
if [ $? -ne 0 ]
then
	echo "Can't load chaincode to the network"
	exit 1
fi


