GOCMD   =  go
GOBUILD =  $(GOCMD)  build
GORUN = $(GOCMD) run
TINYGOCMD   =  tinygo
TINYGOBUILD = $(TINYGOCMD) build

compile-wasm:export GOOS=js
compile-wasm:export GOARCH=wasm
compile-wasm:
	$(GOBUILD) -o main.wasm && rm ./html/main.wasm && cp ./main.wasm ./html/main.wasm
compile-wasm:
	$(GOBUILD) -o main.wasm && rm ./html/main.wasm && cp ./main.wasm ./html/main.wasm	
compile-wasm-tinygo:export GOOS=js
compile-wasm-tinygo:export GOARCH=wasm
compile-wasm-tinygo:
	$(TINYGOBUILD) -o main.wasm && 	rm ./html/main.wasm	&& cp ./main.wasm ./html/main.wasm
run_server:
	cd html && $(GORUN) server.go	