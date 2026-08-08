package crypto_test

import (
	"testing"

	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHasher_HashAndCompare(t *testing.T) {
	h := crypto.NewPasswordHasher()
	hash, err := h.Hash("mypassword")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = h.Compare(hash, "mypassword")
	require.NoError(t, err)

	err = h.Compare(hash, "wrongpassword")
	assert.Error(t, err)
}
