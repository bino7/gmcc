package data

import (
	"bufio"
	"fmt"
	. "github.com/bino7/gmcc/model"
	"github.com/jinzhu/gorm"
	"io"
	"os"
	"text/template"
)

func initTmpl() error {
	if err := initUserTmpl(); err != nil {
		return err
	}
	if err := initRoleTmpl(); err != nil {
		return err
	}
	if err := initOrgTmpl(); err != nil {
		return err
	}
	if err := initAuthorityTmpl(); err != nil {
		return err
	}
	return nil
}

var userT *template.Template

func initUserTmpl() (err error) {
	var funcMap = template.FuncMap{
		"GetOutTime": func(u User) string {
			if u.OutTime.IsZero() {
				return "20060102150405+0800"
			}
			return u.OutTime.Format("20060102150405+0800")
		},
		"GetDn": func(u User, baseDn string) string {
			return fmt.Sprintf("cn=%s,%s", u.LoginId, baseDn)
		},
	}
	userT, err = template.New("user").Funcs(funcMap).Parse(userTempl)
	return
}

var orgT *template.Template

func initOrgTmpl() (err error) {
	var funcMap = template.FuncMap{
		"GetDn": getOrgDn,
	}
	orgT, err = template.New("org").Funcs(funcMap).Parse(orgTmpl)
	return
}

var roleT *template.Template

func initRoleTmpl() (err error) {
	var funcMap = template.FuncMap{
		"GetDn": getRoleDn,
	}
	roleT, err = template.New("role").Funcs(funcMap).Parse(roleTmpl)
	return
}

var authorityT *template.Template

func initAuthorityTmpl() (err error) {
	authorityT, err = template.New("authority").Parse(authorityTmpl)
	return
}
func exportUsers(baseDn, sql, outputFile, dbHost string, dbPort int, dbName, dbUser, dbPassword, sslmode string) error {
	db, err := openDB(sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	if err != nil {
		return err
	}
	defer db.Close()

	var users []*User
	db.Raw(sql).Scan(&users)
	fmt.Printf("exporting %d user to %s\n", len(users), outputFile)

	w, err := getWriter(outputFile)
	if err != nil {
		return err
	}
	defer w.Flush()

	if userT == nil {
		if err := initUserTmpl(); err != nil {
			return err
		}
	}

	data := map[string]interface{}{
		"baseDn": baseDn,
		"users":  users,
	}
	err = userT.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func exportRoles(baseDn, sql, rSql, outputFile, dbHost string, dbPort int, dbName, dbUser, dbPassword, sslmode string) error {
	db, err := openDB(sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	if err != nil {
		return err
	}
	defer db.Close()

	if roleT == nil {
		if err := initRoleTmpl(); err != nil {
			return err
		}
	}

	var roles []Role
	db.Raw(sql).Scan(&roles)

	if rSql == "" {
		fmt.Printf("ecursively exporting roles of %d root roles to %s\n", len(roles), outputFile)
	} else {
		fmt.Printf("exporting roles of %d root roles to %s\n", len(roles), outputFile)
	}

	w, err := getWriter(outputFile)
	if err != nil {
		return err
	}
	defer w.Flush()

	return exportRole(baseDn, roles, rSql, w, db)
}

func exportRole(baseDn string, roles []Role, rSql string, w io.Writer, db *gorm.DB) error {
	data := map[string]interface{}{
		"baseDn": baseDn,
		"roles":  roles,
	}
	err := roleT.Execute(w, data)
	if err != nil {
		return err
	}
	if rSql != "" {
		for _, r := range roles {
			var kRoles []Role
			db.Raw(rSql, r.ID).Scan(&kRoles)
			if len(kRoles) > 0 {
				if err := exportRole(getRoleDn(r, baseDn), kRoles, rSql, w, db); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func exportOrgs(baseDn, sql, rSql, authoritySql, outputFile, dbHost string, dbPort int, dbName, dbUser, dbPassword, sslmode string) error {
	db, err := openDB(sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	if err != nil {
		return err
	}
	defer db.Close()

	var orgs []Org
	db.Raw(sql).Scan(&orgs)

	if rSql == "" {
		fmt.Printf("ecursively exporting %d root orgs to %s\n", len(orgs), outputFile)
	} else {
		fmt.Printf("exporting %d root orgs to %s\n", len(orgs), outputFile)
	}

	w, err := getWriter(outputFile)
	if err != nil {
		return err
	}
	defer w.Flush()

	return exportOrg(baseDn, orgs, rSql, authoritySql, w, db)
}

func getWriter(outputFile string) (*bufio.Writer, error) {
	f, err := os.Create(outputFile)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)
	return w, nil
}

func exportOrg(baseDn string, orgs []Org, rSql, authoritySql string, w io.Writer, db *gorm.DB) error {
	data := map[string]interface{}{
		"baseDn": baseDn,
		"orgs":   orgs,
	}

	for _, o := range orgs {
		err := orgT.Execute(w, data)
		if err != nil {
			return err
		}
		if authoritySql != "" {
			authorities := getAuthorities(authoritySql, o.ID, db)
			if len(authorities) > 0 {
				dn := getOrgDn(o, baseDn)
				if err := exportAuthorityByOrgIdIntenal(authorities, dn, w); err != nil {
					return err
				}
			}
		}

		if rSql != "" {
			var kOrgs []Org
			db.Raw(rSql, o.ID).Scan(&kOrgs)
			if len(kOrgs) > 0 {
				if err := exportOrg(getOrgDn(o, baseDn), kOrgs, rSql, authoritySql, w, db); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func exportAuthorityByOrgId(baseDn, sql, outputFile, dbHost string, dbPort int, dbName, dbUser, dbPassword, sslmode string) error {
	if authorityT == nil {
		initAuthorityTmpl()
	}
	db, err := openDB(sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
	if err != nil {
		return err
	}

	authorities := getAuthorities(sql, "", db)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	defer w.Flush()

	return exportAuthorityByOrgIdIntenal(authorities, baseDn, w)
}

func getAuthorities(sql, oid string, db *gorm.DB) []Authority {
	var authorities []Authority
	if oid == "" {
		db.Raw(sql).Scan(&authorities)
	} else {
		db.Raw(sql, oid).Scan(&authorities)
	}

	return authorities
}

func exportAuthorityByOrgIdIntenal(authorities []Authority, baseDn string, w io.Writer) error {
	mapper := make(map[string][]Authority)
	for _, v := range authorities {
		k := v.RoleId
		if mapper[k] == nil {
			mapper[k] = make([]Authority, 0)
		}
		mapper[k] = append(mapper[k], v)
	}
	data := map[string]interface{}{
		"baseDn":      baseDn,
		"authorities": mapper,
	}
	if err := authorityT.Execute(w, data); err != nil {
		return err
	}
	return nil
}
