# rssd protocol specification

## requests
requests are in plain text, and contains a command, his arguments and a CTRF in the end. a request can only contain spaces if the are surrounded with quotes.

### command list

GET
- description: get a list of items.
- argument: source
- argument description: source can be ALL, READ, UNREAD or a feedname.
- optional argument: limit

READ
- argument: id
- description: marks the especified item as read.

UNREAD
- argument: id
- description: marks the especified item as not read.

DELETE
- argument: id
- description: deletes a feed and his items.

UPDATE
- description: query the feeds to update the database.

FIND
- description: find a item
- arguments: text
- argument description: a text that is on the attributes of the item.

## responses
the responses are in json, and have the obrigatory status attribute.
the status attribute can be "yes" or "not".

### response examples
example GET response
```json 
{
    "status": "yes",
    "response": [
        {
            "id": 1,
            "title", "lorem",
            "description", "lorem",
            "updated", "2026-06-06",
            "content", "lorem ipsum",
            "read", true,
        }
    ]
}
```

example READ response
```json 
{
    "status": "yes",
    "response": "item with id: 1 is read"
}
```

example UNREAD response
```json 
{
    "status": "yes",
    "response": "item with id: 1 is unread"
}
```

example DELETE response
```json 
{
    "status": "yes",
    "response": "item with id: 1 is deleted"
}
```

example UPDATE response
```json 
{
    "status": "yes",
    "response": "the database was updated"
}
```

example FIND response
returns a ordered array of possibly results.
```json 
{
    "status": "yes",
    "response": [1, 4, 3]
}
```
