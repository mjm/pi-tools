module.exports = {
    src: './Sources',
    schema: '../../homebase/schema.graphql',
    language: 'swift',
    artifactDirectory: './__generated__',
    customScalars: {
        Time: 'String'
    }
};
