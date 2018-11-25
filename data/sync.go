package data

import (
	"fmt"
	"github.com/bino7/gmcc/model"
	"github.com/go-ldap/ldap"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
	"strings"
)

func sync(ldapHost, ldapUserDn, ldapPassword, ldapBaseDn, sslmode, dbHost string, dbPort int,
	dbUser, dbName, dbPassword string) (*SyncResult, *SyncResult, *SyncResult, error) {
	db, err := openDB(sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open db failed %v", err)
	}
	defer db.Close()

	ldapConn, err := openLdap(ldapHost, ldapUserDn, ldapPassword)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open ldap falied %v", err)
	}
	defer ldapConn.Close()

	userSyncResult, err := syncUsers(db, ldapConn, ldapBaseDn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sync user failed %v", err)
	}
	roleSyncResult, err := syncRoles(db, ldapConn)
	if err != nil {
		return userSyncResult, nil, nil, fmt.Errorf("sync role failed %v", err)
	}
	orgSyncResult, err := syncOrgs(db, ldapConn)
	if err != nil {
		return userSyncResult, roleSyncResult, nil, fmt.Errorf("sync org failed %v", err)
	}
	return userSyncResult, roleSyncResult, orgSyncResult, nil
}

type SyncResult struct {
	Total    int
	Handled  int
	Success  int
	Failed   int
	Created  int
	Modified int
	Errors   []error
}

func (sr *SyncResult) AddFailed(err error) {
	fmt.Println(err)
	sr.Failed++
	if sr.Errors == nil {
		sr.Errors = make([]error, 0)
	}
	sr.Errors = append(sr.Errors, err)
}

func syncUsers(db *gorm.DB, ldapConn *ldap.Conn, ldapBaseDn string) (*SyncResult, error) {
	syncResult := &SyncResult{}
	db.Table("pure_user").Count(&syncResult.Total)
	if syncResult.Total == 0 {
		return syncResult, nil
	}
	var users []model.User
	db.Table("pure_user").Find(&users)
	fmt.Printf("syncing %d users\n", syncResult.Total)
	for _, u := range users {
		syncResult.Handled++
		dn := getUserDn(u)
		if ok, err := isEntryExisted(dn, ldapConn); err != nil {
			syncResult.AddFailed(err)
			fmt.Printf("user %s search failed \n", u.LoginId)
		} else if ok {
			//modify user

			err := ldapConn.Del(ldap.NewDelRequest(dn, nil))
			if err != nil {
				err = fmt.Errorf("delete user %s for modify failed", u.LoginId)
				syncResult.AddFailed(err)
			} else {
				req := getUserAddRequest(u)
				err := ldapConn.Add(req)
				if err != nil {
					err = fmt.Errorf("user %s modified failed %v \n", u.LoginId, err)
					syncResult.AddFailed(err)
				} else {
					modifyUserPassword(dn, u, "", ldapConn, syncResult)
					syncResult.Modified++
					fmt.Printf("user %s modified \n", u.LoginId)
				}
			}
		} else {
			req := getUserAddRequest(u)
			err := ldapConn.Add(req)
			if err != nil {
				err = fmt.Errorf("user %s created failed %v \n", u.LoginId, err)
				syncResult.AddFailed(err)
			} else {
				syncResult.Created++
				fmt.Printf("user %s created \n", u.LoginId)
			}
		}
	}
	return syncResult, nil
}

func modifyUserPassword(dn string, u model.User, oldPassword string, ldapConn *ldap.Conn, syncResult *SyncResult) {
	req := ldap.NewPasswordModifyRequest(dn, "", u.Password)
	if _, err := ldapConn.PasswordModify(req); err != nil {
		err = fmt.Errorf("user %s password modify failed", u.LoginId)
		syncResult.AddFailed(err)
	} else {
		fmt.Printf("user %s password modify success\n", u.LoginId)
	}
}

func getUserDn(u model.User) string {
	return fmt.Sprintf("cn=%s,ou=people,dc=10086,dc=cn", u.LoginId)
}

