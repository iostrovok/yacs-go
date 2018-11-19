package myconst

const (
	// DocKeyName "@doc": At any level, it is container for processing instructions.
	DocKeyName string = "@doc"
	// JSONRefKeyName "$ref": At any level, it is reference for includes or schemas.
	JSONRefKeyName string = "$ref"
	// LockKeyName "@lock_names": At any level, don't allow values defined at this level to be overwritten.
	LockKeyName string = "@lock_names"
	// ParentKeyName "@parent": Treat this sub-structure as a parent and the containing structure as overrrides.
	ParentKeyName string = "@parent"
	// ResolveKeyName "resolve" is used "@doc".
	ResolveKeyName string = "resolve"
	// SchemaKeyName "@schemas": At any level, it is object with references to schemas.
	SchemaKeyName string = "@schemas"
)
