package main

import (
	"bufio"
	"context"
	"datafilter/charset"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const FILE_ADD_SUFFIX = "_filter"

var OPT = Options{}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	OPT.load()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Greet %s !", name)
}

// GetOptions 获取配置文件JSON
func (a *App) GetOptions() Options {
	return OPT.GetData()
}

// SetOptions 解析传来的配置
func (a *App) SetOptions(val Options) {
	OPT.SetData(&val)
}

type PathStruct struct {
	Directory    string `json:"directory"`
	Filename     string `json:"filename"`
	Ext          string `json:"ext"`
	AddSuffixStr string `json:"addSuffixStr"`
}

// SplitFilePath 拆分文件路径信息
func (a *App) SplitFilePath(filePath string) PathStruct {
	directory, filename, ext := splitFilePath(filePath)
	return PathStruct{
		directory,
		filename,
		ext,
		FILE_ADD_SUFFIX,
	}
}

// 分解文件路径
func splitFilePath(filePath string) (string, string, string) {
	path, filename := filepath.Split(filePath)
	fileNameSp := strings.Split(filename, ".")
	return path, fileNameSp[0], fileNameSp[1]
}

// SelectFile 选择需要处理的文件
func (a *App) SelectFile(title string, filetype string) string {
	if title == "" {
		title = "选择文件"
	}
	if filetype == "" {
		filetype = "*.txt;*.json"
	}
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// FilterFile 处理文件数据
func (a *App) FilterFile(filePath string) string {
	var checkCoding = false
	var isFileEncodingUtf8 bool
	timeStart := time.Now()
	// 读取文件
	fileHandle, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer fileHandle.Close()

	// 保存文件
	path, filename, ext := splitFilePath(filePath)
	newFilePath := path + filename + FILE_ADD_SUFFIX + "." + ext
	fileSaveHandle, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer fileSaveHandle.Close()

	accumulationLine := 0
	bufSaveHandle := bufio.NewWriter(fileSaveHandle)
	lineScanner := bufio.NewScanner(fileHandle)
	for lineScanner.Scan() {
		newStr := lineScanner.Text()
		accumulationLine++
		runtime.EventsEmit(a.ctx, "filter-change", accumulationLine)

		if !checkCoding {
			checkCoding = true
			isFileEncodingUtf8 = charset.IsUtf8(newStr)
		}
		if !isFileEncodingUtf8 {
			// 仅处理GBK情况
			newStr, _ = charset.ToUtf8(charset.GBK, newStr)
		}
		newStr = filterTextData(newStr)
		_, err := bufSaveHandle.WriteString(newStr + "\n")
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
	}
	bufSaveHandle.Flush()

	return time.Now().Sub(timeStart).String()
}

// 文本处理
func filterTextData(str string) string {
	result := str
	if OPT.EraseZeroWidthCharacter {
		result = eraseZeroWidthCharacter(str, "")
	}
	result = OPT.ExecRegex(result)
	return result
}

// 处理零宽字符
func eraseZeroWidthCharacter(str string, val string) string {
	myRegex, err := regexp.Compile("[\u200b-\u200f|\ufefe]")
	if err != nil {
		fmt.Println(err)
		return str
	}
	result := myRegex.ReplaceAllString(str, val)
	return result
}
