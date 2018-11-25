package data

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var CMDS []*cobra.Command

func init() {
	if err := initTmpl(); err != nil {
		panic(err)
	}
	var echoTimes int
	var cmdPrint = &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}
	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [# times] [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}
	var test string
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
	cmdTimes.Flags().StringVarP(&test, "test", "f", "df", "test")
	//"sslmode=disable host=58.213.165.248 port=5432 user=postgres dbname=nokk password=Standard@2017"
	var ldapHost, ldapUserDn, ldapPassword, ldapBaseDn, dbHost, dbUser, dbName, dbPassword string
	var dbPort int
	var sslmode string
	var cmdSync = &cobra.Command{
		Use:   `sync --ldap_host [ldap host] --ldap_user_dn [ldap user dn] --ldap_password [ldap password] --ldap_base_dn [ldap base dn] --sslmode [sslmode] --dbHost [db host] --dbPort [db port] --dbUser [db user] --dbName [db name] --dbPassword [db password]`,
		Short: "sync data from postgres db to ldap",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			userSyncResult, roleSyncResult, orgSyncResult, err := sync(ldapHost, ldapUserDn, ldapPassword, ldapBaseDn, sslmode, dbHost, dbPort, dbUser, dbName, dbPassword)
			if err != nil {
				panic(err)
			}
			fmt.Printf("total %d user,handled %d,success %d,created %d,modifyed %d,failed %d\n", userSyncResult.Total, userSyncResult.Handled, userSyncResult.Success,
				userSyncResult.Created, userSyncResult.Modified, userSyncResult.Failed)
			if len(userSyncResult.Errors) > 0 {
				fmt.Println("errors:")
				for _, err := range userSyncResult.Errors {
					fmt.Println(err)
				}
			}

			fmt.Printf("total %d role,handled %d,success %d,created %d,modifyed %d,failed %d\n", roleSyncResult.Total, roleSyncResult.Handled, roleSyncResult.Success,
				roleSyncResult.Created, roleSyncResult.Modified, roleSyncResult.Failed)
			if len(roleSyncResult.Errors) > 0 {
				fmt.Println("errors:")
				for _, err := range roleSyncResult.Errors {
					fmt.Println(err)
				}
			}

			fmt.Printf("total %d org,handled %d,success %d,created %d,modifyed %d,failed %d\n", orgSyncResult.Total, orgSyncResult.Handled, orgSyncResult.Success,
				orgSyncResult.Created, orgSyncResult.Modified, orgSyncResult.Failed)
			if len(orgSyncResult.Errors) > 0 {
				fmt.Println("errors:")
				for _, err := range orgSyncResult.Errors {
					fmt.Println(err)
				}
			}

		},
	}
	cmdSync.Flags().StringVar(&ldapHost, "ldap_host", "localhost", "ldap host")
	cmdSync.Flags().StringVar(&ldapUserDn, "ldap_user_dn", "cn=admin", "ldap user dn")
	cmdSync.Flags().StringVar(&ldapPassword, "ldap_password", "password", "ldap password")
	cmdSync.Flags().StringVar(&ldapBaseDn, "ldap_base_dn", "dc=example,dc=com", "ldap base dn")
	cmdSync.Flags().StringVar(&sslmode, "sslmode", "disable", "db sslmode")
	cmdSync.Flags().StringVar(&dbHost, "dbhost", "localhost", "db host")
	cmdSync.Flags().IntVar(&dbPort, "dbport", 5432, "db port")
	cmdSync.Flags().StringVar(&dbUser, "dbuser", "admin", "db user")
	cmdSync.Flags().StringVar(&dbName, "dbname", "localhost", "db name")
	cmdSync.Flags().StringVar(&dbPassword, "dbpassword", "password", "db password")

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

	CMDS = []*cobra.Command{cmdPrint, cmdEcho, cmdTimes, cmdSync, cmdExport}
}
