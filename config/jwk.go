package config

import (
	"fmt"
	"net/http"
)

const JWK_VALUE = "MDNGMjU2M0U3RERFQUEwOUUzQUMwQ0NBN0Y1RUY0OEIxNTRDM0IxMw"

func JwkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	{
  "keys": [
    {
      "alg": "HS256",
      "kid": %s
    }
  ]
}
	`, JWK_VALUE)))
}
