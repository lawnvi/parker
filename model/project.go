package model

/**
find in letpub
*/
type Project struct {
	BaseModel
	Leader       string //负责人姓名
	Company      string //单位
	Fund         int    //基金
	Name         string //项目名称
	No           string //项目编号
	Subcategory  string //项目类型
	Belong       string //所属学部
	ApprovalYear int    //批准年份
	Comment      string //学科分类
	CommentS1    string //学科分类1
	CommentS2    string //学科分类2
	CommentS3    string //学科分类3
	CommentCode1 string //学科代码1
	CommentCode2 string //学科代码2
	CommentCode3 string //学科代码3
	StartTime    string //开始时间 年月
	EndTime      string //完成时间 年月
	KeywordCN    string //中文关键字
	KeywordEN    string //英文关键字
	Abstract     string //摘要
}
