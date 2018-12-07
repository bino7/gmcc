package model

//角色
type Role struct {
	ID       string `gorm:"column:role_id;type:varchar(30)"`   //角色编码
	ParentId string `gorm:"column:parent_id;type:varchar(30)"` //上级角色编码
	Name     string `gorm:"column:role_name;type:varchar(50)"` //角色名称
	Memo     string `gorm:"column:memo;type:varchar(500)"`     //注释
	Ord      int    `gorm:"column:ord;type:int4(32)"`          //排序
}
