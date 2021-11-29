#!/bin/bash

if [ $# -ne 2 ]; then
	echo "Arguments are missing. ex) ./cc_ms.sh instantiate 1.0.0"
	exit 1
fi

instruction=$1
version=$2

set -ev

#chaincode install
docker exec cli peer chaincode install -n musicshare -v $version -p github.com/musicshare
#chaincode instatiate
docker exec cli peer chaincode $instruction -n musicshare -v $version -C msharenet -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member")'
sleep 5
# chaincode invoke init musics
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["initmusic"]}'
sleep 5
# chaincode invoke registerMusic
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["register","0006", "Into the Night","YOASOBI","4:23","100"]}'
sleep 5
# chaingcode query showMusicList
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["mquery"]}'
sleep 5
#chaincode invoke setContract
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["donate","0001","Mike","2", "30"]}'
sleep 5
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["donate","0001","Alice","3", "30"]}'
sleep 5
docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["make","0001","Bob", "60"]}'
sleep 5

# docker exec cli peer chaincode invoke -n musicshare -C msharenet -c '{"Args":["fill","0001","Mike","2"]}'
# sleep 5
#chaincode query the contract of 0001
docker exec cli peer chaincode query -n musicshare -C msharenet -c '{"Args":["cquery","0001"]}'
sleep 5
docker exec cli peer chaincode query -n musicshare -C msharenet -c '{"Args":["dquery","0001"]}'
sleep 5
#chaincode query shared profit of 0001
docker exec cli peer chaincode query -n musicshare -C msharenet -c '{"Args":["share","0001", "500000"]}'
sleep 5
#chaincode invoke stake of 0001
docker exec cli peer chaincode query -n musicshare -C msharenet -c '{"Args":["expire","0001"]}'

echo '-------------------------------------END-------------------------------------'