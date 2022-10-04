package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

const FILE_NAME = "options.json"

type Options struct {
	EraseZeroWidthCharacter bool          `json:"eraseZeroWidthCharacter"`
	Regex                   []RegexConfig `json:"regex"`
	eachList                []EachList
}

type RegexConfig struct {
	Regexp string `json:"regexp"`
	Value  string `json:"value"`
}

type EachList struct {
	regexp *regexp.Regexp
	value  string
}

// 载入配置
func (opt *Options) load() {
	cwd, _ := os.Getwd()

	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(content) == 0 {
		opt.EraseZeroWidthCharacter = true
		opt.Regex = []RegexConfig{}
		data, _ := json.Marshal(opt)
		file.WriteString(string(data))
		return
	}

	json.Unmarshal([]byte(content), &opt)
}

// 保存配置文件
func (opt *Options) save() {
	cwd, _ := os.Getwd()
	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	data, _ := json.Marshal(opt)
	file.WriteString(string(data))
}

// 设置值
func (opt *Options) set(key string, val bool) {
	switch key {
	case "EraseZeroWidthCharacter":
		opt.EraseZeroWidthCharacter = val
	}
}

// 添加正则
func (opt *Options) addRegex(key string, val string) {
	opt.Regex = append(opt.Regex, RegexConfig{key, val})
}

// SetData 存储新的配置
func (opt *Options) SetData(newOpt *Options) {
	if reflect.DeepEqual(newOpt, opt) {
		return
	}
	optionByte, _ := json.Marshal(newOpt)
	err := json.Unmarshal(optionByte, &opt)
	if err != nil {
		return
	}
	opt.eachList = opt.eachList[0:0]
	opt.save()
}

// GetData 获取配置信息
func (opt *Options) GetData() Options {
	return *opt
}

// 初始化正则列表
func (opt *Options) initEachList() {
	opt.eachList = []EachList{}
	for _, item := range opt.Regex {
		if item.Regexp != "" {
			regex, err := regexp.Compile(item.Regexp)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			item := EachList{
				regex,
				item.Value,
			}
			opt.eachList = append(opt.eachList, item)
		}
	}
}

// ExecRegex 执行每一个正则替换
func (opt *Options) ExecRegex(str string) string {
	if len(opt.eachList) == 0 {
		opt.initEachList()
	}
	result := str
	for _, item := range opt.eachList {
		result = item.regexp.ReplaceAllString(result, item.value)
	}
	return result
}
