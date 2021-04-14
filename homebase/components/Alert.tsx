import React from "react";
import {ExclamationIcon, XCircleIcon} from "@heroicons/react/solid";

const severityStyles = {
    warning: {
        bg: "bg-yellow-50",
        titleColor: "text-yellow-800",
        bodyColor: "text-yellow-700",
        icon: <ExclamationIcon className="h-5 w-5 text-yellow-400"/>,
    },
    error: {
        bg: "bg-red-50",
        titleColor: "text-red-800",
        bodyColor: "text-red-700",
        icon: <XCircleIcon className="h-5 w-5 text-red-400"/>,
    },
} as const;

export default function Alert({severity = "warning", rounded = true, title, children}: {
    severity?: keyof typeof severityStyles;
    rounded?: boolean;
    title: React.ReactNode;
    children: React.ReactNode;
}) {
    const styles = severityStyles[severity];
    return (
        <div className={`${rounded ? "rounded-md" : ""} ${styles.bg} p-4`}>
            <div className="flex">
                <div className="flex-shrink-0">
                    {styles.icon}
                </div>
                <div className="ml-3">
                    <h3 className={`text-sm leading-5 font-medium ${styles.titleColor}`}>
                        {title}
                    </h3>
                    <div className={`mt-2 text-sm leading-5 ${styles.bodyColor}`}>
                        {children}
                    </div>
                </div>
            </div>
        </div>
    );
}
