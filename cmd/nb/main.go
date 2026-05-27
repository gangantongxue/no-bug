package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var version = "dev"

func init() {
	if version != "dev" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	v := info.Main.Version
	if v != "" && v != "(devel)" {
		version = v
	}
}

// 注释格式映射
var commentFormats = map[string]string{
	// 单行注释语言
	".go":     "//",
	".js":     "//",
	".ts":     "//",
	".java":   "//",
	".cpp":    "//",
	".cxx":    "//",
	".c":      "//",
	".h":      "//",
	".hpp":    "//",
	".cs":     "//",
	".swift":  "//",
	".kt":     "//",
	".rs":     "//",
	".scala":  "//",
	".groovy": "//",
	".dart":   "//",
	".php":    "//",
	".rb":     "#",
	".py":     "#",
	".sh":     "#",
	".bash":   "#",
	".pl":     "#",
	".pm":     "#",
	".lua":    "--",
	".sql":    "--",
	".vue":    "<!--",
	".html":   "<!--",
	".xml":    "<!--",
	".css":    "/*",
	".scss":   "/*",
	".sass":   "/*",
	".less":   "/*",
}

// 佛像图案
const buddha = `
                    _ooo8ooo_
                   o888888888o
                   88"  .  "88
                   (|  -_-  |)
                   0\   =   /0
                 ____/'==='\____
               .' \\|       |// '.
              / \\|||   :   |||// \
             / _|||||  -:-  |||||_ \
            |   | \\\   -   /// |   |
            | \_|  ''\-----/''  |_/ |
            \  .-\__   '-'   __/-.  /
          ___'. .'   /--.--\   '. .'___
       ."" '<   '.___\_<|>_/___.'   >' "".
      | | :   '- \'.:'\ _ /':.'/ -'   : | |
      \  \ '-.    \_ __\ /__ _/    .-' /  /
======='-.____'.____ \_____/ ____.'____.-'=======
                     '=---='
                     佛祖保佑
                     no  bug
`

func main() {
	showVersion := flag.Bool("v", false, "显示版本号")
	dryRun := flag.Bool("d", false, "预览模式，只打印不写入文件")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "nb - 佛祖保佑，no bug！\n\n")
		fmt.Fprintf(os.Stderr, "Usage: nb [options] <file_path>...\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("nb version %s\n", version)
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	for _, filePath := range flag.Args() {
		if err := addBuddhaComment(filePath, *dryRun); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filePath, err)
			continue
		}
		if !*dryRun {
			fmt.Printf("Successfully added Buddha comment to %s\n", filePath)
		}
	}
}

func addBuddhaComment(filePath string, dryRun bool) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	commentPrefix, ok := commentFormats[ext]
	if !ok {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// 生成注释内容
	var comment strings.Builder

	// 处理HTML/XML注释
	if commentPrefix == "<!--" {
		comment.WriteString("<!--")
		comment.WriteString(buddha)
		comment.WriteString("-->")
	} else if commentPrefix == "/*" {
		// 处理CSS等多行注释
		comment.WriteString("/*")
		comment.WriteString(buddha)
		comment.WriteString("*/")
	} else {
		// 处理单行注释
		scanner := bufio.NewScanner(strings.NewReader(buddha))
		for scanner.Scan() {
			line := scanner.Text()
			comment.WriteString(commentPrefix)
			comment.WriteString(line)
			comment.WriteString("\n")
		}
		// 移除最后一个换行符
		commentStr := comment.String()
		if len(commentStr) > 0 {
			commentStr = commentStr[:len(commentStr)-1]
			comment.Reset()
			comment.WriteString(commentStr)
		}
	}

	// 合并内容（处理 shebang）
	rawContent := string(content)
	var newContent string
	if strings.HasPrefix(rawContent, "#!") {
		idx := strings.Index(rawContent, "\n")
		if idx != -1 {
			shebangLine := rawContent[:idx]
			rest := rawContent[idx+1:]
			newContent = shebangLine + "\n\n" + comment.String() + "\n\n" + strings.TrimLeft(rest, "\n")
		} else {
			newContent = comment.String() + "\n\n" + rawContent
		}
	} else {
		newContent = comment.String() + "\n\n" + rawContent
	}

	// 预览模式
	if dryRun {
		fmt.Print(newContent)
		return nil
	}

	// 写入文件
	return os.WriteFile(filePath, []byte(newContent), 0644)
}
