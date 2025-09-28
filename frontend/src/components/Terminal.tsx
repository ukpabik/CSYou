import { useEffect, useState } from "react";

export default function Terminal() {
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    // TODO: replace with Kafka stream
    const interval = setInterval(() => {
      setLogs(prev => [...prev, `Event at ${new Date().toLocaleTimeString()}`]);
    }, 2000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="h-full w-full p-2 font-mono text-sm text-green-400 overflow-y-auto">
      {logs.map((line, i) => (
        <div key={i}>{line}</div>
      ))}
    </div>
  );
}
