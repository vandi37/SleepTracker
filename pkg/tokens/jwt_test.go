package tokens_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vandi37/SleepTracker/pkg/tokens"
)

func TestTokenService_Generate(t *testing.T) {
	secret := "test-secret"
	expiration := 1 * time.Hour
	service := tokens.New(secret, expiration)

	tests := []struct {
		name       string
		sub        string
		data       []tokens.Data
		wantError  bool
		wantClaims func(*testing.T, jwt.MapClaims)
	}{
		{
			name: "Simple token with subject",
			sub:  "user123",
			wantClaims: func(t *testing.T, claims jwt.MapClaims) {
				assert.Equal(t, "user123", claims["sub"])
				assert.Equal(t, []interface{}{"https://github.com/vandi37/Calculator"}, claims["aud"])
				assert.NotEmpty(t, claims["iat"])
				assert.NotEmpty(t, claims["exp"])
				assert.NotEmpty(t, claims["nbf"])
				assert.Equal(t, "https://github.com/vandi37/Calculator", claims["iss"])
			},
		},
		{
			name: "Token with additional data",
			sub:  "user456",
			data: []tokens.Data{
				{Key: "role", Value: "admin"},
				{Key: "email", Value: "user@example.com"},
			},
			wantClaims: func(t *testing.T, claims jwt.MapClaims) {
				assert.Equal(t, "user456", claims["sub"])
				assert.Equal(t, "admin", claims["role"])
				assert.Equal(t, "user@example.com", claims["email"])
			},
		},
		{
			name:      "Empty subject",
			sub:       "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := service.Generate(tt.sub, tt.data...)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenStr)

			// Verify the token can be parsed and has expected claims
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)

			claims, ok := token.Claims.(jwt.MapClaims)
			require.True(t, ok)

			if tt.wantClaims != nil {
				tt.wantClaims(t, claims)
			}

			// Verify timings
			now := time.Now()
			iat := time.Unix(int64(claims["iat"].(float64)), 0)
			exp := time.Unix(int64(claims["exp"].(float64)), 0)

			assert.WithinDuration(t, now, iat, time.Second)
			assert.WithinDuration(t, now.Add(expiration), exp, time.Second)
		})
	}
}

func TestTokenService_Parse(t *testing.T) {
	secret := "test-secret"
	service := tokens.New(secret, 1*time.Hour)

	t.Run("Valid token", func(t *testing.T) {
		tokenStr, err := service.Generate("user123")
		require.NoError(t, err)

		token, err := service.Parse(tokenStr)
		require.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)
		assert.Equal(t, "user123", claims["sub"])
	})

	t.Run("Invalid signature", func(t *testing.T) {
		tokenStr, err := service.Generate("user123")
		require.NoError(t, err)

		// Create another service with different secret
		otherService := tokens.New("different-secret", 1*time.Hour)
		_, err = otherService.Parse(tokenStr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Malformed token", func(t *testing.T) {
		_, err := service.Parse("not.a.valid.token")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token contains an invalid number of segments")
	})

	t.Run("Expired token", func(t *testing.T) {
		// Create service with very short expiration
		shortService := tokens.New(secret, -1*time.Hour)
		tokenStr, err := shortService.Generate("user123")
		require.NoError(t, err)

		_, err = service.Parse(tokenStr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

}

func TestNewTokenService(t *testing.T) {
	t.Run("Default parameters", func(t *testing.T) {
		secret := "test-secret"
		expiration := 1 * time.Hour

		service := tokens.New(secret, expiration,)
		assert.Equal(t, []byte(secret), service.GetSecret())
		assert.Equal(t, expiration, service.GetExpiration())
	})

	t.Run("Zero expiration", func(t *testing.T) {
		service := tokens.New("secret", 0)
		assert.Equal(t, time.Duration(0), service.GetExpiration())
	})
}

func TestTokenService_EdgeCases(t *testing.T) {
	t.Run("Empty secret", func(t *testing.T) {
		service := tokens.New("", 1*time.Hour)
		_, err := service.Generate("user123")
		require.NoError(t, err) // tokens allows empty secret
	})

	t.Run("Very long expiration", func(t *testing.T) {
		service := tokens.New("secret", 365*24*time.Hour)
		tokenStr, err := service.Generate("user123")
		require.NoError(t, err)

		token, err := service.Parse(tokenStr)
		require.NoError(t, err)
		assert.True(t, token.Valid)
	})

	t.Run("Multiple custom claims", func(t *testing.T) {
		service := tokens.New("secret", 1*time.Hour)
		data := []tokens.Data{
			{Key: "name", Value: "John Doe"},
			{Key: "age", Value: 30},
			{Key: "admin", Value: true},
			{Key: "permissions", Value: []string{"read", "write"}},
		}

		tokenStr, err := service.Generate("user123", data...)
		require.NoError(t, err)

		token, err := service.Parse(tokenStr)
		require.NoError(t, err)

		claims := token.Claims.(jwt.MapClaims)
		assert.Equal(t, "John Doe", claims["name"])
		assert.Equal(t, float64(30), claims["age"]) // JSON numbers are float64
		assert.Equal(t, true, claims["admin"])
		assert.Equal(t, []interface{}{"read", "write"}, claims["permissions"])
	})
}
