#!/bin/bash


NAME="crypto-follower"
VERSION="1.0.0"

BUILD=`date +%FT%T%z`
echo ${BUILD}
echo "清理"
rm -rf ../../bin
mkdir -p ../../bin/configs

cd ../../backend
echo "当前目录是：" `pwd`

echo "复制配置文件到bin目录  ====>"

cp -rv ./configs ../bin


echo "编译mac  ====>"
go build -ldflags "-X main.profile=prod" -o ../bin/${NAME}_backend

echo "编译mac完成"


echo "准备复制到window  ====>"
cd ..

echo "当前目录是：" `pwd`


rm -rf window/bin
cp -rv bin window
echo "复制到window完成"