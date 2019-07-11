#!/bin/bash

mkdir -p target
go build -o ./target/yh-radius
cp -r yh-radius startup.sh shutdown.sh attributes config target
cd target
chmod +x *.sh
