package idp_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var (
	HttpClient = &http.Client{
		Timeout: 3 * time.Second,
	}
)

type MyError struct {
	msg string
}

func (error *MyError) Error() string {
	return error.msg
}

func Error(msg string) error {
	return &MyError{msg}
}

type Client struct {
	token, region, url_post, url_get string
}

func NewClient(token, region string) (c *Client) {
	if token == "" {
		Error("Token is required!")
	}
	reg := make(map[string]bool)
	reg["sea"] = true
	reg["test"] = true
	if reg[region] != true {
		Error("Region is limited in ['test','sea']!")
	}
	c = &Client{}
	c.token = token
	c.region = region
	if region == "test" {
		region = ""
	} else {
		region = "-" + region
	}
	c.url_post = "https://idp" + region + ".6estates.com/customer/extraction/fields/async"
	c.url_get = "https://idp" + region + ".6estates.com/customer/extraction/field/async/result/"
	return c
}

func (c *Client) CreateTask(params map[string]string, filePath string) (map[string]interface{}, error) {
	body := &bytes.Buffer{}
	ret := make(map[string]interface{})
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file")
		return ret, err
	}
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("error fetching stat of the file")
		return ret, err
	}
	defer file.Close()
	writer := multipart.NewWriter(body)
	content_type := writer.FormDataContentType()

	formFile, err := writer.CreateFormFile("file", fi.Name())

	_, err = io.Copy(formFile, file)
	if err != nil {
		return ret, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return ret, err
	}
	// fmt.Println(body)
	req, err := http.NewRequest("POST", c.url_post, body)
	if err != nil {
		return ret, err
	}
	// fmt.Println(req)
	//req.Header.Set("Content-Type","multipart/form-data")
	// req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-ACCESS-TOKEN", c.token)
	req.Header.Set("Content-Type", content_type)
	// fmt.Println(req)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return ret, err
	}
	fmt.Println(ret)
	return ret, nil
}

func (c *Client) TaskID(task map[string]interface{}) interface{} {
	return task["data"]
}

func (c *Client) TaskResult(task map[string]interface{}) (map[string]interface{}, error) {
	taskID := fmt.Sprintf("%v", c.TaskID(task))
	ret := make(map[string]interface{})
	req, _ := http.NewRequest("GET", c.url_get+taskID, nil)
	req.Header.Set("X-ACCESS-TOKEN", c.token)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return ret, err
	}
	return ret, nil

}

func (c *Client) ResultStatus(ret map[string]interface{}) string {
	data := ret["data"].(map[string]interface{})
	return fmt.Sprintf("%v", data["taskStatus"])
}

func (c *Client) Poll(task map[string]interface{}) (map[string]interface{}, error) {
	ret, err := c.TaskResult(task)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	for i := 0; i < 200; i++ {
		time.Sleep(3 * time.Second)
		ret, err := c.TaskResult(task)
		if err != nil {
			fmt.Println(err)
			break
		}
		if !(c.ResultStatus(ret) == "Doing" || c.ResultStatus(ret) == "Init") {
			return ret, nil
		}
		fmt.Println(ret)

	}
	return ret, nil
}

func (c *Client) RunSimpleTask(params map[string]string, filePath string) (map[string]interface{}, error) {
	task, err := c.CreateTask(params, "E:\\work\\idp-sdk\\idp_sdk_go\\[UOB]202103_UOB_2222.pdf")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return c.Poll(task)
}
