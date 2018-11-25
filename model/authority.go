package model

type Authority struct {
	UserId   string `gorm:"column:login_id;type:varchar(32)"`
	RoleId   string `gorm:"column:role_code;type:varchar(32)"`
	RoleName string `gorm:column:role_name;`
	OrgId    string `gorm:"column:group_code;type:varchar(32)"`
}
