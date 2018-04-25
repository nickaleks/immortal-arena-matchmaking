package main

import (
	"strconv"
	"sync/atomic"
	"encoding/json"
	"net/http"
)

var (
	lobbies map[int64]Lobby
	last_lobby_id int64
)

type Player struct {
	Address string	`json:"address"`
	Ready bool	`json:"ready"`
}

type Lobby struct {
	ID int64	`json:"id"`
	PlayersCount int `json:"players_count"`
	Players []Player `json:"players"`
}

func createLobby(w http.ResponseWriter, r *http.Request) {
	lobby := Lobby{last_lobby_id, 0, []Player{}}
	lobbies[last_lobby_id] = lobby

	atomic.AddInt64(&last_lobby_id, 1)

	json, err := json.Marshal(&lobby)

	if err != nil {
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(json)

}

func getLobby(w http.ResponseWriter, r *http.Request) {
	lid, ok := r.URL.Query()["lobby_id"]
	if !ok {
		w.Write([]byte("No lobby_id key"))
		return
	}

	lobby_id, err := strconv.ParseInt(lid[0], 10, 32)
	if err != nil {
		w.Write([]byte("lobby_id must be valid int"))
		return
	}

	lobby, ok := lobbies[lobby_id]
	if !ok {
		w.Write([]byte("lobby does not exist"))
	}

	json, err := json.Marshal(&lobby)

	if err != nil {
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(json)
}

func getAllLobbies(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(lobbies)
	if err != nil {
		w.Write([]byte("Internal server error"))
		return
	}
	w.Write(json)
}

func leaveLobby(w http.ResponseWriter, r *http.Request) {
	lid, ok := r.URL.Query()["lobby_id"]
	if !ok {
		w.Write([]byte("No lobby_id key"))
		return
	}

	address, ok := r.URL.Query()["address"]
	if !ok {
		w.Write([]byte("No address key"))
		return
	}

	lobby_id, err := strconv.ParseInt(lid[0], 10, 32)
	if err != nil {
		w.Write([]byte("lobby_id must be valid int"))
		return
	}
	lobby, ok := lobbies[lobby_id]
	if !ok {
		w.Write([]byte("Lobby does not exist"))
		return
	}
	i := 0
	for lobby.Players[i].Address != address[0] && i < len(lobby.Players) {
		i++
	}
	if i == len(lobby.Players) {
		w.Write([]byte("No IP found"))
	}
	lobbies[lobby_id] = Lobby{lobby_id, lobby.PlayersCount - 1, append(lobby.Players[:i], lobby.Players[i+1:]...)}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Removed from lobby"))
}

func joinLobby(w http.ResponseWriter, r *http.Request) {
	lid, ok := r.URL.Query()["lobby_id"]
	if !ok {
		w.Write([]byte("No lobby_id key"))
		return
	}

	address, ok := r.URL.Query()["address"]
	if !ok {
		w.Write([]byte("No address key"))
		return
	}

	lobby_id, err := strconv.ParseInt(lid[0], 10, 32)
	if err != nil {
		w.Write([]byte("lobby_id must be valid int"))
		return
	}
	lobby, ok := lobbies[lobby_id]
	if !ok {
		w.Write([]byte("Lobby does not exist"))
		return
	}
	if lobby.PlayersCount == 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lobby is full"))
		return
	}
	lobbies[lobby_id] = Lobby{lobby_id, lobby.PlayersCount + 1, append(lobby.Players, Player{address[0], false})}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Added to lobby"))
}

func deleteLobby(w http.ResponseWriter, r *http.Request) {
	lid, ok := r.URL.Query()["lobby_id"]
	if !ok {
		w.Write([]byte("No lobby_id key"))
		return
	}

	lobby_id, err := strconv.ParseInt(lid[0], 10, 32)
	if err != nil {
		w.Write([]byte("lobby_id must be valid int"))
		return
	}

	lobby, ok := lobbies[lobby_id]
	if !ok {
		w.Write([]byte("lobby does not exist"))
	}

	json, err := json.Marshal(&lobby)

	if err != nil {
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(json)

	delete(lobbies, lobby_id)
}

func playerReady(w http.ResponseWriter, r *http.Request) {
	lid, ok := r.URL.Query()["lobby_id"]
	if !ok {
		w.Write([]byte("No lobby_id key"))
		return
	}

	address, ok := r.URL.Query()["address"]
	if !ok {
		w.Write([]byte("No address key"))
		return
	}

	lobby_id, err := strconv.ParseInt(lid[0], 10, 32)
	if err != nil {
		w.Write([]byte("lobby_id must be valid int"))
		return
	}
	lobby, ok := lobbies[lobby_id]
	if !ok {
		w.Write([]byte("Lobby does not exist"))
		return
	}
	if lobby.PlayersCount < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lobby is not full"))
		return
	}
	for i := 0; i < len(lobby.Players); i++ {
		if lobby.Players[i].Address == address[0] {
			lobby.Players[i].Ready = true
			break
		}
	}

	lobbies[lobby_id] = lobby
	json, err := json.Marshal(&lobby)

	if err != nil {
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(json)
}

func main() {
	atomic.StoreInt64(&last_lobby_id, 0)
	lobbies = make(map[int64]Lobby)
	http.HandleFunc("/lobbies", getAllLobbies)
	http.HandleFunc("/join", joinLobby)
	http.HandleFunc("/create", createLobby)
	http.HandleFunc("/lobby", getLobby)
	http.HandleFunc("/delete", deleteLobby)
	http.HandleFunc("/leave", leaveLobby)
	http.HandleFunc("/ready", playerReady)

	http.ListenAndServe(":8080", nil)
}
