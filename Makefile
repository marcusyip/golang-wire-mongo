.PHONY: generate test e2e

generate:
	@echo Run go generate
	go generate ./...

test:
	@echo Run go test
	go test ./...

e2e:
	@echo Run E2E test
	go test github.com/marcusyip/golang-wire-mongo/test/e2e -ginkgo.debug -ginkgo.failFast
