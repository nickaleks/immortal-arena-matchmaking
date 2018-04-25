# Matchmaking server

__Restarting server resets its state__

## JSON Format

Lobby:
```
{
    "id":0,
    "players_count":0,
    "players":[
        "127.0.0.1"
    ]
}
```

Map of Lobbies:
```
{
    "1": {
        "id":1,
        "players_count":0,
        "players":[]
    },
    "2": {
        "id":2,
        "players_count":0,
        "players":[]
    }
}
```

## Endpoints

All requests are GET

#### Create Lobby

```/create```

Returns: Lobby json

#### Get Lobby

```/lobby?lobby_id=```

Returns: Lobby json

#### Get Lobbies

```/lobbies```

Returns: Map of lobbies json

#### Join Lobby

```/join?lobby_id=&address=```

Returns: "Added to Lobby" string

#### Leave Lobby

```/leave?lobby_id=&address=```

Returns: "Removed from lobby" string

#### Delete Lobby

```/delete?lobby_id=```

Returns: Lobby json

### Player Ready

```/ready?lobby_id=&address=```

Returns: Lobby json

### Player not Ready

```/ready?lobby_id=&address=```

Returns: Lobby json
