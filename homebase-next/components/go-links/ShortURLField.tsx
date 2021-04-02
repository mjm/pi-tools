import {Field} from "formik";

export default function ShortURLField() {
    return (
        <div className="sm:col-span-4">
            <label htmlFor="shortURL"
                   className="block text-sm font-medium leading-5 text-gray-700">
                Short URL
            </label>
            <div className="mt-1 flex rounded-md shadow-sm">
            <span
                className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 sm:text-sm">
              go/
            </span>
                <Field name="shortURL"
                       type="text"
                       className="flex-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300"/>
            </div>
        </div>
    );
}
