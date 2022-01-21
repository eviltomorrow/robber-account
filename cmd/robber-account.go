package main

import (
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

func init() {
	system.MainVersion = MainVersion
	system.GitSha = GitSha
	system.GitTag = GitTag
	system.GitBranch = GitBranch
	system.BuildTime = BuildTime
}

func main() {
	command.Execute()
}
