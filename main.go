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
	os.Remove("/tmp/rssd.sock");

	// TODO: add option to choose between unixsockets and tcp
	listener, err := net.Listen("unix", "/tmp/rssd.sock");
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on listening: %v\n", err);
		os.Exit(1);
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
			fmt.Fprintf(os.Stderr, "ERROR: on listening: %v\n", err);
			os.Exit(1);
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

