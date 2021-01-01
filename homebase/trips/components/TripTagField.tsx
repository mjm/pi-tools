import React from "react";
import {TripTag} from "com_github_mjm_pi_tools/homebase/trips/components/TripTag";
import {updateTripTags} from "com_github_mjm_pi_tools/homebase/trips/lib/mutate";
import {graphql, useFragment} from "react-relay/hooks";
import {TripTagField_trip$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/TripTagField_trip.graphql";

export function TripTagField({trip}: { trip: TripTagField_trip$key }) {
    const data = useFragment(
        graphql`
            fragment TripTagField_trip on Trip {
                rawID
                tags
            }
        `,
        trip,
    );

    const [draftTags, setDraftTags] = React.useState<string | null>(null);
    const [isSaving, setSaving] = React.useState(false);

    function onEdit() {
        setDraftTags(data.tags.join(", "));
    }

    function onChange(e: React.ChangeEvent<HTMLInputElement>) {
        setDraftTags(e.target.value);
    }

    async function onSave(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        setSaving(true);
        const newTags = draftTags.split(",").map(tag => tag.trim());

        try {
            await updateTripTags(data.rawID, data.tags, newTags);
            setDraftTags(null);
        } catch (e) {
            // TODO surface error
            console.error(e);
        } finally {
            setSaving(false);
        }
    }

    function onCancel() {
        setDraftTags(null);
    }

    return draftTags === null ? (
        <div className="group flex flex-row items-center space-x-3">
            {data.tags.length === 0 ? (
                <span>No tags</span>
            ) : (
                data.tags.map(tag => (
                    <TripTag key={tag}>
                        {tag}
                    </TripTag>
                ))
            )}
            <span className="invisible group-hover:visible inline-flex rounded-md shadow-sm">
  <button type="button"
          onClick={onEdit}
          className="inline-flex items-center px-2.5 py-1.5 border border-gray-300 text-xs leading-4 font-medium rounded text-gray-700 bg-white hover:text-gray-500 focus:outline-none focus:border-blue-300 focus:ring active:text-gray-800 active:bg-gray-50 transition ease-in-out duration-150">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="-ml-0.5 mr-2 h-4 w-4">
  <path
      d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"/>
</svg>
    Edit tags
  </button>
</span>
        </div>
    ) : (
        <form onSubmit={onSave}>
            <div className="flex flex-row space-x-3">
                <div>
                    <label htmlFor="tags" className="sr-only">Tags</label>
                    <div className="relative rounded-md shadow-sm">
                        <input id="tags" type="text"
                               className="max-w-lg block w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:max-w-xs sm:text-xs sm:leading-4 border-gray-300 rounded-md"
                               autoFocus
                               placeholder="Tags, comma-separated"
                               onChange={onChange}
                               value={draftTags}/>
                    </div>
                </div>
                <span className="inline-flex rounded-md shadow-sm">
  <button type="submit"
          disabled={isSaving}
          className="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs leading-4 font-medium rounded text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:border-indigo-700 focus:ring-indigo active:bg-indigo-700 transition ease-in-out duration-150">
    Save
  </button>
</span>
                <button type="button"
                        onClick={onCancel}
                        className="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs leading-4 font-medium rounded text-indigo-700 bg-indigo-100 hover:bg-indigo-50 focus:outline-none focus:border-indigo-300 focus:ring-indigo active:bg-indigo-200 transition ease-in-out duration-150">
                    Cancel
                </button>
            </div>
        </form>
    );
}
