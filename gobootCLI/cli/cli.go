package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// 模板仓库地址,clone 完成后新项目的 module 名会替换掉 oldModule
const (
	templateRepo = "https://github.com/Go2Jv/GoSingleBoot.git"
	oldModule    = "GoSingleBoot"
)

// templateFiles 是 clone 后需要从新项目中删除的模板残留文件/目录。
// .git 必须删除(新项目会重新 git init);README/LICENSE 是模板仓库自身的说明;
// gobootCLI 是 CLI 工具自身,不能随模板一起留给新项目,否则会残留嵌套 module 导致 go mod tidy 失败。
var templateFiles = []string{".git", "README.md", "LICENSE", "gobootCLI"}

// Run 执行完整的创建流程,任何一步失败都直接返回错误并退出
func Run() error {
	in := bufio.NewReader(os.Stdin)

	// 1. 询问项目名称
	projectName, err := ask(in, "Project name: ")
	if err != nil {
		return err
	}
	if err := validateProjectName(projectName); err != nil {
		return err
	}

	// 2. 询问代码存放位置,直接回车默认当前目录
	location, err := ask(in, "Project location (default: current directory): ")
	if err != nil {
		return err
	}
	if location == "" {
		location = "."
	}
	location, err = expandHome(location)
	if err != nil {
		return err
	}
	projectDir := filepath.Join(location, projectName)

	// 目标目录已存在则报错退出,不覆盖
	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("目录 %s 已存在,请换一个项目名", projectDir)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查目录 %s 失败: %w", projectDir, err)
	}
	if err := os.MkdirAll(location, 0o755); err != nil {
		return fmt.Errorf("创建目录 %s 失败: %w", location, err)
	}

	// 3. clone 模板(默认 clone GitHub 仓库,可用环境变量 GOBOT_TEMPLATE_REPO 覆盖,便于测试)
	repo := templateRepo
	if r := os.Getenv("GOBOT_TEMPLATE_REPO"); r != "" {
		repo = r
	}
	fmt.Printf("Cloning template into %s...\n", projectDir)
	clone := exec.Command("git", "clone", "--depth", "1", repo, projectDir)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	if err := clone.Run(); err != nil {
		// 清理 clone 失败留下的半成品目录,避免下次运行时误报"目录已存在"
		os.RemoveAll(projectDir)
		return fmt.Errorf("git clone 失败: %w", err)
	}

	// clone 一完成就清理模板文件,保证后续 go mod edit / tidy 只作用于新项目源码
	if err := cleanClone(projectDir); err != nil {
		return fmt.Errorf("清理模板文件失败: %w", err)
	}
	fmt.Printf("Removed template files: %s\n", strings.Join(templateFiles, ", "))

	// 把新项目 config.json 中的 Application.Name 改为用户输入的项目名
	if err := updateConfig(projectDir, projectName); err != nil {
		return fmt.Errorf("更新 config.json 失败: %w", err)
	}

	// 4. 询问 Go module path,直接回车默认使用项目名
	modulePath, err := ask(in, fmt.Sprintf("Go module path (default: %s): ", projectName))
	if err != nil {
		return err
	}
	if strings.TrimSpace(modulePath) == "" {
		modulePath = projectName
	}

	// 修改新项目的 go.mod
	gomod := filepath.Join(projectDir, "go.mod")
	if _, err := os.Stat(gomod); err != nil {
		return fmt.Errorf("找不到 %s,模板仓库可能缺少 go.mod: %w", gomod, err)
	}
	fmt.Printf("Setting go.mod module to %s...\n", modulePath)
	edit := exec.Command("go", "mod", "edit", "-module", modulePath, gomod)
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr
	if err := edit.Run(); err != nil {
		return fmt.Errorf("go mod edit 失败: %w", err)
	}

	// 同步模板的 go 版本号到用户本机的 Go 版本(静默执行,不打扰用户)。
	// 每个用户电脑上的 Go 版本不同,模板自带的 go 版本号若高于本机版本,
	// 后续 go mod tidy 会直接失败
	if err := syncGoVersion(gomod); err != nil {
		return err
	}

	// 5. 替换源码 import 中的旧 module 路径(只处理 .go 文件,跳过 .git 目录)
	if err := replaceModule(projectDir, modulePath); err != nil {
		return err
	}

	// 6. 在新项目目录中执行 go mod tidy
	fmt.Printf("Running go mod tidy in %s...\n", projectDir)
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = projectDir
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy 失败: %w", err)
	}

	// 7. 询问是否初始化 Git,默认 Yes
	gitInit := true
	answer, err := ask(in, "Initialize Git repository? (Y/n): ")
	if err != nil {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "n", "no":
		gitInit = false
	}

	if gitInit {
		init := exec.Command("git", "init")
		init.Dir = projectDir
		init.Stdout = os.Stdout
		init.Stderr = os.Stderr
		if err := init.Run(); err != nil {
			return fmt.Errorf("git init 失败: %w", err)
		}
	}

	// 8. 完成
	fmt.Println()
	fmt.Println("Project created successfully!")
	fmt.Println()
	fmt.Println("Location:")
	fmt.Println(projectDir)
	fmt.Println()
	fmt.Println("Next:")
	fmt.Println()
	fmt.Printf("cd %s\n", projectDir)
	fmt.Println("go run cmd/main.go")
	if gitInit {
		fmt.Println()
		fmt.Println("Git repository initialized.")
	}
	return nil
}

