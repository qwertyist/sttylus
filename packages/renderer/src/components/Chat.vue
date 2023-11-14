<template>
  <b-sidebar ref="chat" title="Distanstolkning" v-model="show" right z-index=800 :header-class="{ navOpen: nav }">

   <!-- <div style="height:15%">
      <div class="sidebar-field">
        <b-icon icon="person-fill" />
        <span class="float-right">
          Anslutna
        </span>
      </div>
      <div v-for="(interpreter, idx) in interpreters">
        <b-badge>{{ interpreter }}</b-badge>
      </div>
      <div v-for="(user, idx) in users">
        <b-badge></b-badge>
      </div>
    </div>
   -->
    <div style="height: 35%">
      <div class="sidebar-field">
          <b-icon icon="chat-dots-fill" />
        <span class="float-right">
          Chatt
        </span>
      </div>
      <b-list-group v-for="(msg, i) in messages" :key="msg.id + '_' + i">
        <b-list-group-item style="white-space: pre-wrap">
          <small>{{ msg.timestamp }}</small> <b>{{ msg.name }}</b>: {{ msg.message }}
        </b-list-group-item>
      </b-list-group>
      <b-list-group-item style="background-color: #ddd;" ref="lastMessage"></b-list-group-item>
    </div>
    <template #footer="{}">
      <b-form @submit.prevent="send" autocomplete="off">
        <div class="d-flex align-items-center px-3 py-2">
          <kbd>TAB</kbd>
          <b-form-select v-model="form.to" :options="targets"></b-form-select>
        </div>
        <div class="d-flex bg-dark text-light align-items-center px-3 py-2">

          <b-form-input @keydown.tab.prevent="changeTarget" v-model="form.message" ref="input" autofocus placeholder="Skriv ett meddelande..."></b-form-input>
          <b-button type="submit" size="sm">Skicka</b-button>
        </div>
      </b-form>
    </template>
  </b-sidebar>
</template>
<script>
import EventBus from '../eventbus.js'
export default {
  props: {
    nav: Boolean
  },
  data() {
    return {
      id: "",
      messages: [],
      unread: 0,
      show: false,
      focus: false,
      form: {
        index: 0,
        message: "",
        to: null,
      },
      interpreters: [],
      users: [],
    }
  },
  computed: {
    name: function() {
      let name = this.$store.state.userData.name.split()[0]
      return "Tolk " + name

    },
    targets: function () {
      let targets = [
        { value: null, text: "Alla anslutna" },
        { value: "interpreters", text: "Alla tolkar" },
        { value: "users", text: "Alla tolkanvändare" }
      ]
      /*
      this.interpreters.forEach(interpreter => {
        if (!interpreter.name) {
          interpreter.name = interpreter.id
        }
        targets.push({ value: interpreter.id, text: interpreter.name })
      })
      this.users.forEach(user => {
        if (!user.name) {
          user.name = user.id
        }
        targets.push({ value: user.id, text: user.name })
      })
      */
      return targets
    }
  },
  mounted() {
    this.addEventListeners();
  },
  beforeDestroy() {
    this.removeEventListeners();
  },
  methods: {
    focusInput(val) {
      console.log("focus input:", val)
      this.focus = val;
    },
    addEventListeners() {
      this.$refs.chat.$on("shown", this.onShow)
      this.$refs.chat.$on("hidden", this.onHide)
      EventBus.$on("abbModalClosed", this.onShow)
      EventBus.$on("RXChat", this.recv)
      EventBus.$on("toggleChat", this.toggleChat)
      EventBus.$on("clientListUpdated", this.updateClients)
      EventBus.$on("clearChat", this.clearMessages)
      EventBus.$on("recvClientId", (id) => { this.id = id })
    },
    removeEventListeners() {
      this.$refs.chat.$off("shown");
      this.$refs.chat.$off("hidden");
      EventBus.$off("RXChat");
      EventBus.$off("abbModalClosed")
      EventBus.$off("recvClientId")
      EventBus.$off("clientListUpdated")
      EventBus.$off("toggleChat");
      EventBus.$off("clearChat")
    },
    toggleChat() {
      this.show = !this.show;
    },
    onShow() {
      this.$store.commit('setModalOpen', true)
      EventBus.$emit('modalOpened')
      this.unread = 0
      this.form.message = ""
      this.$nextTick( () => {
        this.$refs.lastMessage.scrollIntoView();
      });
      this.$nextTick(() => {
        setTimeout(() => this.$refs.input.focus(), 250);
      })
    },
    onHide() {
      this.$store.commit('setModalOpen', false)
      EventBus.$emit('modalClosed')
      EventBus.$emit("refocus", false)
    },
    updateClients() {
      this.users = []
      this.interpreters = []
      this.$store.state.clients.forEach(c => {
        if (c.interpreter) {
          this.interpreters.push(c)
        }
        if (!c.interpreter) {
          this.users.push(c)
        }

      })
    },
    changeTarget() {
      if (this.form.index < this.targets.length - 1) {
        this.form.index++
        this.form.to = this.targets[this.form.index].value
      } else {
        this.form.index = 0
        this.form.to = this.targets[this.form.index].value
      }
    },
    send() {
      this.$nextTick(() => {
        EventBus.$emit("TXChat", { to: this.form.to, message: this.form.message, name: this.name})
        this.form.message = "";
      })
    },
    recv(msg) {
      console.log("recv:", msg)
      let now = new Date()
       let timestamp = now.toLocaleTimeString().slice(0,5);
      if (msg.chat.message == undefined) { msg.chat.message = "anslöt" }
      this.messages.push(
        {
          timestamp: timestamp,
          id: msg.id,
          name: msg.chat.name,
          message: msg.chat.message
        }
      )
      if(this.show) {
        this.$nextTick( () => {
        this.$refs.lastMessage.scrollIntoView();
        });
      } else {
        this.unread++;
        this.$toast.info("Oläst chatmeddelande (" + this.unread +")")
      }
    },
    clearMessages() {
      this.messages = []
    },
  }
}
</script>
<style>

.sidebar-field {
  font-size: 1.1rem;
  font-weight: bold;
  padding: 0.5rem 1rem;
}

.navOpen {
  padding-top: 2em !important;
}
.bg-grey {
  background-color: var(--grey);
}
</style>
