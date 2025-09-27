# CSYou — Real-Time CS2 Analytics (🚧 In Progress)

CSYou is a self-hosted analytics platform for Counter-Strike 2 that captures **live game telemetry** via CS2 Game State Integration (GSI).  
It uses **Go, Kafka, Redis, and Docker** to build a local event-driven pipeline that powers a real-time GUI dashboard.  

> ⚠️ This project is still **in progress**. Core event collection is working, but the GUI and advanced querying are under development.

---

## ✨ Features
- 🎯 **Real-time player state tracking** (health, armor, money, weapons, etc.)
- 🔫 **Per-kill event logs** (weapon used, headshots, ammo state at time of kill)
- 📊 **Round-by-round analytics** (economy impact, kill timeline, win conditions)
- 🚀 **Event-driven pipeline**:
  - **CS2 GSI → Go collector → Kafka → Redis → GUI**
- 🖥 **Self-hosted Tauri GUI** for querying and visualizing match data
- 📦 **Dockerized deployment** for running Redis, Kafka, and the Go services together

---

## ⚡ Tech Stack
- **Go** — core event collector + processing  
- **Kafka** — event streaming backbone  
- **Redis (ReJSON)** — fast in-memory store for snapshots and logs  
- **Docker** — containerized environment  
- **Tauri (planned)** — lightweight desktop GUI for stats & queries  

---

## 🚀 Usage (coming soon)
