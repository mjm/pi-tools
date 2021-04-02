module.exports = {
    future: {
        webpack5: true
    },

    async rewrites() {
        if (process.env.NODE_ENV !== 'development') {
            return {};
        }

        return {
            beforeFiles: [
                {
                    source: '/graphql',
                    destination: 'http://localhost:6460/graphql',
                }
            ]
        }
    }
}
