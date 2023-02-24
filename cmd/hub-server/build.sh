#!/bin/bash

source ../../env.sh
export CGO_ENABLED=0

echo -e "\e[1;31m VERSION \e[0m" \\t $VERSION
echo -e "\e[1;31m BUILD_DATE \e[0m" \\t $BUILD_DATE
echo -e "\e[1;31m COMMIT \e[0m" \\t $COMMIT
echo -e "\e[1;31m GO_VERSION \e[0m" \\t $GO_VERSION

go build -ldflags "-w -X github.com/calvindc/Web3RpcHub/cmd/hub-server/mainimpl.Appversion=$VERSION
-X github.com/calvindc/Web3RpcHub/cmd/hub-server/mainimpl.Builddate=$BUILD_DATE
-X github.com/calvindc/Web3RpcHub/cmd/hub-server/mainimpl.Commit=$COMMIT
-X github.com/calvindc/Web3RpcHub/cmd/hub-server/mainimpl.Goversion=$GO_VERSION "

#cp metalifeserver $GOPATH/bin