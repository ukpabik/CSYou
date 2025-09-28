import { Input } from "@/components/ui/input";

interface SearchBarProps {
  placeholder?: string;
}

export default function SearchBar({ placeholder }: SearchBarProps) {
  return (
    <div className="flex items-center">
      <Input
        type="text"
        placeholder={placeholder}
        className="w-full bg-zinc-900 border-zinc-700 text-white"
      />
    </div>
  );
}
