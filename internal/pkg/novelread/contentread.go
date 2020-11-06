package novelread

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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
	simplified = map[int]string{0: "零", 1: "一", 2: "二", 3: "三", 4: "四", 5: "五", 6: "六", 7: "七", 8: "八", 9: "九", 10: "十", 11: "百", 12: "千", 13: "万"}
	//繁体
	traditional = map[int]string{0: "零", 1: "壹", 2: "贰", 3: "叁", 4: "肆", 5: "伍", 6: "陆", 7: "柒", 8: "捌", 9: "玖", 10: "拾", 11: "佰", 12: "仟", 13: "萬"}
	//章节 章节集回话 ----- 章节识别关键字
	chapter = map[int]string{0: "章", 1: "节", 2: "集", 3: "回", 4: "话", 5: " "}
	//卷
	reel = "卷"
)

//Chapter 小说章节内容的结构
type Chapter struct {
	Volume  int
	Index   int
	Titles  string
	Content string
}

//SplitChapter 文本流入口
func SplitChapter(input io.Reader) {
	var (
		conts     string
		volumeNum int
		c         = Chapter{}
	)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		cont, vnum, cnum, t := lineTextDiscern(scanner.Text())
		//_, vnum, cnum, t := lineTextDiscern(scanner.Text())
		//lineTextDiscern(scanner.Text())

		//卷号处理，没抓取到就赋值
		if vnum == 0 {
			vnum = volumeNum
		}
		volumeNum = vnum
		//else if vnum == volumeNum || vnum == volumeNum+2 || vnum == volumeNum+1 { //抓取到了判断值是否合理，是否是同卷，或者是下一卷或者第一卷的卷号没写
		//	volumeNum = vnum
		//}
		conts = conts + cont
		if cnum != 0 {
			c.Titles = strings.TrimSpace(strings.Trim(t, "\n\r"))
			c.Volume = volumeNum
			c.Index = cnum
			c.Content = conts
			conts = ""
			//fmt.Println(c.Volume, c.Index, c.Titles)
			ch <- c

		}
	}
	close(ch)
}

//lineTextDiscern 行文本识别 （cont 正文，volume 卷号，index 章节，title 章节标题）
func lineTextDiscern(line string) (cont string, volume int, index int, title string) {
	length := lineLength(line)
	if length { //长度是否超过80 有就返回为正文
		return strings.TrimSpace(line), 0, 0, ""
	}
	b := lineRetractIsContent(line)
	if b { //是否有缩进 有就返回为正文
		return strings.TrimSpace(line), 0, 0, ""
	}
	if !lineGreaterThanSetValueNoNumber(line, 5) { //识别前五位是否有数字，有数字为卷章节，否则为正文
		return strings.TrimSpace(line), 0, 0, ""
	}
	// 提取 卷号 ，章节号，章节名称
	volume, index, title = lineFindNumAtChapterAndVolume(line)
	return
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
	return len(line) > 240
}

//lineGreaterThanSetValueNoNumber 行前五位是否有数字
func lineGreaterThanSetValueNoNumber(line string, set int) bool {
	countSplit := strings.Split(line, "") //切割字符串
	s := countSplit
	if len(countSplit) >= set { //如果行的长度大于设定值，就截取设定的长度
		s = countSplit[0:set]
	}
	for _, v := range s {
		if getTheStringIsNumber(v) != "" { //有预设的识别数字返回为 可能为章节卷
			return true
		}
	}
	return false
}

