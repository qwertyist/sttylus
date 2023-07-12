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

Vue.use(VueToast, { duration: 5000 })
Vue.prototype.$collab = import.meta.env.VITE_STTYLUS_COLLAB_SERVER
console.log(Vue.prototype.$collab)
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
      console.log('setsessionid')
      state.session.id = id
    },
    setName(state, identify) {
      state.name = identify.name
      state.remember = identify.remember
      if (state.remember) {
        VueCookie.set('name', identify.name, { expires: '1M' })
      } else {
        VueCookie.delete('name')
      }
    },
    remember(state) {
      console.log('Should remember:', state.remember)
    },
    loginSuccess(state) {
      state.loginFailed = false
    },
    loginFailed(state) {
      state.loginFailed = true
    },
    setSessionInfo(state, info) {
      console.log('set session info:', info)
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
