"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Search, Filter, Calendar, MapPin, User } from "lucide-react"

export function SearchFilters() {
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedMap, setSelectedMap] = useState("")
  const [selectedPlayer, setSelectedPlayer] = useState("")
  const [dateRange, setDateRange] = useState("")
  const [activeFilters, setActiveFilters] = useState<string[]>([])

  const maps = ["de_dust2", "de_mirage", "de_inferno", "de_cache", "de_overpass", "de_train", "de_nuke"]
  const players = ["Player1", "Player2", "Player3", "Player4", "Player5"]

  const handleSearch = () => {
    // TODO: Handle with backend API
    console.log("Searching with:", { searchQuery, selectedMap, selectedPlayer, dateRange })
  }

  const addFilter = (type: string, value: string) => {
    const filter = `${type}:${value}`
    if (!activeFilters.includes(filter)) {
      setActiveFilters([...activeFilters, filter])
    }
  }

  const removeFilter = (filter: string) => {
    setActiveFilters(activeFilters.filter((f) => f !== filter))
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Search className="h-5 w-5" />
          Search & Filters
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search matches, rounds, or specific events..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          <Button onClick={handleSearch} className="gap-2">
            <Search className="h-4 w-4" />
            Search
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Select value={selectedMap} onValueChange={setSelectedMap}>
            <SelectTrigger className="gap-2">
              <MapPin className="h-4 w-4" />
              <SelectValue placeholder="Select Map" />
            </SelectTrigger>
            <SelectContent>
              {maps.map((map) => (
                <SelectItem key={map} value={map}>
                  {map}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={selectedPlayer} onValueChange={setSelectedPlayer}>
            <SelectTrigger className="gap-2">
              <User className="h-4 w-4" />
              <SelectValue placeholder="Select Player" />
            </SelectTrigger>
            <SelectContent>
              {players.map((player) => (
                <SelectItem key={player} value={player}>
                  {player}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={dateRange} onValueChange={setDateRange}>
            <SelectTrigger className="gap-2">
              <Calendar className="h-4 w-4" />
              <SelectValue placeholder="Date Range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="today">Today</SelectItem>
              <SelectItem value="week">Last Week</SelectItem>
              <SelectItem value="month">Last Month</SelectItem>
              <SelectItem value="custom">Custom Range</SelectItem>
            </SelectContent>
          </Select>

          <Button variant="outline" className="gap-2 bg-transparent">
            <Filter className="h-4 w-4" />
            Advanced
          </Button>
        </div>

        {activeFilters.length > 0 && (
          <div className="flex flex-wrap gap-2">
            <span className="text-sm text-muted-foreground">Active filters:</span>
            {activeFilters.map((filter) => (
              <Badge key={filter} variant="secondary" className="gap-1">
                {filter}
                <button onClick={() => removeFilter(filter)} className="ml-1 hover:text-destructive">
                  ×
                </button>
              </Badge>
            ))}
          </div>
        )}

        <div className="flex flex-wrap gap-2">
          <span className="text-sm text-muted-foreground">Quick filters:</span>
          <Button size="sm" variant="outline" onClick={() => addFilter("weapon", "AWP")} className="h-7 text-xs">
            AWP Kills
          </Button>
          <Button size="sm" variant="outline" onClick={() => addFilter("event", "clutch")} className="h-7 text-xs">
            Clutch Rounds
          </Button>
          <Button size="sm" variant="outline" onClick={() => addFilter("result", "win")} className="h-7 text-xs">
            Won Matches
          </Button>
          <Button size="sm" variant="outline" onClick={() => addFilter("map", "de_dust2")} className="h-7 text-xs">
            Dust2 Only
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
