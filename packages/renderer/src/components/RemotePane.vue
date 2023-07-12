<template>
    <b-sidebar
        style="margin-left: 1em; position: fixed"
        width="20vw"
        no-slide
        v-model="visible"
        right
        z-index="900"
        title="Distanstolkning"
        no-close-on-esc
        @hidden="hide"
    >
        <div class="float-right">Öppna/stäng med <kbd>F3</kbd></div>
        <br />
        Anslutna användare:
        <b-list-group v-for="[id, client] in clients" :key="id">
            <b-list-group-item>
                {{ client.name
                }}<span class="float-right" v-if="client.type == 'user'"
                    ><b-badge>Tolkanvändare</b-badge></span
                ><span class="float-right" v-else><b-badge>Tolk</b-badge></span>
            </b-list-group-item>
        </b-list-group>
        <br />
        <!--<b-checkbox :checked="keepVisible"
      >Håll öppen <kbd>CTRL+F3</kbd></b-checkbox
    >
    <div style="height: 75%">
      <b-list-group v-for="msg in messages" :key="msg.id">
        <b-list-group-item
          style="white-space: pre-wrap"
          :class="{ user: from(msg.from), interpreter: !from(msg.from) }"
        >
          <b>{{ msg.name }}</b
          >: {{ msg.msg }}
        </b-list-group-item>
      </b-list-group>
    </div>
    <b-badge class="interpreter">Meddelanden till/från kollega</b-badge
    ><b-badge class="user">Meddelande till/från tolkanvändare</b-badge>
    <b-form @submit.prevent="onSend">
      <b-form-group label="Skicka meddelande">
        <b-form-input
          @focus="focusInput"
          @blur="blurInput"
          ref="chatInput"
          v-model="msgInput"
        />
      </b-form-group>
      <b-form-group label="Till">
        <b-form-select v-model="to" :options="recipients" />
      </b-form-group>
    </b-form>
    -->
    </b-sidebar>
</template>
<script>
import EventBus from '../eventbus.js'
export default {
    data() {
        return {
            visible: false,
            msgInput: '',
            id: 3,
            to: 0,
            focus: false,
            keepVisible: false,
            recipients: [
                { value: 0, text: 'Tolkanvändare' },
                { value: 1, text: 'Kollega' },
            ],
            messages: [
                {
                    msg: 'hello this is your interpreter typing. How is everything working for you?',
                    id: 0,
                    name: 'David',
                    from: 0,
                },
                { msg: "What's up?", id: 1, name: 'Tolkanvändare', from: 0 },
                {
                    msg: 'Lägger du in blaa=bland annat?',
                    id: 2,
                    name: 'Oskar',
                    from: 1,
                },
            ],
        }
    },
    computed: {
        clients() {
            return this.$store.state.clients
        },
    },
    methods: {
        hide() {
            EventBus.$emit('chatHidden')
        },
        from(from) {
            if (from == 0) {
                return true
            } else {
                return false
            }
        },
        toggleChat() {
            this.visible = !this.visible
            /*
      this.$nextTick(() => {
        if (this.visible && !this.keepVisible) {
          console.log("keepvisible is false, focus input");
          this.addEventListeners();
          this.$refs.chatInput.focus();
        } else if (!this.visible) {
          console.log("Do not keep visible so closing");
          this.removeEventListeners();
          this.focus = false;
          EventBus.$emit("refocus", "");
        } else if (this.visible && this.keepVisible) {
          EventBus.$emit("refocus", "");

          console.log("outside visible logic");
        }
      });
      */
        },
        focusInput() {
            this.focus = true
            console.log('in focus')
        },
        blurInput() {
            this.focus = false
            console.log('blur focus')
        },

        onSend() {
            console.log(this.to)
            EventBus.$emit('sendMsg', {
                from: this.to,
                name: this.$store.state.userData.name,
                msg: this.msgInput,
                id: this.id,
            })
            this.id++
            this.msgInput = ''
        },
        handleMessage(msg) {
            console.log('msg received:', msg)
            this.messages.push(msg)
        },
        eventListeners(e) {
            console.log(e)
            if (e.key == 'F3' && this.visible) {
                this.visible = false
            }
            if (e.key == 'F3' && !this.visible) {
                this.removeEventListeners()
                return
            }
            if (e.key == 'F3' && this.keepVisible && !e.ctrlKey) {
                this.$nextTick(() => {
                    if (this.focus == undefined || !this.focus) {
                        console.log('input not in focus')
                        this.$refs.chatInput.focus()
                    } else {
                        console.log('input is in focus')
                        this.removeEventListeners()

                        EventBus.$emit('refocus', '')
                    }
                })
                return
            }
            if (e.key == 'Tab') {
                e.preventDefault()
                if (this.to == 0) {
                    this.to++
                } else {
                    this.to = 0
                }
                return
            }
            /* if (e.key == "Escape") {
          e.preventDefault();
          this.removeEventListeners();

          EventBus.$emit("refocus", "");
        }*/
        },
        updateClients() {
            this.$forceUpdate()
        },
        addEventListeners() {
            console.log('add event listeners')
            window.addEventListener('keydown', this.eventListeners)
        },
        removeEventListeners() {
            console.log('remove event listeners')
            window.removeEventListener('keydown', this.eventListeners)
        },
    },
    mounted() {
        EventBus.$on('toggleChat', this.toggleChat)
        EventBus.$on('msg', this.handleMessage)
        EventBus.$on('client_connected', this.updateClients)
        EventBus.$on('client_disconnected', this.updateClients)
        EventBus.$on('chatHidden', () => {
            this.visible = false
        })
    },
    beforeDestroy() {
        this.removeEventListeners()
        EventBus.$off('client_connected', this.updateClients)
        EventBus.$off('client_disconnected', this.updateClients)
        EventBus.$off('msg')
        EventBus.$off('toggleChat')
        EventBus.$off('chatHidden')
    },
}
</script>
<style scoped>
::v-deep .b-sidebar-body {
    padding: 1rem;
}
.user {
    background-color: #77c;
}
.interpreter {
    background-color: #7c7;
}
</style>
