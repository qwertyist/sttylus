<template>
  <div>
    <b-row>
    <div style="height: 45vh; overflow-y: auto">
      <b-list-group v-for="(msg, i) in messages" :key="msg.id + '_' + i">
        <b-list-group-item style="white-space: pre-wrap">
          <small>{{ msg.timestamp }}</small> <b>{{ msg.name }}</b>: {{ msg.message }}
        </b-list-group-item>
      </b-list-group>
      <b-list-group-item style="background-color: #ddd;" ref="lastMessage"></b-list-group-item>
    </div>
    </b-row>
    <b-row>
      <b-form @submit.prevent="send" autocomplete="off">
        <div v-if="name == ''" class="d-flex align-items-center px-3 py-2">
          <b-form-input v-model="form.name" ref="input" autofocus placeholder="Välj ett namn för att chatta"></b-form-input>
          <b-button type="submit" size="sm">OK</b-button>
        </div>
        <div v-else class="d-flex bg-dark text-light align-items-center px-3 py-2">
          <b-form-input v-model="form.message" ref="input" autofocus placeholder="Skriv ett meddelande..."></b-form-input>
          <b-button type="submit" size="sm">Skicka</b-button>
        </div>
      </b-form>
    </b-row>
  </div>
</template>
<script>
import EventBus from '../eventbus.ts'
import { mapState }from "vuex";
export default {
  data() {
    return {
      id: "",
      modalOpen: false,
      unread: 0,
      show: false,
      form: {
        index: 0,
        name: "",
        message: "",
        to: null,
      },
      interpreters: [],
      users: [],
    }
  },
  computed: {
    name: function() {
      return this.$store.state.name
    },
    ...mapState(['messages']),
  },
  mounted() {
    this.addEventListeners();
    /*
    if(storedMessages
    this.messages = storedMessages;
      */
  },
  beforeDestroy() {
    this.$store.commit("storeMessages", this.messages)
    this.removeEventListeners();
  },
  methods: {
    addEventListeners() {
      EventBus.$on("modalOpen", val => { this.modalOpen = val  })
      EventBus.$on("toggleChat", this.toggleChat)
      EventBus.$on("recvClientId", (id) => { this.id = id })
    },
    removeEventListeners() {
      EventBus.$off("modalOpen")
      EventBus.$off("recvClientId")
      EventBus.$off("toggleChat");
    },
    onShow() {
      this.unread = 0
      this.form.message = ""
      this.$nextTick(() => {
        setTimeout(() => this.$refs.input.focus(), 250);
      })
    },
    onHide() {
      EventBus.$emit("refocus", false)
    },
    send() {
      if (this.name == "" && this.form.name != "") {
        this.$store.commit("setName", this.form.name)
        return
      }
      this.$nextTick(() => {
        EventBus.$emit("TXChat", { to: this.form.to, message: this.form.message, name: this.name})
        this.form.message = "";
      })
    },
  }
}
</script>
<style>
.navOpen {
  padding-top: 2em !important;
}
.bg-grey {
  background-color: var(--grey);
}
</style>
