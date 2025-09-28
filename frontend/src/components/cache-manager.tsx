"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Trash2, Database, RefreshCw } from "lucide-react"

export function CacheManager() {
  const [cacheSize, setCacheSize] = useState("0 MB")
  const [isClearing, setIsClearing] = useState(false)
  const [lastCleared, setLastCleared] = useState<Date | null>(null)

  useEffect(() => {
    const updateCacheSize = () => {
      const size = Math.floor(Math.random() * 500) + 50
      setCacheSize(`${size} MB`)
    }

    updateCacheSize()
    const interval = setInterval(updateCacheSize, 30000)

    return () => clearInterval(interval)
  }, [])

  const handleClearCache = async () => {
    setIsClearing(true)

    try {
      // TODO: Call backend API to clear cache

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 2000))

      setCacheSize("0 MB")
      setLastCleared(new Date())
    } catch (error) {
      console.error("Failed to clear cache:", error)
    } finally {
      setIsClearing(false)
    }
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2 bg-transparent">
          <Database className="h-4 w-4" />
          Cache: {cacheSize}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-80" align="end">
        <Card className="border-0 shadow-none">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <Database className="h-4 w-4" />
              Cache Management
            </CardTitle>
            <CardDescription>Manage Redis cache for optimal performance</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Current Size:</span>
              <Badge variant="secondary">{cacheSize}</Badge>
            </div>

            {lastCleared && (
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Last Cleared:</span>
                <span className="text-sm">{lastCleared.toLocaleTimeString()}</span>
              </div>
            )}

            <div className="flex gap-2">
              <Button onClick={handleClearCache} disabled={isClearing} className="flex-1 gap-2" variant="destructive">
                {isClearing ? <RefreshCw className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                {isClearing ? "Clearing..." : "Clear Cache"}
              </Button>
            </div>

            <div className="text-xs text-muted-foreground">
              Clearing cache will remove all stored match data and force fresh data retrieval from the database.
            </div>
          </CardContent>
        </Card>
      </PopoverContent>
    </Popover>
  )
}
