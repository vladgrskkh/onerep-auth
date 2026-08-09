package handler

import (
	"crypto/rsa"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func JWKSHandler(publicKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		key, err := jwk.FromRaw(publicKey)
		if err != nil {
			WriteSystemError(w, http.StatusInternalServerError, "JWK_CREATE_FAILED")
			return
		}
		if err := key.Set(jwk.KeyIDKey, "onerep-auth-signing-key"); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, "JWK_KEY_ID_FAILED")
			return
		}
		if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, "JWK_ALGORITHM_FAILED")
			return
		}

		ks := jwk.NewSet()
		if err := ks.AddKey(key); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, "JWK_SET_FAILED")
			return
		}

		WriteJSON(w, http.StatusOK, ks)
	}
}
