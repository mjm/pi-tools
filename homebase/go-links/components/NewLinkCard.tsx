import React from "react";

export function NewLinkCard() {
    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="px-4 py-5 sm:px-6">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                    Add a new link
                </h3>
            </div>
            <div className="bg-gray-50 px-4 py-5 sm:p-6">
                <div className="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
                    <div className="sm:col-span-4">
                        <label htmlFor="short_url"
                               className="block text-sm font-medium leading-5 text-gray-700">
                            Short URL
                        </label>
                        <div className="mt-1 flex rounded-md shadow-sm">
            <span
                className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 sm:text-sm">
              go/
            </span>
                            <input id="short_url"
                                   className="flex-1 form-input block w-full min-w-0 rounded-none rounded-r-md transition duration-150 ease-in-out sm:text-sm sm:leading-5"/>
                        </div>
                    </div>

                    <div className="sm:col-span-6">
                        <label htmlFor="destination_url"
                               className="block text-sm font-medium leading-5 text-gray-700">
                            Destination URL
                        </label>
                        <div className="mt-1 rounded-md shadow-sm">
                            <input id="destination_url" type="url"
                                   className="form-input block w-full transition duration-150 ease-in-out sm:text-sm sm:leading-5"
                                   placeholder="https://www.google.com/"/>
                        </div>
                    </div>

                    <div className="sm:col-span-6">
                        <label htmlFor="description"
                               className="block text-sm font-medium leading-5 text-gray-700">
                            Description
                        </label>
                        <div className="mt-1 rounded-md shadow-sm">
                                        <textarea id="description" rows={3}
                                                  className="form-textarea block w-full transition duration-150 ease-in-out sm:text-sm sm:leading-5"/>
                        </div>
                    </div>
                </div>
            </div>
            <div className="px-4 py-5 sm:px-6 text-right">
            <span className="inline-flex rounded-md shadow-sm">
              <button type="submit"
                      className="inline-flex justify-center py-2 px-4 border border-transparent text-sm leading-5 font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:border-indigo-700 focus:shadow-outline-indigo active:bg-indigo-700 transition duration-150 ease-in-out">
                Create
              </button>
            </span>
            </div>
        </div>
    );
}
