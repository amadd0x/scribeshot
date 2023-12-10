NAME=scribeshot
VERSION=v0.1.1

build: clean
	CGO_ENABLED=0 go build -o build/$(NAME)

clean:
	rm -rf build

run: build
	./build/$(NAME)

release:
	goreleaser release --clean