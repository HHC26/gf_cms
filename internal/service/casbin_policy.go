// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	ICasbinPolicy interface {
		ObjBackend() string
		ObjBackendApi() string
		CheckByRoleId(roleId, obj, act string) bool
		CheckByAccountId(AccountId, obj, act string) bool
		AddByRoleId(roleId, obj, act string) bool
		RemoveByRoleId(roleId, obj, act string) bool
	}
)

var (
	localCasbinPolicy ICasbinPolicy
)

func CasbinPolicy() ICasbinPolicy {
	if localCasbinPolicy == nil {
		panic("implement not found for interface ICasbinPolicy, forgot register?")
	}
	return localCasbinPolicy
}

func RegisterCasbinPolicy(i ICasbinPolicy) {
	localCasbinPolicy = i
}
