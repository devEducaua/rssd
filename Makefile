
TARGET = rssd
TARGETDIR = bin

SRCS := $(wildcard *.go)

PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
CONFIGDIR = $(HOME)/.config/rssd

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$@ .

install:
	cp ./examples/* $(CONFIGDIR)/
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/

clean:
	rm -rf rssd.db $(TARGETDIR)

