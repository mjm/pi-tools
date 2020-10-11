import { nodeResolve } from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import replace from "@rollup/plugin-replace";

console.log(process.env)
const env = process.env.COMPILATION_MODE === 'opt' ? 'production' : 'development';

export default {
    output: {
    },
    plugins: [
        replace({
            // TODO use bazel compilation mode
            'process.env.NODE_ENV': JSON.stringify(env),
        }),
        nodeResolve(),
        commonjs(),
    ],
};
