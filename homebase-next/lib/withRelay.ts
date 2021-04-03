import {RelayOptions, RelayProps, withRelay as originalWithRelay} from "relay-nextjs";
import {ComponentType} from "react";
import {GraphQLTaggedNode} from "relay-runtime";
import {NextPageContext} from "next";
import {getClientEnvironment} from "./environment/client";

type Options = Omit<RelayOptions, "createServerEnvironment" | "createClientEnvironment">;

export default function withRelay<Props extends RelayProps, ServerSideProps>(Component: ComponentType<Props>, query: GraphQLTaggedNode, opts: Options = {}): {
    (props: Props): JSX.Element;
    getInitialProps: (ctx: NextPageContext) => Promise<{
        __wired__server__context: {};
        __wired_error_context?: undefined;
    } | {
        __wired_error_context: {};
        __wired__server__context?: undefined;
    } | {
        __wired__client__context: {};
    }>;
} {
    return originalWithRelay(Component, query, {
        async createServerEnvironment(_ctx, { cookie }: { cookie: string }) {
            const {createServerEnvironment} = await import("./environment/server");
            return createServerEnvironment(cookie);
        },
        createClientEnvironment() {
            return getClientEnvironment();
        },
        async serverSideProps({req}) {
            return { cookie: req.headers.cookie };
        },
        ...opts,
    });
}