//lineFindNumAtChapterAndVolume 在行里查找数字并且返回 章 卷 的值和 章节名称
func lineFindNumAtChapterAndVolume(line string) (int, int, string) {
	var (
		volumeNum       int    // 卷 值
		SectionPosition int    //卷 在这一行的位置
		chapterNum      int    // 章值
		title           string //标题值
		isVolume        = true
	)
	countSplit := strings.Split(line, "") //切割字符串
	if len(countSplit) == 0 {             // 行的长度为0 直接返回
		return 0, 0, ""
	}
	if lineIsPureNumber(countSplit) {
		return 0, getStringNumber(line), ""
	}
	volumeNum, SectionPosition = getVolumeNum(countSplit)                    //卷值
	chapterNum, title, isVolume = getChapterNum(countSplit, SectionPosition) //章值
	if !isVolume {                                                           // 当识别到的卷在章节之后的时候，就不可使用该卷值，设置为0
		volumeNum = 0
	}
	return volumeNum, chapterNum, title
}

//lineIsPureNumber 判断这一行是否为纯数字
func lineIsPureNumber(s []string) bool {
	line := ""
	line = strings.TrimSpace(line)
	for _, v := range s {
		line = line + v
	}
	if line != "" {
		i, _ := numcn.DecodeToInt64(line)
		if int(i) != 0 {
			return true
		} else {
			num, _ := strconv.Atoi(line)
			if num != 0 {
				return true
			}
		}
	}
	return false
}

//getVolumeNum 卷号识别并且提取 卷节号 及其 卷位置
func getVolumeNum(countSplit []string) (int, int) {
	volumeNum := ""
	for k, v := range countSplit {
		if getTheStringIsNumber(v) != "" { // 判断是否是匹配得上数字匹配上了就相加组合在一起
			volumeNum = volumeNum + v
		} else if strings.Contains(v, reel) { //未匹配上数字，就匹配是否是 “卷”
			i := getStringNumber(volumeNum) //卷号
			return i, k
		} else { //以上条件都不满足的，收集到的数字 释放掉
			volumeNum = ""
		}

	}
	return 0, 0
}

//getChapterNum 识别章数并提取 章节号 和章节名称 chapter章节号 titles标题 volume 之前识别的卷号是否可用
func getChapterNum(countSplit []string, SectionPosition int) (chapters int, titles string, volume bool) {
	chapterNum := ""
	for k, v := range countSplit {
		if getTheStringIsNumber(v) != "" {
			chapterNum = chapterNum + v
		} else {
			for _, vv := range chapter {
				if v == vv {
					if chapters := getStringNumber(chapterNum); chapters > 0 {
						for _, v := range countSplit[k+1:] {
							titles = titles + v
						}
						volume = true
						if SectionPosition > k { // 当卷的位置大于章节的位置，认为是标题有卷名称，无效的卷
							volume = false
						}
						return chapters, titles, volume
					}
				}
			}
			chapterNum = ""
		}
	}

	return
}

//getTheStringIsNumber 字符串是否是预设识别需要的数字 预设值 number simplified traditional
func getTheStringIsNumber(s string) string {
	for _, v := range simplified {
		if s == v {
			return s
		}
	}
	for _, v := range number {
		if s == v {
			return s
		}

	}
	for _, v := range traditional {
		if s == v {
			return s
		}

	}
	return ""
}

//getStringNumber 获取字符串中的数字
func getStringNumber(line string) int {
	// 获取卷值
	if i := getNumber(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return i
	} else if i := getSimplified(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return i
	} else if i := getTraditional(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return i
	}
	return 0
}

//getNumber 获取每行的阿拉伯数字
func getNumber(line string) int {
	if len(line) != 0 { //每行的数字
		i, _ := strconv.Atoi(line)
		return i
	}
	return 0
}

//getSimplified 获取简写的数字的值
func getSimplified(line string) int {
	if utf8.RuneCountInString(line) != 0 { //每行的数字
		num, _ := numcn.DecodeToInt64(line)
		return int(num)
	}
	return 0
}

//getTraditional 获取繁体的数字的值
func getTraditional(line string) int {
	if utf8.RuneCountInString(line) != 0 { //每行的数字
		num, _ := numcn.DecodeToInt64(line)
		return int(num)
	}
	return 0
}
