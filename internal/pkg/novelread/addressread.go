package novelread

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"

	"github.com/sirupsen/logrus"
)

var (
	listFilePrefix = "  "
)

//NovelRead 小说读取
func NovelRead() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 当我们取完需要的整数后调用cancel
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
		// 小说读取
		for v := range SplitChapter(ctx, file) {
			if v.Err != nil {
				fmt.Println(v.Err)
				continue
			}
			fmt.Println(v.Chapter.Volume, v.Chapter.Index, v.Chapter.Titles)
		}
		cancel()
		useTime := time.Since(st)
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
