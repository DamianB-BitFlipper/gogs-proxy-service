module gogs.io/main

go 1.14

replace gogs => ./api_client

require (
	github.com/gorilla/mux v1.8.0
	gogs v0.0.0-00010101000000-000000000000
	unknwon.dev/clog/v2 v2.2.0
)
