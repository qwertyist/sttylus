import Vue from "vue"
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist';
import VueCookie from 'vue-cookie';
import EventBus from '../eventbus';
Vue.use(VueCookie);
Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
});

const emptySession = {
  host: false,
  connected: false,
  id: '',
  password: '',
  zoom: {
    apiToken: '',
  },
}

export const store = new Vuex.Store({

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
        apiToken: ""
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
        background: '#000000'
      }
    }
  },
  mutations: {
    setSettings(state, newSettings) {
      console.log("setSettings")
      state.settings = newSettings;
    },
    setFontSettings(state, newFontSettings) {
      console.log("setFontSettings");
      state.settings.font = newFontSettings;
    },
    joinRemoteSession(state, session) {
      console.log("joinRemoteSession")
      state.session = session
    },

    setSessionID(state, id) {
      console.log("setsessionid")
      state.session.id = id
    },
    setName(state, identify) {
      state.name = identify.name
      state.remember = identify.remember
      if (state.remember) {
        VueCookie.set("name", identify.name, { expires: "1M" });
      } else {
        VueCookie.delete("name");
      }
    },
    remember(state) {
      console.log("Should remember:", state.remember)
    },
    loginSuccess(state) {
      state.loginFailed = false
    },
    loginFailed(state) {
      state.loginFailed = true
    },
    setSessionInfo(state, info) {
      console.log("set session info:", info)
      state.session.protected = info.protected
    },
    setSessionConnected(state, connected) {
      state.session.connected = connected;
    }
  }
});