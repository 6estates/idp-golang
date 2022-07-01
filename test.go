package main

import "idp_sdk"
import "fmt"

func main(){
	c:=idp_sdk.NewClient("your-token","your-region")
	params:=map[string]string{"fileType":"type-of-the-file"}
	task,err:=c.CreateTask(params,"path-to-the-file")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(c.Poll(task))
}
