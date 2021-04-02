import React from "react";

const severityStyles = {
    warning: {
        bg: "bg-yellow-50",
        iconColor: "text-yellow-400",
        titleColor: "text-yellow-800",
        bodyColor: "text-yellow-700",
        icon: <path fillRule="evenodd"
                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                    clipRule="evenodd"/>,
    },
    error: {
        bg: "bg-red-50",
        iconColor: "text-red-400",
        titleColor: "text-red-800",
        bodyColor: "text-red-700",
        icon: <path fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"/>,
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
                    <svg className={`h-5 w-5 ${styles.iconColor}`} xmlns="http://www.w3.org/2000/svg"
                         viewBox="0 0 20 20"
                         fill="currentColor">
                        {styles.icon}
                    </svg>
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
