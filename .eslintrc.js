module.exports = {
    extends: [
        'eslint:recommended',
        'plugin:vue/essential',
        'plugin:vue/base',
        'plugin:vue/strongly-recommended',
        'plugin:vue/recommended',
        'prettier'
    ],
    parser: 'vue-eslint-parser',
    rules: {
        'vue/no-unused-vars': 'error',
    },
    env: {
        node: true,
    },
};
