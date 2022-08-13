module game-genshin

go 1.18

require flswld.com/common v0.0.0-incompatible

replace flswld.com/common => ../../common

require flswld.com/logger v0.0.0-incompatible

replace flswld.com/logger => ../../logger

require flswld.com/air-api v0.0.0-incompatible // indirect

replace flswld.com/air-api => ../../air-api

require flswld.com/light v0.0.0-incompatible

replace flswld.com/light => ../../light

require (
	flswld.com/gate-genshin-api v0.0.0-incompatible
	google.golang.org/protobuf v1.28.0
)

replace flswld.com/gate-genshin-api => ../../gate-genshin-api

// mongodb
require go.mongodb.org/mongo-driver v1.8.3

// jwt
require github.com/golang-jwt/jwt/v4 v4.4.0

// csv
require github.com/jszwec/csvutil v1.7.1

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/text v0.3.5 // indirect
)
