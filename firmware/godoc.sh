#!/usr/bin/env bash

GOROOT=/usr/share/go-1.23 GOPROXY="direct" godoc -http=:8900 &
xdg-open http://localhost:8900/pkg/eclair/
