package data

const userTempl = `
{{ $baseDn := .baseDn }}
{{ range $i, $v := .users }}
#{{$i}}
dn: {{GetDn $v $baseDn}}
cn: {{$v.LoginId}}
idcard: {{$v.IdCard}}
mail: {{$v.Email}}
mobile: {{$v.Mobile}}
objectclass: top
objectclass: person
objectclass: organizationalPerson
objectclass: inetOrgPerson
objectclass: gmccUser
outtime: {{GetOutTime $v}}
pwdstate: {{$v.PwdState}}
sex: {{$v.Sex}}
sn: {{$v.Name}}
state: {{$v.State}}
uid: {{$v.ID}}
userpassword: {md5}CY9rzUYh03PK3k6DJie09g==

{{end}}`
