// maybe this can eventually be replaced by calling the relay-compiler in a Bazel rules
module.exports = {
    schema: 'schema.graphql',
    src: '.',
    exclude: ["**/node_modules/**", "api/**"],
    language: 'typescript',
    artifactDirectory: 'api/__generated__',
    customScalars: {
        Cursor: 'String',
        Time: 'String',
    },
};
