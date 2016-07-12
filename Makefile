all:
	go build gist.go
clean:
	rm -rf gist
install: all
	sudo cp gist /usr/local/bin/gist
uninstall:
	sudo rm -f /usr/local/bin/gist
