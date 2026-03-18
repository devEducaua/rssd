package main

import (
    "database/sql"
    _ "modernc.org/sqlite"
)
func SqlConnect() (*sql.DB, error) { const DBPATH = "./rssd.db";

    db, err := sql.Open("sqlite", DBPATH);
    if err != nil {
        return nil, err;
    }
    
    return db, nil;
}

func SqlCreateTablesIfNotExists(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS feeds (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            custom_name TEXT NOT NULL,
            url TEXT NOT NULL,
            description TEXT NOT NULL,
            UNIQUE(url, custom_name)
        )
    `)
    if err != nil {
        return err;
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            updated TEXT NOT NULL,
            content TEXT NOT NULL,
            read BOOLEAN DEFAULT FALSE,
            url TEXT NOT NULL,
            feed_id INTEGER NOT NULL,
            FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
            UNIQUE(feed_id, url)
        )
    `)
    if err != nil {
        return err;
    }

    return nil;
}

func SqlUpdateItemRead(db *sql.DB, id int64, read bool) error {
	_, err := db.Exec("UPDATE items SET read=? WHERE id=?", read, id);
	if err != nil {
		return err;
	}
	return nil;
}

func SqlDeleteFeed(db *sql.DB, url string) error {
	_, err := db.Exec("DELETE FROM feeds WHERE url=?", url);
	if err != nil {
		return err;
	}
	return nil;
}

// TODO: find a better to this function
func SqlGetAllItemsAttributesByCustom(db *sql.DB, limit int64, query string, queryArgs ...any) ([]Item, error) {
    rows, err := db.Query(query, queryArgs);
    if err != nil {
        return nil, err;
    }

    var items []Item;
    for rows.Next() {
        var it Item;
        if err := rows.Scan(&it.Url, &it.Title, &it.Updated, &it.Content, &it.Read); err != nil {
            return nil, err;    
        }
        items = append(items, it);
    }

    if err := rows.Err(); err != nil {
        return nil, err;
    }

    return items, nil;
}

func SqlGetAllItems(db *sql.DB, limit int64) ([]Item, error) {
	items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT * FROM items LIMIT ?", limit);
	if err != nil {
		return nil, err;
	}
	return items, nil;
}

func SqlGetItemsByRead(db *sql.DB, read bool, limit int64) ([]Item, error) {
	items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT * FROM items LIMIT ? WHERE read=?", limit, read);

	if err != nil {
		return nil, err;
	}
	return items, nil;
}

func SqlGetItemsByName(db *sql.DB, name string, limit int64) ([]Item, error) {
	items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT * FROM items LIMIT ? WHERE name=?", limit, name);
    
	if err != nil {
		return nil, err;
	}
	return items, nil;
}
