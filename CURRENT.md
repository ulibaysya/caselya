POST /item-tag-categories

```json
{
    "name": "some-category",
    "public": false,
    "values": [
        "v1", "v2", "v3"
    ]
}
```

POST /items

```json
{
    "name": "Royal Blue",
    "price": 127485,
    "tags": {
        "weapon": ["usp-s"],
        "quality": ["mw"], // minimal wear 
        "testing": ["yes"]
    }
}
```

GET /items/{item} # user API
GET /items/{item} # admin API

DELETE /items/{item}
PATCH /items/{item}

GET /items

POST /items/{item}/image-upload
