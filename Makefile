.PHONY: test

fontFile=$(CURDIR)/resource/fonts/Hiragino_Sans_GB_W3.ttf
all: install

test:
	@AVATAR_FONT=$(fontFile) go test -v ./avatar
  
install:
	go install github.com/holys/initials-avatar/cmd/avatar


