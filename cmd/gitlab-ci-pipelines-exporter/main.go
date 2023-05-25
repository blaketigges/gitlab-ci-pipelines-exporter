package main

import (
	"os"

	"github.com/blaketigges/gitlab-ci-pipelines-exporter/internal/cli"
)

var version = "devel"

func main() {
	cli.Run(version, os.Args)
}
