package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"rssd/internal"
)


func main() {
	config, err := internal.GetConfig();
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
		os.Exit(1);
	}

	db, err := internal.SqlConnect();
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on connecting to the database: %v\n", err);
		os.Exit(1);
	}

	err = internal.SqlCreateTablesIfNotExists(db);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on creating tables: %v\n", err);
		os.Exit(1);
	}
	db.Close();

	// TODO: get from config
	go internal.PeriodicReload(15*60);

	var listener net.Listener;

	if config.Method == "unix" {
		os.Remove(config.UnixPath);

		var err error;
		listener, err = net.Listen("unix", config.UnixPath);
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on listening on unix: %v\n", err);
			os.Exit(1);
		}
	}

	if config.Method == "tcp" {
		var err error;

		port := fmt.Sprintf(":%v", config.TcpPort);
		listener, err = net.Listen("tcp", port);
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on listening on tcp with port: %v: %v\n", port, err);
			os.Exit(1);
		}
	}

	defer listener.Close();
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

	res := internal.ParseCommand(command);

	b, err := json.MarshalIndent(res, "", "    ");
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on marshal the response json: %v\n", err);
		return;	
	}

	fmt.Fprint(conn, string(b));
}

