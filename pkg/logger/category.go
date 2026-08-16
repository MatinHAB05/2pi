package logger

type Category string
type SubCategory string
type ExtraKey string

// Categories
const (
	General    Category = "General"
	IO         Category = "IO"
	Internal   Category = "Internal"
	Postgres   Category = "Postgres"
	Redis      Category = "Redis"
	Validation Category = "Validation"
)

// SubCategories
const (
	// General
	Startup SubCategory = "Startup"

// Postgres

// Internal

// Validation

// IO

// Redis

)

// Extra Keys
const (
	// System & Metadata

	// Request & Response

	// User & Business Context

	// Error Details
	ErrorMessage ExtraKey = "ErrorMessage"
)
