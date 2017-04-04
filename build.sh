#!/bin/sh

cd fetcher
go fmt
go build
cd ..

cd parser
go fmt
go build
cd ..

cd apis
go fmt
go build
cd ..

cd database
go fmt
go build
cd ..

cd apiclient
go fmt
go install
cd ..
