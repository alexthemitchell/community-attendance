GO_TEST_FLAGS := -v -cover -count=1 
GO_DEEP_TEST_FLAGS := $(GO_TEST_FLAGS) -race
GO_TEST_TARGET := ./...
test:
	go test $(GO_TEST_FLAGS) $(GO_TEST_TARGET)

deep_test:
	go test $(GO_DEEP_TEST_FLAGS) $(GO_TEST_TARGET)

ci_test:
	go test $(GO_DEEP_TEST_FLAGS) -coverprofile=coverage.txt -covermode=atomic $(GO_TEST_TARGET)

cli:
	make -C cli binary
	mv cli/bin/* ./bin/

.PHONY: cli test
