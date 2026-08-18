package tokencontext

type contextKey string

type TokenContext struct {
	Authentication *AuthenticationContextToken
}
