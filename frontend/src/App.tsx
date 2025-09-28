import { useState } from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import Realtime from "./pages/Realtime";
import Historical from "./pages/Historical";
import Terminal from "./components/Terminal";
import { Button } from "@/components/ui/button";
import './App.css'

function App() {
  const [terminalOpen, setTerminalOpen] = useState(false);

  return (
    <div className="h-screen w-screen flex flex-col bg-zinc-950 text-white">
      {/* Navbar */}
      <div className="flex items-center justify-between px-4 py-2 border-b border-zinc-800 bg-zinc-900">
        <h1 className="text-xl font-bold">CSYou Analytics</h1>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={() => setTerminalOpen(!terminalOpen)}
          >
            {terminalOpen ? "Hide Terminal" : "Show Terminal"}
          </Button>
          <Button variant="destructive">
            Clear Cache (120MB) {/* TODO: wire up */}
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="realtime" className="flex-1 flex flex-col">
        <TabsList className="bg-zinc-800 px-4 py-2 border-b border-zinc-700">
          <TabsTrigger value="realtime">Realtime (Redis)</TabsTrigger>
          <TabsTrigger value="historical">Historical (ClickHouse)</TabsTrigger>
        </TabsList>

        <TabsContent value="realtime" className="flex-1 overflow-y-auto">
          <Realtime />
        </TabsContent>
        <TabsContent value="historical" className="flex-1 overflow-y-auto">
          <Historical />
        </TabsContent>
      </Tabs>

      {/* Terminal */}
      {terminalOpen && (
        <div className="h-40 border-t border-zinc-700 bg-black">
          <Terminal />
        </div>
      )}
    </div>
  );
}

export default App;
