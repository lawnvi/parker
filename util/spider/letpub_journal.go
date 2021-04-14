package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"parker/config/log"
	"parker/model"
	"strconv"
	"strings"
)

type LetPubJournal struct {
	domain string //主网址
	HttpSpider
}

type journalBuf struct {
	Name      string
	ISSN      string
	ShortName string
	DetailUrl string
}

func NewLetPubJournal() JournalSpider {
	var l = LetPubJournal{
		domain: "https://www.letpub.com.cn/index.php?page=journalapp",
	}
	l.InitSpider()
	return l
}

//https://www.letpub.com.cn/index.php?page=journalapp&view=search
//searchname=IEEE&searchissn=&searchfield=&searchimpactlow=&searchimpacthigh=&searchscitype=&view=search&searchcategory1=&searchcategory2=&searchjcrkind=&searchopenaccess=&searchsort=relevance
//https://www.letpub.com.cn/index.php?journalid=3342&page=journalapp&view=detail
func (l LetPubJournal) visitDetail(url string) (journal model.Journal) {
	l.Colly.OnHTML("#yxyz_content > table:nth-child(12) > tbody", func(element *colly.HTMLElement) {
		journal = l.toJournal(*element)
	})
	l.Colly.OnError(func(response *colly.Response, e error) {
		fmt.Println(e)
	})
	err := l.Colly.Visit(url)
	if err != nil {
		log.E(tag, fmt.Sprintf("letpub visit journal error: %v", err))
	}
	return
}

func (l LetPubJournal) doSearch(url string) (journals []journalBuf) {
	l.Colly.OnHTML("#yxyz_content > table.table_yjfx > tbody", func(element *colly.HTMLElement) {
		journals = l.toJournalsBuf(*element)
	})
	l.Colly.OnError(func(response *colly.Response, e error) {
		fmt.Println(e)
	})
	err := l.Colly.Visit(url)
	if err != nil {
		log.E(tag, fmt.Sprintf("letpub visit journal error: %v", err))
	}
	return
}

func (l LetPubJournal) SearchByISSN(issn string) model.Journal {
	var params = "&searchname=&searchfield=&searchimpactlow=&searchimpacthigh=&searchscitype=&view=search&searchcategory1=&searchcategory2=&searchjcrkind=&searchopenaccess=&searchsort=relevance"
	var url = fmt.Sprintf("%s&searchissn=%s%s", l.domain, issn, params)
	var journals = l.doSearch(url)
	if len(journals) < 1 || journals[0].ISSN != issn{
		log.I(tag, "journals size is < 1, not find journal by issn: "+ issn)
		return model.Journal{}
	}
	return l.visitDetail("https://www.letpub.com.cn"+ journals[0].DetailUrl[1:len(journals[0].DetailUrl)])
}

func (l LetPubJournal) SearchByName(name string) []model.Journal {
	var params = "&searchissn=&searchfield=&searchimpactlow=&searchimpacthigh=&searchscitype=&view=search&searchcategory1=&searchcategory2=&searchjcrkind=&searchopenaccess=&searchsort=relevance"
	var url = fmt.Sprintf("%s&searchname=%s%s", l.domain, name, params)
	var bufs = l.doSearch(url)
	var journals []model.Journal
	for b := range bufs{
		journals = append(journals, bufs[b].toJournal())
	}
	return journals
}

