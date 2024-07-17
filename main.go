package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <path to ini file> <path to output file>")
		return
	}

	inputFilePath := os.Args[1]
	outputFilePath := os.Args[2]
	config, err := ParseIniFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading ini file:", err)
		return
	}

	json1 := TransformToJson1(config)
	json2 := TransformToJson2(config)

	result := map[string]interface{}{
		"output1": json1,
		"output2": json2,
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = WriteToFile(outputFilePath, jsonResult)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Output written to", outputFilePath)
}

func ParseIniFile(filename string) (map[string]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]map[string]string)
	var section string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			config[section] = make(map[string]string)
		} else {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[section][key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func TransformToJson1(config map[string]map[string]string) map[string]string {
	result := make(map[string]string)
	for section, items := range config {
		for key, _ := range items {
			jsonKey := key
			if section != "" {
				jsonKey = strings.ToUpper(section) + "_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
			}
			result[key] = "${" + jsonKey + "}"
		}
	}
	return result
}

func TransformToJson2(config map[string]map[string]string) map[string]string {
	result := make(map[string]string)
	for section, items := range config {
		for key, value := range items {
			jsonKey := strings.ToUpper(key)
			if section != "" {
				jsonKey = strings.ToUpper(section) + "_" + strings.ReplaceAll(jsonKey, ".", "_")
			}
			result[jsonKey] = value
		}
	}
	return result
}

func WriteToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
