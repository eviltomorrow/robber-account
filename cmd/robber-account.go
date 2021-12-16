package main

import (
	"log"

	"github.com/eviltomorrow/robber-account/internal/command"
	"github.com/eviltomorrow/robber-core/pkg/system"
)

var (
	GitSha      = ""
	GitTag      = ""
	GitBranch   = ""
	BuildTime   = ""
	MainVersion = "v3.0"
)

func main() {
	setupVersion()
	setupEnv()
	command.Execute()
}

func setupEnv() {
	if err := system.InitEnv(); err != nil {
		log.Fatalf("[Fatal] Robber-account init basic env failure, nest error: %v\r\n", err)
	}
}

func setupVersion() {
	system.MainVersion = MainVersion
	system.GitSha = GitSha
	system.GitTag = GitTag
	system.GitBranch = GitBranch
	system.BuildTime = BuildTime
}
