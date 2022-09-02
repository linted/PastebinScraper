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
mkdir -p $(pwd)/../go_dir
export GOPATH=$(pwd)/../go_dir:$(pwd):$GOPATH
export PKG_CONFIG_PATH=$(pwd)/libyara/libyara/:$(go env PKG_CONFIG_PATH)
export CGO_CFLAGS="-I$(pwd)/libyara/libyara/include/ $(go env CGO_CFLAGS)" 
export CGO_LDFLAGS="-L$(pwd)/libyara/libyara/.libs -lm $(go env CGO_LDFLAGS) -lcrypto"
if ! [ -f $(pwd)/go.mod ]; then
    go mod init example.com/m/v2
fi
go get -tags yara_static github.com/hillu/go-yara
go get github.com/bwmarrin/discordgo
echo "running go-yara tests"
go test github.com/hillu/go-yara
if [ $? -ne 0 ]; then
    echo "__________________________BUILD ERROR__________________________"
    go env
fi

go build -race cmd/*.go
echo "Done"