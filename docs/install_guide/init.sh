#!/bin/bash

echo "input node index"
read nodeIndex
node="node$nodeIndex"

if [ ! -f $node/genesis.json ]
then
	mkdir -p $node/geth/
	cp setup/genesis.json $node/genesis.json
fi

if [ ! -f $node/static-nodes.json ]
then
	cp setup/static-nodes.json $node/static-nodes.json
fi

if [ ! -f $node/geth/nodekey ]
then
	cp setup/$node/nodekey $node/geth/
fi

if [ ! -f $node/bsn.json ]
then
	cp setup/bsn.json $node/bsn.json
fi

./geth init $node/genesis.json --datadir $node
