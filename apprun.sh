#!/usr/bin/env sh

cd ./application/go

go mod tidy
if [ $? -ne 0 ]
then
	echo "tidy error"
	exit 1
fi

go build
if [ $? -ne 0 ]
then
	echo "build error"
	exit 1
fi

rm -r ./wallet ./keystore

./application $1 $2
