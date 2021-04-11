package spider

import (
	"github.com/gocolly/colly"
	"parker/model"
)

//项目爬虫接口
type ProjectSpider interface {
	//搜索
	SearchByNo(no string) []model.Project
	SearchByName(name string) []model.Project
	//researcher name+company
	SearchByResearcher(researcher model.Researcher) []model.Project
	//解析成对象
	toProjects(html colly.HTMLElement) []model.Project
}

//期刊爬虫接口
type JournalSpider interface {

}

//论文爬虫接口
type ArticleSpider interface {

}