module github.com/zacscoding/echo-gorm-realworld-app-it

go 1.16

require (
	github.com/gavv/httpexpect/v2 v2.3.1
	github.com/google/uuid v1.2.0
	github.com/spf13/cast v1.3.0
	github.com/zacscoding/echo-gorm-realworld-app v0.0.0-00010101000000-000000000000
)

replace github.com/zacscoding/echo-gorm-realworld-app => ../
