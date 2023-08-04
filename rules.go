package auth_library

import "errors"

const (
	VIEW       Method = "view"
	MODIFY     Method = "modify"
	MEMBERSHIP Method = "member"
)

type Method string

type Rules interface {
	Authorize(permissionApi PermissionApi, subjects RequestOwner) error
}

type ViewRule struct {
	Profile string
}
type ModifyRule struct {
	Profile string
}
type MembershipRule struct {
}

type DenyRule struct {
}

func (rule DenyRule) Authorize(permissionApi PermissionApi, subjects RequestOwner) error {
	return errors.New("DenyRule")
}

func (rule ViewRule) Authorize(permissionApi PermissionApi, subjects RequestOwner) error {
	user := "profile:" + subjects.UsersProfile
	object := "profile_info:" + rule.Profile
	return permissionApi.Check(user, string(VIEW), object)
}

func (rule ModifyRule) Authorize(permissionApi PermissionApi, subjects RequestOwner) error {
	user := "profile:" + subjects.UsersProfile
	object := "profile_info:" + rule.Profile
	return permissionApi.Check(user, string(MODIFY), object)
}

func (rule MembershipRule) Authorize(permissionApi PermissionApi, subjects RequestOwner) error {
	user := "user:" + subjects.User
	object := "profile:" + subjects.UsersProfile
	return permissionApi.Check(user, string(MEMBERSHIP), object)
}
