import React from "react";
import {useHistory} from "react-router-dom";

export function TransitionLink({to, className, children}: {
    to: string;
    className?: string;
    children: React.ReactNode;
}) {
    const history = useHistory();
    // @ts-ignore
    const [startTransition, isPending] = React.unstable_useTransition();

    function onClick(e) {
        e.preventDefault();
        startTransition(() => {
            history.push(to);
        });
    }

    return (
        <a className={className} href={to} onClick={onClick}>
            {children}
        </a>
    );
}
