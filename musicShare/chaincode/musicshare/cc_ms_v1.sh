#!/bin/bash

# if [ $# -ne 2 ]; then
# 	echo "Arguments are missing. ex) ./cc_ms.sh instantiate 1.0.0"
# 	exit 1
# fi

instruction=$1
version=$2

set -ev

#chaincode install
docker exec cli peer chaincode install -n musicshare -v $version -p github.com/musicshare
#chaincode instatiate
docker exec cli peer chaincode $instruction -n musicshare -v $version -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 5
#chaincode invoke music1
docker exec cli peer chaincode invoke -n musicshare -C mychannel -c '{"Args":["register","Into the Night","YOASOBI","4:23","100"]}'
sleep 5
#chaincode query music1
docker exec cli peer chaincode query -n musicshare -C mychannel -c '{"Args":["readRating","user2"]}'

#chaincode invoke add rating
docker exec cli peer chaincode invoke -n musicshare -C mychannel -c '{"Args":["addRating","user2","p1","5.0"]}'
sleep 5

#chaincode query user1
docker exec cli peer chaincode query -n musicshare -C mychannel -c '{"Args":["readRating","user2"]}'

#chaincode query user1
docker exec cli peer chaincode query -n musicshare -C mychannel -c '{"Args":["getHistory","user2"]}'

echo '-------------------------------------END-------------------------------------'