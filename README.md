# 6Estates idp-golang

A Golang SDK for communicating with the 6Estates Intelligent Document Processing(IDP) Platform.

## Documentation
The documentation for the 6Estates IDP API can be found via https://idp-sea.6estates.com/docs


## Setup
    to be filled
    
## Usage
1. To Extract Fields in Synchronous Way
If you just need to do one file at a time

```go
    package main

    import "github.com/6estates/idp-golang/idp_sdk"
    import "fmt"

    func main(){
      c:=idp_sdk.NewClient("your-token","your-region")
      params:=map[string]string{"fileType":"type-of-the-file"}
      result,err:=c.RunSimpleTask(params,"path-to-the-file")
      if err!=nil{
        fmt.Println(err)
      }
      fmt.Println(result)
    }
 ```

2. To Extract Fields in Asynchronous Way
If you need to do a batch of files
```go
    package main

    import "github.com/6estates/idp-golang/idp_sdk"
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
 ```
