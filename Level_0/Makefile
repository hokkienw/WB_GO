.PHONY: all sub pub sub_sc tidy_sc pub_sc
all:

sub:
	clear
	go run cmd/sub/main.go

pub:
	clear
	go run cmd/pub/main.go

sub_sc:
	clear
	/opt/goinfre/hokkienw/brew/bin/go run cmd/sub/main.go

tidy_sc:
	clear
	/opt/goinfre/hokkienw/brew/bin/go mod init level_0
	/opt/goinfre/hokkienw/brew/bin/go mod tidy

	

pub_sc:
	clear
	/opt/goinfre/hokkienw/brew/bin/go run cmd/pub/main.go
