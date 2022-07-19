package rbac

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

const (
	ResourceAbsentForm = "absent_form"
	ResourceAbsentList = "absent_list"
)

type Action string

const (
	ActionCreateAny = "create_any"
	ActionReadAny   = "view_any"
	ActionUpdateAny = "update_any"
	ActionDeleteAny = "delete_any"
)

type ResourceAction struct {
	Resource Resource
	Action   Action
}

var _permissions = map[ResourceAction][]Role{
	{ResourceAbsentForm, ActionCreateAny}: {RoleAdmin},
	{ResourceAbsentList, ActionCreateAny}: {RoleAdmin, RoleMember},
}

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
