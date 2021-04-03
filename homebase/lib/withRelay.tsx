import {RelayOptions, RelayProps, withRelay as originalWithRelay} from "relay-nextjs";
import {ComponentType} from "react";
import {GraphQLTaggedNode} from "relay-runtime";
import {NextPageContext} from "next";
import {getClientEnvironment} from "./environment/client";
import Loading from "../components/Loading";

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
        async createServerEnvironment(_ctx, {cookie, user}: { cookie: string; user: string }) {
            const {createServerEnvironment} = await import("./environment/server");
            return createServerEnvironment(cookie, user);
        },
        createClientEnvironment() {
            return getClientEnvironment();
        },
        async serverSideProps({req}) {
            return {
                cookie: req.headers.cookie,
                user: req.headers["X-Auth-Request-User"],
            };
        },
        fallback: <Loading />,
        ...opts,
    });
}
