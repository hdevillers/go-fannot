INSTALL_DIR = /usr/local/bin

ifdef prefix
	INSTALL_DIR = $(prefix)
endif

build:
	go build -o bin/uniprot-count ./cmd/uniprot-count/main.go
	go build -o bin/uniprot-subset ./cmd/uniprot-subset/main.go
	go build -o bin/uniprot-create-refdb ./cmd/uniprot-create-refdb/main.go
	go build -o bin/fasta-create-refdb ./cmd/fasta-create-refdb/main.go
	go build -o bin/uniprot-prune ./cmd/uniprot-prune/main.go
	go build -o bin/uniprot-split ./cmd/uniprot-split/main.go
	go build -o bin/uniprot-download ./cmd/uniprot-download/main.go
	go build -o bin/fannot-run ./cmd/fannot-run/main.go
	go build -o bin/refdb-info ./cmd/refdb-info/main.go

test:
	go test -v tools/tools_test.go
	go test -v refdb/refdb.go refdb/refdb_test.go
	go test -v fannot/rule.go fannot/rule_test.go
	go test -v fannot/besthit.go fannot/besthit_test.go
	go test -v fannot/fields.go fannot/fields_test.go
	go test -v fannot/fields.go fannot/description.go fannot/description_test.go
	go test -v fannot/fields.go fannot/description.go fannot/format.go fannot/format_test.go
	go test -v fannot/rule.go fannot/param.go fannot/param_test.go
	go test -v fannot/rule.go fannot/param.go fannot/fannot.go fannot/fields.go fannot/description.go fannot/format.go fannot/result.go fannot/besthit.go fannot/fannot_test.go
	

install:
	cp bin/uniprot-count $(INSTALL_DIR)/uniprot-count
	cp bin/uniprot-subset $(INSTALL_DIR)/uniprot-subset
	cp bin/uniprot-create-refdb $(INSTALL_DIR)/uniprot-create-refdb
	cp bin/fasta-create-refdb $(INSTALL_DIR)/fasta-create-refdb
	cp bin/uniprot-prune $(INSTALL_DIR)/uniprot-prune
	cp bin/uniprot-split $(INSTALL_DIR)/uniprot-split
	cp bin/uniprot-download $(INSTALL_DIR)/uniprot-download
	cp bin/fannot-run $(INSTALL_DIR)/fannot-run
	cp bin/refdb-info $(INSTALL_DIR)/refdb-info

uninstall:
	rm -f $(INSTALL_DIR)/uniprot-count
	rm -f $(INSTALL_DIR)/uniprot-subset
	rm -f $(INSTALL_DIR)/uniprot-create-refdb
	rm -f $(INSTALL_DIR)/fasta-create-refdb
	rm -f $(INSTALL_DIR)/uniprot-prune
	rm -f $(INSTALL_DIR)/uniprot-split
	rm -f $(INSTALL_DIR)/uniprot-download
	rm -f $(INSTALL_DIR)/fannot-run
	rm -f $(INSTALL_DIR)/refdb-info