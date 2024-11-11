build: clean init cp_client
	go build -o ./dist src/main.go

run: export HOSTNAME=0.0.0.0:3000
run: export LOG_LEVEL=debug
run: build
	dist/main

clean:
	rm -rf ./dist

init:
	mkdir -p ./dist/static

cp_client:
	cp ../client/src/* ./dist/static

