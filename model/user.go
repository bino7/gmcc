package model

import "time"

//用户
type User struct {
	ID         uint      `gorm:"column:user_id;type:bigint"`        //用户ID
	Name       string    `gorm:"column:user_name;type:varchar(30)"` //用户姓名
	LoginId    string    `gorm:"column:login_id;type:varchar(30)"`  //登录ID
	Password   string    `gorm:"column:password;type:varchar(255)"` //密码
	Sex        int       `gorm:"column:sex;type:int"`               //性别
	Email      string    `gorm:"column:email;type:varchar(100)"`    //邮箱地址
	Mobile     string    `gorm:"column:mobile;type:varchar(20)"`    //移动电话
	State      int       `gorm:"column:state;type:int"`             //状态
	PwdState   int       `gorm:"column:pwd_state;type:int"`         //密码状态
	Memo       string    `gorm:"column:memo;type:varchar(500)"`     //备注信息
	RegDate    time.Time `gorm:"column:reg_date;type:timestamp"`    //注册时间
	UpdateDate time.Time `gorm:"column:update_date;type:timestamp"` //更新时间
	CreaterID  uint      `gorm:"column:creater_id;type:bigint"`     //创建人
	OrgID      string    `gorm:"column:org_id;type:varchar(255)"`   //机构编码
	IdCard     string    `gorm:"column:id_card;type:varchar(255)"`  //身份证号码
	OutTime    time.Time `gorm:"column:out_time;type:timestamp"`    //过期时间
}
