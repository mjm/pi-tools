import React from "react";

const severityStyles = {
    warning: {
        color: "yellow",
        icon: <path fillRule="evenodd"
                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                    clipRule="evenodd"/>,
    },
    error: {
        color: "red",
        icon: <path fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"/>,
    },
} as const;

export function Alert({severity = "warning", rounded = true, title, children}: {
    severity?: keyof typeof severityStyles;
    rounded?: boolean;
    title: React.ReactNode;
    children: React.ReactNode;
}) {
    const styles = severityStyles[severity];
    return (
        <div className={`${rounded ? "rounded-md" : ""} bg-${styles.color}-50 p-4`}>
            <div className="flex">
                <div className="flex-shrink-0">
                    <svg className={`h-5 w-5 text-${styles.color}-400`} xmlns="http://www.w3.org/2000/svg"
                         viewBox="0 0 20 20"
                         fill="currentColor">
                        {styles.icon}
                    </svg>
                </div>
                <div className="ml-3">
                    <h3 className={`text-sm leading-5 font-medium text-${styles.color}-800`}>
                        {title}
                    </h3>
                    <div className={`mt-2 text-sm leading-5 text-${styles.color}-700`}>
                        {children}
                    </div>
                </div>
            </div>
        </div>
    );
}
