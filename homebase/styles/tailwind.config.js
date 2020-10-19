const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
    future: {
        removeDeprecatedGapUtilities: true,
        purgeLayersByDefault: true,
    },
    purge: [],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter var', ...defaultTheme.fontFamily.sans],
            },
        },
    },
    variants: {
        display: ['responsive', 'group-hover'],
        visibility: ['responsive', 'group-hover'],
    },
    plugins: [
        require('@tailwindcss/ui'),
    ],
}
