import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

interface ChartCardProps {
  title: string;
  className?: string;
}

export default function ChartCard({ title, className }: ChartCardProps) {
  return (
    <Card className={`bg-zinc-900 border-zinc-800 ${className}`}>
      <CardHeader>
        <CardTitle className="text-white">{title}</CardTitle>
      </CardHeader>
      <CardContent className="h-60 flex items-center justify-center text-zinc-400">
        <p>Chart Placeholder</p>
      </CardContent>
    </Card>
  );
}
