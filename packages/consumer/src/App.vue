<template>
  <b-overlay :show="qr.show" @click="toggleQRCode()" rounded="sm">
    <template #overlay>
      <div class="text-center">
        <qrcode-vue :value="qr.url" :size="350" level="H" />
      </div>
    </template>
    <div id="app" @mousedown="handleMouseDown" @dbclick="handleDblClick">
      <Consumer />
      <ConsumerSettings />
    </div>
  </b-overlay>
</template>
<script>
//import Login from "./components/modals/Login.vue"

import ConsumerSettings from "./components/modals/ConsumerSettings.vue"
import EventBus from "./eventbus"

import QrcodeVue from 'qrcode.vue'
export default {
    components: {
        ConsumerSettings,
        QrcodeVue
    },
    data() {
        return {
            qr: {
              show: true,
              url: "",
            },
            clicks: 0,
            timer: null,
            screenLock: {}
        }
    },
    mounted() {

      EventBus.$on("toggleQRCode", this.toggleQRCode)
      EventBus.$on("setQRCodeURL", this.setQRCode)
      this.getScreenLock();
    },
    beforeDestroy() {
      EventBus.$off("toggleQRCode")
      EventBus.$off("setQRCodeURL")
      if(typeof this.screenLock !== "undefined" && this.screenLock != null) {
          this.screenLock.release()
          .then(() => {
            console.log("Lock released 🎈");
            this.screenLock = null;
          });
      }
    },
    methods: {
        setQRCode(val) {
          this.qr.url = val
        },
        toggleQRCode() {
          this.qr.show = !this.qr.show; 
        },
        getScreenLock() {
            navigator.wakeLock.request("screen")
            .then(lock => {
              this.screenLock = lock
            })
            .catch(err => {
              console.log(err.name, err.message)});
        },
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
