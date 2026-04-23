PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1
CONFIGDIR = $(HOME)/.config/rssd

TARGET = rssd
TARGETDIR = bin
MANPAGE = rssd.1

SRCS := $(wildcard *.go)

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$@ .

install:
	cp ./examples/* $(CONFIGDIR)/
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/
	cp ./$(MANPAGE) $(MANDIR)/
	gzip -f $(MANDIR)/$(MANPAGE)

uninstall:
	rm $(BINDIR)/$(TARGET)
	rm -f $(MANDIR)/$(MANPAGE).gz

clean:
	rm -rf rssd.db $(TARGETDIR)

