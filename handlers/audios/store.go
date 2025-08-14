package audios

import (
	"aitring/model"
	"aitring/services"
	AudioServ "aitring/services/audios"
	"aitring/utils"
	"io"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"encoding/json"
		"encoding/base64"

 "strings"

	"fmt"
	
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

type UploadResponse struct {
    Status  string `json:"status"`
    ChunkID string `json:"chunk_id"`
}
func (h *audioHandler) Upload(w http.ResponseWriter, r *http.Request) {
    userID := r.FormValue("user_id")
    sessionID := r.FormValue("session_id")
    tsStr := r.FormValue("timestamp")

    if userID == "" || sessionID == "" {
        http.Error(w, "missing user_id or session_id", http.StatusBadRequest)
        return
    }

    // Parse timestamp if provided, else use now
    ts := time.Now()
    if tsStr != "" {
        if parsed, err := time.Parse(time.RFC3339, tsStr); err == nil {
            ts = parsed
        }
    }

    var data []byte
    // Detect multipart vs raw
    ct := strings.ToLower(r.Header.Get("Content-Type"))
    if strings.HasPrefix(ct, "multipart/form-data") {
        if err := r.ParseMultipartForm(25 << 20); err != nil {
            http.Error(w, "parse multipart failed", http.StatusBadRequest)
            return
        }

        files := r.MultipartForm.File["files"]
        if len(files) == 0 {
            http.Error(w, "no files provided", http.StatusBadRequest)
            return
        }

        var responses []UploadResponse
        for _, fileHeader := range files {
            file, err := fileHeader.Open()
            if err != nil {
                http.Error(w, "open file failed", http.StatusBadRequest)
                return
            }
            data, err := io.ReadAll(file)
            file.Close()
            if err != nil {
                http.Error(w, "read file failed", http.StatusBadRequest)
                return
            }

            chunkID := uuid.NewString()
            ackCh := make(chan model.ChunkMeta, 1)
            raw := model.RawChunk{
                ChunkID:   chunkID,
                SessionID: sessionID,
                UserID:    userID,
                Data:      data,
                Received:  ts,
                AckCh:     ackCh,
            }

            // Push into pipeline
            if ok, _ := h.audioService.UploadAudio(r.Context(), raw); !ok {
                http.Error(w, "pipeline backpressure: rejected", http.StatusTooManyRequests)
                return
            }

            responses = append(responses, UploadResponse{
                Status:  "accepted",
                ChunkID: chunkID,
            })
        }

        // Respond with all chunk IDs
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusAccepted)
        json.NewEncoder(w).Encode(responses)
        return
    
    } else {
        defer r.Body.Close()
        b, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "read body failed", http.StatusBadRequest)
            return
        }
        data = b
    }

    // Create chunk
    chunkID := uuid.NewString()
    ackCh := make(chan model.ChunkMeta, 1)
    raw := model.RawChunk{
        ChunkID:   chunkID,
        SessionID: sessionID,
        UserID:    userID,
        Data:      data,
        Received:  ts,
        AckCh:     ackCh,
    }

	
    // Push into pipeline
    if ok, _ := h.audioService.UploadAudio(r.Context(), raw); !ok {
        http.Error(w, "pipeline backpressure: rejected", http.StatusTooManyRequests)
        return
    }

    // Respond with 202 Accepted
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(UploadResponse{
        Status:  "accepted",
        ChunkID: chunkID,
    })
}

func (h audioHandler) GetChunkByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	meta, err := h.audioService.GetAudioMetadata(id)
	if err != nil {
		utils.ErrorResponse(w,"audio metadata not found", http.StatusNotFound)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, meta)
}

func (h audioHandler) GetChunksByUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	chunks, err := h.audioService.GetAudioChunks(userID)
	if err != nil {
		// http.Error(w, "failed to get audio chunks", http.StatusInternalServerError)
		utils.ErrorResponse(w,"audio metadata not found", http.StatusNotFound)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, chunks)
}

func (h audioHandler) WSHandler(w http.ResponseWriter, r *http.Request) {
userID := r.URL.Query().Get("user_id")
	sessionID := r.URL.Query().Get("session_id")
	if userID == "" || sessionID == "" { http.Error(w, "user_id & session_id required", http.StatusBadRequest); return }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil { http.Error(w, "upgrade failed", http.StatusBadRequest); return }
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil { return }

		// Support both binary frames and JSON text frames with base64
		var data []byte
		var ts = time.Now()
		if mt == websocket.BinaryMessage {
			data = msg
		} else {
			var payload struct {
				AudioB64  string `json:"audio_b64"`
				Timestamp string `json:"timestamp,omitempty"`
			}
			if err := json.Unmarshal(msg, &payload); err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid json"}`))
				continue
			}
			if payload.Timestamp != "" { if parsed, err := time.Parse(time.RFC3339, payload.Timestamp); err == nil { ts = parsed } }
			if payload.AudioB64 == "" {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"missing audio_b64"}`))
				continue
			}
			b, err := base64.StdEncoding.DecodeString(payload.AudioB64)
			if err != nil { _ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid base64"}`)); continue }
			data = b
		}

		chunkID := uuid.NewString()
		ackCh := make(chan model.ChunkMeta, 1)
		raw := model.RawChunk{ChunkID: chunkID, SessionID: sessionID, UserID: userID, Data: data, Received: ts, AckCh: ackCh}
		ok, _ := h.audioService.UploadAudio(r.Context(), raw)

		// Immediate ack
		ack := map[string]any{"type": "ack", "chunk_id": chunkID, "accepted": ok}
		bAck, _ := json.Marshal(ack)
		_ = conn.WriteMessage(websocket.TextMessage, bAck)
		if !ok { continue }

		// Stream metadata when processed (or timeout)
		select {
		case meta := <-ackCh:
			metaEv := map[string]any{"type": "metadata", "chunk_id": chunkID, "meta": meta}
			b, _ := json.Marshal(metaEv)
			_ = conn.WriteMessage(websocket.TextMessage, b)
		case <-time.After(5 * time.Second):
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"type":"metadata","chunk_id":"%s","status":"pending"}`, chunkID)))
		}
}
}
