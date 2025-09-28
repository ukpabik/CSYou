"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from "recharts"
import { Activity, Target, Trophy, Clock, Crosshair, Shield, Zap } from "lucide-react"

interface AnalyticsChartsProps {
  dataSource: "redis" | "clickhouse"
}

// Mock data - replace this with real API calls
const killData = [
  { time: "00:30", kills: 2, deaths: 1, assists: 0 },
  { time: "01:15", kills: 5, deaths: 2, assists: 1 },
  { time: "02:00", kills: 8, deaths: 3, assists: 2 },
  { time: "02:45", kills: 12, deaths: 4, assists: 3 },
  { time: "03:30", kills: 15, deaths: 6, assists: 4 },
]

const roundData = [
  { round: 1, score: 1, economy: 4000, winner: "CT" },
  { round: 2, score: 2, economy: 3200, winner: "T" },
  { round: 3, score: 2, economy: 5500, winner: "CT" },
  { round: 4, score: 3, economy: 2800, winner: "T" },
  { round: 5, score: 4, economy: 4800, winner: "CT" },
]

const weaponData = [
  { weapon: "AK-47", kills: 45, percentage: 35 },
  { weapon: "M4A4", kills: 38, percentage: 30 },
  { weapon: "AWP", kills: 25, percentage: 20 },
  { weapon: "Glock", kills: 12, percentage: 10 },
  { weapon: "Other", kills: 6, percentage: 5 },
]

const COLORS = ["#8b5cf6", "#06d6a0", "#f72585", "#ffbe0b", "#fb8500"]

// TODO: Update this with a real data source
export function AnalyticsCharts({ dataSource }: AnalyticsChartsProps) {
  return (
    <div className="grid gap-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Kills</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,247</div>
            <p className="text-xs text-muted-foreground">+12% from last match</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">K/D Ratio</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1.34</div>
            <p className="text-xs text-muted-foreground">+0.08 improvement</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Win Rate</CardTitle>
            <Trophy className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">67%</div>
            <p className="text-xs text-muted-foreground">Last 30 matches</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Match Time</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">34m</div>
            <p className="text-xs text-muted-foreground">-2m faster than avg</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Crosshair className="h-5 w-5" />
              Kill Timeline
            </CardTitle>
            <CardDescription>Real-time kill tracking during matches</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={killData}>
                <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
                <XAxis dataKey="time" stroke="hsl(var(--muted-foreground))" />
                <YAxis stroke="hsl(var(--muted-foreground))" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "hsl(var(--card))",
                    border: "1px solid hsl(var(--border))",
                    borderRadius: "8px",
                  }}
                />
                <Line type="monotone" dataKey="kills" stroke="hsl(var(--chart-1))" strokeWidth={2} />
                <Line type="monotone" dataKey="deaths" stroke="hsl(var(--chart-4))" strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5" />
              Round Economy
            </CardTitle>
            <CardDescription>Economic performance by round</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={roundData}>
                <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
                <XAxis dataKey="round" stroke="hsl(var(--muted-foreground))" />
                <YAxis stroke="hsl(var(--muted-foreground))" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "hsl(var(--card))",
                    border: "1px solid hsl(var(--border))",
                    borderRadius: "8px",
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="economy"
                  stroke="hsl(var(--chart-2))"
                  fill="hsl(var(--chart-2))"
                  fillOpacity={0.3}
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Zap className="h-5 w-5" />
              Weapon Usage
            </CardTitle>
            <CardDescription>Most effective weapons by kill count</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={weaponData}
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="kills"
                  label={({ payload }: any) => `${payload.weapon} (${payload.percentage}%)`}
                >
                  {weaponData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: "hsl(var(--card))",
                    border: "1px solid hsl(var(--border))",
                    borderRadius: "8px",
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Trophy className="h-5 w-5" />
              Match Performance
            </CardTitle>
            <CardDescription>Recent match results and trends</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[
                { map: "de_dust2", result: "Win", score: "16-12", kd: "1.8" },
                { map: "de_mirage", result: "Loss", score: "14-16", kd: "1.2" },
                { map: "de_inferno", result: "Win", score: "16-8", kd: "2.1" },
                { map: "de_cache", result: "Win", score: "16-10", kd: "1.6" },
              ].map((match, index) => (
                <div key={index} className="flex items-center justify-between p-3 rounded-lg bg-accent/50">
                  <div className="flex items-center gap-3">
                    <img
                      src={`/.jpg?height=32&width=32&query=${match.map}`}
                      alt={match.map}
                      className="w-8 h-8 rounded"
                    />
                    <div>
                      <p className="font-medium">{match.map}</p>
                      <p className="text-sm text-muted-foreground">{match.score}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <Badge variant={match.result === "Win" ? "default" : "destructive"}>{match.result}</Badge>
                    <p className="text-sm text-muted-foreground mt-1">K/D: {match.kd}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
