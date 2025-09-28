import ChartCard from "../components/ChartCard";
import SearchBar from "../components/SearchBar";

export default function Realtime() {
  return (
    <div className="p-4 grid grid-cols-12 gap-4">
      {/* Search */}
      <div className="col-span-12">
        <SearchBar placeholder="Search live matches, rounds, players..." />
      </div>

      {/* Charts */}
      <ChartCard title="Kill by Kill Analysis" className="col-span-6" />
      <ChartCard title="Round by Round Analysis" className="col-span-6" />
      <ChartCard title="Match by Match Analysis" className="col-span-12" />
    </div>
  );
}
