package jwt_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	infrajwt "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManager_IssueAndValidate(t *testing.T) {
	key, err := infrajwt.GenerateKeyPair()
	require.NoError(t, err)

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := infrajwt.NewTokenManager(string(privPEM), 15*time.Minute)
	require.NoError(t, err)

	userID := uuid.Must(uuid.NewV7())
	email := "test@example.com"

	pair, err := tm.IssueTokenPair(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, 900, pair.ExpiresIn)

	id, err := infrajwt.GetUserIDFromToken(pair.AccessToken, tm.PublicKeyPEM())
	require.NoError(t, err)
	assert.Equal(t, userID, id)
}
