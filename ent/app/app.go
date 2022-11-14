// Code generated by ent, DO NOT EDIT.

package app

const (
	// Label holds the string label denoting the app type in the database.
	Label = "app"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldProjectId holds the string denoting the projectid field in the database.
	FieldProjectId = "project_id"
	// FieldAppTypeID holds the string denoting the app_type_id field in the database.
	FieldAppTypeID = "app_type_id"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// EdgeAppsTypes holds the string denoting the apps_types edge name in mutations.
	EdgeAppsTypes = "apps_types"
	// Table holds the table name of the app in the database.
	Table = "apps"
	// AppsTypesTable is the table that holds the apps_types relation/edge.
	AppsTypesTable = "apps"
	// AppsTypesInverseTable is the table name for the AppType entity.
	// It exists in this package in order to avoid circular dependency with the "apptype" package.
	AppsTypesInverseTable = "apps_types"
	// AppsTypesColumn is the table column denoting the apps_types relation/edge.
	AppsTypesColumn = "app_type_id"
)

// Columns holds all SQL columns for app fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldProjectId,
	FieldAppTypeID,
	FieldActive,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultActive holds the default value on creation for the "active" field.
	DefaultActive bool
)
