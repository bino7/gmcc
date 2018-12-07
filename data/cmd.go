package data

import (
	"fmt"
	"github.com/spf13/cobra"
)

var CMDS []*cobra.Command

func init() {
	if err := initTmpl(); err != nil {
		panic(err)
	}
	var dbHost, dbUser, dbName, dbPassword string
	var dbPort int
	var sslmode string
	var sql, recursivelySql, authoritySql, outputFile, baseDn string
	var cmdExport = &cobra.Command{
		Use:   `export --sslmode [sslmode] --dbHost [db host] --dbPort [db port] --dbUser [db user] --dbName [db name] --dbPassword [db password]`,
		Short: "sync data from postgres db to ldap",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			t := args[0]
			var err error
			switch t {
			case "user":
				err = exportUsers(baseDn, sql, outputFile, dbHost, dbPort, dbName, dbUser, dbPassword, sslmode)
			case "role":
				err = exportRoles(baseDn, sql, recursivelySql, outputFile, dbHost, dbPort, dbName, dbUser, dbPassword, sslmode)
			case "org":
				err = exportOrgs(baseDn, sql, recursivelySql, authoritySql, outputFile, dbHost, dbPort, dbName, dbUser, dbPassword, sslmode)
			case "authority":
				err = exportAuthorityByOrgId(baseDn, sql, outputFile, dbHost, dbPort, dbName, dbUser, dbPassword, sslmode)
			default:
				fmt.Printf("unknow export type %s, should be one of user,role,org\n", t)
				return
			}
			if err != nil {
				fmt.Printf("export failed, %v\n", err)
			} else {
				fmt.Println("export success!")
			}
		},
	}

	cmdExport.Flags().StringVarP(&sslmode, "sslmode", "m", "disable", "db sslmode")
	cmdExport.Flags().StringVarP(&dbHost, "dbhost", "H", "58.213.165.248", "db host")
	cmdExport.Flags().IntVarP(&dbPort, "dbport", "p", 5432, "db port")
	cmdExport.Flags().StringVarP(&dbUser, "dbuser", "u", "postgres", "db user")
	cmdExport.Flags().StringVarP(&dbName, "dbname", "n", "nokk", "db name")
	cmdExport.Flags().StringVarP(&dbPassword, "dbpassword", "w", "Standard@2017", "db password")
	cmdExport.Flags().StringVarP(&sql, "sql", "s", "", "sql")
	cmdExport.Flags().StringVarP(&outputFile, "outputFile", "o", "", "output file")
	cmdExport.Flags().StringVarP(&recursivelySql, "recursivelySql", "r", "", "recursively sql")
	cmdExport.Flags().StringVarP(&authoritySql, "authoritySql", "a", "", "authority sql")
	cmdExport.Flags().StringVarP(&baseDn, "baseDn", "b", "", "base dn")

	CMDS = []*cobra.Command{cmdExport}
}
