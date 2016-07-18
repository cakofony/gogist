all:
	mkdir -p build
	go build -o build/gist gist.go
clean:
	rm -rf build
install: all
	sudo cp build/gist /usr/local/bin/gist
uninstall:
	sudo rm -f /usr/local/bin/gist
