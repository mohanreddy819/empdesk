package appcontext

type contextKey string

const (
	UserIDKey   = contextKey("userID")
	UserRoleKey = contextKey("role")
)
// we define them globally so as to 