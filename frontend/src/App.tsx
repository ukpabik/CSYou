"use client"

import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Database, Zap, Terminal, Target } from "lucide-react"
import { AnalyticsCharts } from "@/components/analytics-charts"
import { LiveTerminal } from "@/components/live-terminal"
import { SearchFilters } from "@/components/search-filters"
import { CacheManager } from "@/components/cache-manager"
import "./App.css"

export default function CSGOAnalytics() {
  const [activeTab, setActiveTab] = useState("redis")
  const [terminalOpen, setTerminalOpen] = useState(false)

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card/50 backdrop-blur-sm">
        <div className="flex h-16 items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Target className="h-8 w-8 text-primary" />
              <h1 className="text-xl font-semibold text-foreground">CS:GO Analytics</h1>
            </div>
            <Badge variant="secondary" className="text-xs">
              Live Data
            </Badge>
          </div>

          <div className="flex items-center gap-3">
            <CacheManager />
            <Button
              variant={terminalOpen ? "default" : "outline"}
              size="sm"
              onClick={() => setTerminalOpen(!terminalOpen)}
              className="gap-2"
            >
              <Terminal className="h-4 w-4" />
              Terminal
            </Button>
          </div>
        </div>
      </header>

      <div className="flex">
        <main className="flex-1 p-6">
          <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
            <TabsList className="grid w-full max-w-md grid-cols-2">
              <TabsTrigger value="redis" className="gap-2">
                <Zap className="h-4 w-4" />
                Redis Data
              </TabsTrigger>
              <TabsTrigger value="clickhouse" className="gap-2">
                <Database className="h-4 w-4" />
                ClickHouse
              </TabsTrigger>
            </TabsList>

            <TabsContent value="redis" className="space-y-6">
              <div className="grid gap-6">
                <SearchFilters />
                <AnalyticsCharts dataSource="redis" />
              </div>
            </TabsContent>

            <TabsContent value="clickhouse" className="space-y-6">
              <div className="grid gap-6">
                <Card className="border-dashed">
                  <CardContent className="flex flex-col items-center justify-center py-12">
                    <Database className="h-12 w-12 text-muted-foreground mb-4" />
                    <h3 className="text-lg font-semibold mb-2">ClickHouse Integration</h3>
                    <p className="text-muted-foreground text-center max-w-md">
                      ClickHouse data source will be available soon. Configure your connection to start analyzing
                      historical match data.
                    </p>
                    <Button className="mt-4" disabled>
                      Configure ClickHouse
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </Tabs>
        </main>

        {terminalOpen && (
          <aside className="w-96 border-l border-border bg-card/30">
            <LiveTerminal />
          </aside>
        )}
      </div>
    </div>
  )
}
