package main

import (
	"os"
	"strings"
	"fmt"
	"strconv"
)

type FeedFile struct {
	Name string
	Url string
}

type ConfigFile struct {
	Method string
	UnixPath string
	TcpPort int
	QueryLimit int64
	ReloadTime int
}

type Config struct {
	Config ConfigFile
	Feeds []FeedFile
}

func getConfig() Config {
	config, err := parseConfig();
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error());
		os.Exit(1);
	}
	feeds, err := parseFeed();
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error());
		os.Exit(1);
	}

	return Config{config, feeds};
}

func parseConfig() (ConfigFile, error) {
	// change to $XDG_CONFIG_HOME/rssd
	var path = "./examples/config";

	cont, err := readFile(path);
	if err != nil {
		return ConfigFile{}, err;
	}

	var c ConfigFile;
	lines := strings.Split(cont, "\n");
	for _,l := range lines {
		l = strings.TrimSpace(l);
		if l != "" && !strings.HasPrefix(l, "//") {
			parts := strings.SplitN(l, ":", 2);

			key := strings.TrimSpace(parts[0]);
			value := strings.TrimSpace(parts[1]);

			switch key {
			case "method":
				c.Method = value;

			case "unix-path":
				if c.Method != "unix" {
					return ConfigFile{}, fmt.Errorf("ERROR: unix-path requires the method to be unix, method: `%v` is not compatible", c.Method);
				}
				c.UnixPath = value
				c.TcpPort = 0

			case "tcp-port":
				if c.Method != "tcp" {
					return ConfigFile{}, fmt.Errorf("ERROR: tcp-port requires the method to be tcp, method: `%v` is not compatible", c.Method);
				}
				converted, err := strconv.Atoi(value);
				if err != nil {
					return ConfigFile{}, err;
				}
				c.TcpPort = converted;
				c.UnixPath = "";

			case "default-query-limit":
				converted, err := strconv.ParseInt(value, 10, 64);
				if err != nil {
					return ConfigFile{}, err;
				}
				c.QueryLimit = converted;

			case "reload-time":
				converted, err := strconv.Atoi(value);
				if err != nil {
					return ConfigFile{}, err;
				}
				c.ReloadTime = converted;

			default:
				return ConfigFile{}, fmt.Errorf("ERROR: unknown option: `%v`", key);
			}
		}
	}

	return c, nil;
}

func parseFeed() ([]FeedFile, error) {
	// change to $XDG_CONFIG_HOME/rssd/
	var path = "./examples/feeds";

	cont, err := readFile(path);
	if err != nil {
		return nil, err;
	}

	var feeds []FeedFile;

	lines := strings.Split(cont, "\n");
	for _,l := range lines {
		if l != "" {
			var f FeedFile;
			parts := strings.SplitN(l, " ", 2);
			f.Name = parts[0];
			f.Url = parts[1];
			feeds = append(feeds, f);
		}
	}

	return feeds, nil;
}

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path);
	if err != nil {
		return "", err;
	}
	return string(dat), nil;
}

