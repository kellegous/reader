package main

import (
	"github.com/kellegous/glue/logging/yarder/zap"
	"github.com/kellegous/poop"
	"github.com/kellegous/reader/internal/cmd"
)

func main() {
	if err := zap.Register(zap.WithApp("reader")); err != nil {
		poop.HitFan(err)
	}
	cmd.Execute()
}
