package scrapper

import (
	"fmt"
	"log"
	"net/http"
)

func getHttp(url string) *http.Response {
	fmt.Println("[GET]Request : ", url)
	req, rErr := http.NewRequest("GET", url, nil)
	checkErr(rErr)
	req.Header.Add("X-Naver-Client-Id", "83JLs9mo_QftyJIcMY8C")
	req.Header.Add("X-Naver-Client-Secret", "XqSZPxw7qV")

	client := &http.Client{}
	res, err := client.Do(req)
	checkErr(err)
	checkCode(res)

	return res
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status : ", res.StatusCode)
	}
}
