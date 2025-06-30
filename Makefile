EXEC_NAME = goposc
BIN_DIR = bin

all: build

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(EXEC_NAME) cmd/tui/main.go

clean:
	rm  -rf bin/$(EXEC_NAME)
