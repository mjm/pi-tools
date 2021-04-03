import Spinner from "./Spinner";

export default function Loading() {
    return (
        <div className="flex flex-col items-center p-8 space-y-2 text-gray-700">
            <Spinner className="w-12 h-12" />
            <div className="font-medium">Loading…</div>
        </div>
    )
}
