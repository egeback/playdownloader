#!/usr/bin/env bash
#version=0.1.3
version=`cat VERSION`
time=$(date)
go build -o mediadownloader -ldflags="-X 'github.com/egeback/mediadownloader/internal/version.BuildTime=$time' -X 'github.com/egeback/mediadownloader/internal/version.BuildVersion=$version'" ./cmd/mediadownloader/main.go
