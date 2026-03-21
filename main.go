package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

type Feed struct {
	Title string
	Name string
	Description string
	Url string
	Items []Item
}

type Item struct {
	Title string
	Updated string
	Content string
	Read bool
	Url string
}

func main() {
	c := getConfig();

	config := c.Config;

	var listener net.Listener;
	var err error;

	if config.Method == "unix" {
		os.Remove(config.UnixPath);

		listener, err = net.Listen("unix", config.UnixPath);
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on listening on unix: %v\n", err);
			os.Exit(1);
		}
	}

	if config.Method == "tcp" {
		port := fmt.Sprintf(":%v", config.TcpPort);
		listener, err = net.Listen("tcp", port);
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on listening on tcp with port: %v: %v\n", port, err);
			os.Exit(1);
		}
	}

	defer listener.Close();

	db, err := SqlConnect();
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on connecting to the database: %v\n", err);
		os.Exit(1);
	}

	err = SqlCreateTablesIfNotExists(db);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on creating tables: %v\n", err);
		os.Exit(1);
	}

	db.Close();

	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on accepting: %v\n", err);
		}

		go handleConnection(conn);
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close();

	reader := bufio.NewReader(conn);
	command, err := reader.ReadString('\n');
	if err != nil {
		if err == io.EOF {
			return;
		}
		fmt.Fprintf(os.Stderr, "ERROR: on reading the command: %v\n", err);
		return;
	}

	res := parseCommand(command);

	b, err := json.MarshalIndent(res, "", "    ");
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on marshal the response json: %v\n", err);
		return;	
	}

	fmt.Fprintf(conn, string(b));
}


