import { user } from './models/user.js'
import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'
import VueCookie from 'vue-cookie'
import { router } from '../router.js'
import api from '../api/api.js'
import EventBus from '../eventbus.js'
Vue.use(VueCookie)
Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
})

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
  key: 'offline',
  state: {
    users: null,
    mode: 'test',
    connection: 'connecting',
    machineId: '',
    licenseKey: '',
    user: user,
    userData: {
      id: '',
    },
    lastLogin: {},
    lastSync: null,
    doc: '',
    focus: 'text',
    lastFocus: 'text',
    modalOpen: false,
    parsedDoc: '',
    stagedNominations: [],
    missedAbbs: new Map(),
    selectedWord: '',
    targetList: '',
    conflictingAbbs: [],
    sharedList: {
      id: '',
      base: '',
    },
    sharedAbbs: [],
    local: {
      connected: false,
    },
    session: {
      host: false,
      connected: false,
      id: '',
      password: '',
      zoom: {
        apiToken: '',
      },
    },
    clients: new Map(),
    settings: {
      font: {
        size: 32,
        family: 'Times New Roman',
        colorID: 2,
        lineHeight: 1.25,
        foreground: '#ffff00',
        background: '#000000',
        margins: {
          top: 10,
          right: 10,
          bottom: 10,
          left: 10,
        },
        customColors: {
          valid: false,
          foreground: '#ffffff',
          background: '#000000',
        },
      },
      behaviour: {
        capitalizeOnNewLine: true,
      },
      selectedLists: {
        standard: '',
        addon: [],
      },
      selectedManuscripts: [],
    },
    //Det här minns jag inte vad det handlade om, men misstänker 
    //att det har att göra med att
    //skriva i samma text samtidigt?
    lookup: "",
    cached: false,
    multiplayer: [],
  },

  mutations: {
    AUTH_REQUEST: (state) => {
      state.status = 'loading'
    },
    AUTH_SUCCESS: (state, token) => {
      state.status = 'success'
      state.token = token
      state.lastLogin = new Date()
    },
    AUTH_ERROR: (state) => {
      state.status = 'error'
    },
    ADD_USER: (state, data) => {
      console.log(state.users)
      state.users.set(data.id, data)
    },
    initState(state) {
      state.token = localStorage.getItem('user-token')
      state.missedAbbs = new Map()
      state.sharedList = { id: '', base: '' }
      state.users = new Map()
      if (!state.lastSync) {
        console.log('No last sync')
        state.lastSync = new Date(0)
        console.log(state.lastSync)
      }
      if (import.meta.env.VITE_STTYLUS_MODE == 'desktop') {
        state.local.connected = false
        api
          .getUsers()
          .then((resp) => {
            if (resp.data != null) {
              console.log(resp.data.length, 'local users')
            }
          })
          .catch((err) => {
            console.error("Store couldn't get local users:", err)
          })
      }

      state.session = {
        host: false,
        connected: false,
        id: '',
        zoom: {
          apiToken: '',
        },
      }

      if (state.token) {
        console.log('prior login stored')
        if (api.checkForToken() == undefined) {
          api.setIdToken(state.token)
        }
        api
          .getSettings()
          .then((resp) => {
            console.log('received settings:', resp.data)
            EventBus.$emit('stateReady', '')
          })
          .catch((err) => {
            console.error('init state couldnt get settings:', err)
          })
        let session = JSON.parse(VueCookie.get('session'))
        if (session) {
          state.session = session
        }
      } else {
        console.log('no prior login stored')
        state.userData = user
      }
    },
    setMode(state, mode) {
      state.mode = mode
    },
    setConnectionStatus(state, status) {
      state.connection = status
    },
    setMachineId(state, id) {
      state.machineId = id
    },
    setLicenseKey(state, key) {
      state.licenseKey = key
    },
    setLastSync(state, lastSync) {
      state.lastSync = lastSync
    },
    setUserData(state, data) {
      state.userData = {}
      state.userData.id = data.id
      state.userData.email = data.email
      state.userData.role = data.role
      state.userData.name = data.name
      state.userData.subscriptions = [].push(data.subscriptions.global_lists)
      state.licenseKey = data.license_key
    },
    setModalOpen(state, open) {
      if (open == true) {
        state.lastFocus = state.focus
        state.focus = 'modal'
      } else {
        state.focus = state.lastFocus
      }

      state.modalOpen = open
    },
    setFocus(state, focus) {
      state.focus = focus
    },
    setSelectedStandard(state, list) {
      state.settings.selectedLists.standard = list
      state.cached = false
    },
    setSelectedAddons(state, selectedLists) {
      state.cached = false
      state.settings.selectedLists.addon = selectedLists
    },
    unsetSelectedAddon(state, list) {
      state.cached = false
      let index = state.settings.selectedLists.addon.indexOf(list.id)
      if (index !== -1) {
        console.log('Found list')
        state.settings.selectedLists.addon.splice(index, 1)
      }
    },
    setLookupPhrase(state, phrase) {
      console.log("state,lookup =", phrase.toLowerCase())
      state.lookup = phrase.toLowerCase()
    },
    saveDoc(state, doc) {
      state.doc = doc
    },
    setSelectedManuscripts(state, manuscriptAbbs) {
      state.settings.selectedManuscripts = manuscriptAbbs

      EventBus.$emit('setSelectedManuscripts')
    },
    setParsedDocContent(state, content) {
      state.parsedDoc = content
    },
    setStagedNominations(state, abbs) {
      state.stagedNominations = abbs
    },
    clearStagedNominations(state) {
      state.stagedNominations = []
    },
    createMissedAbbsMap(state) {
      state.missedAbbs = new Map()
    },
    incrementMissedAbb(state, abb) {
      console.log('look for missed abb in map', state.missedAbbs)
      if (state.missedAbbs.has(abb.abb)) {
        let missed = state.missedAbbs.get(abb.abb)
        missed.counter++
        state.missedAbbs.set(abb.abb, missed)
      } else {
        let missed = { ...abb, counter: 1 }
        state.missedAbbs.set(abb.abb, missed)
      }
    },
    forgetMissedAbb(state, abb) {
      state.missedAbbs.delete(abb)
    },
    setSelectedWord(state, word) {
      state.selectedWord = word
    },
    setTargetList(state, targetList) {
      console.log(targetList)
      state.targetList = targetList
    },
    setSharedList(state, sharedList) {
      state.sharedList = sharedList
    },
    setSharedAbbs(state, abbs) {
      state.sharedAbbs = abbs
    },
    setConflictingAbbs(state, abbs) {
      console.log('setConflictingAbbs:', abbs)
      state.conflictingAbbs = abbs
    },
    setSettings(state, newSettings) {
      state.settings = newSettings
      if (newSettings.selectedManuscripts == null) {
        console.log('No manuscripts')
        state.settings.selectedManuscripts = []
      }
      console.log('set settings to:', newSettings)
      if (
        newSettings.selectedLists.standard == '' ||
        newSettings.selectedLists.standard == null ||
        newSettings.selectedLists.standard == undefined
      ) {
        api
          .getUserLists()
          .then((resp) => {
            console.log('user lists:')
            if (resp.data.length > 0) {
              resp.data.map((list) => {
                if (list.type == 0) {
                  state.settings.selectedLists.standard = list.id
                  state.targetList = list.id
                }
              })
            }
          })
          .catch((err) => {
            console.log('state couldnt get user lists', err)
          })
      }
    },
    setFontSettings(state, newFontSettings) {
      console.log('Font settings updated')
      state.settings.font = newFontSettings.font
      state.settings.behaviour = newFontSettings.behaviour
      EventBus.$emit('fontSettingsUpdated')
    },
    setFontSize(state, fontSize) {
      state.settings.font.size = parseInt(fontSize)
    },
    setFontColor(state, colors) {
      console.log('set state color ID to', colors.colorID)
      state.settings.font.colorID = colors.colorID
      state.settings.font.background = colors.background
      state.settings.font.foreground = colors.foreground
    },
    setCustomColors(state, customColors) {
      state.settings.font.foreground = customColors.foreground
      state.settings.font.background = customColors.background
      state.settings.font.colorID = 6
      state.settings.font.customColors = customColors
    },
    setLocalSession(state, local) {
      state.local = local
    },
    joinRemoteSession(state, session) {
      VueCookie.set('session', JSON.stringify(session), { expires: '8h' })
      state.session = session
    },
    closeRemoteSession(state) {
      state.session = {}
    },
    setSessionId(state, id) {
      state.session.id = id
    },
    setSessionConnected(state, connected) {
      state.session.connected = connected
      if (!connected) {
        VueCookie.delete('session')
        state.session = emptySession
      }
    },
    setSessionPassword(state, password) {
      state.session.password = password
    },
    clearClients(state) {
      state.clients = 0
      EventBus.$emit('clientListUpdated', null)
    },
    updateClients(state, clients) {
      state.clients = clients
      EventBus.$emit('clientListUpdated', null)
    },
    storeSettings(state) {
      api
        .saveSettings(state.settings)
        .then((resp) => {
          console.log('saved settings:', resp.data)
        })
        .catch((err) => {
          console.error("couldn't save settings on logout", err)
          this.$toast.error('Kunde inte spara dina inställningar.')
        })
    },
    getStoredSettings(state) {
      api.getSettings().then((resp) => {
        state.settings = resp.data
        if (!state.settings.selectedLists.addon) {
          state.settings.selectedLists.addon = []
        }
      })
    },
    subscribeList() {
      /* console.log("subscribe to list:", state.userData.subscriptions, id)
       if (state.userData.subscriptions == undefined) {
         state.userData.subscriptions = [].push(id)
         return
       }
       const index = state.userData.subscriptions.indexOf(id)
       if (index > -1) {
         return
       }
       state.userData.subscriptions.push(id)
       */
    },
    unsubscribeList() {
      /* console.log("unsubscribe from list:", state.userData.subscriptions)

       const index = state.userData.subscriptions.indexOf(id)
       if (index > -1) {
         state.userData.subscriptions.splice(index)
       }*/
    },
    addInterpreter() {},
    removeInterpreter() {},
    reset(state) {
      state.settings.selectedManuscripts = []
      state.settings.selectedLists.addon = []
      state.userData.subscriptions = []
      EventBus.$emit('resetStore')
      state.session = {
        host: false,
        connected: false,
        id: '',
        zoom: {
          apiToken: '',
        },
      }
      api.setIdToken(state.token)
    },
    logout(state) {
      this.commit('storeSettings')

      state.status = ''
      state.doc = ''
      state.token = ''
      state.userData = ''
      state.sharedAbbs = []
      state.sharedList = { id: '', base: null }
      state.session = {
        host: false,
        connected: false,
        id: '',
        zoom: {
          apiToken: '',
        },
      }

      state.licenseKey = ''
      state.settings = {}
      state.targetList = {}
      localStorage.removeItem('user-token')

      router.push('/login')
    },
    beforeClose(state) {
      state.session = {
        host: false,
        connected: false,
        id: '',
        zoom: {
          apiToken: '',
        },
      }
    },
  },
  getters: {
    getUsers(state) {
      return state.users.values()
    },
    getUserId(state) {
      return state.userData.id
    },
    getFocus(state) {
      return state.focus
    },
    getModalOpen(state) {
      return state.modalOpen
    },
    selectedManuscripts(state) {
      return state.settings.selectedManuscripts
    },
    subscriptions(state) {
      return state.userData.subscriptions
    },
    isAuthenticated: (state) => !!state.token,
    getSettings(state) {
      return state.settings
    },
  },
  actions: {
    AUTH_REQUEST: ({ commit }, user) => {
      return new Promise((resolve, reject) => {
        commit('AUTH_REQUEST')
        api
          .login(user)
          .then((resp) => {
            const token = resp.data.id
            const id = resp.data.id
            VueCookie.set('user-id', id, 1)

            localStorage.setItem('user-token', token)

            commit('AUTH_SUCCESS', token)
            //dispatch("USER_REQUEST")

            commit('setUserData', resp.data)
            api.setIdToken(token)
            commit('setSettings', resp.data.settings)
            console.log('load user settings')
            commit('ADD_USER', resp.data)
            setTimeout(() => resolve(resp), 250)
          })
          .catch((err) => {
            commit('AUTH_ERROR', err)
            console.log('error:', err)
            VueCookie.delete('user-id')

            localStorage.removeItem('user-token')
            reject(err)
          })
      })
    },
  },
})
