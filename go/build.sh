#!/bin/bash

#build libyara
cd libyara
./bootstrap.sh
./configure --disable-shared --enable-static --without-crypto
make clean
make -j 8
cd ..

#start building go-yara
export GOPATH=$(pwd):$GOPATH
export export PKG_CONFIG_PATH=$(pwd)/libyara/libyara/:$PKG_CONFIG_PATH
export CGO_CFLAGS="-I$(pwd)/libyara/libyara/include/ $(go env CGO_CFLAGS)" 
export CGO_LDFLAGS="-L$(pwd)/libyara/libyara/.libs -lm $(go env CGO_LDFLAGS)"
go get github.com/hillu/go-yara
if [ $? -ne 0 ]; then
    echo "__________________________BUILD ERROR__________________________"
    go env
fi