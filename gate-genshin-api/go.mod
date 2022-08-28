module gate-genshin-api

go 1.18

require flswld.com/common v0.0.0-incompatible // indirect

replace flswld.com/common => ../common

require flswld.com/logger v0.0.0-incompatible

require github.com/BurntSushi/toml v0.3.1 // indirect

replace flswld.com/logger => ../logger

// protobuf
require google.golang.org/protobuf v1.28.0
