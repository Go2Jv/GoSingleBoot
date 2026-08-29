package docs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/logger"
)

// respec 安装包路径(不含 @latest 之外的版本约束,始终安装最新版)
const respecPkg = "github.com/Zachacious/go-respec/cmd/respec@latest"

// respecName 返回当前平台下 respec 的可执行文件名。
// Windows 下可执行文件带 .exe 后缀,macOS/Linux 没有
func respecName() string {
	if runtime.GOOS == "windows" {
		return "respec.exe"
	}
	return "respec"
}

// getGoPath 通过 go env GOPATH 获取实际 GOPATH,而不是假设 ~/go。
// 输出会带尾部换行,需要 TrimSpace
func getGoPath() (string, error) {
	out, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return "", fmt.Errorf("执行 go env GOPATH 失败: %w", err)
	}
	gopath := strings.TrimSpace(string(out))
	if gopath == "" {
		return "", fmt.Errorf("go env GOPATH 返回为空")
	}
	return gopath, nil
}

// findRespec 依次在 PATH、GOPATH/bin 中查找 respec 的绝对路径,找不到返回空字符串。
// 优先 PATH(兼容已把工具目录加入 PATH 的环境),
// 再查 GOPATH/bin(go install 的默认安装位置,即使不在 PATH 中也能找到)
func findRespec() string {
	if path, err := exec.LookPath(respecName()); err == nil {
		return path
	}

	gopath, err := getGoPath()
	if err != nil {
		logger.Logger.Sugar().Warnf("获取 GOPATH 失败,无法在 GOPATH/bin 中查找 respec: %v", err)
		return ""
	}
	candidate := filepath.Join(gopath, "bin", respecName())
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		return candidate
	}
	return ""
}

// installRespec 安装 respec 并返回安装后的绝对路径。
// 安装只是兜底,不应该成为每次启动的固定步骤
func installRespec() (string, error) {
	cmd := exec.Command("go", "install", respecPkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go install %s 失败: %w", respecPkg, err)
	}

	// 安装完成后当前进程的 PATH 不会自动更新,
	// 重新通过 GOPATH/bin 确定 respec 的实际绝对路径
	gopath, err := getGoPath()
	if err != nil {
		return "", err
	}
	bin := filepath.Join(gopath, "bin", respecName())
	if _, err := os.Stat(bin); err != nil {
		return "", fmt.Errorf("go install 完成但在 %s 未找到 respec: %w", bin, err)
	}
	return bin, nil
}

// GenerateOpenAPI 在开发环境(ApplicationCfg.Text 开启)下生成 openapi.yaml。
// 属于辅助能力,任何一步失败只记录日志并返回,不影响主服务启动
func GenerateOpenAPI() {
	if !config.Config.ApplicationCfg.Text {
		return
	}

	// 先查找 respec,确实找不到时才安装
	respecPath := findRespec()
	if respecPath == "" {
		logger.Logger.Info("未找到 respec,执行 go install ...")
		var err error
		respecPath, err = installRespec()
		if err != nil {
			logger.Logger.Sugar().Errorf("安装 respec 失败: %v", err)
			return
		}
	}

	// 使用绝对路径执行,不依赖 PATH 中注册了 GOPATH/bin
	cmd := exec.Command(respecPath, ".", "-o", "openapi.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Logger.Sugar().Errorf("生成 openapi.yaml 失败: %v", err)
		return
	}

	logger.Logger.Info("openapi 文件生成成功")
}
