package novelread

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"

	"github.com/sirupsen/logrus"
)

var (
	listFilePrefix = "  "
)

//NovelRead 小说读取
func NovelRead() {
	wgr := sync.WaitGroup{}
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
		c, err := SplitChapter(file)
		if err != nil { //读取错误
			logrus.Println(err)
		}
		wgr.Add(1)
		go getBookCatalogue(c, &wgr)

		wgr.Wait()
		useTime := time.Since(st)
		logrus.Printf("用时为：%s", useTime)
	}
	useAllTime := time.Since(st)
	logrus.Printf("用时为：%s", useAllTime)
}
func getBookCatalogue(c chan Chapter, wgr *sync.WaitGroup) {
	defer wgr.Done()
	for v := range c {
		fmt.Println(v.Volume, v.Index, v.Titles)
		//select {
		//case v := <-c:
		//	fmt.Println(v.Volume, v.Index, v.Titles)
		//}
	}
	return

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
