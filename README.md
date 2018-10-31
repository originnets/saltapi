# SaltApi介绍
  实现了salt的salt任务执行并返回执行id值,通过id值返回执行结果
1. 创建客户端
  
		  var config = Cfg{
		  url,
		  Auth{name, password, "pam"},
		  }
2. 创建客户端

		client, err := New(&config)
		if err != nil {
			fmt.Print(err)
		}
3. 推送指令并返回id

		id, err := client.Execute("local", "test.ping", "", "*", "glob")
		if err != nil {
			fmt.Print(err)
		}
		
4. 读取id返回执行结果

		job, err := client.Job(id)
		if err != nil {
			fmt.Print(err)
		}
		for _, value := range job.Minions {
			fmt.Print(job.Result[value].Return)
		}
