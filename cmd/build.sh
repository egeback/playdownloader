#!/usr/bin/env bash
#version=0.1.3
version=`cat VERSION`
time=$(date)
swag init --output ./internal/docs --parseInternal --exclude ./internal/models --generalInfo ./cmd/mediadownloader/main.go
go build -o mediadownloader -ldflags="-X 'github.com/egeback/mediadownloader/internal/version.BuildTime=$time' -X 'github.com/egeback/mediadownloader/internal/version.BuildVersion=$version'" ./cmd/mediadownloader/main.go
