package novelread

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkumza/numcn"
)

//有缩进就算正文
//超过80个字算正文
//从左边数超过5个字都没有数字算正文
//数字之后，如果存在“章节集回话”这几个字，算识别到了章节
//要识别“卷”，卷可能是单独一行，也可能跟章节放到一起，如果章节那一行没有包含卷信息，使用之前识别到的卷信息

var (
	//数字
	number = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9"}
	//简写
	simplified = map[int]string{0: "零", 1: "一", 2: "二", 3: "三", 4: "四", 5: "五", 6: "六", 7: "七", 8: "八", 9: "九", 10: "十", 11: "百", 12: "千"}
	//繁体
	traditional = map[int]string{0: "零", 1: "壹", 2: "贰", 3: "叁", 4: "肆", 5: "伍", 6: "陆", 7: "柒", 8: "捌", 9: "玖", 10: "拾", 11: "佰", 12: "仟"}

	//章节 章节集回话 ----- 章节识别关键字
	chapter = map[int]string{0: "章", 1: "节", 2: "集", 3: "回", 4: "话", 5: " "}
	//卷
	reel = "卷"
)

//Chapter 小说结构
type Chapter struct {
	Titles  string
	Volume  int
	Index   int
	Content string
}

//GetChapter 获取结构体
func GetChapter() *Chapter {
	return &Chapter{}
}

//TrimFile 小说文件处理
func TrimFile(filePath string) {
	s, err := getFileContentAsStringLines(filePath)
	if err != nil {
		log.Println(err)
	}
	c := GetChapter()
	var conts string
	for _, v := range s {
		cont, vnum, cnum, t := lineTextDiscern(v)
		conts = conts + cont
		if cnum != 0 {
			c.Titles = strings.Trim(t, "\n\r")
			c.Volume = vnum
			c.Index = cnum
			//c.Content = conts
			//conts = ""
			fmt.Printf("%+v\n", c)
		}

	}
	wg.Done()
}

//getFileContentAsStringLines 读取行
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
		//lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	log.Printf("get file content as lines: %v, size: %v", filePath, len(result))
	return result, nil
}

//lineTextDiscern 行文本识别
func lineTextDiscern(line string) (string, int, int, string) {
	b := lineRetractIsContent(line)
	length := lineLength(line)
	//赋值逻辑
	if b { //是否有缩进 有就返回为正文
		return strings.TrimSpace(line), 0, 0, ""
	} else if length { //长度是否超过80 有就返回为正文
		return strings.TrimSpace(line), 0, 0, ""
	}
	// 提取 卷号 ，章节号，章节名称
	v, ch, t := lineFindNumAtChapterAndVolume(line, 5)
	return "", v, ch, t

}

//lineRetractIsContent判断是否有缩进是否为正文
func lineRetractIsContent(line string) bool {
	countSplit := strings.Split(line, "") //切割字符串
	if len(countSplit) > 2 {
		for k, r := range line {
			if k < 2 {
				return unicode.IsSpace(r)
			}

		}
	}
	return false
}

//lineLengthIsContent 根据行的长度判断是否正文
func lineLength(line string) bool {
	if len(line) > 80 {
		return true
	}
	return false
}

//lineFindNumAtChapterAndVolume 在行里查找数字并且返回章卷的值
func lineFindNumAtChapterAndVolume(line string, seat int) (int, int, string) {
	var (
		volumeNum  int      // 卷 值
		chapterNum int      // 章值
		title      []string //标题值
	)
	countSplit := strings.Split(line, "") //切割字符串
	s := countSplit
	if len(countSplit) >= seat { //如果行的长度大于设定值，就截取设定的长度
		s = countSplit[0:seat]
	}
	if len(s) == 0 { // 行的长度为0 直接返回
		return 0, 0, ""
	}

	// 获取卷值
	if i := getNumber(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		volumeNum = getVolumeNum(line) //卷值
	} else if i := getSimplified(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		volumeNum = getVolumeNum(line) //卷值
	} else if i := getTraditional(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		volumeNum = getVolumeNum(line) //卷值
	}
	// 获取章值
	if i := getNumber(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		chapterNum, title = getChapterNum(line) //章值
	} else if i := getSimplified(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		chapterNum, title = getChapterNum(line) //章值
	} else if i := getTraditional(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		chapterNum, title = getChapterNum(line) //章值
	}
	t := ""
	for _, v := range title {
		t = t + v
	}
	return volumeNum, chapterNum, t
}

