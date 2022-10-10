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
	isOauth bool
}

type FileType struct {
	bank_statement, invoice, cheque, credit_bureau_singapore, receipt, payslip, packing_list, bill_of_lading, air_waybill, kartu_tanda_penduduk,
	hong_kong_annual_return, purchase_order, delivery_order string
	list map[string]bool
}

func FileTypeList() (f *FileType) {
	f = &FileType{
		bank_statement:          "CBKS",
		invoice:                 "CINV",
		cheque:                  "CHQ",
		credit_bureau_singapore: "CBS",
		receipt:                 "RCPT",
		payslip:                 "PS",
		packing_list:            "PL",
		bill_of_lading:          "BL",
		air_waybill:             "AWBL",
		kartu_tanda_penduduk:    "KTP",
		hong_kong_annual_return: "HKAR",
		purchase_order:          "PO",
		delivery_order:          "DO",
		list: map[string]bool{
			"CBKS": true,
			"CINV": true,
			"CHQ":  true,
			"CBS":  true,
			"RCPT": true,
			"PS":   true,
			"PL":   true,
			"BL":   true,
			"AWBL": true,
			"KTP":  true,
			"HKAR": true,
			"PO":   true,
			"DO":   true,
		},
	}
	return f
}

func (c *Client) ResultFields(result map[string]interface{}) interface{} {
	return result["data"].(map[string]interface{})["fields"]
}

func (c *Client) FieldCode(field map[string]interface{}) string {
	return fmt.Sprintf("%v", field["field_code"])
}

func (c *Client) FieldName(field map[string]interface{}) string {
	return fmt.Sprintf("%v", field["field_name"])
}

func (c *Client) FieldValue(field map[string]interface{}) string {
	return fmt.Sprintf("%v", field["value"])
}

func (c *Client) FieldType(field map[string]interface{}) string {
	return fmt.Sprintf("%v", field["type"])
}

func NewClient(token, region string, isOauth ...bool) (c *Client) {
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
	c.isOauth = false
	if len(isOauth) > 0 {
    		c.isOauth = isOauth[0]
        }
	
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
	ret := make(map[string]interface{})
	if FileTypeList().list[params["fileType"]] != true {
		return ret, Error("File type is invalid")
	}
	body := &bytes.Buffer{}

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
	
	req, err := http.NewRequest("POST", c.url_post, body)
	if err != nil {
		return ret, err
	}

	if c.isOauth{
		req.Header.Set("Authorization", c.token)
	}else {
		req.Header.Set("X-ACCESS-TOKEN", c.token)
	}
	
	req.Header.Set("Content-Type", content_type)
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
	if c.isOauth{
		req.Header.Set("Authorization", c.token)
	}else {
		req.Header.Set("X-ACCESS-TOKEN", c.token)
	}
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
	task, err := c.CreateTask(params, filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return c.Poll(task)
}


func  OauthUtil(authorization, region string) (map[string]interface{}, error) {
        if authorization  == "" {
		Error("Authorization  is required!")
	}
	reg := make(map[string]bool)
	reg["sea"] = true
	reg["test"] = true
	if reg[region] != true {
		Error("Region is limited in ['test','sea']!")
	}

	if region == "test" {
		region = "-onp"
	} else {
		region = "-" + region
	}
	
	url_post := "https://oauth" + region + ".6estates.com/oauth/token?grant_type=client_bind"
        ret := make(map[string]interface{})
	req, _ := http.NewRequest("POST", url_post, nil) 
	req.Header.Set("Authorization", authorization)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return  ret, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return  ret, err
	}
	
	isExpired := ret["data"].(map[string]interface{})["expired"].(bool)

	if isExpired {
	    return ret, fmt.Errorf("This IDP Authorization is expired, please re-send the request to get new IDP Authorization.")
	}
	return ret, nil
}
