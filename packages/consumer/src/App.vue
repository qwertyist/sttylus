<template>
  <b-overlay :show="qr.show" @click="toggleQRCode()" rounded="sm">
    <template #overlay>
      <div class="text-center">
        <qrcode-vue :value="qr.url" :size="350" level="H" />
      </div>
    </template>
    <div id="app" @mousedown="handleMouseDown" @dbclick="handleDblClick">
      <Consumer />
      <ConsumerSettings :visisble="modalOpen" />
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
              show: false,
              url: "",
            },
            clicks: 0,
            timer: null,
            unread: 0,
            screenLock: {},
            modalOpen: false,
        }
    },
    mounted() {

      window.addEventListener('beforeunload', this.leaving);
      console.log("mobile:", this.$mobile);
      EventBus.$on("RXChat", this.recv)
      EventBus.$on("toggleQRCode", this.toggleQRCode)
      EventBus.$on("setQRCodeURL", this.setQRCode)
      EventBus.$on("modalOpen", val => this.modalOpen = val)
      this.getScreenLock();
    },
    beforeDestroy() {
      EventBus.$off("modalOpen")
      EventBus.$off("modalClosed")
      EventBus.$off("RXChat")
      EventBus.$off("toggleQRCode")
      EventBus.$off("setQRCodeURL")
      if(typeof this.screenLock !== "undefined" && this.screenLock != null) {
          this.screenLock.release()
          .then(() => {
            this.screenLock = null;
          });
      }
    },
    methods: {
        leaving() {
          this.$store.commit("clearMessages")
        },
        recv(msg) {
          if (msg.chat.message == null) { return }
          let now = new Date()
          let timestamp = now.toLocaleTimeString().slice(0,5);
          if(msg.id != this.$store.getters.socketId) {
            this.unread++;
          }
          this.$store.commit("addMessage",
            {
              timestamp: timestamp,
              id: msg.id,
              name: msg.chat.name,
              message: msg.chat.message
            }
          )
          if(this.modalOpen) {
            this.unread = 0
            EventBus.$emit("recv")
          } else {
            if (this.unread > 0) {
              this.$toast.info(
                "Oläst chatmeddelande (" + this.unread +")",
                { duration: 750}
              )
            }
          }
        },
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
              console.error(err.name, err.message)});
        },
        handleDblClick() {
        },
        handleMouseDown() {
            if (this.$mobile) {
                this.clicks++;
                if (this.clicks === 1) {
                    this.timer = setTimeout(() => {
                        this.clicks = 0;
                    }, 700)
                } else {
                    clearTimeout(this.timer);
                    this.clicks = 0;
                    this.unread = 0
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