//getNumber 获取每行的阿拉伯数字
func getNumber(line []string) int {
	ss := ""
	for _, v := range line {
		for _, vv := range number { //循环数字做匹配
			if v == vv {
				ss = ss + v
			}
		}
	}

	if len(ss) != 0 { //每行的数字
		i, _ := strconv.Atoi(ss)
		return i
	}
	return 0
}

//getSimplified 获取简写的数字的值
func getSimplified(line []string) int64 {
	ss := ""
	for _, v := range line {
		for _, vv := range simplified { //循环数字做匹配
			if v == vv {
				ss = ss + v
			}
		}
	}
	if len(ss) != 0 { //每行的数字
		num, _ := numcn.DecodeToInt64(ss)
		return num
	}
	return 0
}

//getTraditional 获取繁体的数字的值
func getTraditional(line []string) int64 {
	ss := ""
	for _, v := range line {
		for _, vv := range traditional { //循环数字做匹配
			if v == vv {
				ss = ss + v
			}
		}
	}
	if len(ss) != 0 { //每行的数字
		num, _ := numcn.DecodeToInt64(ss)
		return num
	}
	return 0
}

//getVolumeNum 卷号识别并且提取
func getVolumeNum(s string) int {
	countSplit := strings.Split(s, "")
	for k, v := range countSplit {
		if strings.Contains(v, reel) {
			stk := 0
			if k >= 10 {
				stk = k - 10
			}
			edk := k + 1
			result := countSplit[stk:edk]      //截取卷前10以内的字符
			if i := getNumber(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return i
			} else if i := getSimplified(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return int(i)
			} else if i := getTraditional(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return int(i)
			}
		}
	}
	return 0
}

//TODO:需要拆分优化,循环判断太多了，需要拆分处理
//getChapterNum 识别章数并提取
func getChapterNum(s string) (int, []string) {
	countSplit := strings.Split(s, "")
	for k, v := range countSplit {
		if strings.Contains(v, reel) { //判断是否有章节前缀
			result := countSplit[k+1:]   //去掉章节前缀之后的
			for rk, rv := range result { //循环卷之后的结果
				for _, vv := range chapter {
					if strings.Contains(rv, vv) { //章节集回话
						stk := 0 //章节号起始
						if rk >= 10 {
							stk = rk - 10
						}
						edk := rk + 1                    //章节号结束
						chapterResult := result[stk:edk] //截取卷前10以内的字符
						title := result[edk:]
						if i := getNumber(chapterResult); i > 0 { //返回的数字大于0 有可能是章节目录有卷
							return i, title
						} else if i := getSimplified(chapterResult); i > 0 { //返回的数字大于0 有可能是章节目录有卷
							return int(i), title
						} else if i := getTraditional(chapterResult); i > 0 { //返回的数字大于0 有可能是章节目录有卷
							return int(i), title
						}
					}
				}
			}
		}

		for _, vv := range chapter {
			if strings.Contains(v, vv) {
				stk := 0 //章节号起始
				if k >= 10 {
					stk = k - 10
				}
				edk := k + 1                  //章节号结束
				result := countSplit[stk:edk] //截取卷前10以内的字符
				title := countSplit[edk:]
				if i := getNumber(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
					return i, title
				} else if i := getSimplified(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
					return int(i), title
				} else if i := getTraditional(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
					return int(i), title
				}
			}
		}
	}
	return 0, nil
}
