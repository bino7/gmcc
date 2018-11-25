package model

//组织机构
type Org struct {
	ID       string `gorm:"column:org_id;type:varchar(255)"`    //机构ID
	Name     string `gorm:"column:org_name;type:varchar(255)"`  //机构名称
	Path     string `gorm:"column:path;type:varchar(255)"`      //路径（/开头，使用/表示上下级）
	Ord      int    `gorm:"column:ord;type:int4(32)"`           //排序
	ParentId string `gorm:"column:parent_id;type:varchar(255)"` //上级编码
}
