#!/bin/bash
set -e
#build libyara
cd libyara
./bootstrap.sh
./configure --disable-shared --enable-static #--enable-magic
if [[ $1 = "-c" ]]; then
    make clean
fi
make -j 8
cd ..

#start building go-yara
export GOPATH=$(pwd):$GOPATH
export PKG_CONFIG_PATH=$(pwd)/libyara/libyara/:$(go env PKG_CONFIG_PATH)
export CGO_CFLAGS="-I$(pwd)/libyara/libyara/include/ $(go env CGO_CFLAGS)" 
export CGO_LDFLAGS="-L$(pwd)/libyara/libyara/.libs -lm $(go env CGO_LDFLAGS) -lcrypto"
go get -tags yara_static github.com/hillu/go-yara
echo "running go-yara tests"
go test github.com/hillu/go-yara
if [ $? -ne 0 ]; then
    echo "__________________________BUILD ERROR__________________________"
    go env
fi

if [ ! -d $GOPATH/src/github.com/elastic/go-elasticsearch ]; then
    git clone --branch master https://github.com/elastic/go-elasticsearch.git $GOPATH/src/github.com/elastic/go-elasticsearch
fi

go build -o scraper -race cmd/*.go
echo "Done"