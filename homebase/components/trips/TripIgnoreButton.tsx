import {useRouter} from "next/router";
import {useIgnoreTrip} from "../../mutations/IgnoreTrip";

export default function TripIgnoreButton({id}: { id: string }) {
    const router = useRouter();
    const [commit, isInFlight] = useIgnoreTrip();

    async function onIgnore() {
        try {
            await commit(id);

            // return to the trips page upon successful ignore
            await router.push("/trips");
        } catch (e) {
            console.error(e);
        }
    }

    return (
        <span className="inline-flex rounded-md shadow-sm">
<button type="button"
        disabled={isInFlight}
        onClick={onIgnore}
        className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm leading-5 font-medium rounded-md text-gray-700 bg-white hover:text-gray-500 focus:outline-none focus:ring-blue focus:border-blue-300 active:bg-gray-50 active:text-gray-800">
<svg className="-ml-1 mr-2 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg"
     viewBox="0 0 20 20"
     fill="currentColor">
<path fillRule="evenodd"
      d="M10 18a8 8 0 100-16 8 8 0 000 16zM7 9a1 1 0 000 2h6a1 1 0 100-2H7z"
      clipRule="evenodd"/>
</svg>
<span>Ignore</span>
</button>
</span>
    );
}
