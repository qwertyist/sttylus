<template>
    <div>CONSUME</div>
</template>

<script>
import TextView from './components/TextView.vue'
export default {
    name: 'Consumer',
    components: {
        TextView,
    },
    data() {
        return {}
    },
    computed: {},
    methods: {},
}
</script>

<style>
body,
html {
    overflow: hidden;
}
#app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
}
#navigation {
    position: relative;
    width: 100%;
    z-index: 1000;
}
.navbar {
    height: 5vh !important;
}
</style>
