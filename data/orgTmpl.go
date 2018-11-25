package data

var orgTmpl = `
{{$baseDn := .baseDn}}
{{range .orgs }}
dn: ou={{.ID}},{{$baseDn}}
ord: {{.Ord}}
ouid: {{.ID}}
ou: {{.ID}}
objectClass: top
objectClass: organizationalUnit
objectClass: gmccUnit
description: {{.Name}}
{{end}}`
