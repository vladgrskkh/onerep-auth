package crypto_test

import (
	"testing"

	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHasher_HashAndCompare(t *testing.T) {
	h := crypto.NewPasswordHasher()
	hash, err := h.Hash("mypassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = h.Compare(hash, "mypassword")
	assert.NoError(t, err)

	err = h.Compare(hash, "wrongpassword")
	assert.Error(t, err)
}
