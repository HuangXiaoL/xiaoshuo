package novelread

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/HuangXiaoL/xiaoshuo/internal/pkg/config"
)

var (
	listFilePrefix = "  "
	wg             sync.WaitGroup
)

//NovelRead 小说读取
func NovelRead() {
	src := config.Get().FileAddress.Address
	srcDir := src
	pathSeparator := string(os.PathSeparator)
	level := 1
	fileName := listAllFileByName(level, pathSeparator, srcDir)
	for _, v := range fileName {
		fileAddres := src + v
		fmt.Println(fileAddres)
		wg.Add(1)
		go TrimFile(fileAddres)
		wg.Wait()
	}

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
