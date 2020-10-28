package novelread

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var (
	//数字
	number = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9"}
	//简写
	simplified = map[int]string{0: "零", 1: "一", 2: "二", 3: "三", 4: "四", 5: "五", 6: "六", 7: "七", 8: "八", 9: "九", 10: "十", 11: "百", 12: "千"}
	//繁体
	traditional = map[int]string{0: "零", 1: "壹", 2: "贰", 3: "叁", 4: "肆", 5: "伍", 6: "陆", 7: "柒", 8: "捌", 9: "玖", 10: "拾", 11: "佰", 12: "仟"}

	//章节 章节集回话 ----- 章节识别关键字
	chapter = map[int]string{0: "章", 1: "节", 2: "集", 3: "回", 4: "话"}
	//卷
	reel = "卷"
)

//Chapter 小说结构
type Chapter struct {
	Title   string
	Index   int
	Volume  int
	Content string
}

//TrimFile 小说文件处理
func TrimFile(filePath string) {
	s, err := getFileContentAsStringLines(filePath)
	if err != nil {
		log.Println(err)
	}
	for _, v := range s {
		lineTextDiscern(v)
	}
	wg.Done()
}
func getFileContentAsStringLines(filePath string) ([]string, error) {

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

//lineTextDiscern 行文本识别
func lineTextDiscern(line string) {
	c := Chapter{}
	c.Volume = 1   //卷号
	c.Index = 1    //章节号
	c.Title = ""   //章节标题
	c.Content = "" //章节内容
	chapterVolume(line)
}

//chapterVolume 卷号识别并且提取
func chapterVolume(s string) int {
	countSplit := strings.Split(s, "")
	for k, v := range countSplit {
		if strings.Contains(v, reel) {
			fmt.Println(countSplit[k])
		}
	}
	return 0
}
