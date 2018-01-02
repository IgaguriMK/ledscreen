
.PHONY: build
build:
	go build viewchar.go
	go build runeid.go
	go build genimage.go

.PHONY: deps
deps:
	go get github.com/go-yaml/yaml
