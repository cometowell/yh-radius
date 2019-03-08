#!/bin/bash

mkdir -p target
go build -o ./target/go-rad
cp -r go-rad startup.sh shutdown.sh attributes config target
cd target
chmod +x *.sh
