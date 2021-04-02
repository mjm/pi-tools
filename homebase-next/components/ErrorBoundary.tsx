import React from "react";
import Alert from "./Alert";

interface State {
    error?: Error;
}

export default class ErrorBoundary extends React.Component<{
    children: React.ReactNode;
    fallback?: (error: Error) => React.ReactNode;
}, State> {
    state: State = {};

    static getDerivedStateFromError(error: Error) {
        return {error};
    }

    render() {
        if (this.state.error) {
            if (this.props.fallback) {
                return this.props.fallback(this.state.error);
            }

            return (
                <Alert title="An error occurred" severity="error">
                    {this.state.error.toString()}
                </Alert>
            );
        }

        return this.props.children;
    }
}
