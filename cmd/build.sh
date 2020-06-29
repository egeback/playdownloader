#!/usr/bin/env bash
version=`cat VERSION`
time=$(date)
swag init --output ./internal/docs --parseInternal --generalInfo ./cmd/mediadownloader/main.go
go build -o mediadownloader -ldflags="-X 'github.com/egeback/mediadownloader/internal/version.BuildTime=$time' -X 'github.com/egeback/mediadownloader/internal/version.BuildVersion=$version'" ./cmd/mediadownloader/main.go
