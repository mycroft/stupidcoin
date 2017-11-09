all:
	go build

test:
	# go test
	./stupidcoin -create-key
	./stupidcoin -create-key
	./stupidcoin -list-keys

clean:
	rm -f stupidcoin