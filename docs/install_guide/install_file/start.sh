#!/bin/bash

echo "input node index"
read nodeIndex
node="node$nodeIndex"

startP2PPort=30300
startRPCPort=8545
mod=`expr $nodeIndex`

port=`expr $startP2PPort + $mod`
rpcport=`expr $startRPCPort + $mod`

coinbases=(0x1B7347c655d7B09aCB3cA13dfed1585235FA187E 0x296c57c333676C5cDB68bBB5d45Eca40b2D40b90 0x38Da2a6DEa3519ecadf934B6F80Dd48145811311 0x3baBf6AC5b776dBDB07420FcB0D79Ae28669Fde7 0x7f9f4af365ae4d1EcCe9fBbC605a7e8a00482343 0xC502652b4aD530d0E1e62884d48F26eD615204FA 0xaF1A5CB144f608B9411B9D75d4BA8eD93eEE7ad9)
miner=${coinbases[$nodeIndex]}
echo "$node and miner is $miner, rpc port $rpcport, p2p port $port"

nohup ./geth --mine --miner.threads 1 \
--miner.etherbase=$miner \
--identity=$node \
--maxpeers=100 \
--syncmode full \
--allow-insecure-unlock \
--datadir $node \
--txlookuplimit 0 \
--networkid 10898 \
--http.corsdomain "*" \
--http.api admin,eth,debug,miner,net,txpool,personal,web3 \
--http --http.addr 0.0.0.0 --http.port $rpcport --http.vhosts "*" \
--port $port \
--verbosity 5 \
--nodiscover >> $node/node.log 2>&1 &

sleep 3s;
ps -ef|grep geth|grep mine|grep -v grep;
