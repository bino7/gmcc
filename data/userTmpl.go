package data

const userTempl = `
{{ $baseDn := .baseDn }}
{{ range $i, $v := .users }}
#{{$i}}
dn: {{GetDn $v $baseDn}}
objectclass: top
objectclass: person
objectclass: organizationalPerson
objectclass: inetOrgPerson
objectclass: gmccUser
cn: {{$v.LoginId}}
idcard: {{if eq $v.IdCard ""}}n/a{{else}}{{$v.IdCard}}{{end}}
mail: {{if eq $v.Email ""}}n/a{{else}}{{$v.Email}}{{end}}
mobile: {{if eq $v.Mobile ""}}n/a{{else}}{{$v.Mobile}}{{end}} 
outtime: {{GetOutTime $v}}
pwdstate: {{$v.PwdState}}
sex: {{$v.Sex}}
sn: {{$v.Name}}
state: {{$v.State}}
uid: {{$v.ID}}
userpassword: {md5}{{GetPassword $v.Password}}
orgId: {{$v.OrgId}}

{{end}}`
