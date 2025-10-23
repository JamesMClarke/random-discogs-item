# Makefile
make:
	go build .

test:
	go run . james

testDebug:
	go run . -debug james 

install:
	go install .

clean:
	rm -f random-discogs-item