package util

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)


func SaveMapYaml(data interface{}){
	yamlData,err:= yaml.Marshal(&data)
	if err!=nil{
		panic("sibla")
	}
	filepath:= "netData.yaml"

	err =os.WriteFile(filepath, yamlData, 0644) // 0644 for read/write by owner, read-only by others
	if err != nil {
		log.Fatalf("Error writing YAML to file: %v", err)
	}
}

func ReadMapNetYaml(loadedData map[string]string)error{
	filePath := "netData.yaml" // 읽어올 YAML 파일 경로

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlData, &loadedData)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML: %v", err)
	}
	return nil
}