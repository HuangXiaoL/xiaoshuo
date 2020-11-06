package novelread

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

type Expect struct {
	Title *Chapter
	OK    bool
	Err   error
}

var (
	testCases = []struct {
		Line   io.Reader
		Expect Expect
	}{

		{
			Line: strings.NewReader("十一"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "",
					Index:  11,
				},
				OK:  true,
				Err: nil,
			},
		},
		{
			Line: strings.NewReader("十 一"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "一",
					Index:  10,
				},
				OK:  true,
				Err: nil,
			},
		},

		{
			Line: strings.NewReader("第二卷第五十章 包含了卷的标题"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "包含了卷的标题",
					Index:  50,
					Volume: 2,
				},
				OK:  true,
				Err: nil,
			},
		},
		{
			Line: strings.NewReader("第二卷 第五十章 包含了卷的标题"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "包含了卷的标题",
					Index:  50,
					Volume: 2,
				},
				OK:  true,
				Err: nil,
			},
		},
		{
			Line: strings.NewReader("第四十八-四十九章"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "",
					Index:  49,
				},
				OK:  true,
				Err: nil,
			},
		},
		{
			Line: strings.NewReader("第十卷 单独的卷标题"),
			Expect: Expect{
				Title: &Chapter{
					Titles: "",
					Index:  0,
					Volume: 10,
				},
				OK:  true,
				Err: nil,
			},
		},
		// 异常测试
		{
			Line: strings.NewReader(" 第十卷 用空格缩进的文字"),
			Expect: Expect{
				Title: nil,
				OK:    false,
				Err:   nil,
			},
		},
		{
			Line: strings.NewReader(" 第十卷 用tab缩进的文字"),
			Expect: Expect{
				Title: nil,
				OK:    false,
				Err:   nil,
			},
		},
		{
			Line: strings.NewReader("测试前面的文字非常长的情况 第一章 三卷天书"),
			Expect: Expect{
				Title: nil,
				OK:    false,
				Err:   nil,
			},
		},
	}
)

func TestSplitChapter(t *testing.T) {
	for _, v := range testCases {
		got, _ := SplitChapter(v.Line)
		for _, vv := range got {
			fmt.Println(vv.Volume)
			fmt.Println(v.Expect.Title.Volume)
			fmt.Println(vv.Index)
			fmt.Println(v.Expect.Title.Index)
			fmt.Println(vv.Titles)
			fmt.Println(v.Expect.Title.Titles)
			if ok := reflect.DeepEqual(vv.Volume, v.Expect.Title.Volume); !ok {
				t.Fatalf("期望得到Volume%v，实际得到Volume%v", v.Expect.Title.Volume, vv.Volume)
			}

			if ok := reflect.DeepEqual(vv.Index, v.Expect.Title.Index); !ok {
				t.Fatalf("期望得到Index%v，实际得到Index%v", v.Expect.Title.Index, vv.Index)
			}
			if ok := reflect.DeepEqual(vv.Titles, v.Expect.Title.Titles); !ok {
				t.Fatalf("期望得到Titles%s，实际得到Titles%s", v.Expect.Title.Titles, vv.Titles)
			}
		}
	}
}
