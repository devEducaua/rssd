package main

import (
	"os"
	"strings"
)

func parseFeedsFile(path string) (map[string]string, error) {

	dat, err := os.ReadFile(path);
	if err != nil {
		return nil, err;
	}
	m := make(map[string]string);

	contents := string(dat);	

	lines := strings.Split(contents, "\n");

	for _,l := range lines {
		if strings.TrimSpace(l) != "" {
			parts := strings.SplitN(l, " ", 2);
			m[parts[0]] = parts[1];
		}
	}

	return m, nil;
}
