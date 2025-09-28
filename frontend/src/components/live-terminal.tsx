"use client"

import { useState, useEffect, useRef } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Terminal, Play, Pause, Trash2, Wifi } from "lucide-react"

interface LogEntry {
  id: string
  timestamp: string
  eventType: string
  message: string
}

export function LiveTerminal() {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [isPaused, setIsPaused] = useState(false)
  const [windowHeight, setWindowHeight] = useState(0)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const updateWindowHeight = () => {
      setWindowHeight(window.innerHeight)
    }

    updateWindowHeight()

    window.addEventListener("resize", updateWindowHeight)

    return () => {
      window.removeEventListener("resize", updateWindowHeight)
    }
  }, [])

  const terminalMaxHeight = Math.max(200, windowHeight - 300)

  const getEventDetails = (eventType: string) => {
    switch (eventType) {
      case "EventPlayerWeaponUse":
        return { icon: "🔫", msg: "Player fired weapon", color: "text-red-400" }
      case "PlayerWeaponReloadStarted":
        return { icon: "🔄", msg: "Reload started", color: "text-yellow-400" }
      case "PlayerWeaponReloadFinished":
        return { icon: "✅", msg: "Reload finished", color: "text-green-400" }
      case "PlayerWeaponChanged":
        return { icon: "🔀", msg: "Switched weapon", color: "text-blue-400" }
      case "PlayerWeaponAdded":
        return { icon: "➕", msg: "Picked up weapon", color: "text-cyan-400" }
      case "PlayerWeaponRemoved":
        return { icon: "➖", msg: "Dropped weapon", color: "text-orange-400" }
      case "PlayerActiveWeaponSwitched":
        return { icon: "🎯", msg: "Switched active weapon", color: "text-purple-400" }
      case "EventPlayerHealthChanged":
        return { icon: "❤️", msg: "Health changed", color: "text-pink-400" }
      case "EventPlayerArmourChanged":
        return { icon: "🛡️", msg: "Armor changed", color: "text-gray-400" }
      case "EventPlayerAlivenessChanged":
        return { icon: "💀", msg: "Player died/respawned", color: "text-red-500" }
      case "EventBombPlanted":
        return { icon: "💣", msg: "Bomb planted", color: "text-orange-500" }
      case "EventBombDefused":
        return { icon: "🛡️", msg: "Bomb defused", color: "text-cyan-500" }
      case "EventBombExploded":
        return { icon: "🔥", msg: "Bomb exploded", color: "text-red-600" }
      case "EventPlayerPaused":
        return { icon: "⏸️", msg: "Player paused", color: "text-muted-foreground" }
      case "EventPlayerPlaying":
        return { icon: "▶️", msg: "Player resumed", color: "text-green-500" }
      case "EventPlayerInTextInput":
        return { icon: "⌨️", msg: "Player typing", color: "text-blue-500" }
      case "HeartBeat":
        return { icon: "❤️‍🔥", msg: "Heartbeat", color: "text-muted-foreground" }
      default:
        return { icon: "📝", msg: eventType, color: "text-foreground" }
    }
  }

  useEffect(() => {
    if (!isConnected) return

    const ws = new WebSocket("ws://localhost:8080/ws")

    ws.onopen = () => {
      console.log("WebSocket connected")
    }

    ws.onmessage = (event) => {
      if (isPaused) return

      try {
        const log = JSON.parse(event.data) as { event_type: string; time: string }
        const details = getEventDetails(log.event_type)

        setLogs((prev) => [
          ...prev.slice(-49),
          {
            id: `${Date.now()}-${Math.random()}`,
            timestamp: new Date(log.time).toLocaleTimeString() || log.time,
            eventType: log.event_type,
            message: details.msg,
          },
        ])
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error)
      }
    }

    ws.onclose = () => {
      console.log("WebSocket disconnected")
    }

    ws.onerror = (error) => {
      console.error("WebSocket error:", error)
    }

    return () => {
      ws.close()
    }
  }, [isConnected, isPaused])

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [logs])

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3 flex-shrink-0">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Terminal className="h-5 w-5" />
            Live Events
          </CardTitle>
          <div className="flex items-center gap-2">
            <Badge variant={isConnected ? "default" : "secondary"} className="gap-1">
              <Wifi className="h-3 w-3" />
              {isConnected ? "Connected" : "Disconnected"}
            </Badge>
          </div>
        </div>

        <div className="flex items-center gap-2 pt-2">
          <Button size="sm" variant="outline" onClick={() => setIsConnected(!isConnected)} className="gap-2">
            {isConnected ? <Pause className="h-3 w-3" /> : <Play className="h-3 w-3" />}
            {isConnected ? "Disconnect" : "Connect"}
          </Button>

          <Button
            size="sm"
            variant="outline"
            onClick={() => setIsPaused(!isPaused)}
            disabled={!isConnected}
            className="gap-2"
          >
            {isPaused ? <Play className="h-3 w-3" /> : <Pause className="h-3 w-3" />}
            {isPaused ? "Resume" : "Pause"}
          </Button>

          <Button size="sm" variant="outline" onClick={() => setLogs([])} className="gap-2">
            <Trash2 className="h-3 w-3" />
            Clear
          </Button>
        </div>
      </CardHeader>

      <CardContent className="flex-1 p-0 min-h-0">
        <div
          ref={scrollRef}
          style={{ maxHeight: `${terminalMaxHeight}px` }}
          className="overflow-y-auto px-4 pb-4 scrollbar-thin scrollbar-thumb-border scrollbar-track-transparent"
        >
          <div className="space-y-2 font-mono text-sm">
            {logs.length === 0 ? (
              <div className="flex items-center justify-center h-32 text-muted-foreground">
                {isConnected ? "Waiting for events..." : "Connect to start streaming events"}
              </div>
            ) : (
              logs.map((log) => {
                const details = getEventDetails(log.eventType)
                return (
                  <div
                    key={log.id}
                    className="flex items-start gap-2 p-2 rounded bg-accent/30 hover:bg-accent/50 transition-colors"
                  >
                    <span className="text-xs text-muted-foreground mt-0.5 min-w-[60px]">{log.timestamp}</span>
                    <span className="mt-0.5">{details.icon}</span>
                    <span className={`flex-1 ${details.color}`}>{details.msg}</span>
                  </div>
                )
              })
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
