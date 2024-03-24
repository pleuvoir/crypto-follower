#!/bin/bash

NAME="crypto-follower"
VERSION="1.0.0"
BUILD=`date +%FT%T%z`
echo ${BUILD}
echo "-X backend/app.Version=${VERSION} -X backend/app.Build=${BUILD}"
echo "清理"
rm -rf ../bin/${NAME}_backend
mkdir ../bin

cd ../backend
echo "编译mac"
go build -ldflags "-X backend/app.Version=${VERSION} -X backend/app.Build=${BUILD}" -o ../bin/${NAME}_backend