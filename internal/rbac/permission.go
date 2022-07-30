package rbac

import (
	"errors"
	"strings"
)

// Role define what roles are available in this app
type Role string

// Role availables here. Should match the enum in the db
const (
	// RoleMember member role
	RoleMember = "MEMBER"

	// RoleAdmin admin role
	RoleAdmin = "ADMIN"
)

// Resource define what resource is available in this app
type Resource string

// Define any resource type available in this app
const (
	ResourceAbsentForm Resource = "absent_form"
	ResourceAbsentList Resource = "absent_list"
	ResourceUser       Resource = "user"
)

// Action define type for action
type Action string

// Define any action that may be used against any resource
const (
	ActionCreateAny Action = "create_any"
	ActionReadAny   Action = "read_any"
	ActionUpdateAny Action = "update_any"
	ActionDeleteAny Action = "delete_any"

	// ActionReadAll differ from ActionReadAny. This Action specifically used
	// to indicate an action to view all the resources. e.g. seeing all absent form
	ActionReadAll Action = "read_all"
	ActionInvite  Action = "create_invitation"
)

// ResourceAction represent a pair of resource and action can be performed
type ResourceAction struct {
	Resource Resource
	Action   Action
}

var _permissions = map[ResourceAction][]Role{
	{ResourceAbsentForm, ActionCreateAny}: {RoleAdmin},
	{ResourceAbsentForm, ActionReadAll}:   {RoleAdmin},
	{ResourceAbsentForm, ActionUpdateAny}: {RoleAdmin},

	{ResourceAbsentList, ActionCreateAny}: {RoleAdmin, RoleMember},

	{ResourceUser, ActionInvite}: {RoleAdmin},
}

// HasAccess detect if the given user role have access / permission in the supplied resource and action
func HasAccess(userRole Role, resource Resource, action Action) bool {
	if string(userRole) == "" {
		return false
	}

	roles := _permissions[ResourceAction{
		Resource: resource,
		Action:   action,
	}]

	for _, role := range roles {
		if role == userRole {
			return true
		}
	}

	return false
}

// ParseStringToRole parse given string to it's equivalent rbac.Role
func ParseStringToRole(str string) (Role, error) {
	role := strings.ToUpper(str)
	switch role {
	default:
		return RoleMember, errors.New("invalid role")
	case RoleMember:
		return RoleMember, nil
	case RoleAdmin:
		return RoleAdmin, nil
	}
}
