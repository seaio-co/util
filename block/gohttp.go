package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var baseUrl = "https://www.qidian.com/rank/yuepiao?page=";
var ch = make(chan int)

func main() {
	//
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go toSpider(i, ch)
	}
	for i := 0; i < 5; i++ {
		fmt.Println("finish", <-ch, "finish read")
	}
}

func toSpider(page int, ch chan int) {
	fmt.Println("begin read" + strconv.Itoa(page) + "page")
	url := baseUrl
	result := httpGet(url)
	toSave(result)

	ch <- page
}
func toSave(result string){

	ret := regexp.MustCompile(`<h4>(?s:(.*?))</h4>`)

	alls :=ret.FindAllStringSubmatch(result,-1)
	for _, tmpTitle := range alls {
		title := tmpTitle[1]
		title = strings.Replace(title, "\t", "", -1)
		fmt.Println(title)
		break
	}
}

func httpGet(url string) (result string) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	doError("connect error", err)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		doError("read error", err)
		//n=0
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return result
}


func doError(name string, err error) {
	if (err != nil && err != io.EOF) {
		fmt.Println(name, ":", err)
	}
}
