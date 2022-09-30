<template>
    <div id="app" @mousedown="handleMouseDown" @dbclick="handleDblClick">
     <Consumer />
     <ConsumerSettings />
    </div>
</template>
<script>
import Login from "./components/modals/Login.vue"
import Consumer from "./components/Consumer.vue";
import ConsumerSettings from "./components/modals/ConsumerSettings.vue"
export default {
    components: {
    ConsumerSettings
},
    data() {
        return {
            clisk: 0,
            timer: null,
        }
    },
    mounted() {
        console.log("MOBILE:", this.$mobile)
    },
    methods: {
        handleDblClick() {
            console.log("Open settings");
        },
        handleMouseDown() {
            if (this.$mobile) {
                this.clicks++;
                if (this.clicks === 1) {
                    this.timer = setTimeout(() => {
                        this.clicks = 0;
                    }, 700)
                } else {
                    console.log("open settings")
                    clearTimeout(this.timer);
                    this.clicks = 0;
                    this.$bvModal.show("consumerSettings")
                }
            }
        }
    }
}
</script>
<style scoped>
.html,
body {
    height: 100%;
}
#app {
    width: 100%;
    overflow: hidden;
}
</style>
