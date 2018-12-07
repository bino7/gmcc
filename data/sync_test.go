package data

import (
	"fmt"
	"github.com/bino7/gmcc/model"
	"github.com/go-ldap/ldap"
	"github.com/jinzhu/gorm"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	db, err := gorm.Open("postgres", "sslmode=disable host=58.213.165.248 port=5432 "+
		"user=postgres dbname=nokk password=Standard@2017")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var count int
	db.Table("pure_user").Count(&count)
	if count != 6632 {
		t.Fail()
	}
}

func TestIsEntryExisted(t *testing.T) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "bino", 389))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	username := "cn=admin,dc=10086,dc=cn"
	password := "wtf5560#@*"
	err = conn.Bind(username, password)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer conn.Close()
	ok, err := isEntryExistedByFilter("(outTime!=20181130145227+0800)", "ou=people,dc=10086,dc=cn", conn)
	fmt.Println(ok, err)
	if !ok {
		t.Fail()
	}
}

func TestFilter(t *testing.T) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "bino", 389))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	username := "cn=admin,dc=10086,dc=cn"
	password := "wtf5560#@*"
	err = conn.Bind(username, password)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer conn.Close()

	rs, err := conn.Search(ldap.NewSearchRequest("ou=People,dc=10086,dc=cn", ldap.ScopeWholeSubtree, 0, 0, 0,
		false, "((outTime<=CurrentTimestamp))", []string{"cn,outTime,CurrentTimestamp"}, nil))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(len(rs.Entries))
	rs.PrettyPrint(2)
	rs.Print()

}

func TestDN(t *testing.T) {
	dn := "cn=admin,dc=10086,dc=cn"
	t.Log(strings.SplitN(dn, ",", 2))
}

func TestSync(t *testing.T) {
	fmt.Println(sync("bino", "cn=admin,dc=10086,dc=cn", "wtf5560#@*",
		"dc=10086,dc=cn", "disable", "58.213.165.248", 5432, "postgres",
		"nokk", "Standard@2017"))
}

func TestGetRootRole(t *testing.T) {
	db, err := gorm.Open("postgres", "sslmode=disable host=58.213.165.248 port=5432 "+
		"user=postgres dbname=nokk password=Standard@2017")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var roles []model.Role
	db.Table("pure_role").Find(&roles, "role_id='root'")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(roles)
	/*for rows.Next() {
		var role model.Role
		fmt.Println(rows.Columns())
		rows.Scan(&role.ID,&role.ParentId,&role.Name,&role.Memo,&role.Ord)
		fmt.Println(role)
	}*/
}

func TestPasswordModify(t *testing.T) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "bino", 389))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	username := "cn=admin,dc=10086,dc=cn"
	password := "wtf5560#@*"
	err = conn.Bind(username, password)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer conn.Close()
	conn.Debug = true
	rs, err := conn.PasswordModify(ldap.NewPasswordModifyRequest("cn=test,ou=people,dc=10086,dc=cn", "test",
		"test"))
	t.Log(rs, err)
}
