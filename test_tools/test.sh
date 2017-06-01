#!/bin/bash
gocmd="$GOROOT/bin/go test"

for i in `ls cases`; do
	(cd cases/$i && $gocmd)
done
