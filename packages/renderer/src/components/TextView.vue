<template>
  <div class="quillWrapper" :style="wrapper">
    <slot name="toolbar"></slot>
    <div
      ref="quillContainer"
      :class="{ 'ql-container': true }"
      :style="settings.font"
      spellcheck="false"
    ></div>
    <AddAbbreviation ref="addAbb" />
  </div>
</template>
<script>
import api from "../api/api";
import EventBus from "../eventbus";
import Quill from "quill";
import keyboard from "./tabula/keyboard";
import Preview from "./tabula/preview";
import AddAbbreviation from "./modals/AddAbbreviation.vue";
import Text from "./tabula/text.js";
import wsConnection from "./tabula/websocket";
export default {
  components: {
    AddAbbreviation,
  },
  props: {
    nav: Boolean,
  },
  data() {
    return {
      quill: null,
      websocket: null,
      presentation: false,
      child: null,
      name: "",
      capitalize: true,
      version: 0,
      settings: {
        font: {
          fontFamily: "Arial",
          fontSize: "45px !important",
          backgroundColor: "#000000",
          color: "#ffff00",
          lineHeight: "1.2",
          height: "100vh",
          colorID: 2,
          backgroundColor: "#000000",
          customColors: {},
        },
        behaviour: {
          capitalizeOnNewLine: true,
        },
      },
      margins: null,
      userColor: "red",
    };
  },
  computed: {
    wrapper() {
      return {
        backgroundColor: this.settings.font.backgroundColor,
        padding: 20 + "px",
        height: (100 - (this.nav ? 5 : 0)) + "vh",
      };
    },
  },
  mounted() {
    api
      .cacheAbbs()
      .then((resp) => { })
      .catch((err) => {
        console.log("couldn't create cache", err);
      });
    this.websocket = null;
    this.initializeEditor();
    this.quill.version = 0;
    /* this.quill.setText(
      "We have been subordinate to our limitations until now. The time has come to cast aside these bonds and to elevate our consciousness to a higher plane. It is time to become a part of all things.",
      "init"
    );*/
    this.addEventListeners();
    setTimeout(() => {
      this.focus();
    }, 500)

  },
  beforeDestroy() {
    console.log("destroying textview");
    if (this.websocket) {
      this.websocket.close();
    }
    this.websocket = null;
    this.removeEventListeners();
    Text.saveTextSettings(this.settings);
  },
  methods: {
    addEventListeners() {
      EventBus.$on("sharedAbb", this.sharedAbbs)
      EventBus.$on("addAbbreviation", this.openAddModal);
      EventBus.$on("setSelectedManuscripts", this.loadManuscripts)
      EventBus.$on("fontSettingsUpdated", this.loadTextSettings)
      EventBus.$on("refocus", this.focus);
      EventBus.$on("clear", this.clear);
      EventBus.$on("newLine", this.newline);
      EventBus.$on("sizeChange", this.changeTextSize);
      EventBus.$on("colorChange", this.changeColor);

      EventBus.$on("joinRemoteSession", this.joinSession);
      EventBus.$on("createSession", this.createSession);
      EventBus.$on("leaveRemoteSession", this.leaveSession)
      EventBus.$on("sendSessionData", this.sendSessionData)

      EventBus.$on("connectLocal", this.connectLocal)
      EventBus.$on("disconnectLocal", this.disconnectLocal)
      EventBus.$on("startPresentation", this.startPresentation)
      EventBus.$on("stopPresentation", this.stopPresentation)
      EventBus.$on("websocketConnected", this.websocketConnected)
      EventBus.$on("websocketClosed", this.websocketClosed)
      EventBus.$on("websocketDropped", this.websocketDropped)
      EventBus.$on("websocketFailed", this.websocketFailed)
      EventBus.$on("websocketReconnecting", this.websocketReconnecting)
      EventBus.$on("clientConnected", this.clientConnected)
      EventBus.$on("clientDisconnected", this.clientDisconnected)
    },
    removeEventListeners() {
      EventBus.$off("sharedAbb")
      EventBus.$off("addAbbreviation");
      EventBus.$off("setSelectedManuscripts")
      EventBus.$off("fontSettingsUpdated")
      EventBus.$off("refocus");
      EventBus.$off("clear");
      EventBus.$off("newLine");
      EventBus.$off("sizeChange");
      EventBus.$off("colorChange");
      EventBus.$off("joinRemoteSession");
      EventBus.$off("createSession");
      EventBus.$off("connectLocal")
      EventBus.$off("disconnectLocal")
      EventBus.$off("websocketConnected")
      EventBus.$off("websocketClosed")
      EventBus.$off("websocketDropped")
      EventBus.$off("websocketFailed")
      EventBus.$off("websocketReconnecting")
      EventBus.$off("startPresentation")
      EventBus.$off("stopPresentation")
      EventBus.$off("clientConnected")
      EventBus.$off("clientDisconnected")
    },
    createSession() {
      this.websocket.createsession();
    },
    joinSession(id) {
      console.log("join Session with id:", id);
      this.clear()
      this.websocket = new wsConnection(this.quill, this.$collabServer + "conn/" + id);
    },
    websocketConnected() {
      this.$toast.success("Du kopplades upp")
    },
    websocketFailed(err) {
      this.$toast.error("Anslutningen misslyckades")
      if(this.websocket) {
        this.websocket = null;
      }
    },
    websocketDropped(err) {
      this.$toast.error("Anslutningen till tolkningen bröts")
      if (err) {
        this.$toast.warning("Fel:", err)
      }
    },
    websocketClosed(msg) {
      this.$toast.info("Du kopplades ner", msg)
    },
    websocketReconnecting(tries = 1) {
      console.log(tries, "försöket...")
        const msg = "Försöker ansluta igen ... (#" + tries + ")"
        this.$toast.info(msg, tries)
        
    },
    clientConnected(c) {
      if (c == "user") {
        this.$toast.success("En tolkanvändare anslöt")
        return
      }
      this.$toast.success("En tolk anslöt")
    },
    clientDisconnected(c) {
      if (c == "user") {
        this.$toast.success("En tolkanvändare kopplade ner")
        return
      }
      this.$toast.success("En tolk kopplade ner")
    },
    sendSessionData(data) {
      console.log("should send session data")
      this.websocket.sendSessionData(data)
    },
    startPresentation(child) {
      this.presentation = true
      console.log("emit to child")
      nw.Window.getAll(windowList => { this.child = windowList[1] })
      setTimeout(() => {
        const msg = { type: "init", delta: this.quill.getContents(), version: this.quill.version }
        this.child.window.postMessage(JSON.stringify(msg), "*")
      }, 125)
    },
    stopPresentation() {
      this.presentation = false
      const msg = { type: "kill" }
      this.child.window.postMessage(JSON.stringify(msg), "*")
      this.child.close()
      this.child = null

    },
    connectLocal(reconnect = false) {
      if (!reconnect && this.websocket) {
        this.$toast.warning("Du måste koppla ner först")
        return
      }
      this.clear()
      this.websocket = new wsConnection(this.quill, "ws://127.0.0.1:80/conn/local")
    },
    disconnectLocal() {
      if (!this.websocket) {
        return
      }
      this.websocket.close()
      this.websocket = null;
      this.quill.version = 0
    },
    leaveSession() {
      if (!this.websocket) {
        return
      }
      this.websocket.close()
      this.websocket = null;
      this.quill.version = 0
      this.quill.setText("");
    },
    changeTextSize({ inc, send }) {
      console.log("presentation:", this.presentation, "send:", send)
      if (this.presentation && send) {
        console.log("send size change", inc ? "inc" : "dec")
        const msg = { type: inc ? "inc" : "dec" }
        this.child.window.postMessage(JSON.stringify(msg), "*")

        return
        //TODO: IF-sats kolla om inställningarna delas mellan skärmar
      }

      this.settings.font.fontSize = Text.changeTextSize(
        inc,
        this.settings.font.fontSize
      );
      this.$store.commit("setFontSize", this.settings.font.fontSize.replace("px", ""))
    },
    changeColor(send) {
      if (this.presentation && send) {
        const msg = { type: "color" }
        this.child.window.postMessage(JSON.stringify(msg), "*")
      }

      let colors = Text.changeColor(this.settings.font.colorID);
      this.settings.font.backgroundColor = colors.background;
      this.settings.font.color = colors.foreground;
      this.settings.font.colorID = colors.colorID;
      this.$store.commit("setFontColor", colors)
    },
    openAddModal(phrase) {
      this.$store.commit("setSelectedWord", phrase);
      this.$bvModal.show("addAbb");
    },
    addAbb() {

    },
    sharedAbbs(abb) {
      if (abb.me) {
        this.websocket.sendSharedAbb(abb)
      } else {
        console.log("RXAbb:", abb)
        let baseListId = this.$store.state.sharedList.base

        if (abb.create) {
          this.$toast.info('"' + abb.abb + '" → "' + abb.word + '" skapades')
          if (baseListId) {
            api.createAbb(baseListId, abb).then(resp => {
              console.log("other user created abb");
            }).catch(err => {
              console.error("other user couldnt create abb", err)
            })
          }
        }
        if (abb.delete) {
          this.$toast.info('"' + abb.abb + '" togs bort')
          api.deleteAbb(baseListId, abb).then(resp => {
            console.log("other user deleted abb");
          }).catch(err => {
            console.error("other user couldnt delete abb", err)
          })
        }
        if (abb.override) {
          this.$toast.info('"' + abb.abb + '" skrevs över')
        }
        EventBus.$emit("sharedAbbEvent")
      }
    },
    focus() {
      Text.initText()
      let settings = Text.loadTextSettings();
      this.settings.font = settings.font;

      /*
      this.$nextTick(() => {
        this.quill.focus();
      })
      */
      let editor = document.querySelector(".ql-editor")
      setTimeout(() => {
        editor.focus();
        this.quill.focus();
        this.quill.setSelection(this.quill.getText().length);
      }, 25)
    },
    clear() {
      window.scrollTo(0, 0);
      if (this.websocket) {
        this.websocket.sendClear()
        this.quill.version = 0
      }
    },
    newline(scroll) {
    },
    initializeEditor() {
      this.setupEditor();
      this.$emit("ready", this.quill);
    },
    loadManuscripts() {
      Text.initText()
    },
    loadTextSettings() {
      let settings = Text.loadTextSettings();
      console.log("settings", settings)
      this.settings.font = settings.font;
      this.settings.behaviour = settings.behaviour;
    },
    setupEditor() {
      this.loadTextSettings();
      Text.initText()

      const editorConfig = {
        debug: false,
        theme: "snow",
        modules: {
          toolbar: null,
          keyboard: {
            capitalizeOnNewLine: this.settings.behaviour.capitalizeOnNewLine,
            manuscriptEditor: false,
          },
        },
      };

      Quill.register("modules/keyboard", keyboard, true);
      Quill.register("formats/preview", Preview, true);
      //    Quill.register("modules/keyboard", keyboard, true);
      this.quill = new Quill(this.$refs.quillContainer, editorConfig);

      this.quill.on("text-change", (delta, oldDelta, source) => {
        if (this.presentation && this.child) {
          const msg = {
            type: "delta", delta: delta, version: this.quill.version
          }
          this.child.window.postMessage(JSON.stringify(msg), "*")
        }
        if (this.websocket && source != "collab" && source != "init") {
          this.websocket.sendDelta(delta);
        }
      });
      this.quill.clipboard.addMatcher(Node.ELEMENT_NODE, (node, delta) => {
        delta.ops = delta.ops.map(op => {
          if (typeof op.insert !== "string") { return { insert: "" } }
          return {
            insert: op.insert
          }
        })
        return delta
      })
    },
    saveSettings() {
      let settings = this.$store.state.settings;
      console.log("save settings:", settings);
    },
  },
};
</script>
<style src="./tabula/quill.scss"></style>
<style src="./tabula/quill.snow.css"></style>
<style>
/*.ql-editor {
  outline: 1px solid red !important
}
*/

.preRender {
  opacity: 0.6;
}
.navMargin {
  position: fixed;
  top: 5vh;
  width: 100vw;
}
.wrapper {
  height: 100vh;
}
</style>
