module.exports = {
    extends: [
        'eslint:recommended',
        'plugin:vue/essential',
        'plugin:vue/base',
        'plugin:vue/strongly-recommended',
        'plugin:vue/recommended',
        'prettier',
        'plugin:@typescript-eslint/recommended',
    ],
    parser: 'vue-eslint-parser',
    parserOptions: {
        ecmaVersion: 'latest',
    },
    rules: {
        'vue/no-unused-vars': 'error',
        'no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
        '@typescript-eslint/no-unused-vars': [
            'error',
            { argsIgnorePattern: '^_' },
        ],
        indent: 'off',
        semi: ['error', 'never'],
    },
    env: {
        node: true,
    },
    globals: {
        __APP_VERSION__: true,
    },
}
