package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"strings"
	"github.com/emacsist/alfred3/utils"
	"os"
	"bufio"
)

func main() {
	fundCode := utils.GetQuery()
	//query := "dsp   "
	fundCodes := strings.Fields(fundCode)

	alfredResponse := utils.NewAlfredResponse()
	defer alfredResponse.WriteOutput()

	if len(fundCodes) == 0 {
		//如果没有，则获取从文件中获取默认的基金代码
		fundCodes = getDefaultFundCodes()
	}

	for _, code := range fundCodes {
		doc, err := goquery.NewDocument(getUrl(code))
		if err != nil {
			alfredResponse.AddItemWithSubTitle("获取基金信息出错", code)
			continue
		}
		// Find the review items
		doc.Find(".merchandiseDetail").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			fundTitle := s.Find(".fundDetail-header > .fundDetail-tit > div").Text() + ")"
			// 净值
			netValue := s.Find(".fundDetail-main").Find("#gz_gsz").Text()

			floatValue := s.Find(".fundDetail-main").Find("#gz_gszze").Text()

			// 变动率
			floatRate := s.Find(".fundDetail-main").Find("#gz_gszzl").Text()

			alfredResponse.AddItemWithSutTitleAndArg(fundTitle, netValue + " | " + floatValue + " | " + floatRate, getUrl(code))
		})
	}
}

// fundCode : 基金代码
func getUrl(fundCode string) string {
	return fmt.Sprintf("http://fund.eastmoney.com/%v.html", fundCode)
}

func getDefaultFundCodes() []string {
	var data []string
	file, err := os.Open("./fund.txt")
	if err != nil {
		log.Fatal(err)
		return data
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 非注释行，则为基金代码
		if !strings.HasPrefix(line, "#") {
			data = append(data, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}