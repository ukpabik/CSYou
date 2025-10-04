"use client"

import { useEffect, useState, useRef } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  ComposedChart,
} from "recharts"
import {
  Activity,
  Target,
  Clock,
  Crosshair,
  RefreshCw,
  DollarSign,
  Heart,
  Users,
  Inbox,
  AlertCircle,
  Search,
} from "lucide-react"

interface AnalyticsChartsProps {
  dataSource: "redis" | "clickhouse"
  pollInterval?: number
}

interface KillEvent {
  match_id: string
  round: number
  map: string
  team: string
  steamid: string
  name: string
  mode: string
  active_gun: {
    name: string
    type: string
    ammo: number
    reserve: number
    skin: string
    headshot: boolean
  }
  timestamp: number
}

interface PlayerEvent {
  match_id: string
  round: number
  map: string
  team: string
  steamid: string
  name: string
  mode: string
  health: number
  armor: number
  helmet: boolean
  money: number
  equip_value: number
  round_kills: number
  round_killhs: number
  kills: number
  assists: number
  deaths: number
  mvps: number
  score: number
  timestamp: number
  win_team: string
}

export function AnalyticsCharts({ dataSource, pollInterval = 2000 }: AnalyticsChartsProps) {
  const [killEvents, setKillEvents] = useState<KillEvent[]>([])
  const [playerEvents, setPlayerEvents] = useState<PlayerEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date())
  const [isRefreshing, setIsRefreshing] = useState(false)
  const intervalRef = useRef<NodeJS.Timeout | null>(null)

  // query filters
  const [filters, setFilters] = useState({
    match_id: "",
    round: "",
    weapon_name: "",
    headshot: "",
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters({ ...filters, [e.target.name]: e.target.value })
  }

  const fetchData = async (showRefreshing = false) => {
    try {
      if (showRefreshing) setIsRefreshing(true)

      const baseUrl =
        dataSource === "redis" ? "http://localhost:8080/redis" : "http://localhost:8080/db"

      const queryString = new URLSearchParams(
        Object.entries(filters).filter(([_, v]) => v !== "")
      ).toString()

      const killUrl =
        queryString.length > 0
          ? `${baseUrl}/kill-events/params?${queryString}`
          : `${baseUrl}/kill-events`

      const playerUrl =
        queryString.length > 0
          ? `${baseUrl}/player-events/params?${queryString}`
          : `${baseUrl}/player-events`

      const [killResponse, playerResponse] = await Promise.all([
        fetch(killUrl),
        fetch(playerUrl),
      ])

      if (!killResponse.ok) throw new Error("Failed to fetch kill events")
      if (!playerResponse.ok) throw new Error("Failed to fetch player events")

      const [rawKillData, rawPlayerData] = await Promise.all([
        killResponse.json(),
        playerResponse.json(),
      ])

      let killData: KillEvent[]
      let playerData: PlayerEvent[]

      const toMs = (v: any) => {
        const n = typeof v === "string" ? parseInt(v, 10) : Number(v)
        if (!Number.isFinite(n)) return Date.now()
        return n < 1e12 ? n * 1000 : n
      }

      if (dataSource === "redis") {
        killData = rawKillData.map((e: any) => ({ ...e, timestamp: toMs(e.timestamp) }))
        playerData = rawPlayerData.map((e: any) => ({ ...e, timestamp: toMs(e.timestamp) }))
      } else {
        const wn = (e: any) => e.weapon_name ?? e.WeaponName ?? e.weapon?.weapon_name ?? ""
        const wt = (e: any) => e.weapon_type ?? e.WeaponType ?? e.weapon?.weapon_type ?? ""
        const wa = (e: any) => e.weapon_ammo ?? e.WeaponAmmo ?? e.weapon?.weapon_ammo ?? 0
        const wr = (e: any) => e.weapon_reserve ?? e.WeaponReserve ?? e.weapon?.weapon_reserve ?? 0
        const ws = (e: any) => e.weapon_skin ?? e.WeaponSkin ?? e.weapon?.weapon_skin ?? ""
        const wh = (e: any) => e.weapon_headshot ?? e.WeaponHeadshot ?? e.weapon?.weapon_headshot ?? false
        const ts = (e: any) =>
          e.timestamp ?? e.Timestamp ?? e.event_timestamp ?? e.EventTS

        killData = rawKillData.map((e: any) => ({
          match_id: e.match_id ?? e.MatchId,
          round: e.round ?? e.Round,
          map: e.map ?? e.Map,
          team: e.team ?? e.Team,
          steamid: e.steamid ?? e.SteamID,
          name: e.name ?? e.Name,
          mode: e.mode ?? e.Mode,
          active_gun: {
            name: wn(e),
            type: wt(e),
            ammo: Number(wa(e)) || 0,
            reserve: Number(wr(e)) || 0,
            skin: ws(e),
            headshot: Boolean(wh(e)),
          },
          timestamp: toMs(ts(e)),
        }))

        playerData = rawPlayerData.map((e: any) => ({
          match_id: e.match_id ?? e.MatchId,
          round: e.round ?? e.Round,
          map: e.map ?? e.Map,
          team: e.team ?? e.Team,
          steamid: e.steamid ?? e.SteamID,
          name: e.name ?? e.Name,
          mode: e.mode ?? e.Mode,
          health: Number(e.health ?? e.Health) || 0,
          armor: Number(e.armor ?? e.Armor) || 0,
          helmet: Boolean(e.helmet ?? e.Helmet),
          money: Number(e.money ?? e.Money) || 0,
          equip_value: Number(e.equip_value ?? e.EquipValue) || 0,
          round_kills: Number(e.round_kills ?? e.RoundKills) || 0,
          round_killhs: Number(e.round_killhs ?? e.RoundKillHS) || 0,
          kills: Number(e.kills ?? e.Kills) || 0,
          assists: Number(e.assists ?? e.Assists) || 0,
          deaths: Number(e.deaths ?? e.Deaths) || 0,
          mvps: Number(e.mvps ?? e.MVPs) || 0,
          score: Number(e.score ?? e.Score) || 0,
          timestamp: toMs(e.event_timestamp ?? e.EventTS ?? e.timestamp ?? e.Timestamp),
          win_team: e.win_team ?? e.WinTeam ?? "",
        }))
      }

      setKillEvents(Array.isArray(killData) && killData.length > 0 ? killData : [])
      setPlayerEvents(Array.isArray(playerData) && playerData.length > 0 ? playerData : [])

      setLastUpdate(new Date())
      setError(null)

      if (loading) setLoading(false)
      if (showRefreshing) setTimeout(() => setIsRefreshing(false), 300)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch data")
      if (loading) setLoading(false)
      if (showRefreshing) setIsRefreshing(false)
    }
  }

  useEffect(() => {
    fetchData()
    intervalRef.current = setInterval(() => fetchData(true), pollInterval)
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current)
    }
  }, [pollInterval, dataSource, filters])

  const calculateStats = () => {
    if (killEvents.length === 0 && playerEvents.length === 0) {
      return {
        totalKills: 0,
        headshotPercentage: "0",
        roundCount: 0,
        killsPerRound: "0.0",
        currentMoney: 0,
        currentHealth: 100,
        kdr: "0.00",
      }
    }

    const totalKills = killEvents.length
    const headshotKills = killEvents.filter((e) => e.active_gun.headshot).length
    const headshotPercentage = totalKills > 0 ? ((headshotKills / totalKills) * 100).toFixed(1) : "0"

    const uniqueRounds = new Set(killEvents.map((e) => e.round))
    const roundCount = uniqueRounds.size || 1

    const latestPlayerEvent =
      playerEvents.length > 0
        ? playerEvents.reduce((latest, current) =>
            current.timestamp > latest.timestamp ? current : latest
          )
        : null

    return {
      totalKills,
      headshotPercentage,
      roundCount,
      killsPerRound: (totalKills / roundCount).toFixed(1),
      currentMoney: latestPlayerEvent?.money || 0,
      currentHealth: latestPlayerEvent?.health || 100,
      kdr: latestPlayerEvent
        ? (latestPlayerEvent.kills / Math.max(latestPlayerEvent.deaths, 1)).toFixed(2)
        : "0.00",
    }
  }

  const processKillTimeline = () => {
    const timeline: { [key: string]: { kills: number; time: string } } = {}
    killEvents.forEach((event) => {
      const date = new Date(event.timestamp)
      const timeKey = `${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}`
      if (!timeline[timeKey]) {
        timeline[timeKey] = { time: timeKey, kills: 0 }
      }
      timeline[timeKey].kills++
    })
    return Object.values(timeline)
      .sort((a, b) => a.time.localeCompare(b.time))
      .slice(-10)
  }

  const processEconomyData = () => {
    const rounds = Array.from(new Set(playerEvents.map((e) => e.round))).sort((a, b) => a - b)
    return rounds.slice(-10).map((round) => {
      const roundPlayerEvents = playerEvents.filter((e) => e.round === round)
      const roundKillEvents = killEvents.filter((e) => e.round === round)

      const avgMoney =
        roundPlayerEvents.length > 0
          ? Math.round(roundPlayerEvents.reduce((sum, e) => sum + e.money, 0) / roundPlayerEvents.length)
          : 0

      const avgEquipValue =
        roundPlayerEvents.length > 0
          ? Math.round(roundPlayerEvents.reduce((sum, e) => sum + e.equip_value, 0) / roundPlayerEvents.length)
          : 0

      return {
        round,
        money: avgMoney,
        equipValue: avgEquipValue,
        kills: roundKillEvents.length,
        ecoRound: avgMoney < 3000,
      }
    })
  }

  const processHealthArmorData = () => {
    const rounds = Array.from(new Set(playerEvents.map((e) => e.round))).sort((a, b) => a - b)
    return rounds.slice(-10).map((round) => {
      const roundEvents = playerEvents.filter((e) => e.round === round)
      const avgHealth =
        roundEvents.length > 0
          ? Math.round(roundEvents.reduce((sum, e) => sum + e.health, 0) / roundEvents.length)
          : 100

      const avgArmor =
        roundEvents.length > 0 ? Math.round(roundEvents.reduce((sum, e) => sum + e.armor, 0) / roundEvents.length) : 0

      return {
        round,
        health: avgHealth,
        armor: avgArmor,
      }
    })
  }

  const processPerformanceData = () => {
    const rounds = Array.from(new Set(playerEvents.map((e) => e.round))).sort((a, b) => a - b)
    return rounds.slice(-10).map((round) => {
      const roundEvents = playerEvents.filter((e) => e.round === round)
      const latestInRound =
        roundEvents.length > 0
          ? roundEvents.reduce((latest, current) => (current.timestamp > latest.timestamp ? current : latest))
          : null

      return {
        round,
        kills: latestInRound?.kills || 0,
        deaths: latestInRound?.deaths || 0,
        assists: latestInRound?.assists || 0,
        score: latestInRound?.score || 0,
      }
    })
  }

  const processRecentKills = () => {
    return [...killEvents]
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, 5)
      .map((kill) => ({
        ...kill,
        timeAgo: getTimeAgo(new Date(kill.timestamp)),
      }))
  }

  const getTimeAgo = (date: Date) => {
    const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
    if (seconds < 60) return `${seconds}s ago`
    const minutes = Math.floor(seconds / 60)
    if (minutes < 60) return `${minutes}m ago`
    const hours = Math.floor(minutes / 60)
    return `${hours}h ago`
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="flex items-center gap-2">
          <RefreshCw className="h-4 w-4 animate-spin text-white" />
          <span className="text-white">Loading analytics data...</span>
        </div>
      </div>
    )
  }

  // Show empty state if no data
  if (killEvents.length === 0 && playerEvents.length === 0 && !loading) {
    return (
      <div className="flex items-center justify-center min-h-[500px]">
        <Card className="bg-gray-900 border-gray-700 max-w-md">
          <CardContent className="pt-6">
            <div className="flex flex-col items-center text-center space-y-4">
              <div className="rounded-full bg-gray-800 p-4">
                <Inbox className="h-12 w-12 text-gray-400" />
              </div>
              <div className="space-y-2">
                <h3 className="text-xl font-semibold text-white">No Data Available</h3>
                <p className="text-sm text-gray-400">
                  There's no analytics data to display yet. Start playing some matches to see your stats here!
                </p>
              </div>
              {error && (
                <div className="flex items-center gap-2 text-xs text-amber-400 bg-amber-950/20 px-3 py-2 rounded-md">
                  <AlertCircle className="h-4 w-4" />
                  <span>Connection issue</span>
                </div>
              )}
              <div className="text-xs text-gray-500">
                Auto-refreshing every {pollInterval / 1000}s...
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (error && killEvents.length === 0 && playerEvents.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-red-400">Error: {error}</div>
      </div>
    )
  }

  const stats = calculateStats()
  const killTimeline = processKillTimeline()
  const economyData = processEconomyData()
  const recentKills = processRecentKills()
  const healthArmorData = processHealthArmorData()
  const performanceData = processPerformanceData()

  return (
    <div className="grid gap-6">
      <Card className="bg-gray-900 border-gray-700">
        <CardContent className="flex flex-wrap gap-2 py-4">
          <input
            name="match_id"
            value={filters.match_id}
            onChange={handleChange}
            placeholder="Match ID"
            className="px-2 py-1 rounded bg-gray-800 text-white text-sm"
          />
          <input
            name="round"
            value={filters.round}
            onChange={handleChange}
            placeholder="Round"
            className="px-2 py-1 rounded bg-gray-800 text-white text-sm"
          />
          <input
            name="weapon_name"
            value={filters.weapon_name}
            onChange={handleChange}
            placeholder="Weapon"
            className="px-2 py-1 rounded bg-gray-800 text-white text-sm"
          />
          <input
            name="headshot"
            value={filters.headshot}
            onChange={handleChange}
            placeholder="Headshot (true/false)"
            className="px-2 py-1 rounded bg-gray-800 text-white text-sm"
          />
          <button
            onClick={() => fetchData(true)}
            className="flex items-center gap-1 px-3 py-1 bg-blue-600 text-white rounded text-sm"
          >
            <Search className="h-4 w-4" /> Search
          </button>
        </CardContent>
      </Card>
      <div className="flex items-center justify-between text-sm">
        <div className="flex items-center gap-2">
          <div className={`h-2 w-2 rounded-full ${isRefreshing ? "bg-yellow-400 animate-pulse" : "bg-green-400"}`} />
          <span className="text-gray-300">
            Auto-refresh every {pollInterval / 1000}s • Last update: {lastUpdate.toLocaleTimeString()}
          </span>
        </div>
        {error && <span className="text-xs text-red-400">No events found... Try a different query</span>}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-6 gap-4">
        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">Total Kills</CardTitle>
            <Target className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.totalKills}</div>
            <p className="text-xs text-gray-400">{stats.killsPerRound} per round</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">K/D Ratio</CardTitle>
            <Users className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.kdr}</div>
            <p className="text-xs text-gray-400">Kill/Death ratio</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">Headshot %</CardTitle>
            <Crosshair className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.headshotPercentage}%</div>
            <p className="text-xs text-gray-400">Precision rating</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">Current Money</CardTitle>
            <DollarSign className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">${stats.currentMoney.toLocaleString()}</div>
            <p className="text-xs text-gray-400">Available funds</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">Health</CardTitle>
            <Heart className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.currentHealth}</div>
            <p className="text-xs text-gray-400">Current HP</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-white">Rounds</CardTitle>
            <Clock className="h-4 w-4 text-gray-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.roundCount}</div>
            <p className="text-xs text-gray-400">Current session</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="bg-gray-900 border-gray-700">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-white">
              <Crosshair className="h-5 w-5" />
              Kill Timeline
            </CardTitle>
            <CardDescription className="text-gray-400">Kills over time (last 10 intervals)</CardDescription>
          </CardHeader>
          <CardContent>
            {killTimeline.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={killTimeline}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="time" stroke="#9ca3af" />
                  <YAxis stroke="#9ca3af" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#1f2937",
                      border: "1px solid #374151",
                      borderRadius: "8px",
                      color: "#ffffff",
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="kills"
                    stroke="#60a5fa"
                    strokeWidth={3}
                    dot={{ fill: "#60a5fa", strokeWidth: 2, r: 4 }}
                    animationDuration={300}
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-gray-400">No kill data available</div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-white">
              <DollarSign className="h-5 w-5" />
              Economy Analysis
            </CardTitle>
            <CardDescription className="text-gray-400">Money and equipment value by round</CardDescription>
          </CardHeader>
          <CardContent>
            {economyData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <ComposedChart data={economyData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="round" stroke="#9ca3af" />
                  <YAxis stroke="#9ca3af" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#1f2937",
                      border: "1px solid #374151",
                      borderRadius: "8px",
                      color: "#ffffff",
                    }}
                    formatter={(value: any, name: string) => [
                      name.includes("money") || name.includes("Value") ? `$${value.toLocaleString()}` : value,
                      name === "money" ? "Money" : name === "equipValue" ? "Equipment Value" : name,
                    ]}
                  />
                  <Area
                    type="monotone"
                    dataKey="money"
                    stroke="#34d399"
                    fill="#34d399"
                    fillOpacity={0.3}
                    strokeWidth={2}
                  />
                  <Line
                    type="monotone"
                    dataKey="equipValue"
                    stroke="#fbbf24"
                    strokeWidth={2}
                    dot={{ fill: "#fbbf24", strokeWidth: 2, r: 3 }}
                  />
                </ComposedChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-gray-400">No economy data available</div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-white">
              <Heart className="h-5 w-5" />
              Health & Armor
            </CardTitle>
            <CardDescription className="text-gray-400">Player survivability by round</CardDescription>
          </CardHeader>
          <CardContent>
            {healthArmorData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={healthArmorData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="round" stroke="#9ca3af" />
                  <YAxis stroke="#9ca3af" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#1f2937",
                      border: "1px solid #374151",
                      borderRadius: "8px",
                      color: "#ffffff",
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="health"
                    stroke="#f87171"
                    strokeWidth={2}
                    dot={{ fill: "#f87171", strokeWidth: 2, r: 3 }}
                  />
                  <Line
                    type="monotone"
                    dataKey="armor"
                    stroke="#60a5fa"
                    strokeWidth={2}
                    dot={{ fill: "#60a5fa", strokeWidth: 2, r: 3 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-gray-400">
                No health/armor data available
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-white">
              <Users className="h-5 w-5" />
              Performance Tracking
            </CardTitle>
            <CardDescription className="text-gray-400">KDA progression over rounds</CardDescription>
          </CardHeader>
          <CardContent>
            {performanceData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={performanceData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="round" stroke="#9ca3af" />
                  <YAxis stroke="#9ca3af" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#1f2937",
                      border: "1px solid #374151",
                      borderRadius: "8px",
                      color: "#ffffff",
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="kills"
                    stroke="#34d399"
                    strokeWidth={2}
                    dot={{ fill: "#34d399", strokeWidth: 2, r: 3 }}
                  />
                  <Line
                    type="monotone"
                    dataKey="deaths"
                    stroke="#f87171"
                    strokeWidth={2}
                    dot={{ fill: "#f87171", strokeWidth: 2, r: 3 }}
                  />
                  <Line
                    type="monotone"
                    dataKey="assists"
                    stroke="#fbbf24"
                    strokeWidth={2}
                    dot={{ fill: "#fbbf24", strokeWidth: 2, r: 3 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-gray-400">
                No performance data available
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card className="bg-gray-900 border-gray-700">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-white">
            <Activity className="h-5 w-5" />
            Recent Kills
          </CardTitle>
          <CardDescription className="text-gray-400">Latest kill feed</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {recentKills.length > 0 ? (
              recentKills.map((kill, index) => (
                <div
                  key={`${kill.timestamp}-${index}`}
                  className="flex items-center justify-between p-2 rounded bg-gray-800 text-sm"
                >
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-white">{kill.name}</span>
                    {kill.active_gun.headshot && (
                      <Badge variant="outline" className="text-xs border-red-400 text-red-400">
                        HS
                      </Badge>
                    )}
                    <span className="text-gray-400">killed with</span>
                    <span className="text-blue-400">{kill.active_gun.name || "Unknown"}</span>
                  </div>
                  <span className="text-xs text-gray-500">{kill.timeAgo}</span>
                </div>
              ))
            ) : (
              <div className="text-center text-gray-400 py-4">No recent kills</div>
            )}
          </div>
        </CardContent>
      </Card>

      <Card className="bg-gray-900 border-gray-700">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-white">
            <Activity className="h-5 w-5" />
            Event Viewer
          </CardTitle>
          <CardDescription className="text-gray-400">
            Showing {killEvents.length} kill events and {playerEvents.length} player events
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-h-[400px] overflow-y-auto">
            {/* Kill Events */}
            <div>
              <h4 className="text-white font-semibold mb-2">Kill Events</h4>
              {killEvents.length > 0 ? (
                <ul className="space-y-1 text-sm">
                  {killEvents.map((e, i) => (
                    <li key={i} className="p-2 rounded bg-gray-800 flex justify-between">
                      <span className="text-white">
                        [{e.round}] {e.name} ({e.team}) → {e.active_gun.name}{" "}
                        {e.active_gun.headshot && (
                          <Badge variant="outline" className="text-xs border-red-400 text-red-400 ml-1">
                            HS
                          </Badge>
                        )}
                      </span>
                      <span className="text-gray-400 text-xs">
                        {new Date(e.timestamp).toLocaleTimeString()}
                      </span>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-gray-400 text-sm">No kill events found</p>
              )}
            </div>

            {/* Player Events */}
            <div>
              <h4 className="text-white font-semibold mb-2">Player Events</h4>
              {playerEvents.length > 0 ? (
                <ul className="space-y-1 text-sm">
                  {playerEvents.map((e, i) => (
                    <li key={i} className="p-2 rounded bg-gray-800 flex justify-between">
                      <span className="text-white">
                        [{e.round}] {e.name} ({e.team}) • HP:{e.health} • ${e.money}
                      </span>
                      <span className="text-gray-400 text-xs">
                        {new Date(e.timestamp).toLocaleTimeString()}
                      </span>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-gray-400 text-sm">No player events found</p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}