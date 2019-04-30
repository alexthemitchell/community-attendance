GO_TEST_FLAGS := -v -race -cover -count=1 
test:
	go test $(GO_TEST_FLAGS) ./...

ci_test:
	go test $(GO_TEST_FLAGS) -coverprofile=coverage.txt -covermode=atomic ./...

cli:
	make -C cli binary
	mv cli/bin/* ./bin/

.PHONY: cli test
