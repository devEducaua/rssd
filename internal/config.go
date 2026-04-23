package internal

import (
	"os"
	"strings"
	"fmt"
	"strconv"
	"path/filepath"
)

type FeedConfig struct {
	Name string
	Url string
}

type Config struct {
	Method string
	UnixPath string
	TcpPort int
	QueryLimit int64
	ReloadTime int
}

func GetConfig() (Config, error){
	config := Config{
		Method: "unix",
		UnixPath: "/tmp/rssd.sock",
		QueryLimit: 100,
		ReloadTime: 900,
	}

	configFile, err := parseConfigFile();
	if err != nil {
		return config, err;
	}

	if configFile.Method == "tcp" {
		config.Method = "tcp";
		config.TcpPort = configFile.TcpPort;
		config.UnixPath = "";
	}

	if configFile.QueryLimit != 0 {
		config.QueryLimit = configFile.QueryLimit;
	}

	if configFile.ReloadTime != 0 {
		config.ReloadTime = configFile.ReloadTime;
	}

	return config, nil;
}


func parseConfigFile() (Config, error) {
	baseDir, err := getBaseDir();
	if err != nil {
		return Config{}, err;
	}

	var path = filepath.Join(baseDir, "config");

	cont, err := readFile(path);
	if err != nil {
		return Config{}, fmt.Errorf("failed to read the config file: %v", err);
	}

	var c Config;
	lines := strings.Split(cont, "\n");
	for _,l := range lines {
		l = strings.TrimSpace(l);
		if l != "" && !strings.HasPrefix(l, "#") {
			parts := strings.SplitN(l, ":", 2);

			key := strings.TrimSpace(parts[0]);
			value := strings.TrimSpace(parts[1]);

			switch key {
			case "method":
				c.Method = value;

			case "unix-path":
				if c.Method != "unix" {
					return Config{}, fmt.Errorf("ERROR: unix-path requires the method to be unix, method: `%v` is not compatible", c.Method);
				}
				c.UnixPath = value
				c.TcpPort = 0

			case "tcp-port":
				if c.Method != "tcp" {
					return Config{}, fmt.Errorf("ERROR: tcp-port requires the method to be tcp, method: `%v` is not compatible", c.Method);
				}
				converted, err := strconv.Atoi(value);
				if err != nil {
					return Config{}, err;
				}
				c.TcpPort = converted;
				c.UnixPath = "";

			case "query-limit":
				converted, err := strconv.ParseInt(value, 10, 64);
				if err != nil {
					return Config{}, err;
				}
				c.QueryLimit = converted;

			case "reload-time":
				converted, err := strconv.Atoi(value);
				if err != nil {
					return Config{}, err;
				}
				c.ReloadTime = converted;

			default:
				return Config{}, fmt.Errorf("ERROR: unknown option: `%v`", key);
			}
		}
	}

	return c, nil;
}


func getFeedsConfig() ([]FeedConfig, error) {
	feeds, err := parseFeedsFile();
	if err != nil {
		return nil, fmt.Errorf("could not read the feeds file: %v\n", err);
	}

	return feeds, nil;
}

func parseFeedsFile() ([]FeedConfig, error) {
	baseDir, err := getBaseDir();
	if err != nil {
		return nil, err;
	}

	var path = filepath.Join(baseDir, "feeds");

	cont, err := readFile(path);
	if err != nil {
		return nil, err;
	}

	var feeds []FeedConfig
	lines := strings.SplitSeq(cont, "\n");
	for l := range lines {
		if l != "" {
			var f FeedConfig;
			parts := strings.SplitN(l, " ", 2);
			f.Name = parts[0];
			f.Url = parts[1];
			feeds = append(feeds, f);
		}
	}

	return feeds, nil;
}

func getBaseDir() (string, error) {
	home, err := os.UserHomeDir();
	if err != nil {
		return "", err;
	}
	path := filepath.Join(home, ".config", "rssd");

	err = os.MkdirAll(path, 0755);
	if err != nil {
		return "", err;
	}

	return path, nil;
}

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path);
	if err != nil {
		return "", err;
	}
	return string(dat), nil;
}

