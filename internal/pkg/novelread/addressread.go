package novelread

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"

	"github.com/sirupsen/logrus"
)

var (
	listFilePrefix = "  "
	//wg             sync.WaitGroup
)

//NovelRead 小说读取
func NovelRead() {
	src := config.Get().FileAddress.Address
	//src := "/www/xiaoshuo/Theoriginalnovel/"
	srcDir := src
	pathSeparator := string(os.PathSeparator)
	level := 1
	fileName := listAllFileByName(level, pathSeparator, srcDir)
	st := time.Now()

	for _, v := range fileName {
		fileAddres := src + v
		//_, err := os.Open(fileAddres)
		//fmt.Println(fileAddres)
		file, err := os.Open(fileAddres)
		if err != nil {
			panic(err)
		}
		c, err := SplitChapter(file)
		//_, _ = SplitChapter(file)
		useTime := time.Since(st)
		if err != nil {
			logrus.Println(err)
		}
		for _, v := range c {
			fmt.Println(v.Volume, v.Index, v.Titles)
		}

		logrus.Printf("用时为：%s", useTime)
	}
	useAllTime := time.Since(st)
	logrus.Printf("用时为：%s", useAllTime)
}

// listAllFileByName 文件列表
func listAllFileByName(level int, pathSeparator, fileDir string) map[int]string {
	var (
		num      = 1                    //计数器
		fileName = make(map[int]string) //文件名称

	)
	files, _ := ioutil.ReadDir(fileDir)
	tmpPrefix := ""
	for i := 1; i < level; i++ {
		tmpPrefix = tmpPrefix + listFilePrefix
	}
	for _, o := range files {
		if o.IsDir() {
			fmt.Printf("\033[34m %s %s \033[0m \n", tmpPrefix, o.Name())
			listAllFileByName(level+1, pathSeparator, fileDir+pathSeparator+o.Name())
		} else {
			fileName[num] = tmpPrefix + o.Name()
		}
		num++
	}

	return fileName
}
