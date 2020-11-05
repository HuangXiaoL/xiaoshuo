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
	//屏蔽词
	shield = map[int]string{0: "京", 1: "两"}
	//章节 章节集回话 ----- 章节识别关键字
	chapter = map[int]string{0: "章", 1: "节", 2: "集", 3: "回", 4: "话", 5: " ", 6: "-"}
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
func SplitChapter(input io.Reader) ([]Chapter, error) {
	var (
		conts     string
		o         = []Chapter{}
		c         = Chapter{}
		volumeNum int
	)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		cont, vnum, cnum, t := lineTextDiscern(scanner.Text())

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
			o = append(o, c)
		}
	}
	return o, nil
}

//lineTextDiscern 行文本识别
func lineTextDiscern(line string) (string, int, int, string) {
	length := lineLength(line)
	if length { //长度是否超过80 有就返回为正文
		return strings.TrimSpace(line), 0, 0, ""
	}
	b := lineRetractIsContent(line)
	//赋值逻辑
	if b { //是否有缩进 有就返回为正文
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
	return utf8.RuneCountInString(line) > 80
}

//lineIsPureNumber 判断这一行是否为纯数字
func lineIsPureNumber(s []string) bool {
	line := ""
I:
	for _, v := range s {
		for _, vv := range shield {
			if v == vv {
				break I
			}

		}
		line = line + v
	}
	line = strings.TrimSpace(line)

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

//lineFindNumAtChapterAndVolume 在行里查找数字并且返回 章 卷 的值和 章节名称
func lineFindNumAtChapterAndVolume(line string, seat int) (int, int, string) {
	var (
		volumeNum       int      // 卷 值
		SectionPosition int      //卷 在这一行的位置
		chapterNum      int      // 章值
		title           []string //标题值
	)
	countSplit := strings.Split(line, "") //切割字符串
	if len(countSplit) == 0 {             // 行的长度为0 直接返回
		return 0, 0, ""
	}
	s := countSplit
	if len(countSplit) >= seat { //如果行的长度大于设定值，就截取设定的长度
		s = countSplit[0:seat]
	}
	if lineIsPureNumber(s) {
		return 0, getStringNumber(s), ""
	}
	// 获取卷值
	if i := getStringNumber(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		volumeNum, SectionPosition = getVolumeNum(line) //卷值
	}
	// 获取章值
	if i := getStringNumber(s); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		chapterNum, title = getChapterNum(line, SectionPosition) //章值
	}
	//标题
	t := ""
	for _, v := range title {
		t = t + v
	}

	return volumeNum, chapterNum, t
}

//getVolumeNum 卷号识别并且提取 卷节号 及其 卷位置
func getVolumeNum(s string) (int, int) {
	countSplit := strings.Split(s, "")
	for k, v := range countSplit {
		if strings.Contains(v, reel) {

			result := countSplit[:k] //截取卷前的字符
			//判断卷前面是否是 数字预防title里有卷， 标题有数字和卷的组合 ，返回卷为0
			if getStringNumber(result[k-1:k]) == 0 {
				k = 0
			}
			for kk, v := range result {
				for _, vv := range chapter {
					if strings.Contains(v, vv) {
						return 0, kk
					}
				}
			}
			if i := getNumber(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return i, k
			} else if i := getSimplified(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return int(i), k
			} else if i := getTraditional(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
				return int(i), k
			}
		}
	}
	return 0, 0
}

//getChapterNum 识别章数并提取 章节号 和章节名称
func getChapterNum(s string, SectionPosition int) (int, []string) {

	countSplit := strings.Split(s, "")
	for k, v := range countSplit {

		for _, vv := range chapter {
			if strings.Contains(v, vv) { //是否 有章,节,集，回等.... 做判断
				sk := k - 1
				if k > 1 {
					sk = k - 2
				}
				if getStringNumber(countSplit[sk:k]) == 0 {
					return 0, nil
				}

				if SectionPosition > k {
					SectionPosition = k - k
				}

				result := countSplit[SectionPosition:k] //截取卷后到章之前的字符

				title := countSplit[k+1:]
				if i := getStringNumber(result); i > 0 { //返回的数字大于0 有可能是章节目录有卷
					return i, title
				}
			}
		}
	}
	return 0, nil
}

//getStringNumber 获取字符串中的数字
func getStringNumber(line []string) int {

	// 获取卷值
	if i := getNumber(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return i
	} else if i := getSimplified(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return int(i)
	} else if i := getTraditional(line); i > 0 { //返回的数字大于0 有可能是章节目录有卷
		return int(i)
	}
	return 0
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
				ss = ss + strings.TrimSpace(v)
			}
		}
		for _, vv := range shield {
			if v == vv {
				return 0
			}
		}
	}
	if utf8.RuneCountInString(ss) != 0 { //每行的数字
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
				ss = ss + strings.TrimSpace(v)
			}
		}
	}

	if utf8.RuneCountInString(ss) != 0 { //每行的数字
		num, _ := numcn.DecodeToInt64(ss)
		return num
	}
	return 0
}
