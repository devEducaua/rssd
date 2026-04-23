package internal

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ItemDB struct {
    Id int64
    Title string
    Updated string
    Content string
    Read bool
    Url string
	FeedId int64
}

type FeedDB struct {
    Id int64
    Title string
    Name string
    Description string
    Url string
}

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
            custom_name TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL UNIQUE,
            description TEXT NOT NULL
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

func SqlGetFeed(db *sql.DB, url string) (FeedDB, error) {
    row := db.QueryRow("SELECT id, title, custom_name, url, description FROM feeds WHERE url=?", url);
    var f FeedDB;
    err := row.Scan(&f.Id, &f.Title, &f.Name, &f.Url, &f.Description);

    if err != nil {
        return FeedDB{}, err;
    }
    
    return f, nil;
}

func SqlUpsertFeed(db *sql.DB, feed Feed) (int64, error) {
    _, err := db.Exec("INSERT INTO feeds (title, custom_name, description, url) VALUES (?, ?, ?, ?) ON CONFLICT(url) DO UPDATE SET title=excluded.title, description=excluded.description, custom_name=excluded.custom_name", 
		feed.Title, feed.Name, feed.Description, feed.Url);

	if err != nil {
		return -1, fmt.Errorf("failed on upsert: %v", err);
	}

	feedQuery, err := SqlGetFeed(db, feed.Url);
	if err != nil {
		return -1, fmt.Errorf("failed to get feed on upsert: %v", err);
	}

	return feedQuery.Id, nil;
}


func SqlDeleteFeed(db *sql.DB, id int64) error {
    _, err := db.Exec("DELETE FROM feeds WHERE id=?", id);
    if err != nil {
        return err;
    }
    return nil;
}

func SqlSaveFeedItems(db *sql.DB, items []Item, feedId int64) (int, error) {
	var inserted int;

    for _,it := range items {
        result, err := db.Exec("INSERT OR IGNORE INTO items (title, updated, content, read, url, feed_id) VALUES (?, ?, ?, ?, ?, ?)", it.Title, it.Updated, it.Content, it.Read, it.Url, feedId);
        if err != nil {
            return inserrted, err;
        }

		affected, err := result.RowsAffected();
		if err != nil {
			return inserted, err;
		}

		if affected > 0 {
			inserted++;
		}
    }
    return inserted, nil;
}

func SqlGetItem(db *sql.DB, id int64) (ItemDB, error) {
    row := db.QueryRow("SELECT id, title, content, updated, url, read, feed_id FROM items WHERE id=?", id);

	var item ItemDB;
    err := row.Scan(&item.Id, &item.Title, &item.Content, &item.Updated, &item.Url, &item.Read, &item.FeedId);
    if err != nil {
        return ItemDB{}, err;
    }
	return item, nil;
}

func SqlGetItemsByName(db *sql.DB, name string, limit int64) ([]ItemDB, error) {
    row := db.QueryRow("SELECT id FROM feeds WHERE custom_name=?", name);

    var id int64;

    err := row.Scan(&id);
    if err != nil {
        return nil, err;
    }
    
    items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT id, url, title, updated, content, read, feed_id FROM items WHERE feed_id=? LIMIT ?", id, limit);
    
    if err != nil {
        return nil, err;
    }
    return items, nil;
}

// TODO: find a better to this function
func SqlGetAllItemsAttributesByCustom(db *sql.DB, limit int64, query string, queryArgs ...any) ([]ItemDB, error) {
    rows, err := db.Query(query, queryArgs...);
    if err != nil {
        return nil, err;
    }

    var items []ItemDB;
    for rows.Next() {
        var it ItemDB;
        if err := rows.Scan(&it.Id, &it.Url, &it.Title, &it.Updated, &it.Content, &it.Read, &it.FeedId); err != nil {
            return nil, err;    
        }
        items = append(items, it);
    }

    if err := rows.Err(); err != nil {
        return nil, err;
    }

    return items, nil;
}

func SqlGetAllItems(db *sql.DB, limit int64) ([]ItemDB, error) {
    items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT id, url, title, updated, content, read, feed_id FROM items LIMIT ?", limit);
    if err != nil {
        return nil, err;
    }
    return items, nil;
}

func SqlGetItemsByRead(db *sql.DB, read bool, limit int64) ([]ItemDB, error) {
    items, err := SqlGetAllItemsAttributesByCustom(db, limit, "SELECT id, url, title, updated, content, read, feed_id FROM items WHERE read=? LIMIT ?", read, limit);

    if err != nil {
        return nil, err;
    }
    return items, nil;
}

func SqlSearchItem(db *sql.DB, text string, limit int64) ([]int64, error) {
    text = "%" + text + "%";
    rows, err := db.Query("SELECT id FROM items WHERE title LIKE ? OR url LIKE ? OR content LIKE ? LIMIT ?", text, text, text, limit);
    if err != nil {
        return nil, err;
    }

    var ids []int64;
    for rows.Next() {
        var id int64;
        if err := rows.Scan(&id); err != nil {
            return nil, err;    
        }
        ids = append(ids, id);
    }

    if err := rows.Err(); err != nil {
        return nil, err;
    }

    return ids, nil;
}

func SqlUpdateItemRead(db *sql.DB, id int64, read bool) error {
    _, err := db.Exec("UPDATE items SET read=? WHERE id=?", read, id);
    if err != nil {
        return err;
    }
    return nil;
}

