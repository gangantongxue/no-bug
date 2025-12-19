package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
                  _ooOoo_
                 o8888888o
                 88" . "88
                 (| -_- |)
                 O\\  =  /O
              ____/` + "`" + `---'____
            .'  \\|     |//  '.
           /  \|||  :  |||//  \
          /  _||||| -:- |||||-  \
          |   | \\  -  /// |   |
          | \_|  ''\---/''  |   |
          \\  .-\\__  '-'  ___/-. /
        ___'. .'  /--.--\  '. . __
     ."" '<  '.___\\_<|>_/___.'  >'"".
    | | :  '- \'.;'\ _ /';.'/ - ' : | |
    \\  \\ '-.   \_ __\\ /__ _/   .-' /  /
======='-.____'-.___\\_____/___.-'____.-'=======
                  '=---='

                  佛祖保佑
                  no  bug
`

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: nb <file_path>...")
		os.Exit(1)
	}

	for _, filePath := range flag.Args() {
		if err := addBuddhaComment(filePath); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			continue
		}
		fmt.Printf("Successfully added Buddha comment to %s\n", filePath)
	}
}

func addBuddhaComment(filePath string) error {
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

	// 合并内容
	newContent := comment.String() + "\n\n" + string(content)

	// 写入文件
	return os.WriteFile(filePath, []byte(newContent), 0644)
}
