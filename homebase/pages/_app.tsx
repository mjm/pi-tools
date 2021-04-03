import {RelayEnvironmentProvider} from "react-relay/hooks";
import {getInitialPreloadedQuery, getRelayProps} from "relay-nextjs/app";
import {AppProps} from "next/app";
import Head from "next/head";
import NavigationBar from "../components/NavigationBar";

import "../styles/globals.css";
import {getClientEnvironment} from "../lib/environment/client";

const initialPreloadedQuery = getInitialPreloadedQuery({
    createClientEnvironment: () => getClientEnvironment(),
});

function MyApp({Component, pageProps}: AppProps) {
    const relayProps = getRelayProps(pageProps, initialPreloadedQuery);
    const env = relayProps.preloadedQuery?.environment ?? getClientEnvironment();

    return (
        <>
            <Head>
                <title>Homebase</title>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
            </Head>
            <div>
                <NavigationBar/>
                <RelayEnvironmentProvider environment={env}>
                    <Component {...pageProps} {...relayProps} />
                </RelayEnvironmentProvider>
            </div>
        </>
    );
}

export default MyApp;
