module github.com/jurgen-kluft/go-home/sensor-server

replace jurgen-kluft/go-home/sensor-server/gnet => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet

replace jurgen-kluft/go-home/sensor-server/gnet/internal/gfd => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/gfd

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/bs => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/bs

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/elastic => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/elastic

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/linkedlist => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/linkedlist

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/ring => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/buffer/ring

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/errors => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/errors

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/io => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/io

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/logging => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/logging

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/math => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/math

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/netpoll => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/netpoll

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/bytebuffer => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/bytebuffer

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/byteslice => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/byteslice

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/goroutine => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/goroutine

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/ringbuffer => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/pool/ringbuffer

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/queue => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/queue

replace jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/socket => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/gnet/internal/pkg/socket

replace jurgen-kluft/go-home/sensor-server/ants => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/ants

replace jurgen-kluft/go-home/sensor-server/ants/pkg/sync => /Users/obnosis5/dev.go/src/github.com/jurgen-kluft/go-home/sensor-server/ants/pkg/sync

go 1.23.0

toolchain go1.23.8

require (
	github.com/stretchr/testify v1.10.0
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.14.0
	golang.org/x/sys v0.33.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
