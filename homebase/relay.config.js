module.exports = {
    src: '.',
    schema: 'schema.graphql',
    exclude: ['**/node_modules/**', '**/__mocks__/**', '**/__generated__/**'],
    extensions: ['ts', 'tsx'],
    language: 'typescript',
    artifactDirectory: '__generated__',
    customScalars: {
        Cursor: 'String',
        Time: 'String',
    },
};
