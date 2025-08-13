# Audio Processing Service in Go

## 📌 Overview
This project is a **high-performance, modular Go service** designed to process incoming audio chunks from both **HTTP** and **WebSocket** clients.  
It implements a **multi-stage concurrent processing pipeline** that:
- Ingests audio
- Validates and transforms it
- Extracts metadata
- Persists results in **memory** and on **disk**
- Exposes REST and WebSocket APIs for retrieval and live streaming

The system is built with **production-grade patterns** such as worker pools, backpressure handling, thread-safe storage, and observability hooks.

---

## ⚙️ Features
- **Audio Ingestion**
  - `POST /upload` — accepts audio via HTTP (raw or multipart-form)
  - `WebSocket /ws` — supports bidirectional streaming of audio chunks
- **Concurrent Processing Pipeline**
  - Stages: `Ingestion → Validation → Transformation → Metadata Extraction → Storage`
  - Worker pool at each stage for controlled concurrency
  - Channels for stage-to-stage communication
  - Context-based cancellation support
  - Simulated transformations (checksum, FFT mock, fake transcript)
  - Backpressure handling (queue size limits, job dropping or rejection)
- **Storage**
  - Thread-safe in-memory store (`sync.RWMutex` or `sync.Map`)
  - Persistent storage to disk (JSON, BadgerDB, or SQLite)
  - Lookup by `chunk_id`, `session_id`, or `user_id`
- **APIs**
  - `POST /upload` — submit audio chunk
  - `GET /chunks/{id}` — retrieve processed metadata
  - `GET /sessions/{user_id}` — list chunks for a user's session
  - WebSocket `/ws` — send audio chunks and receive metadata/transcripts in real time

---

## 🛠 Tech Stack
- **Language:** Go
- **Concurrency:** Goroutines, Channels, Worker Pools
- **Storage:** In-memory (thread-safe map) + Disk persistence
- **Networking:** net/http, Gorilla WebSocket
- **Data Encoding:** JSON
- **Optional DB:** SQLite or BadgerDB

---

## 📂 Project Structure
├── cmd/ # Main application entry point
├── ai-tring/
│ ├── app/ # Routes
│ ├── handlers/# HTTP and WebSocket handlers
│ ├── services/ # Multi-stage processing pipeline
│ ├── store/  # In-memory and persistent storage
│ └── utils/ # Helper functions (checksum, mock FFT, etc.)
├── data/ # Persistent storage (JSON/DB files)
└── README.md # Project documentation
