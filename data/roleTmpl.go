package data

var roleTmpl = `
{{$baseDn := .baseDn}}
{{range .roles }}
dn: ou={{.ID}},{{$baseDn}}
ord: {{.Ord}}
ouid: {{.ID}}
ou: {{.ID}}
objectClass: top
objectClass: organizationalUnit
objectClass: gmccUnit
description: {{.Name}}

{{end}}`
