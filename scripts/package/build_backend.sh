#!/bin/bash

NAME="crypto-follower"
VERSION="1.0.0"

BUILD=`date +%FT%T%z`
echo ${BUILD}
echo "清理"
rm -rf ../../bin/${NAME}_backend
mkdir -p ../../bin

cd ../../backend
echo "当前目录是：" `pwd`

echo "编译mac  ====>"
go build -ldflags "-X main.profile=prod" -o ../bin/${NAME}_backend

echo "编译mac完成"


echo "准备复制到window  ====>"
cd ..

echo "当前目录是：" `pwd`


rm -rf window/bin
cp -rv bin window
echo "复制到window完成"