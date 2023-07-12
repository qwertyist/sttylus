import { defineConfig } from 'vite'
import { createVuePlugin } from 'vite-plugin-vue2'
import ViteComponents from 'vite-plugin-components'
import pkg from '../../package.json'
import visualizer from 'rollup-plugin-visualizer'

// https://vitejs.dev/config/
export default defineConfig({
    mode: process.env.NODE_ENV,
    root: __dirname,
    plugins: [
        visualizer(),
        createVuePlugin(),
        ViteComponents({ transformer: 'vue2' }),
    ],
    resolve: {
        alias: {
            vue: 'vue/dist/vue.common',
        },
    },
    base: './',
    build: {
        /*commonjsOptions: {
			include: [/node_modules/]
	},*/
        sourcemap: true,
        outDir: '../../dist/renderer',
    },
    server: {
        port: pkg.env.PORT,
    },
    define: {
        __APP_VERSION__: JSON.stringify(process.env.npm_package_version),
    },
})
