package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hyunki08/newsscrapper/scrapper"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("검색어를 입력해주세요.")
		return
	}
	query := args[0]
	fmt.Println(query + "를 검색합니다.")

	maxPage := 40
	if len(args) == 1 {
		fmt.Println("최대 40페이지까지 결과를 csv 파일에 출력합니다.")
	} else if len(args) == 2 {
		conv, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("최대 40페이지까지 결과를 csv 파일에 출력합니다.")
		} else {
			maxPage = conv
			fmt.Println(strconv.Itoa(maxPage) + "페이지까지 결과를 csv 파일에 출력합니다.")
		}
	}

	scrapper.Scrapper(query, maxPage)
}
