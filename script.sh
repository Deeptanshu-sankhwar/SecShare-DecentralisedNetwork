#!/bin/bash

CURR=$PWD

for f in $PWD/src/*
do
	cd $f
	# echo $HOME/concurrency-decentralized-network/docs/${PWD##*/}/index.html
	godoc -html -goroot=$CURR cmd/${PWD##*/} > $CURR/docs/${PWD##*/}/index.html

done