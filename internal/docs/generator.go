package docs

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/logger"
	"os"
	"os/exec"
)

func GenerateOpenAPI() {
	if !config.Config.ApplicationCfg.Text {
		return
	}

	// go install ...
	command := exec.Command(
		"go",
		"install",
		"github.com/Zachacious/go-respec/cmd/respec@latest",
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		panic(err)
	}

	// respec . -o openapi.yaml
	command = exec.Command(
		"respec",
		".",
		"-o",
		"openapi.yaml",
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		panic(err)
	}

	logger.Logger.Info("openapi 文件生成成功")
}
