package context

const UserIDKey contextKey = "user_id"
const AccountIDKey contextKey = "account_id"
const AccountOwnerIDKey contextKey = "account_owner_id"

type AuthenticationContextToken struct {
	UserId         int64
	AccountID      *int64
	AccountOwnerID *int64
}

//TODO : get Token function for context
