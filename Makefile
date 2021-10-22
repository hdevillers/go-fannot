INSTALL_DIR = /usr/local/bin

ifdef prefix
	INSTALL_DIR = $(prefix)
endif

build:
	go build -o bin/swiss-count ./cmd/swiss-count/main.go
	go build -o bin/swiss-subset ./cmd/swiss-subset/main.go
	go build -o bin/swiss-create-refdb ./cmd/swiss-create-refdb/main.go
	go build -o bin/fannot-run ./cmd/fannot-run/main.go
	go build -o bin/refdb-info ./cmd/refdb-info/main.go

test:
	go test -v swiss/swiss_test.go

install:
	cp bin/swiss-count $(INSTALL_DIR)/swiss-count
	cp bin/swiss-subset $(INSTALL_DIR)/swiss-subset
	cp bin/swiss-create-refdb $(INSTALL_DIR)/swiss-create-refdb
	cp bin/fannot-run $(INSTALL_DIR)/fannot-run
	cp bin/refdb-info $(INSTALL_DIR)/refdb-info

uninstall:
	rm -f $(INSTALL_DIR)/swiss-count
	rm -f $(INSTALL_DIR)/swiss-subset
	rm -f $(INSTALL_DIR)/swiss-create-refdb
	rm -f $(INSTALL_DIR)/fannot-run
	rm -f $(INSTALL_DIR)/refdb-info