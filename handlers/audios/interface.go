package audios

import "net/http"


type Handler interface {
	Upload(w http.ResponseWriter, r *http.Request)
GetChunkByID(w http.ResponseWriter, r *http.Request)
GetChunksByUser(w http.ResponseWriter, r *http.Request) 
WSHandler(w http.ResponseWriter, r *http.Request)
}