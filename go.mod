module gogs.io/main

go 1.14

require (
	github.com/cloudflare/cfssl v1.5.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fxamacker/cbor/v2 v2.2.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	gogs v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	unknwon.dev/clog/v2 v2.2.0
	webauthn v0.0.0-00010101000000-000000000000
)

replace gogs => ./api_client

replace webauthn => ./webauthn
