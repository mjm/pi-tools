import React from "react";

export default function DescriptionField({label, children, offset}: {
    label: React.ReactNode;
    children: React.ReactNode;
    offset?: boolean
}) {
    return (
        <div className={`${offset ? "bg-gray-50" : "bg-white"} px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6`}>
            <dt className="text-sm leading-5 font-medium text-gray-500">
                {label}
            </dt>
            <dd className="mt-1 text-sm leading-5 text-gray-900 sm:mt-0 sm:col-span-2">
                {children}
            </dd>
        </div>
    );
}
