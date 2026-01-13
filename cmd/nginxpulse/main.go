package main

import (
	"os"

	"github.com/likaia/nginxpulse/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := app.Run(); err != nil {
		logrus.WithError(err).Error("服务启动失败")
		os.Exit(1)
	}
}
