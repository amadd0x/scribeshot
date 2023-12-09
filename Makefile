NAME=scribeshot
VERSION=v0.1.0

build: clean
	CGO_ENABLED=0 go build -o build/$(NAME)

clean:
	rm -rf build

run: build
	./build/$(NAME)