func isEntryExisted(dn string, ldapConn *ldap.Conn) (bool, error) {
	rdns := strings.SplitN(dn, ",", 2)
	if len(rdns) != 2 {
		return false, fmt.Errorf("wrong dn %s", dn)
	}
	filter := fmt.Sprintf("(&(%s))", rdns[0])
	baseDn := rdns[1]
	return isEntryExistedByFilter(filter, baseDn, ldapConn)
}

func isEntryExistedByFilter(filter, baseDn string, ldapConn *ldap.Conn) (bool, error) {
	return isExistedByFilter(filter, baseDn, ldapConn, ldap.ScopeSingleLevel)
}

func isExistedByFilter(filter, baseDn string, ldapConn *ldap.Conn, scope int) (bool, error) {
	req := &ldap.SearchRequest{
		baseDn,
		scope,
		ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{},
		nil,
	}
	rs, err := ldapConn.Search(req)
	if err != nil {
		return false, err
	}
	return len(rs.Entries) > 0, nil
}

func getUserAddRequest(u model.User) *ldap.AddRequest {
	dn := fmt.Sprintf("cn=%s,ou=people,dc=10086,dc=cn", u.LoginId)
	addReq := ldap.NewAddRequest(dn, nil)
	addReq.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "inetOrgPerson", "gmccUser"})
	id := fmt.Sprintf("%d", u.ID)
	addReq.Attribute("uid", []string{id})
	addReq.Attribute("cn", []string{u.LoginId})
	//addReq.Attribute("userPassword ", []string{fmt.Sprintf("{MD5}%s",u.Password)})
	addReq.Attribute("sn", []string{u.Name})
	sex := fmt.Sprintf("%d", u.Sex)
	addReq.Attribute("sex", []string{sex})
	addReq.Attribute("mail", []string{u.Email})
	addReq.Attribute("mobile", []string{u.Mobile})
	state := fmt.Sprintf("%d", u.State)
	addReq.Attribute("state", []string{state})
	ps := fmt.Sprintf("%d", u.PwdState)
	addReq.Attribute("pwdState", []string{ps})
	//addReq.Attribute("description",[]string{u.Memo})
	addReq.Attribute("idCard", []string{u.IdCard})
	addReq.Attribute("outTime", []string{"20181130145227+0800"})
	return addReq
}

func isUserExisted(userId string, ldapConn *ldap.Conn) (bool, error) {
	return isExistedByFilter(fmt.Sprintf("(cn=%s)", userId), "ou=people,dc=10086,dc=cn", ldapConn, ldap.ScopeWholeSubtree)
}

func syncRoles(db *gorm.DB, ldapConn *ldap.Conn) (*SyncResult, error) {
	syncResult := &SyncResult{}
	db.Table("pure_role").Count(&syncResult.Total)
	if syncResult.Total == 0 {
		return syncResult, nil
	}
	var roles []model.Role
	db.Table("pure_role").Find(&roles, "role_id='root'")
	if len(roles) == 0 {
		return syncResult, fmt.Errorf("role root not found")
	}
	root := roles[0]
	fmt.Printf("syncing %d roles \n", syncResult.Total)
	err := syncRole(root, "ou=roles,dc=10086,dc=cn", db, ldapConn, syncResult)
	return syncResult, err
}

func syncRole(r model.Role, parentDn string, db *gorm.DB, ldapConn *ldap.Conn, syncResult *SyncResult) error {
	syncResult.Handled++
	if ok, err := isEntryExisted(getRoleDn(r, parentDn), ldapConn); err != nil {
		syncResult.AddFailed(fmt.Errorf("role %s search failed %v", r.ID, err))
		return err
	} else if ok {
		//modify user
		syncResult.Modified++
		fmt.Printf("role %s modified \n", r.ID)
	} else {
		dn := getRoleDn(r, parentDn)
		req := getRoleAddRequest(r, dn)
		err := ldapConn.Add(req)
		if err != nil {
			syncResult.AddFailed(fmt.Errorf("role %s created failed %v", r.ID, err))
			return err
		} else {
			syncResult.Created++
			fmt.Printf("role %s created \n", r.ID)
		}
	}

	kids := getRoleKids(r, db)
	dn := getRoleDn(r, parentDn)
	for _, k := range kids {
		syncRole(k, dn, db, ldapConn, syncResult)
	}
	return nil
}

