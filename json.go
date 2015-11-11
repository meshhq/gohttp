package gohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func JSONData(params interface{}) (*bytes.Buffer, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(jsonData)
	return buffer, nil
}

func ParseJSON(data io.Reader) (map[string]interface{}, error) {
	jsonData, err := ioutil.ReadAll(data)
	if err != nil {
		fmt.Printf("Failed to read data with error %v\n", err)
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Printf("Failed to parse json with error %v\n", err)
		return nil, err
	}
	JSONMap := response.(map[string]interface{})
	return JSONMap, err
}

func PrettyPrint(item interface{}) {
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)
}
