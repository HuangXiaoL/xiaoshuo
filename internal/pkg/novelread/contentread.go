package novelread

import (
	"io/ioutil"
	"log"
	"strings"
)

func GetFileContentAsStringLines(filePath string) ([]string, error) {

	log.Printf("get file content as lines: %v", filePath)
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("read file: %v error: %v", filePath, err)
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	log.Printf("get file content as lines: %v, size: %v", filePath, len(result))
	return result, nil
}
