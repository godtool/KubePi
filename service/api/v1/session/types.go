package session

import v1 "k8s.io/api/rbac/v1"

type LoginCredential struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	AuthMethod string `json:"authMethod"`
}
type MfaCredential struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
	Code     string `json:"code"`
}

type PasswordSetter struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type ProfileSetter struct {
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	Language string `json:"language"`
}

type UserProfile struct {
	Name                string              `json:"name"`
	NickName            string              `json:"nickName"`
	Email               string              `json:"email"`
	Language            string              `json:"language"`
	ResourcePermissions map[string][]string `json:"resourcePermissions"`
	IsAdministrator     bool                `json:"isAdministrator"`
	Mfa                 Mfa                 `json:"mfa"`
}

type ClusterUserProfile struct {
	UserProfile
	ClusterRoles []v1.ClusterRole `json:"clusterRoles"`
}

type Mfa struct {
	Enable   bool   `json:"enable"`
	Secret   string `json:"secret"`
	Approved bool   `json:"approved"`
}

type UserInfo struct {
	DomainAccount string `json:"loginName"`  // 域账号
	EmpID         string `json:"empId"`      // 工号，用户ID(empID)
	Name          string `json:"lastName"`   // 真名（不唯一）
	NickName      string `json:"nickNameCn"` // 花名,并非所有人都有
	UserType      string `json:"userType"`   //员工类型：R,正式; O,外包; W,部门公共账号
	HrStatus      string `json:"hrStatus"`   //在职状态：A,在职; I,离职
	Available     string `json:"available"`  // 账号状态：T,有效; F,无效
	Email         string `json:"emailAddr"`  // 常用邮箱
	CellPhone     string `json:"cellPhone"`  // 手机号
	Department    string `json:"depDesc"`    // 部门
	AvatarURL     string `json:"avatarURL"`  // 头像地址
	Token         string `json:"token"`      // 登录会话SSO_TOKEN，有效期7d，可以换取SSO_TICKET,可以用于心跳请求，心跳后原TOKEN在1分钟后失效
	DisplayName   string // 显示名称（唯一），优先显示花名，然后真实姓名
	PicURL        string // 头像
}
