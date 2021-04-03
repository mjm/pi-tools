import {Field} from "formik";

export default function DestinationURLField() {
    return (
        <div className="sm:col-span-6">
            <label htmlFor="destinationURL"
                   className="block text-sm font-medium leading-5 text-gray-700">
                Destination URL
            </label>
            <div className="mt-1">
                <Field name="destinationURL"
                       type="url"
                       className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                       placeholder="https://www.google.com/"
                />
            </div>
        </div>
    );
}
