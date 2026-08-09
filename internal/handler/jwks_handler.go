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
			WriteSystemError(w, http.StatusInternalServerError, err)
			return
		}
		if err := key.Set(jwk.KeyIDKey, "onerep-auth-signing-key"); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, err)
			return
		}
		if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, err)
			return
		}

		ks := jwk.NewSet()
		if err := ks.AddKey(key); err != nil {
			WriteSystemError(w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(w, http.StatusOK, ks)
	}
}
