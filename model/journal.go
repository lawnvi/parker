package model

type Journal struct {
	BaseModel
	Name          string //刊名
	ShortName     string //简称
	Issn          string
	EIssn         string
	ImpactFactor  float32 //IF
	SelfCiting    float32 //自引率
	HIndex        int     //h指数
	CiteScore     float32 //引用分数
	SJR           float32
	SNIP          float32
	Website       string //刊网站
	PostAddress   string //刊投递地址
	IsOA          bool
	Communication string //通讯方式
	Publisher     string //出版商
	Specialism    string //研究方向
	Area          string //地区
	PublishCycle  string //出版周期
	PublishYear   int    //出版年份
	NumberPerYear int    //年出版数
	Warning       bool   //学业预警
	PubMedCentral string //PMC链接
	Speed         string //平均审稿速度
	Difficulty    string //录用难度

	//基础SCI
	BasicSCIMain      string
	//BasicSCISubject   string
	BasicSCIMainLevel int
	//BasicSCISubLevel  int
	BasicSCIIsTop     bool
	BasicSCIIsReview  bool //综述期刊

	//试行SCI
	TestSCIMain      string
	//TestSCISubject   string
	TestSCIMainLevel int
	//TestSCISubLevel  int
	TestSCIIsTop     bool
	TestSCIIsReview  bool //综述期刊
}
