import {Field} from "formik";

export default function DescriptionField() {
    return (
        <div className="sm:col-span-6">
            <label htmlFor="description"
                   className="block text-sm font-medium leading-5 text-gray-700">
                Description
            </label>
            <div className="mt-1">
                <Field as="textarea"
                       name="description"
                       rows={3}
                       className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                />
            </div>
        </div>
    );
}
