package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Response interface{} `json:"response"`
}

func parseCommand(command string) Response {
	parts := strings.Split(command, " ");

	var r Response;

	var msg string;
	var items []ItemDB;
	var findIds []int64;

	var err error;

	r.Status = "yes"

	switch parts[0] {
		case "GET":
			items, err = getCommand(parts);
			r.Response = items;
		case "UPDATE":
			msg, err = updateCommand(parts);
			r.Response = msg;
		case "READ":
			msg, err = readCommand(parts);
			r.Response = msg;	
		case "UNREAD":
			msg, err = unreadCommand(parts);
			r.Response = msg;	
		case "DELETE":
			msg, err = deleteCommand(parts);
			r.Response = msg;
		case "FIND":
			findIds, err = findCommand(parts);
			r.Response = findIds;
		default:
			err = fmt.Errorf("command: %v doesn't exists", parts[0]);
	}

	if err != nil {
		r.Status = "no";
		r.Response = err;
	}

	return r;
}

func getCommand(command []string) ([]ItemDB, error) {
	if len(command) != 2 {
		return nil, fmt.Errorf("invalid syntax on the `UPDATE` command: `UPDATE` only accepts one argument");
	}

	var limit int64 = 100;
	if len(command) != 3 {
		limit, _ = strconv.ParseInt(command[2], 10, 64);
	}

	arg := command[1];

	db, err := SqlConnect();
	if err != nil {
		return nil, err;
	}	
	defer db.Close();

	var items []ItemDB;

	switch arg {
	case "ALL":
		items, err = SqlGetAllItems(db, limit);
	case "UNREAD":
		items, err = SqlGetItemsByRead(db, false, limit);
	case "READ":
		items, err = SqlGetItemsByRead(db, true, limit);
	default:
		items, err = SqlGetItemsByName(db, arg, limit);
	}

	if err != nil {
		return nil, err;
	}

	return items, nil;
}

func updateCommand(command []string) (string, error) {
	if len(command) != 2 {
		return "", fmt.Errorf("invalid syntax on the `UPDATE` command: `UPDATE` only accepts one argument");
	}

	// TODO: add feedname support
	db, err := SqlConnect();
	if err != nil {
		return "", err;
	}
	defer db.Close();



	return fmt.Sprintf("the database was updated"), nil;
}

func changeRead(stringId string, read bool) (int64, error) {
	id, err := strconv.ParseInt(stringId, 10, 64);
	if err != nil {
		return -1, err;
	}

	db, err := SqlConnect();
	if err != nil {
		return -1, err;
	}
	defer db.Close();

	err = SqlUpdateItemRead(db, id, read);
	if err != nil {
		return -1, err;
	}

	return id, nil;
}

func readCommand(command []string) (string, error) {
	if len(command) != 2 {
		return "", fmt.Errorf("invalid syntax on the `READ` command: `READ` only accepts one argument");
	}

	id, err := changeRead(command[1], true);
	if err != nil {
		return "", err;
	}

	return fmt.Sprintf("item with id: %v is read", id), nil;
}

func unreadCommand(command []string) (string, error) {
	if len(command) != 2 {
		return "", fmt.Errorf("invalid syntax on the `UNREAD` command: `UNREAD` only accepts one argument");
	}

	id, err := changeRead(command[1], false);
	if err != nil {
		return "", err;
	}

	return fmt.Sprintf("item with id: %v is unread", id), nil;
}

func deleteCommand(command []string) (string, error) {
	if len(command) != 2 {
		return "", fmt.Errorf("invalid syntax on the `DELETE` command: `DELETE` only accepts one argument");
	}

	db, err := SqlConnect();
	if err != nil {
		return "", err;
	}
	defer db.Close();

	url := command[1];

	err = SqlDeleteFeed(db, url);
	if err != nil {
		return "", err;
	}

	return fmt.Sprintf("feed with url: %v deleted", url), nil;
}

func findCommand(command []string) ([]int64, error) {
	if len(command) > 3 || len(command) < 2 {
		return nil, fmt.Errorf("invalid syntax on the `FIND` command: `FIND` only accepts two arguments");
	}

	var limit int64 = 3;
	if len(command) == 3 {
		limit, _ = strconv.ParseInt(command[2], 10, 64);
	}

	db, err := SqlConnect();
	if err != nil {
		return nil, err;
	}
	defer db.Close();

	ids, err := SqlSearchItem(db, command[1], limit);
	if err != nil {
		return nil, err;
	}

	return ids, nil;
}

