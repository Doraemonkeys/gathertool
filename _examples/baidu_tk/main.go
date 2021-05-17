/*
	百度题库抓取
 */

package main

import (
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main(){
	getZJ()
}

// １．获取每一章节对应的目录
func getZJ(){
	caseUrl := "https://tiku.baidu.com/tikupc/chapterlist/1bfd700abb68a98271fefa04-16-knowpoint"
	c, _ := gt.Get(caseUrl)
	c.Do()
	//log.Println(string(c.RespBody))

	dom,err := gt.NewGoquery(string(c.RespBody))
	if err != nil{
		log.Println(err)
		return
	}

	dom.Find("div[class=detail]").Each(func(i int, div *goquery.Selection){
		chapter := div.Find("div[class=detail-chapter]")
		//log.Println(chapter.Html())
		chapterHtml, _ := chapter.Html()

		// 章节名称
		zjTitle := gt.RegHtmlH(chapterHtml, "3")
		log.Println("章节名称 ==> ", zjTitle)

		chapter.Each(func(i int, div2 *goquery.Selection){
			kpoint := div2.Find("div[class=detail-kpoint-1]")
			//log.Println(kpoint.Html())
			kpointHtml, _ := kpoint.Html()

			// 小节名称
			xjTitle := gt.RegHtmlH(kpointHtml, "4")
			log.Println("小节名称 ==> ", xjTitle)

			// 获取课程
			kpoint.Each(func(i int, div3 *goquery.Selection){
				kpoint2 := div3.Find("div[class=detail-kpoint-2]")
				//log.Println(kpoint2.Html())
				kpoint2Html,_ := kpoint2.Html()

				//课程名称
				kcTitle := gt.RegHtmlHTxt(kpoint2Html, "5")
				log.Println("课程名称 ==> ", kcTitle)
				//课程链接
				kcLink := gt.RegHtmlHrefTxt(kpoint2Html)
				if len(kcLink) > 0 {
					link := "https://tiku.baidu.com"+kcLink[0]
					log.Println("课程链接 ==> ", link)
				}
				log.Println("目录 ==> ", zjTitle, xjTitle, kcTitle)

				log.Println("-------------------------------")

			})
		})

	})
}