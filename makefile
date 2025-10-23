# Makefile
test:
	go run . james

testDebug:
	go run . -debug james 

make:
	go build .

install:
	go install .

clean:
	rm -f random-discogs-item