# 6Estates idp-golang

A Golang SDK for communicating with the 6Estates Intelligent Document Processing(IDP) Platform.

## Documentation
The documentation for the 6Estates IDP API can be found via https://idp-sea.6estates.com/docs


## Setup
    go get github.com/6estates/idp-golang/idp_sdk
    
## Usage
### 1. Initialize the 6Estates IDP Client
6E API Access Token(Deprecated)
```go
    package main

    import "github.com/6estates/idp-golang/idp_sdk"

    func main(){
      c:=idp_sdk.NewClient("your-token","your-region")
    }
 ```
 
6E API Authorization based on oauth 2.0
```go
    package main

    import "github.com/6estates/idp-golang/idp_sdk"
    import "fmt"
    
    func main(){
      ret, err :=idp_sdk.OauthUtil("your-authorization"，"your-region")
      if err != nil {
		fmt.Println(err)
	  }
      oauth:=ret["data"].(map[string]interface{})["value"].(string)
      isOauth:=true
      c:=idp_sdk.NewClient(oauth, "your-region", isOauth)
    }
 ```

### 2. To Extract Fields in Synchronous Way
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

### 3. To Extract Fields in Asynchronous Way<br>
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
