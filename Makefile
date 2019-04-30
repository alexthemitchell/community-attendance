test:
	go test -v -cover -count=1 ./...

cli:
	make -C cli binary
	mv cli/bin/* ./bin/

.PHONY: cli test
