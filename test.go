package main

import "idp_sdk"
import "fmt"

func main(){
	c:=idp_sdk.NewClient("QWhPM7Wxqr3xc4dCQFmXhH8xYD8CTq3N41XnvV38OblJQpTw5R9DyKHA0coN5m81","test")
	params:=map[string]string{"fileType":"CBKS"}
	task,err:=c.CreateTask(params,"E:\\work\\idp-sdk\\idp_sdk_go\\[UOB]202103_UOB_2222.pdf")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(c.Poll(task))
}