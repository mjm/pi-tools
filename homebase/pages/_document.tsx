import {createRelayDocument, RelayDocument} from "relay-nextjs/document";
import NextDocument, {DocumentContext, Head, Html, Main, NextScript} from "next/document";

interface DocumentProps {
    relayDocument: RelayDocument;
}

class MyDocument extends NextDocument<DocumentProps> {
    static async getInitialProps(ctx: DocumentContext) {
        const relayDocument = createRelayDocument();

        const renderPage = ctx.renderPage;
        ctx.renderPage = () =>
            renderPage({
                enhanceApp: (App) => relayDocument.enhance(App),
            });

        const initialProps = await NextDocument.getInitialProps(ctx);

        return {
            ...initialProps,
            relayDocument,
        };
    }

    render() {
        const {relayDocument} = this.props;

        return (
            <Html>
                <Head>
                    <link rel="stylesheet" href="https://rsms.me/inter/inter.css"/>
                    <relayDocument.Script/>
                </Head>
                <body>
                    <Main/>
                    <NextScript/>
                </body>
            </Html>
        );
    }
}

export default MyDocument;
