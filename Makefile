
all: rssd

rssd:
	go build -o $@

clean:
	rm rssd.db rssd

