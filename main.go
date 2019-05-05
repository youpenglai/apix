package main

import (
	"github.com/youpenglai/apix/mgr"
	"github.com/youpenglai/goutils/logger"
)

func main() {
	l := logger.GetLogger("ApiX")
	l.Info("ApiX Started")
	mgr.RunManagerServer()
}
