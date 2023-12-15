<template>
  <b-sidebar ref="chat" title="Distanstolkning" v-model="show" right z-index=800 :header-class="{ navOpen: nav }" no-close-on-esc>

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
   <div :class="{ chatFocused: focused, textFocused: !focused }" @click="onFocus" style="height: 100%; overflow-y:hidden;">
      <div class="sidebar-field">
          <b-icon icon="chat-dots-fill" />
        <span class="float-right">
          Chatt &mdash; {{ time }}
        </span>
      </div>
     <div style="height: 93%; overflow-y: scroll;">
      <b-list-group v-for="(msg, i) in messages" :key="msg.id + '_' + i">
        <b-list-group-item style="white-space: pre-wrap" :class="{bgOther: !msg.me}">
            <small><i>[{{ msg.timestamp }}]</i> &ndash; {{ msg.name }}:</small> <b>{{ msg.message }}</b>
        </b-list-group-item>
      </b-list-group>
      <b-list-group-item style="background-color: #ddd;" ref="lastMessage"></b-list-group-item>
     </div>
    </div>
    <template #footer="{}">
      <b-form @click="onFocus" @submit.prevent="send" autocomplete="off">
        <div class="d-flex align-items-center px-3 py-2">
          <kbd>TAB</kbd>
          <b-form-select v-model="form.to" :options="targets"></b-form-select>
        </div>
        <div class="d-flex bg-dark text-light align-items-center px-3 py-2">
          <b-form-input @focus="focused = true" @blur="focused = false" @keydown.esc.prevent="hideNav" @keydown.tab.prevent="changeTarget" v-model="form.message" ref="input" autofocus placeholder="Skriv ett meddelande..."></b-form-input>
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
      focused: false,
      time: "16:14",
      timerID: null,
      updateInterval: null,
      form: {
        index: 0,
        message: "",
        to: "interpreters",
      },
      interpreters: [],
      users: [],
    }
  },
  computed: {
    name: function() {
      let name = this.$store.state.userData.name.split()[0]
      return name
    },
    targets: function () {
      let targets = [
        { value: null, text: "Alla anslutna" },
        { value: "interpreters", text: "Alla tolkar" },
        { value: "users", text: "Alla tolkanvändare" }
      ]
      return targets
    }
  },
  watch: {
    focused: (newVal, oldVal) => {
      if (oldVal == newVal) return;
      if (newVal == true) {
        EventBus.$emit("chatFocused")
        console.log("Chat.Vue - chat focused")
      }
    }
  },
  mounted() {
    this.timerID = setInterval(this.updateTime, 1000);
    this.updateTime();
    this.updateInterval = setInterval(
        () => this.update()
    ,100);

    this.addEventListeners();
  },
  beforeDestroy() {
    clearInterval(this.updateInterval)
    clearInterval(this.timerID)
    this.removeEventListeners();
  },
  methods: {
    onFocus() {
      console.log("chat focused")
      this.focus()
      this.$store.commit("setFocus","chat")
    },
    onBlur() {
      console.log("chat blurred")
      this.$store.commit("setFocus","text")
    },
    update() {
      if(this.focused) {
        EventBus.$emit("scrollDown", "")
      }
    },
    addEventListeners() {
      this.$refs.chat.$on("shown", this.onShow)
      this.$refs.chat.$on("hidden", this.onHide)
      EventBus.$on("RXChat", this.recv)
      EventBus.$on("toggleChat", this.toggleChat)
      EventBus.$on("showChat", this.showChat)
      EventBus.$on("hideChat", this.hideChat)
      EventBus.$on("focusChat", this.focus)
      EventBus.$on("clearChat", this.clearMessages)
      EventBus.$on("myClientId", (id) => { console.log("myClientId:", id); this.id = id })
    },
    removeEventListeners() {
      this.$refs.chat.$off("shown");
      this.$refs.chat.$off("hidden");
      EventBus.$off("RXChat");
      EventBus.$off("recvClientId")
      EventBus.$off("toggleChat");
      EventBus.$off("showChat")
      EventBus.$off("hideChat")
      EventBus.$off("focusChat");
      EventBus.$off("clearChat")
    },
    hideNav() {
      console.log("hidenav")
      EventBus.$emit("closeNav")
    },
    toggleChat() {
      this.show = !this.show;
    },
    showChat() {
      this.show = true
    },
    hideChat() {
      this.show = false
    },
    onShow() {
      this.$store.commit('setModalOpen', true)
      EventBus.$emit('modalOpened')
      EventBus.$emit('chatOpened')
      this.unread = 0
      this.form.message = ""
      this.$refs.input.focus()
    },
    onHide() {
      this.$store.commit('setModalOpen', false)
      EventBus.$emit('modalClosed')
      EventBus.$emit('chatClosed')
      EventBus.$emit("refocus", false)
    },
    focus() {
      this.show = true
      this.$nextTick( () => {
        this.$refs.lastMessage.scrollIntoView();
      });
      this.$nextTick(() => {
        setTimeout(() => this.$refs.input.focus(), 150);
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
      console.log(this.id)
      let me = msg.id == this.id
      let now = new Date()
       let timestamp = now.toLocaleTimeString().slice(0,5);
      if (msg.chat.message == undefined) {
        msg.chat.message = "anslöt"
        return
      }
      this.messages.push(
        {
          timestamp: timestamp,
          id: msg.id,
          name: msg.chat.name,
          message: msg.chat.message,
          me: me
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
    updateTime() {
      var cd = new Date();
      this.time = this.zeroPadding(cd.getHours(), 2) + ':' + this.zeroPadding(cd.getMinutes(), 2) + ':' + this.zeroPadding(cd.getSeconds(), 2);
    },
    zeroPadding(num, digit) {
      var zero = '';
      for(var i = 0; i < digit; i++) {
          zero += '0';
      }
      return (zero + num).slice(-digit);
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

.chatFocused {
  border-left: 12px solid green;
}

.textFocused {
  border-left: 12px solid red;
}

.navOpen {
  padding-top: 2em !important;
}
.bg-grey {
  background-color: var(--grey);
}
.bgOther {
  background-color: #d5e6e6
}
</style>
