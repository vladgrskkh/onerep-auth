package handler

import (
	"crypto/rsa"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func JWKSHandler(publicKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := jwk.FromRaw(publicKey)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create JWK")
			return
		}
		key.Set(jwk.KeyIDKey, "gym-auth-signing-key")
		key.Set(jwk.AlgorithmKey, "RS256")

		ks := jwk.NewSet()
		ks.AddKey(key)

		writeJSON(w, http.StatusOK, ks)
	}
}
