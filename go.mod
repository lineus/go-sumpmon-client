module github.com/lineus/go-sump-client

go 1.15

replace github.com/lineus/go-sumpmon => ../../src/github.com/lineus/go-sumpmon

replace github.com/lineus/go-sqlitelogs => ../../src/github.com/lineus/go-sqlitelogs

replace github.com/lineus/go-notify => ../../src/github.com/lineus/go-notify

replace github.com/lineus/go-loadaws => ../../src/github.com/lineus/go-loadaws

require (
	github.com/lineus/go-loadaws v0.0.0-00010101000000-000000000000
	github.com/lineus/go-notify v0.0.0-00010101000000-000000000000
	github.com/lineus/go-sqlitelogs v0.0.0-00010101000000-000000000000
	github.com/lineus/go-sumpmon v0.0.0-00010101000000-000000000000
)
