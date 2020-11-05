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
