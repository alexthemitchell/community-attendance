test:
	go test -count=1 ./...

cli:
	make -eC cli binary
	mv cli/bin/* ./bin/

.PHONY: cli test
