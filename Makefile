
SRCS=$(wildcard *go)

all: rssd

rssd: $(SRCS)
	go build -o $@

clean:
	rm rssd.db rssd

