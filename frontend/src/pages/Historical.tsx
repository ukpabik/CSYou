import ChartCard from "../components/ChartCard";
import SearchBar from "../components/SearchBar";

export default function Historical() {
  return (
    <div className="p-4 grid grid-cols-12 gap-4">
      {/* Search */}
      <div className="col-span-12">
        <SearchBar placeholder="Search past matches (e.g. Dust2 on 2025-09-27)" />
      </div>

      {/* Charts */}
      <ChartCard title="Kill by Kill Trends" className="col-span-6" />
      <ChartCard title="Round by Round Trends" className="col-span-6" />
      <ChartCard title="Match History Overview" className="col-span-12" />
    </div>
  );
}
