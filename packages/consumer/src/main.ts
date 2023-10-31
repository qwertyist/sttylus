import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'

import './assets/custom.scss'
import('typeface-roboto-mono')
import 'typeface-roboto-mono/index.css'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'

import bootstrap from './bootstrap.js'
Vue.use(bootstrap)

import VueToast from 'vue-toast-notification'
import './assets/toast.css'

import Consumer from './components/Consumer.vue'

Vue.use(VueRouter)
const router = new VueRouter({
  mode: 'hash',
  routes: [
    { path: '/:id?', component: Consumer },
    { path: '/:id/:password?', component: Consumer },
  ],
})

Vue.use(VueToast, { duration: 3000, pauseOnHover: false })
Vue.prototype.$collab = import.meta.env.VITE_STTYLUS_COLLAB_SERVER
Vue.prototype.$mobile =
  /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent
  )
    ? true
    : false

Vue.use(Vuex)
const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
})

const store = new Vuex.Store({
  plugins: [vuexLocal.plugin],
  state: {
    mode: 'test',
    connection: 'connecting',
    machineId: '',
    licenseKey: '',
    userData: {},
    name: '',
    messages: [],
    session: {
      host: false,
      connected: false,
      id: '',
      public: false,
      local: false,
      zoom: {
        apiToken: '',
      },
    },
    localIP: '',
    sessionId: '',
    mySessionId: '',
    settings: {
      font: {
        size: 32,
        family: 'Times New Roman',
        fontColorID: 2,
        foreground: '#ffff00',
        background: '#000000',
      },
    },
  },
  getters: {
    storedMessages: (state) => {
      if (!state.messages) {
        state.messages = []
      }
      return state.messages
    },
  },
  mutations: {
    setSettings(state, newSettings) {
      console.log('setSettings')
      state.settings = newSettings
    },
    setFontSettings(state, newFontSettings) {
      console.log('setFontSettings')
      state.settings.font = newFontSettings
    },
    joinRemoteSession(state, session) {
      console.log('joinRemoteSession')
      state.session = session
    },

    setSessionID(state, id) {
      state.session.id = id
    },
    setName(state, name) {
      state.name = name
    },
    addMessage(state, message) {
      state.messages.push(message)
    },
    storeMessages(state, messages) {
      if (!messages) {
        state.messages = []
      } else {
        state.messages = messages
      }
    },
    clearMessages(state) {
      state.messages = []
    },
    remember(state) {
    },
    loginSuccess(state) {
      state.loginFailed = false
    },
    loginFailed(state) {
      state.loginFailed = true
    },
    setSessionInfo(state, info) {
      state.session.protected = info.protected
    },
    setSessionConnected(state, connected) {
      state.session.connected = connected
    },
  },
})

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app')
