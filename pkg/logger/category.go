package logger

type Category string
type SubCategory string
type ExtraKey string

// Categories
const (
	Panic      Category = "Panic"
	General    Category = "General"
	IO         Category = "IO"
	Internal   Category = "Internal"
	Postgres   Category = "Postgres"
	Redis      Category = "Redis"
	Validation Category = "Validation"
	Service    Category = "Service"
)

// SubCategories
const (
	MPanic          SubCategory = "Recovery-Middleware"
	LoggerM         SubCategory = "Logger-Middleware"
	ExternalService SubCategory = "External"

	// General
	Startup SubCategory = "Startup"

	// Service
	UserService          SubCategory = "UserService"
	TargetAccountService SubCategory = "TargetAccountService"
	RBACService          SubCategory = "RBACService"
)

// Extra Keys
const (
	// Business Context
	UserID          ExtraKey = "UserID"
	TargetAccountID ExtraKey = "TargetAccountID"
	OwnerID         ExtraKey = "OwnerID"
	Role            ExtraKey = "Role"
	Action          ExtraKey = "Action"
	Username        ExtraKey = "Username"
	Limit           ExtraKey = "Limit"
	Offset          ExtraKey = "Offset"

	// Error Details
	ErrorMessage ExtraKey = "ErrorMessage"
)
