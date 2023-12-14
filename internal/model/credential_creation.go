package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-webauthn/webauthn/protocol"
)

func MarshalCredentialCreation(c protocol.CredentialCreation) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		b, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		if _, err := w.Write(b); err != nil {
			log.Println(err)
			panic(err)
		}
	})
}

func UnmarshalCredentialCreation(v interface{}) (protocol.CredentialCreation, error) {
	return protocol.CredentialCreation{}, fmt.Errorf("%T is not a parsable", v)
}
