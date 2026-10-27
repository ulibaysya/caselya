# User API

## Authentication

### Steam

- POST /auth/steam/url
- GET /auth/steam/callback

### Opaque tokens

- POST /auth/opaque/exchange

## Items

- GET /items/{item}

## Container

- GET /containers

List

- GET /containers/by-alias

Resolve container ID by its alias

- POST /containers/{container}/opening

Open container

## User

- GET /users/me
- GET /users/{user}

### Items

- GET /users/{user}/items
- GET /users/{user}/items/{item}
- POST /users/{user}/items/{item}/withdrawal
- POST /users/{user}/items/{item}/selling

# Admin API

## Items

- POST /item-tag-categories

- GET /items

Filtered search

- POST /items

```json
{
    "name": "Royal Blue",
    "price": 127485,
    "tags": {
        "weapon": ["usp"],
        "quality": ["mw"], // minimal wear 
        "testing": ["yes"]
    }
}
```
 
Some tags are predefined and cannot be modified.

The system allows multiple tags to be assigned to some categories (e.g.
`"weapon_type": ["usp", "ak47"]`), which may not make sense.

The system also allows some tag categories to be specified while others are omitted (e.g.
specifying weapon_type but omitting weapon_rarity), which may not make sense either.

- GET /items/{item}
- DELETE /items/{item}
- PATCH /items/{item}
- POST /items/{item}/image-upload

Returns S3 presign url, where user uploads image

## Container

- POST /containers
- POST /containers/{container}/image-upload

Returns S3 presign url, where user uploads image

## User

- POST /users
- POST /users/{user}/opaque-tokens

### Items

- GET /users/{user}/items
- POST /users/{user}/items
- DELETE /users/{user}/items/{item}
- POST /users/{user}/items/{item}/withdrawal
- POST /users/{user}/items/{item}/selling

### Balance

- PUT /users/{user}/balance

Updates balance entirely. Can't be lower than 0.

- PATCH /users/{user}/balance

Adjusts balance. Body is delta - can be any number (but 0 is no-op)
