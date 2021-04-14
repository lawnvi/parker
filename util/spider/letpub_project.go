package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"parker/config/log"
	"parker/model"
	"strconv"
	"strings"
	"time"
)

var tag = "letpub_proj"

//可查找项目与期刊信息
//datakind=excel excel方式下载
type LetPubProject struct {
	domain string //主网址
	HttpSpider
}

func NewLetPubProject() ProjectSpider {
	var l = LetPubProject{
		domain: "https://www.letpub.com.cn/nsfcfund_search.php?mode=advanced&datakind=list&currentpage=1",
	}
	l.InitSpider()
	return l
}

func (l LetPubProject)visit(url string) (projects []model.Project) {
	l.Colly.OnHTML("#dict", func(element *colly.HTMLElement) {
		//#dict > center > div > b
		var size = element.ChildText("center > div > b")
		fmt.Printf("find projects size: %s\n", size)
		projects = l.toProjects(*element)
	})
	l.Colly.OnError(func(response *colly.Response, e error) {
		fmt.Println(e)
	})
	err := l.Colly.Visit(url)
	if err != nil {
		log.E(tag, fmt.Sprintf("letpub visit project error: %v", err))
	}
	return
}

func (l LetPubProject) SearchByNo(no string) []model.Project{
	var params = "page=&name=&person=&company=&addcomment_s1=&addcomment_s2=&addcomment_s3=&addcomment_s4=&money1=&money2=&startTime=1997&province_main=&subcategory=&searchsubmit=true"
	var url = fmt.Sprintf("%s&no=%s&endTime=%d&%s", l.domain, no, time.Now().Year() - 2, params)
	return l.visit(url)
}

func (l LetPubProject) SearchByName(name string) []model.Project{
	var params = "page=&no=&person=&company=&addcomment_s1=&addcomment_s2=&addcomment_s3=&addcomment_s4=&money1=&money2=&startTime=1997&province_main=&subcategory=&searchsubmit=true"
	var url = fmt.Sprintf("%s&name=%s&endTime=%d&%s", l.domain, name, time.Now().Year() - 2, params)
	return l.visit(url)
}

func (l LetPubProject) SearchByResearcher(researcher model.Researcher) []model.Project{
	var params = "page=&no=&name=&addcomment_s1=&addcomment_s2=&addcomment_s3=&addcomment_s4=&money1=&money2=&startTime=1997&province_main=&subcategory=&searchsubmit=true"
	var url = fmt.Sprintf("%s&person=%s&company=%s&endTime=%d&%s", l.domain, researcher.Name, researcher.Company, time.Now().Year() - 2, params)
	return l.visit(url)
}

func (l LetPubProject) toProjects(element colly.HTMLElement) []model.Project{
	var projects []model.Project
	//if node size > 2, do next
	var le = len(element.DOM.Children().Nodes)
	println(le)
	var p *model.Project
	element.ForEach("table > tbody > tr", func(i int, element *colly.HTMLElement) {
		if i > 1 {
			if len(element.DOM.Children().Nodes) == 7{
				projects = append(projects, model.Project{})
				p = &projects[len(projects)-1]
				p.Name = element.ChildText("td:nth-child(1)")
				p.Company = element.ChildText("td:nth-child(2)")
				p.Fund, _ = strconv.Atoi(element.ChildText("td:nth-child(3)"))
				p.No = element.ChildText("td:nth-child(4)")
				p.Belong = element.ChildText("td:nth-child(5)")
				p.Subcategory = element.ChildText("td:nth-child(6)")
				p.ApprovalYear, _ = strconv.Atoi(element.ChildText("td:nth-child(7)"))
			}
			var attr = element.ChildText("td:nth-child(1)")
			var value = element.ChildText("td:nth-child(2)")
			switch attr {
			case "题目": p.Name = value
			case "学科分类": {
				var list = strings.Split(value, "，")
				p.CommentS1 = strings.Split(list[0], "级：")[1]
				p.CommentS2 = strings.Split(list[1], "级：")[1]
				p.CommentS3 = strings.Split(list[2], "级：")[1]
			}
			case "学科代码": {
				var list = strings.Split(value, "，")
				p.CommentCode1 = strings.Split(list[0], "级：")[1]
				p.CommentCode2 = strings.Split(list[1], "级：")[1]
				p.CommentCode3 = strings.Split(list[2], "级：")[1]
			}
			case "执行时间":{
				var list = strings.Split(value, " 至 ")
				p.StartTime = list[0]
				p.EndTime = list[1]
			}
			case "中文关键词": p.KeywordCN = value
			case "英文关键词": p.KeywordEN = value
			case "结题摘要": p.Abstract = value
			}
		}
	})
	log.I(tag, fmt.Sprintf("find projects size: %d\n%v", len(projects), projects[0]))
	return projects
}