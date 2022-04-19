#!/bin/bash

echo "input node index"
read nodeIndex
node="node$nodeIndex"

kill -s SIGINT $(ps aux|grep geth|grep mine|grep $node|grep -v grep|awk '{print $2}');

sleep 3s;
ps -ef|grep geth|grep mine|grep -v grep;