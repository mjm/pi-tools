import React from "react";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";

interface State {
    error?: Error;
}

export class ErrorBoundary extends React.Component<{
    children: React.ReactNode;
}, State> {
    state: State = {};

    static getDerivedStateFromError(error: Error) {
        return {error};
    }

    render() {
        if (this.state.error) {
            return (
                <main className="mb-8">
                    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                        <Alert title="An error occurred" severity="error">
                            {this.state.error.toString()}
                        </Alert>
                    </div>
                </main>
            );
        }

        return this.props.children;
    }
}