func (l LetPubJournal) toJournal(html colly.HTMLElement) model.Journal {
	var journal model.Journal
	html.ForEach("#yxyz_content > table:nth-child(12) > tbody > tr", func(i int, element *colly.HTMLElement) {
		var title = element.ChildText("#yxyz_content > table:nth-child(12) > tbody > tr > td:nth-child(1)")
		var value = element.ChildText("#yxyz_content > table:nth-child(12) > tbody > tr > td:nth-child(2)")
		switch {
		case title == "期刊名字":
			{
				journal.ShortName = element.ChildText("td:nth-child(2) > span > font")
				journal.Name = element.ChildText("td:nth-child(2) > span > a")
			}
		case title == "期刊ISSN": journal.Issn = value
		case strings.Contains(title, "自引率"):
			{
				if f, err := strconv.ParseFloat(value[0 : strings.Count(value, "") - 12], 32); err == nil{
					journal.SelfCiting = float32(f)
				}
				fmt.Println(journal.SelfCiting)
			}
		case title == "h-index": journal.HIndex, _ = strconv.Atoi(value)
		case title == "CiteScore":{
			cite := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(1)")
			sjr := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(2)")
			snip := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(3)")
			if f, e := strconv.ParseFloat(cite, 32); e == nil{
				journal.CiteScore = float32(f)
			}
			if f, e := strconv.ParseFloat(sjr, 32); e == nil{
				journal.SJR = float32(f)
			}
			if f, e := strconv.ParseFloat(snip, 32); e == nil{
				journal.SNIP = float32(f)
			}
		}
		case title == "期刊官方网站": journal.Website = value
		case title == "期刊投稿网址": journal.PostAddress = value
		case title == "是否OA开放访问": journal.IsOA = value != "No"
		case title == "通讯方式": journal.Communication = value
		case title == "出版商": journal.Publisher = value
		case title == "涉及的研究方向": journal.Specialism = value
		case title == "出版国家或地区": journal.Area = value
		case title == "出版周期": journal.PublishCycle = value
		case title == "出版年份": journal.PublishYear, _ = strconv.Atoi(value)
		case title == "年文章数":journal.NumberPerYear, _ = strconv.Atoi(value[0 : strings.Count(value, "") - 12])
		case title == "PubMed Central (PMC)链接": journal.PubMedCentral = value
		case title == "平均审稿速度": journal.Speed = value
		case title == "平均录用比例": journal.Difficulty = value
		case strings.Contains(title, "预警名单"): journal.Warning = !strings.Contains(value, "不在")
		case strings.Contains(title, "最新基础版"):{
			//todo display to check which is right level
			level := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(1) > span:nth-child(2)")
			isTop := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(3)")
			isReview := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(4)")
			journal.BasicSCIMainLevel, _ = strconv.Atoi(level[0:1])
			journal.BasicSCIIsTop = isTop == "是"
			journal.BasicSCIIsReview = isReview == "是"
		}
		case strings.Contains(title, "最新升级版"):{
			//todo display to check which is right level
			level := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(1) > span:nth-child(2)")
			isTop := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(3)")
			isReview := element.ChildText("td:nth-child(2) > table > tbody > tr:nth-child(2) > td:nth-child(4)")
			journal.TestSCIMainLevel, _ = strconv.Atoi(level[0:1])
			journal.TestSCIIsTop = isTop == "是"
			journal.TestSCIIsReview = isReview == "是"
		}
		}
	})
	//todo get Impact Factor, the old url does not work now
	return journal
}



func (l LetPubJournal) toJournalsBuf(html colly.HTMLElement) []journalBuf {
	var list []journalBuf
	html.ForEach("#yxyz_content > table.table_yjfx > tbody > tr", func(i int, element *colly.HTMLElement) {
		if i > 1 &&  len(element.DOM.Children().Nodes) == 12{
			var issn = element.ChildText("td:nth-child(1)")
			var name = element.ChildText("td:nth-child(2) > a")
			var short = element.ChildText("td:nth-child(2) > font")
			var url = element.ChildAttr("td:nth-child(2) > a", "href")
			list = append(list, journalBuf{
				Name:      name,
				ISSN:      issn,
				ShortName: short,
				DetailUrl: url,
			})
		}
	})
	return list
}

func (b journalBuf) toJournal() model.Journal {
	return model.Journal{
		Name: b.Name,
		Issn: b.ISSN,
		ShortName: b.ShortName,
	}
}