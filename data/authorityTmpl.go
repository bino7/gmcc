package data

var authorityTmpl = `{{$baseDn := .baseDn}}
{{range $k,$v := .authorities }}
dn: cn={{$k}},{{$baseDn}}
{{range $v}}
{{if $k eq 0 }}
objectclass: top
objectclass: groupOfNames
cn: {{.RoleId}}
description: {{.RoleName}}
{{end}}
member: {{.UserId}}
{{end}}
{{end}}`
