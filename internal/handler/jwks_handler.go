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
			WriteError(w, http.StatusInternalServerError, "failed to create JWK")
			return
		}
		err = key.Set(jwk.KeyIDKey, "gym-auth-signing-key")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to set key ID")
			return
		}
		err = key.Set(jwk.AlgorithmKey, "RS256")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to set algorithm")
			return
		}

		ks := jwk.NewSet()
		err = ks.AddKey(key)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to add key to set")
			return
		}

		WriteJSON(w, http.StatusOK, ks)
	}
}
