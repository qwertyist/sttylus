<template>
  <div id="app">
    <Navigation :key="editorKey" :view="view" v-show="showNav" id="navigation" />
    <div @dblclick="dbclick">
      <Settings v-show="view == 'settings'" />
      <TextView
        v-show="view == 'tabula'"
        :key="editorKey"
        :nav="showNav"
        :chat="showChat"
        ref="tabula"
      />
      <RemotePane />
    </div>
  </div>
</template>

<script>
import EventBus from "./eventbus.js";
import Settings from "./Settings.vue";
import Navigation from "./components/Navigation.vue";
import TextView from "./components/TextView.vue";
import RemotePane from "./components/RemotePane.vue";
import api from "./api/api.js";
export default {
  name: "Tabula",
  components: {
    Navigation,
    TextView,
    RemotePane,
    Settings,
  },
  data() {
    return {
      showNav: true,
      showChat: false,
      editorKey: "0",
      view: "tabula",
      //view: "settings",
    };
  },
  computed: {
    foreground() {
      return this.$store.state.settings.font.foreground;
    },
    background() {
      return this.$store.state.settings.font.background;
    },
  },
  methods: {
    dbclick() {
      if (this.view != "settings") {
        this.showNav = !this.showNav;
      }
    },
    reload() {
      this.editorKey =
        Math.floor(Math.random() * (999999 - 100000 + 1)) + 100000;
    },
    hotkeys(e) {
      if (e.key == "F1") {
        e.preventDefault();
        this.$bvModal.show("support");
      }
      if (e.key == "F5") {
        e.preventDefault();
        if (this.view == "tabula") {
          this.openSettings();
        } else {
          this.openTextView();
        }
        //this.$router.push("/settings");
      }
    },
    openSettings() {
      this.showChat = false;
      if (this.view == "tabula") {
        this.$store.commit("setModalOpen", true)
        api.saveSettings(null).then(() => {
        }).catch(err => {
          console.error("couldn't save settings", err)
          this.$toast.error("Dina inställningar kunde inte sparas");
        })
        window.scrollTo(0, 0);
        this.view = "settings";
        this.showNav = true;

        EventBus.$emit("chatHidden")
        EventBus.$emit("openSettings")
      }
    },
    openTextView() {
      if (this.view == "settings") {
        api.saveSettings(null).then(() => {
          this.cacheAbbs();
          this.view = "tabula";
          this.showNav = true;
          this.$store.commit("setModalOpen", false)
          EventBus.$emit("closeManuscriptEditor");
          EventBus.$emit("refocus");
        }).catch(err => {
          console.error("couldn't save settings", err)
          this.$toast.error("Dina inställningar kunde inte sparas");
        })
      }
    },
    cacheAbbs() {
      api
        .cacheAbbs()
        .then((resp) => { })
        .catch((err) => {
          console.log("couldn't create cache", err);
        });
    },
  },
  mounted() {
    //TEMP
    //this.$bvModal.show("support");
    //ENDTEMP
    this.cacheAbbs();
    window.addEventListener("keydown", this.hotkeys);

    EventBus.$on("reloadEditor", this.reload);
    EventBus.$on("openSettings", this.openSettings);
    EventBus.$on("openTextView", this.openTextView);
    EventBus.$on("chatHidden", () => {
      this.showChat = false;
    });
    EventBus.$on("toggleNav", () => {
      if (this.view != "settings") {
        this.showNav = !this.showNav;
      }
    });
    EventBus.$on("toggleChat", () => {
      this.showChat = !this.showChat;
    });
    EventBus.$on("closeNav", () => {
      if (this.view != "settings") {
        this.showNav = false;
      }
    });
  },
  destroyed() { },

  beforeDestroy() {
    window.removeEventListener("keydown", this.openSettings);
    EventBus.$off("reloadEditor");
    EventBus.$off("toggleNav");
    EventBus.$off("closeNav");
    EventBus.$off("toggleChat");
    EventBus.$off("chatHidden");
    EventBus.$off("openSettings");
    EventBus.$off("openTextView");
  },
};
</script>

<style>
body,
html {
  overflow: hidden;
}
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
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
