package internal

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Response struct {
	Ok bool `json:"ok"`
	Data any `json:"data"`
	Error error `json:"error"`
}

type UpdateDataResponse struct {
	Updated int;
}

func ParseCommand(command string) Response {
	parts := strings.Split(command, " ");

	var r = Response{
		Ok: true,
		Data: nil,
		Error: nil,
	};

	var data any;
	var err error;

	switch parts[0] {
		case "GET":
			data, err = getCommand(parts);
		case "UPDATE":
			data, err = updateCommand(parts);
		case "READ":
			err = readCommand(parts);
			data = nil;
		case "UNREAD":
			err = unreadCommand(parts);
			data = nil;
		case "DELETE":
			err = deleteCommand(parts);
			data = nil;
		case "FIND":
			data, err = findCommand(parts);
		case "OPEN":
			err = openCommand(parts);
			data = nil;
		case "ADD":
			err = addCommand(parts);
			data = nil
		default:
			err = fmt.Errorf("command: %v doesn't exists", parts[0]);
	}
	if err != nil {
		r.Ok = false;
		r.Error = err;
	}

	r.Data = data;

	return r;
}

func getCommand(command []string) ([]ItemDB, error) {
	if len(command) < 2 {
		return nil, fmt.Errorf("invalid syntax on the `GET` command: `GET` needs one argument");
	}

	config, err := GetConfig();
	if err != nil {
		return nil, err;
	}

	var limit int64 = config.QueryLimit;
	if len(command) == 3 {
		limit, _ = strconv.ParseInt(strings.TrimSpace(command[2]), 10, 64);
	}

	arg := strings.TrimSpace(command[1]);

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
	case "FEED":
		items, err = SqlGetItemsByName(db, arg, limit);
	case "ID":
		if len(command) < 3 {
			return nil, fmt.Errorf("subcommand ID needs an argument");
		}
		id, err := strconv.ParseInt(command[2], 10, 64);
		if err != nil {
			return nil, fmt.Errorf("argument of the subcommand ID needs to be a integer");
		}

		var item ItemDB;
		item, err = SqlGetItem(db, id);
		items = []ItemDB{item};
	default:
		err = fmt.Errorf("invalid subcommand: %v", arg);
	}

	if err != nil {
		return nil, err;
	}

	return items, nil;
}


func updateCommand(command []string) (UpdateDataResponse, error) {
	var data UpdateDataResponse;

	if len(command) != 2 {
		return data, fmt.Errorf("invalid syntax on the `UPDATE` command: `UPDATE` only accepts one argument");
	}

	db, err := SqlConnect();
	if err != nil {
		return data, err;
	}
	defer db.Close();

	arg := strings.TrimSpace(command[1]);

	feeds, err := getFeedsConfig();
	if err != nil {
		return data, err;
	}

	var inserted int;

	// do paralelization here
	if arg == "ALL" {
		for _,v := range feeds {
			inserted, err = updateOneFeed(db, v.Name, v.Url);
			if err != nil {
				return data, err;
			}
		}
	} else {
		var feedUrl string;
		for _,v := range feeds {
			if v.Name == arg {
				feedUrl = v.Url;
			}
		}

		if feedUrl == "" {
			return data, fmt.Errorf("feeds name not found: `%v`", arg);
		}

		inserted, err = updateOneFeed(db, arg, feedUrl);
		if err != nil {
			return data, err;
		}
	}
	data.Updated = inserted;

	return data, nil;
}

func updateOneFeed(db *sql.DB, name string, url string) (int, error) {
	feed, err := getFeedFromWeb(url);
	if err != nil {
		return 0, err;
	}

	feed.Name = name;

	id, err := SqlUpsertFeed(db, feed);
	if err != nil {
		return 0, err;
	}

	inserted, err := SqlSaveFeedItems(db, feed.Items, id);
	if err != nil {
		return 0, err;
	}

	return inserted, nil;
}

func changeRead(stringId string, read bool) error {
	id, err := strconv.ParseInt(stringId, 10, 64);
	if err != nil {
		return err;
	}

	db, err := SqlConnect();
	if err != nil {
		return err;
	}
	defer db.Close();

	err = SqlUpdateItemRead(db, id, read);
	if err != nil {
		return err;
	}

	return nil;
}

func readCommand(command []string) error {
	if len(command) != 2 {
		return fmt.Errorf("invalid syntax on the `READ` command: `READ` only accepts one argument");
	}

	arg := strings.TrimSpace(command[1]);
	err := changeRead(arg, true);
	if err != nil {
		return err;
	}

	return nil;
}

func unreadCommand(command []string) error {
	if len(command) != 2 {
		return fmt.Errorf("invalid syntax on the `UNREAD` command: `UNREAD` only accepts one argument");
	}

	arg := strings.TrimSpace(command[1]);
	err := changeRead(arg, false);
	if err != nil {
		return err;
	}

	return nil;
}

func deleteCommand(command []string) error {
	if len(command) != 2 {
		return fmt.Errorf("invalid syntax on the `DELETE` command: `DELETE` only accepts one argument");
	}

	db, err := SqlConnect();
	if err != nil {
		return err;
	}
	defer db.Close();

	id, err := strconv.ParseInt(strings.TrimSpace(command[1]), 10, 64);
	if err != nil {
		return err;
	}

	err = SqlDeleteFeed(db, id);
	if err != nil {
		return err;
	}

	return nil;
}

func findCommand(command []string) ([]ItemDB, error) {
	if len(command) > 3 || len(command) < 2 {
		return nil, fmt.Errorf("invalid syntax on the `FIND` command: `FIND` only accepts two arguments");
	}

	config, err := GetConfig();
	if err != nil {
		return nil, err;
	}

	var limit int64 = config.QueryLimit;
	if len(command) == 3 {
		limit, _ = strconv.ParseInt(strings.TrimSpace(command[2]), 10, 64);
	}

	db, err := SqlConnect();
	if err != nil {
		return nil, err;
	}
	defer db.Close();

	items, err := SqlSearchItem(db, strings.TrimSpace(command[1]), limit);
	if err != nil {
		return nil, err;
	}

	return items, nil;
}

func openCommand(command []string) error {
	if len(command) < 2 {
		return fmt.Errorf("invalid syntax on the `OPEN` command: `OPEN` only accepts one arguments");
	}

	arg := strings.TrimSpace(command[1]);

	id, err := strconv.ParseInt(arg, 10, 64);
	if err != nil {
		return err;
	}

	db, err := SqlConnect();
	if err != nil {
		return err;
	}
	defer db.Close();

	item, err := SqlGetItem(db, id);
	if err != nil {
		return err;
	}

	cmd := exec.Command("xdg-open", item.Url);

	err = cmd.Run();
	if err != nil {
		return err;
	}

	return nil;
}

func addCommand(command []string) error {

	return fmt.Errorf("TODO: not implemented");
}
