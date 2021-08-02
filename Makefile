# GOLOOP
#
# Load data about the package:
#	NAME
#		Name of the GoLang's module;
#	VERSION
#		Current version.
# 	REPOSITORY
# 		The name of the repository where the package is stored,
# 		for example: github.com/goloop;
MODULE_NAME:=$(shell cat go.mod | grep module | awk '{split($$2,v,"/"); print v[3]}')
MODULE_VERSION:=$(shell cat doc.go | grep "const version" | awk '{gsub(/"/, "", $$4); print $$4}')
REPOSITORY_NAME:=$(shell cat go.mod | grep module | awk '{split($$2,v,"/"); print v[1] "/" v[2]}')
 
# Help information.
define MSG_HELP
Go-package's manager

Commands:
	help
		Show this help information
	go.test
		Run tests
	go.test.cover
		Check test coverage
	go.lint
		Check code with GoLints
		 
		Requires `golangci-lint`, install as:
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0
	readme
		Create readme from the GoLang code
		
		Requires `godocdown`, install as:
		go get github.com/robertkrimen/godocdown/godocdown
	git.commit
		Update readme, create commit and update tag from the .module file.

		Usage as: make git.commit am="Commit message"
endef

# Constants.
export MSG_HELP
REPOSITORY_PATH=${REPOSITORY_NAME}/${MODULE_NAME}

all: help
help:
	@echo "$$MSG_HELP"
go.test:
	@go clean -testcache; go test ${REPOSITORY_PATH}
go.test.cover:
	@go test -cover ${REPOSITORY_PATH} && \
		go test -coverprofile=/tmp/coverage.out ${REPOSITORY_PATH} && \
		go tool cover -func=/tmp/coverage.out && \
		go tool cover -html=/tmp/coverage.out
go.lint:
ifeq (, $(shell which golangci-lint))
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0
endif
	golangci-lint run --no-config --issues-exit-code=0 --timeout=30m \
		--disable-all --enable=deadcode  --enable=gocyclo --enable=golint \
		--enable=varcheck --enable=structcheck --enable=maligned \
		--enable=gosec --enable=megacheck --enable=ineffassign \
		--enable=interfacer --enable=unconvert \
		--enable=goconst #--enable=errcheck --enable=dupl
readme:
ifeq (, $(shell which godocdown))
	@go get github.com/robertkrimen/godocdown/godocdown
endif
	@godocdown -plain=true -template=.godocdown.md ./ | \
		sed -e 's/\.ModuleVersion/v${MODULE_VERSION}/g' > README.md
git.commit: readme
ifeq ($(am),)
	@echo "You must provide a message to commit as: make commit am='Commit message'"
else
	@git add . && git commit -am "${am}" && \
		git tag v${MODULE_VERSION} && \
		git push -u origin --all && \
		git push -u origin --tag
endif