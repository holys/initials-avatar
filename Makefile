all: install
  
install:
	go get ./...
	go install github.com/holys/initials-avatar/cmd/avatar


