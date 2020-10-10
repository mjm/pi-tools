import { nodeResolve } from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import replace from "@rollup/plugin-replace";

export default {
    output: {
    },
    plugins: [
        replace({
            // TODO use bazel compilation mode
            'process.env.NODE_ENV': JSON.stringify('production'),
        }),
        nodeResolve(),
        commonjs(),
    ],
};
