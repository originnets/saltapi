package main

import (
	"fmt"
)

func mian() {
	url := "http://192.168.1.218:8888"
	name := "root"
	password := "123456"
	// 配置文件
	var config = Cfg{
		url,
		Auth{name, password, "pam"},
	}
	// 创建客户端
	client, err := New(&config)
	if err != nil {
		fmt.Print(err)
	}
	// 推送指令并返回id
	id, err := client.Execute("local", "test.ping", "", "*", "glob")
	if err != nil {
		fmt.Print(err)
	}
	//读取id返回执行结果
	job, err := client.Job(id)
	if err != nil {
		fmt.Print(err)
	}
	for _, value := range job.Minions {
		fmt.Print(job.Result[value].Return)
	}
}
