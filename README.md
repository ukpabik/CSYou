# CSYou — Real-Time CS2 Analytics

CSYou is a self-hosted analytics platform for Counter-Strike 2 that captures **live game telemetry** via CS2's Game State Integration (GSI) to provide real-time insights into your gameplay.

<img width="1908" height="999" alt="image" src="https://github.com/user-attachments/assets/890a1dcb-b7aa-4def-b98c-ee1c79edafa3" />

---

## ✨ Features

-   🎯 **Real-time Player State Tracking**: Monitor your health, armor, money, weapons, and ammunition as you play.
-   🔫 **Per-Kill Event Logs**: Get detailed logs for each kill, including the weapon used, headshot status, and your ammo state at the time of the kill.
-   📊 **Round-by-Round Analytics**: Analyze round outcomes, economy impact, kill timelines, and win conditions.
-   🚀 **Event-Driven Architecture**: Built on a modern, scalable event pipeline:
    -   **CS2 GSI → Go Collector → Kafka → Redis → GUI**
-   🖥 **Self-Hosted Tauri GUI**: A cross-platform desktop application for querying and visualizing your match data.
-   📦 **Dockerized Deployment**: Includes a `docker-compose` setup to easily run Redis, Kafka, and other services.

---

## ⚙️ How It Works

The platform works by listening for HTTP POST requests sent directly from the CS2 game client.

1.  **CS2 Game State Integration (GSI)**: You configure your CS2 client to send JSON payloads containing game state data to a local endpoint.
2.  **Go Collector**: A lightweight Go service listens on this endpoint (`http://127.0.0.1:3000`), validates the incoming data, and publishes it as a raw event to a Kafka topic.
3.  **Kafka**: Acts as a durable and scalable message bus, decoupling the data ingestion from processing.
4.  **Redis**: A fast in-memory store that holds the latest game state, allowing the GUI to display live data with minimal latency.
5.  **Tauri GUI**: The frontend application reads directly from Redis to provide a real-time view of your current match statistics.

*(In the future, a separate processor service will consume events from Kafka, enrich them, and store them in ClickHouse for historical analysis and long-term storage.)*

---

## 🔧 Getting Started

Follow these steps to set up and run CSYou on your local machine for development.

### Prerequisites

-   **Go** (≥1.22)
-   **Node.js** (≥18) + **npm**
-   **Rust** + [Tauri system dependencies](https://tauri.app/v1/guides/getting-started/prerequisites)
-   **Docker** + **Docker Compose**
-   **Counter-Strike 2**

---

### 1. Configure CS2 Game State Integration

You must tell CS2 where to send its telemetry data. Create a file named `gamestate_integration_cs2.cfg` inside your CS2 `cfg/` folder. This is typically located at:
`Steam/steamapps/common/Counter-Strike Global Offensive/game/csgo/cfg/`

Paste the following configuration into the file:

```txt
"CS2 Analytics Integration"
{
    "uri" "[http://127.0.0.1:3000](http://127.0.0.1:3000)"
    "timeout" "5.0"
    "heartbeat" "0.25"

    "auth"
    {
        "token" "cs2_secret_token"
    }

    "data"
    {
        "provider" "1"            // game + client info
        "map" "1"                 // map info (name, round, mode)
        "round" "1"               // round state (phase, win team)
        "player_id" "1"           // Steam ID
        "player_state" "1"        // health, armor, alive
        "player_weapons" "1"      // weapons, ammo
        "player_match_stats" "1"  // kills, assists, deaths, mvps, score
        "bomb" "1"                // bomb events (planted, defused, exploded)
    }
}
```
### 2. Configure Your Player ID

The backend needs to know your Steam ID to isolate your player data from the GSI stream.

In the `backend/` directory, copy the example configuration file:

```bash
cp backend/config.example.json backend/config.json
```

Now, open `backend/config.json` and replace the placeholder with your **Steam64 ID**. You can find your ID using a service like [SteamID Finder](https://steamid.io/).

```json
{
  "steam_id": "76561197960287930" // <-- Replace with your Steam64 ID
}
```

> **Note**: At this stage, only your own player telemetry is tracked.

---

### 3. Run the Backend Services

With Docker running, start the required services (Redis, Kafka, and ClickHouse) using Docker Compose. From the project root:

```bash
docker-compose up -d
```

This will run the containers in the background.

---

### 4. Run the Go Collector

Navigate to the backend directory and run the Go application. This server will listen for data from your CS2 client.

```bash
cd backend
go run cmd/main.go
```

---

### 5. Run the Frontend GUI

Finally, navigate to the frontend directory, install dependencies, and launch the Tauri desktop application.

```bash
cd frontend
npm install
npm run tauri dev
```

This will open the CSYou desktop app. Once you join a match in CS2, events will begin streaming through the pipeline, and your live stats will appear in the GUI.
