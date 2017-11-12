all:
	go build

test:
	go test

coverage:
	go test -cover
	go test -coverprofile=coverage.out 
	go tool cover -func=coverage.out
	# go tool cover -html=coverage.out

clean:
	rm -f stupidcoin