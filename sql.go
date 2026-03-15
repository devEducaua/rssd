package main

import (
    "database/sql"
    _"modernc.org/sqlite"
)

func SqlCreate() (*sql.DB, error) {
    const DBPATH = "./rssd.db";

    db, err := sql.Open("sqlite", DBPATH);
    if err != nil {
        return nil, err;
    }
    
    return db, nil;
}

func SqlCreateTables(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS feeds (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            url TEXT NOT NULL,
            description TEXT NOT NULL,
            UNIQUE(url)
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
            FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE NOT NULL,
            UNIQUE(feed_id, url)
        )
    `)
    if err != nil {
        return err;
    }

    return nil;
}

func SqlGetItemsByFeed(db *sql.DB, feedId int64) ([]Item, error) {
    rows, err := db.Query("SELECT title, updated, content, url, read FROM items WHERE feed_id = ?", feedId);
    if err != nil {
        return nil, err;
    }
    defer rows.Close();

    var items []Item;

    for rows.Next() {
        var it Item;
        if err := rows.Scan(&it.Title, &it.Updated, &it.Content, &it.Url, &it.Read); err != nil {
            return items, err
        }
        items = append(items, it);
    }
    if err := rows.Err(); err != nil {
        return nil, err;
    }

    return items, nil;
}

func SqlSaveFeed(db *sql.DB, feed Feed) (int64, error) {
    _, err := db.Exec(`INSERT OR IGNORE INTO feeds (title, url, description) VALUES (?, ?, ?)`, feed.Title, feed.Url, feed.Description);
    if err != nil {
        return -1, err;
    }

    row := db.QueryRow(`SELECT id FROM feeds WHERE url = ?`, feed.Url);

    var feedId int64;
    if err := row.Scan(&feedId); err != nil {
        return -1, err;
    }

    for _,f := range feed.Items {
        _, err = db.Exec(`INSERT OR IGNORE INTO items (title, updated, content, url, feed_id) VALUES (?, ?, ?, ?, ?)`, f.Title, f.Updated, f.Content, f.Url, feedId);    
        if err != nil {
            return -1, err;
        }
    }
    return feedId, nil;
}

func SqlRemoveFeed(db *sql.DB, feedId int64) error {
    _, err := db.Exec(`DELETE FROM feeds WHERE id=?`, feedId);
    if err != nil {
        return err;
    }
    return nil;
}

func SqlUpdateFeed(db *sql.DB, feed Feed, feedId int64) error {
    _, err := db.Exec(`UPDATE feeds SET title=?, url=?, description=? WHERE id=?`, feed.Title, feed.Url, feed.Description, feedId);
    if err != nil {
        return err;
    }
    return nil;
}

func SqlCustomQuery(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.Query(query, args);
	if err != nil {
		return nil, err;
	}

	return rows, nil; 
}