func getRoleKids(r model.Role, db *gorm.DB) []model.Role {
	var kids []model.Role
	db.Table("pure_role").Find(&kids, fmt.Sprintf("parent_id='%s'", r.ID))
	return kids
}

func getRoleDn(r model.Role, parentDn string) string {
	return fmt.Sprintf("ou=%s,%s", r.ID, parentDn)
}

func getRoleAddRequest(r model.Role, dn string) *ldap.AddRequest {
	addReq := ldap.NewAddRequest(dn, nil)
	addReq.Attribute("objectClass", []string{"top", "organizationalUnit", "gmccUnit"})
	addReq.Attribute("ou", []string{r.ID})
	addReq.Attribute("ouid", []string{r.ID})
	addReq.Attribute("ord", []string{strconv.Itoa(r.Ord)})
	addReq.Attribute("description", []string{r.Name})
	return addReq
}

func isRoleExisted(rid string, ldapConn *ldap.Conn) (bool, error) {
	return isExistedByFilter(fmt.Sprintf("(ou=%s)", rid), "ou=roles,dc=10086,dc=cn", ldapConn, ldap.ScopeWholeSubtree)
}

func getRoleNameById(rid string, ldapConn *ldap.Conn) (string, error) {
	req := &ldap.SearchRequest{
		"ou=roles,dc=10086,dc=cn",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(ou=%s)", rid),
		[]string{"description"},
		nil,
	}
	rs, err := ldapConn.Search(req)
	if err != nil {
		return "", fmt.Errorf("search role error %v", err)
	}
	if len(rs.Entries) == 0 {
		return "", fmt.Errorf("role with id %s not found", rid)
	}
	if len(rs.Entries) > 1 {
		return "", fmt.Errorf("found more than one role with id %s", rid)
	}
	return rs.Entries[0].GetAttributeValue("description"), nil

}

func syncOrgs(db *gorm.DB, ldapConn *ldap.Conn) (*SyncResult, error) {
	syncResult := &SyncResult{}
	db.Table("pure_org").Count(&syncResult.Total)
	if syncResult.Total == 0 {
		return syncResult, nil
	}
	var orgs []model.Org
	db.Table("pure_org").Find(&orgs, "org_id='020'")
	if len(orgs) == 0 {
		return syncResult, fmt.Errorf("org root not found")
	}
	root := orgs[0]
	fmt.Printf("syncing %d orgs \n", syncResult.Total)
	err := syncOrg(root, "ou=orgs,dc=10086,dc=cn", db, ldapConn, syncResult)
	return syncResult, err
}

