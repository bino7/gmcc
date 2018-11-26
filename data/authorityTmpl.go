package data

var authorityTmpl = `{{$baseDn := .baseDn}}
{{range $k,$v := .authorities }}
dn: cn={{$k}},{{$baseDn}}
objectclass: top
objectclass: groupOfNames
{{$f:=index $v 0}}cn: {{$f.RoleId}}
description:{{$f.RoleName}}{{range $v}}
member: cn={{.UserId}},ou=people,dc=10086,dc=cn{{end}}
{{end}}`
