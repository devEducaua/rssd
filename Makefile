
SRCS=$(wildcard *.go)

all: rssd

rssd: $(SRCS)
	go build -o $@ .

client:
	go build -o client ./client/main.go

clean:
	rm rssd.db rssd