func syncOrg(o model.Org, parentDn string, db *gorm.DB, ldapConn *ldap.Conn, syncResult *SyncResult) error {
	syncResult.Handled++
	if ok, err := isEntryExisted(getOrgDn(o, parentDn), ldapConn); err != nil {
		syncResult.AddFailed(fmt.Errorf("org %s search failed %v", o.ID, err))
		return err
	} else if ok {
		//modify user
		syncResult.Modified++
		fmt.Printf("org %s modified \n", o.ID)
	} else {
		dn := getOrgDn(o, parentDn)
		req := getOrgAddRequest(o, dn)
		err := ldapConn.Add(req)
		if err != nil {
			syncResult.AddFailed(fmt.Errorf("org %s created failed %v", o.ID, err))
			return err
		} else {
			syncResult.Created++
			fmt.Printf("org %s created \n", o.ID)
		}
	}

	userOrgAndRoles := make([]model.Authority, 0)
	rows, err := db.Table("pure_usergroupandrole").Select("login_id,role_code,group_code").Where("group_code=?", o.ID).Group("login_id,role_code,group_code").Rows()
	if err != nil {
		syncResult.AddFailed(fmt.Errorf("org %s search pure_usergroupandrole failed %v", o.ID, err))
	}
	for rows.Next() {
		data := model.Authority{}
		rows.Scan(&data.UserId, &data.RoleId, &data.OrgId)
		userOrgAndRoles = append(userOrgAndRoles, data)
	}
	rows.Close()
	mapper := make(map[string][]model.Authority)
	for _, v := range userOrgAndRoles {
		ok1, _ := isRoleExisted(v.RoleId, ldapConn)
		ok2, _ := isUserExisted(v.UserId, ldapConn)
		if ok1 && ok2 {
			k := fmt.Sprintf("%s-%s", v.RoleId, v.OrgId)
			if mapper[k] == nil {
				mapper[k] = make([]model.Authority, 0)
			}
			mapper[k] = append(mapper[k], v)
		}
	}
	dn := getOrgDn(o, parentDn)
	for _, v := range mapper {
		if len(v) == 0 {
			continue
		}
		rid := v[0].RoleId
		rName, err := getRoleNameById(rid, ldapConn)
		if err != nil {
			syncResult.AddFailed(fmt.Errorf("org %s add memberOfNames failed %v", o.ID, err))
			continue
		}
		req := getAuthoritysAddRequest(dn, rid, rName, v)
		if ok, _ := isEntryExisted(req.DN, ldapConn); ok {
			if err := ldapConn.Del(ldap.NewDelRequest(req.DN, nil)); err != nil {
				syncResult.AddFailed(fmt.Errorf("org %s delete memberOfNames failed %v", o.ID, err))
				continue
			}
		}
		if err := ldapConn.Add(req); err != nil {
			syncResult.AddFailed(fmt.Errorf("org %s add memberOfNames failed %v", o.ID, err))
		} else {
			fmt.Printf("memberOfNames %s added\n", req.DN)
		}
	}

	kids := getOrgKids(o, db)
	for _, k := range kids {
		syncOrg(k, dn, db, ldapConn, syncResult)
	}
	return nil
}

func getAuthoritysAddRequest(baseDn, rid, rName string, userOrgAndRoles []model.Authority) *ldap.AddRequest {
	dn := fmt.Sprintf("cn=%s,%s", rid, baseDn)
	member := make([]string, len(userOrgAndRoles))
	for i, v := range userOrgAndRoles {
		member[i] = fmt.Sprintf("cn=%s,ou=people,dc=10086,dc=cn", v.UserId)
	}
	addReq := ldap.NewAddRequest(dn, nil)
	addReq.Attribute("objectClass", []string{"top", "groupOfNames"})
	addReq.Attribute("cn", []string{rid})
	addReq.Attribute("member", member)
	addReq.Attribute("description", []string{rName})
	return addReq
}

func getOrgKids(o model.Org, db *gorm.DB) []model.Org {
	var kids []model.Org
	db.Table("pure_org").Find(&kids, fmt.Sprintf("parent_id='%s'", o.ID))
	return kids
}

func getOrgDn(o model.Org, parentDn string) string {
	return fmt.Sprintf("ou=%s,%s", o.ID, parentDn)
}

func getOrgAddRequest(o model.Org, dn string) *ldap.AddRequest {
	addReq := ldap.NewAddRequest(dn, nil)
	addReq.Attribute("objectClass", []string{"top", "organizationalUnit", "gmccUnit"})
	addReq.Attribute("ou", []string{o.ID})
	addReq.Attribute("ouid", []string{o.ID})
	addReq.Attribute("ord", []string{strconv.Itoa(o.Ord)})
	addReq.Attribute("description", []string{o.Name})
	return addReq
}

func openDB(sslmode string, dbHost string, dbPort int, dbUser, dbName, dbPassword string) (*gorm.DB, error) {
	connStr := fmt.Sprintf("sslmode=%s host=%s port=%d user=%s dbname=%s password=%s", sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		err = fmt.Errorf("open db failed %v", err)
	}
	return db, err
}

func openLdap(ldapHost, ldapUserDn, ldapPassword string) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "bino", 389))
	if err != nil {
		return nil, fmt.Errorf("can't open ldap conn %v", err)
	}
	err = conn.Bind(ldapUserDn, ldapPassword)
	if err != nil {
		return nil, fmt.Errorf("ldap user or password is wrong")
	}
	return conn, nil
}
