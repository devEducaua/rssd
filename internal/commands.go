package internal

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Response any `json:"response"`
}

func ParseCommand(command string) Response {
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
		case "OPEN":
			msg, err = openCommand(parts);
			r.Response = msg;
		default:
			err = fmt.Errorf("command: %v doesn't exists", parts[0]);
	}

	if err != nil {
		r.Status = "no";
		r.Response = err.Error();
	}

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
	default:
		id, err := strconv.ParseInt(arg, 10, 64);

		// TODO: refactor this to 2 separate functions
		if err != nil {
			items, err = SqlGetItemsByName(db, arg, limit);

		} else {
			// TODO: return just the item without array
			var item ItemDB;
			item, err = SqlGetItem(db, id);
			items = []ItemDB{item};
		}
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

	db, err := SqlConnect();
	if err != nil {
		return "", err;
	}
	defer db.Close();

	arg := strings.TrimSpace(command[1]);

	feeds, err := getFeedsConfig();
	if err != nil {

		return "", err;
	}

	// do paralelization here
	if arg == "ALL" {
		for _,v := range feeds {
			err := updateOneFeed(db, v.Name, v.Url);
			if err != nil {
				return "", err;
			}
		}
	} else {

		// TODO: turn CONFIG.feeds on a hash map to better performance
		var feedUrl string;
		for _,v := range feeds {
			if v.Name == arg {
				feedUrl = v.Url;
			}
		}

		if feedUrl == "" {
			return "", fmt.Errorf("feeds name not found: `%v`", arg);
		}

		err = updateOneFeed(db, arg, feedUrl);
		if err != nil {
			return "", err;
		}
	}

	return "feeds are updated", nil;
}

func updateOneFeed(db *sql.DB, name string, url string) error {
	feed, err := getFeedFromWeb(url);
	if err != nil {
		return err;
	}

	feed.Name = name;

	id, err := SqlUpsertFeed(db, feed);
	if err != nil {
		return err;
	}

	err = SqlSaveFeedItems(db, feed.Items, id);
	if err != nil {
		return err;
	}

	return nil;
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

	arg := strings.TrimSpace(command[1]);
	id, err := changeRead(arg, true);
	if err != nil {
		return "", err;
	}

	return fmt.Sprintf("item with id: %v is read", id), nil;
}

func unreadCommand(command []string) (string, error) {
	if len(command) != 2 {
		return "", fmt.Errorf("invalid syntax on the `UNREAD` command: `UNREAD` only accepts one argument");
	}

	arg := strings.TrimSpace(command[1]);
	id, err := changeRead(arg, false);
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

	id, err := strconv.ParseInt(strings.TrimSpace(command[1]), 10, 64);
	if err != nil {
		return "", err;
	}

	err = SqlDeleteFeed(db, id);
	if err != nil {
		return "", err;
	}

	return fmt.Sprintf("feed with id: %v deleted", id), nil;
}

func findCommand(command []string) ([]int64, error) {
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

	ids, err := SqlSearchItem(db, strings.TrimSpace(command[1]), limit);
	if err != nil {
		return nil, err;
	}

	return ids, nil;
}

func openCommand(command []string) (string, error) {
	if len(command) < 2 {
		return "", fmt.Errorf("invalid syntax on the `OPEN` command: `OPEN` only accepts one arguments");
	}

	arg := strings.TrimSpace(command[1]);

	id, err := strconv.ParseInt(arg, 10, 64);
	if err != nil {
		return "", err;
	}

	db, err := SqlConnect();
	if err != nil {
		return "", err;
	}
	defer db.Close();

	item, err := SqlGetItem(db, id);
	if err != nil {
		return "", err;
	}

	cmd := exec.Command("xdg-open", item.Url);

	err = cmd.Run();
	if err != nil {
		return "", err;
	}

	return "external program finished", nil;
}

