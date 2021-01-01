import {nodeResolve} from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import replace from "@rollup/plugin-replace";
import styles from "rollup-plugin-styles";
import html, {makeHtmlAttributes} from "@rollup/plugin-html";
import babel from "@rollup/plugin-babel";

const env = process.env.COMPILATION_MODE === 'opt' ? 'production' : 'development';

export default {
    output: {
        entryFileNames: 'assets/[name]-[hash].js',
    },
    plugins: [
        replace({
            'process.env.NODE_ENV': JSON.stringify(env),
        }),
        nodeResolve(),
        commonjs(),
        babel({
            babelHelpers: 'bundled',
            plugins: [
                ['relay', {
                    eagerESModules: true,
                    artifactDirectory: 'node_modules/com_github_mjm_pi_tools/homebase/api/__generated__',
                }]
            ],
        }),
        styles({
            mode: 'extract',
        }),
        html({
            publicPath: '/',
            template: htmlTemplate,
        }),
    ],
};

async function htmlTemplate({attributes, files, meta, publicPath, title}) {
    const scripts = (files.js || [])
        .map(({fileName}) => {
            const attrs = makeHtmlAttributes(attributes.script);
            return `<script src="${publicPath}${fileName}"${attrs}></script>`;
        })
        .join('\n');

    const links = (files.css || [])
        .map(({fileName}) => {
            const attrs = makeHtmlAttributes(attributes.link);
            return `<link href="${publicPath}${fileName}" rel="stylesheet"${attrs}>`;
        })
        .join('\n');

    const metas = meta
        .map((input) => {
            const attrs = makeHtmlAttributes(input);
            return `<meta${attrs}>`;
        })
        .join('\n');

    return `
<!doctype html>
<html${makeHtmlAttributes(attributes.html)}>
  <head>
    ${metas}
    <title>${title}</title>
    ${links}
  </head>
  <body class="bg-gray-100">
    <div id="root"></div>
    ${scripts}
  </body>
</html>`;
}