// ask 打印提示并读取一行用户输入,返回去除首尾空白后的内容
func ask(in *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// validateProjectName 校验项目名:
// 不能为空、不能包含路径分隔符、不能以 . 开头(隐藏目录,且可能与 .git 等模板清理逻辑冲突)
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("项目名不能为空")
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("项目名 %q 不能包含路径分隔符", name)
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("项目名 %q 不能以 . 开头", name)
	}
	return nil
}

// expandHome 把以 ~ 开头的路径展开为用户主目录,便于输入 ~/projects 这类位置
func expandHome(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// replaceModule 遍历新项目中的 .go 文件,
// 把 import 里的旧 module 路径(GoSingleBoot/ 开头)替换成新的 module path
func replaceModule(dir, modulePath string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// 跳过 .git 目录(含嵌套仓库)。模板的 gobootCLI 子目录已在 cleanClone 中删除,
			// 这里不再按名字跳过其他目录,避免项目名恰好是 gobootCLI 时误跳过整个项目
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// 只匹配 import 路径开头("GoSingleBoot/ 形式),避免误伤注释等位置
		oldImport := `"` + oldModule + `/`
		newImport := `"` + modulePath + `/`
		if !strings.Contains(string(data), oldImport) {
			return nil
		}

		updated := strings.ReplaceAll(string(data), oldImport, newImport)
		if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
			return err
		}
		fmt.Printf("Updated imports in %s\n", path)
		return nil
	})
}

// syncGoVersion 把新项目 go.mod 的 go 版本号改为用户本机的 Go 版本。
// 整个步骤静默执行:不打印任何内容,也不让 go 命令的输出出现在终端,
// 用户只需要知道结果(go.mod 已适配本机环境)
func syncGoVersion(gomod string) error {
	verOut, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		return fmt.Errorf("获取本机 Go 版本失败: %w", err)
	}

	// go1.26.4 这类输出只保留主.次版本,小版本号不是语言版本约束
	m := regexp.MustCompile(`^go(\d+)\.(\d+)`).FindStringSubmatch(strings.TrimSpace(string(verOut)))
	if m == nil {
		return fmt.Errorf("无法解析本机 Go 版本: %q", strings.TrimSpace(string(verOut)))
	}
	ver := m[1] + "." + m[2]

	edit := exec.Command("go", "mod", "edit", "-go", ver, gomod)
	if out, err := edit.CombinedOutput(); err != nil {
		return fmt.Errorf("修改 go.mod 的 go 版本失败: %w\n%s", err, out)
	}
	return nil
}

// updateConfig 把新项目 config.json 中的 Application.Name 改为项目名,
// 只做最小编码修改,保留原文件的其他内容(缩进、注释、键顺序)
func updateConfig(dir, name string) error {
	configPath := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg struct {
		Application struct {
			Name string `json:"Name"`
		} `json:"Application"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析 %s: %w", configPath, err)
	}
	if cfg.Application.Name == "" {
		return fmt.Errorf("%s 中缺少 Application.Name", configPath)
	}
	if cfg.Application.Name == name {
		return nil
	}

	oldVal := `"Name": "` + cfg.Application.Name + `"`
	// 用 json.Marshal 转义,防止项目名含引号、反斜杠等字符时写坏 JSON
	newVal := `"Name": ` + string(mustJSONString(name))
	updated := strings.Replace(string(data), oldVal, newVal, 1)
	if updated == string(data) {
		return fmt.Errorf("在 %s 中未找到 Application.Name", configPath)
	}
	if err := os.WriteFile(configPath, []byte(updated), 0o644); err != nil {
		return err
	}
	fmt.Printf("Updated Application.Name in %s\n", configPath)
	return nil
}

// mustJSONString 把字符串编码为 JSON 字符串字面量(含双引号)。
// name 是普通字符串,json.Marshal 不会失败,忽略错误
func mustJSONString(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		return []byte(`""`)
	}
	return b
}

// cleanClone 删除 clone 时新项目里遗留的模板文件(见 templateFiles)
func cleanClone(dir string) error {
	for _, name := range templateFiles {
		if err := os.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}
	return nil
}
