<template>
  <div id="app">
    <Dashboard />
    <Navigation
      :key="editorKey"
      :view="view"
      v-show="showNav"
      id="navigation"
    />
    <div @dblclick="dbclick()">
      <Settings v-show="view == 'settings'" />
      <TextView
        @click="focusText()"
        v-show="view == 'tabula'"
        :key="editorKey"
        :nav="showNav"
        :chat="showChat"
        ref="tabula"
      />
    </div>
    <Chat :nav="showNav" ref="chat" />
  </div>
</template>

<script>
import EventBus from './eventbus.js'
import Settings from './Settings.vue'
import Navigation from './components/Navigation.vue'
import TextView from './components/TextView.vue'
import Chat from './components/Chat.vue'
import Dashboard from './components/modals/DashboardModal.vue'
import api from './api/api.js'
import db from './store/db.js'
export default {
  name: 'Tabula',
  components: {
    Navigation,
    Dashboard,
    TextView,
    Settings,
    Chat,
  },
  data() {
    return {
      focused: null,
      showNav: true,
      showChat: false,
      editorKey: '0',
      view: 'tabula',
      //view: "settings",
    }
  },
  computed: {
    foreground() {
      return this.$store.state.settings.font.foreground
    },
    background() {
      return this.$store.state.settings.font.background
    },
  },
  mounted() {
    this.setupDexieStorage() 
    this.$bvModal.show('dashboard')
    //TEMP
    //this.$bvModal.show("support");
    //ENDTEMP
    this.$store.commit('setModalOpen', false)
    window.addEventListener('keydown', this.hotkeys)
    EventBus.$on('cacheAbbs', this.cacheAbbs)
    EventBus.$on('reloadEditor', this.reload)
    EventBus.$on('openSettings', this.openSettings)
    EventBus.$on('openTextView', this.openTextView)
    // TODO: Rename
    EventBus.$on('modalClosed', this.abbModalClosed)
    EventBus.$on('abbModalClosed', this.abbModalClosed)
    EventBus.$on('toggleCollab', this.toggleCollab)
    EventBus.$on('chatFocused', this.chatFocused)
    EventBus.$on('chatBlurred', this.chatBlurred)
    EventBus.$on('chatOpened', () => {
      this.showChat = true
      this.$store.commit('setFocus', 'chat')
    })
    EventBus.$on('chatClosed', () => {
      this.focusText()
      this.showChat = false
    })
    this.$store.commit('setFocus', 'text')
    EventBus.$on('chatHidden', () => {
      this.showChat = false
      this.focusText()
      this.$store.commit('setFocus', 'text')
    })
    EventBus.$on('toggleNav', () => {
      if (this.view != 'settings') {
        this.showNav = !this.showNav
      }
    })
    EventBus.$on('closeNav', () => {
      if (this.view != 'settings') {
        this.showNav = false
      }
    })

    this.$nextTick(() => {
      this.focused = 'text'
    })
  },
  destroyed() {},

  beforeDestroy() {
    window.removeEventListener('keydown', this.hotkeys)
    EventBus.$off('toggleCollab')
    EventBus.$off('chatClosed')
    EventBus.$off('chatOpened')
    EventBus.$off('chatFocused')
    EventBus.$off('reloadEditor')
    EventBus.$off('toggleNav')
    EventBus.$off('closeNav')
    EventBus.$off('chatHidden')
    EventBus.$off('openSettings')
    EventBus.$off('openTextView')
    EventBus.$off('abbModalClosed')
  },
  methods: {
    setupDexieStorage() {
      db.syncData()
    },
    dbclick() {
      if (this.view != 'settings') {
        this.showNav = !this.showNav
      }
    },
    reload() {
      this.editorKey =
        Math.floor(Math.random() * (999999 - 100000 + 1)) + 100000
    },
    hotkeys(e) {
      if (e.ctrlKey && this.view != 'tabula') {
        if (['1', '2', '3', '4', '5'].indexOf(e.key) != -1) {
          EventBus.$emit('changeStandardList', e.key)
          e.preventDefault()
        }
      }

      if (e.key == 'F1') {
        this.$bvModal.show('support')
      }

      if (e.key == 'F3') {
        e.preventDefault()
        EventBus.$emit('toggleCollab', '')
      }
      if (e.key == 'F9') {
        e.preventDefault()
        this.$bvModal.show('remote')
      }

      if (e.key == 'F10') {
        if (e.shiftKey) {
          e.preventDefault()
          console.log('Close chat')
          EventBus.$emit('hideChat')
          this.focusText()
          return
        }
        e.preventDefault()
        this.toggleFocus()
      }

      if (e.key == 'F5') {
        e.preventDefault()
        if (this.view == 'tabula') {
          this.$bvModal.hide('support')
          this.$bvModal.hide('addAbb')
          this.openSettings()
        } else {
          EventBus.$emit('showTextView')
          this.$bvModal.hide('support')
          this.$bvModal.hide('addAbb')
          this.openTextView()
        }
        //this.$router.push("/settings");
      }
    },
    openSettings() {
      if (this.showChat) {
        EventBus.$emit('hideChat')
      }
      if (this.view == 'tabula') {
        this.$store.commit('setModalOpen', true)
        EventBus.$emit('showSettings')
        api
          .saveSettings(null)
          .then(() => {})
          .catch((err) => {
            console.error("couldn't save settings", err)
            this.$toast.error('Dina inställningar kunde inte sparas')
          })
        window.scrollTo(0, 0)
        this.view = 'settings'
        this.showNav = true

        EventBus.$emit('chatHidden')
        EventBus.$emit('openSettings')
      }
    },
    openTextView() {
      this.$store.commit('setFocus', 'text')
      if (this.showChat) {
        EventBus.$emit('showChat')
      }
      if (this.view == 'settings') {
        api
          .saveSettings(null)
          .then(() => {
            this.cacheAbbs()
          })
          .catch((err) => {
            console.error("couldn't save settings", err)
            this.$toast.error('Dina inställningar kunde inte sparas')
          })

        this.view = 'tabula'
        this.showNav = true
        this.$store.commit('setModalOpen', false)
        EventBus.$emit('closeManuscriptEditor')
        EventBus.$emit('openTextView')
        EventBus.$emit('refocus', true)
      }
    },
    cacheAbbs() {
      console.log('cached: ', this.$store.state.cached)
      if (!this.$store.state.cached) {
        api
          .cacheAbbs()
          .then(() => {
            EventBus.$emit('getAbbCache')
          })
          .catch((err) => {
            console.log("couldn't create cache", err)
          })
      }
    },
    focusText() {
      this.focused = 'text'
      EventBus.$emit('refocus')
      this.$store.commit('setFocus', 'text')
    },
    focusChat() {
      this.focused = 'chat'
      EventBus.$emit('focusChat')
      this.$store.commit('setFocus', 'chat')
    },
    chatFocused() {
      this.focused = 'chat'
    },
    chatBlurred() {
      this.focused = 'text'
    },
    toggleFocus() {
      console.log('this.showChat', this.showchat)
      if (this.showChat) {
        if (this.focused == 'text') {
          this.focusChat()
        } else if (this.focused == 'chat') {
          this.focusText()
        } else {
          this.focusText()
        }
      } else {
        this.showChat = true
        EventBus.$emit('focusChat')
      }
    },
    toggleCollab() {
      EventBus.$emit('sendReadySignal')
    },
    abbModalClosed() {
      console.log(this.focused)
      if (this.focused == 'text') {
        this.focusText()
      } else if (this.showChat && this.focused == 'chat') {
        this.focusChat()
      }
    },
  },
}
</script>

<style>
body,
html {
  overflow: hidden;
}
#app {
  zoom: 1;
  -moz-transform: scale(1);
  -moz-transform-origin: 0 0;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
#navigation {
  position: relative;
  width: 100%;
  z-index: 1000;
}
.navbar {
  height: 5vh !important;
}
</style>
