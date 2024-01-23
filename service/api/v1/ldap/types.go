package ldap

import (
	v1Ldap "github.com/godtool/kubeone/service/model/v1/ldap"
	"github.com/godtool/kubeone/service/model/v1/user"
)

type Ldap struct {
	v1Ldap.Ldap
}

type TestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ImportRequest struct {
	Users []user.ImportUser `json:"users"`
}
