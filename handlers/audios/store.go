package audios

import (
	"aitring/model"
	"aitring/services"
	AudioServ "aitring/services/audios"
	"aitring/utils"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
)

func New(store services.Store) Handler {
	return &audioHandler{
		audioService: store.AudioServ,
	}
}

type audioHandler struct {
	audioService AudioServ.Service
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000" || origin == "https://tringai.com"
		      // Allow all origins (for dev)
        return true
	},
}

func (h *audioHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	sessionID := r.URL.Query().Get("session_id")
	if userID == "" || sessionID == "" {
		http.Error(w, "user_id and session_id required", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	chunkID := uuid.NewString()
	meta := model.ChunkMeta{
		ChunkID:    chunkID,
		SessionID:  sessionID,
		UserID:     userID,
		Timestamp:  time.Now(),
		Size:       len(data),
		Transcript: "processing...", // placeholder
	}
	err = h.audioService.UploadAudio(chunkID, data, meta)
	if err != nil {
		http.Error(w, "failed to upload audio", http.StatusInternalServerError)
	}

	utils.ReturnResponse(w, http.StatusOK, meta)
}

func (h audioHandler) GetChunkByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	meta, err := h.audioService.GetAudioMetadata(id)
	if err != nil {
		http.Error(w, "audio metadata not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.ReturnResponse(w, http.StatusOK, meta)
}

func (h audioHandler) GetChunksByUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	chunks, err := h.audioService.GetAudioChunks(userID)
	if err != nil {
		http.Error(w, "failed to get audio chunks", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunks)
}

func (h audioHandler) WSHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	sessionID := r.URL.Query().Get("session_id")
	if userID == "" || sessionID == "" {
		http.Error(w, "user_id and session_id required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		chunkID := uuid.NewString()
		meta := model.ChunkMeta{
			ChunkID:    chunkID,
			SessionID:  sessionID,
			UserID:     userID,
			Timestamp:  time.Now(),
			Size:       len(msg),
			Transcript: "simulated transcript text",
		}
		err = h.audioService.UploadAudio(chunkID, msg, meta) // Acknowledge
		if err != nil {
			utils.ErrorResponse(w, "failed to upload audio", http.StatusInternalServerError)
		}
		resp := map[string]interface{}{
			"type":     "ack",
			"chunk_id": chunkID,
			"meta":     meta,
		}
if err := conn.WriteJSON(resp); err != nil {
    log.Println("write json error:", err)
}
	}
}
