#!/bin/sh

echo "Generating code from GraphQL schema..."

cd "$(dirname "$0")"

go get -d github.com/99designs/gqlgen@v0.12.0
go run github.com/99designs/gqlgen

sleep 5
