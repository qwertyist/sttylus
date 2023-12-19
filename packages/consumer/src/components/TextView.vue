<template>
   <div @dblclick="doubleClickHandler" class="quillWrapper" :style="wrapper">
      <slot name="toolbar"></slot>
      <div
        v-show="!captions"
        ref="quillContainer"
        id="quill-container"
        :class="{ 'ql-container': true }"
        :style="settings.font"

        spellcheck="false"
      ></div>
      <Captions v-if="captions" :id="id" />
    </div>

</template>
<script>
import EventBus from "../eventbus";
import Quill from "quill";
import Preview from "./tabula/preview";
import Protype from "./tabula/protype";
import Text from "./tabula/text"
import wsConnection from "./tabula/websocket";
export default {
  props: {
    nav: Boolean,
  },
  data() {
    return {
      id: "",
      quill: null,
      websocket: null,
      captions: false,
      password: "",
      name: "",
      capitalize: true,
      version: 0,
      waiting: true,
      settings: {
        font: {
          fontFamily: "Arial",
          fontSize: "45px !important",
          backgroundColor: "#000000",
          color: "#ffff00",
          lineHeight: "1.2",
          height: "100vh",
          colorID: 2,
          padding: 20 + "px",
          backgroundColor: "#000000",
          customColors: {},
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
        height: 100 + "vh",
      };
    },
  },
  mounted() {

    const id = this.$route.params.id
    this.password = this.$route.params.password
    this.websocket = null;
    this.initializeEditor();
    this.quill.version = 0;
    /* this.quill.setText(
      "We have been subordinate to our limitations until now. The time has come to cast aside these bonds and to elevate our consciousness to a higher plane. It is time to become a part of all things.",
      "init"
    );*/
    this.addEventListeners();
    this.quill.focus();
    setTimeout(() => {
      if (id !== "") {
        this.$store.commit("setSessionID", id);
        //this.$store.commit("setSessionPassword", this.password);
        this.joinSession(id, this.password)
        return
      } else {
        this.$store.commit("setSessionID", "local")
        this.joinSession("local")
      }
    }, 500)
  },
  beforeDestroy() {
    if (this.websocket) {
      this.websocket.close();
    }
    this.websocket = null;
    this.removeEventListeners();
    Text.saveTextSettings(this.settings);
  },
  methods: {
    addEventListeners() {
      var editor = document.getElementsByClassName("ql-editor")
      editor[0].addEventListener("scroll", (e) => {
        editor[0].scrollTop = editor[0].scrollHeight;
      })
      document.addEventListener('keypress', (e) => {
        if (e.key == "~") {
          e.preventDefault();

        }
      });
      EventBus.$on("addAbbreviation", this.openAddModal);
      EventBus.$on("refocus", this.focus);
      EventBus.$on("clear", this.clear);
      EventBus.$on("newLine", this.newline);
      EventBus.$on("sizeChange", this.changeTextSize);
      EventBus.$on("colorChange", this.changeColor);
      EventBus.$on("joinRemoteSession", this.joinSession);
      EventBus.$on("joinedEmptySession", this.joinedEmptySession)
      EventBus.$on("joinedRunningSession", this.joinedRunningSession)
      EventBus.$on("createSession", this.createSession);
      EventBus.$on("setSessionData", this.setSessionData)
      EventBus.$on("leaveRemoteSession", this.leaveSession)
      EventBus.$on("notAuthorized", this.enterPassword)

      EventBus.$on("websocketConnected", this.websocketConnected)
      EventBus.$on("websocketClosed", this.websocketClosed)
      EventBus.$on("websocketDropped", this.websocketDropped)
      EventBus.$on("websocketFailed", this.websocketFailed)
      EventBus.$on("websocketReconnecting", this.websocketReconnecting)

      EventBus.$on("TXChat", this.sendChat)
      EventBus.$on("toggleCaptions", this.toggleCaptions)

    },
    removeEventListeners() {
      EventBus.$off("addAbbreviation");
      EventBus.$off("refocus");
      EventBus.$off("clear");
      EventBus.$off("newLine");
      EventBus.$off("sizeChange");
      EventBus.$off("colorChange");
      EventBus.$off("joinRemoteSession");
      EventBus.$off("joinedEmptySession")
      EventBus.$off("joinedRunningSession")
      EventBus.$off("createSession");
      EventBus.$off("setSessionData")
      EventBus.$off("notAuthorized")

      EventBus.$off("websocketConnected")
      EventBus.$off("websocketClosed")
      EventBus.$off("websocketDropped")
      EventBus.$off("websocketFailed")
      EventBus.$off("websocketReconnecting")

      EventBus.$off("TXChat")
      EventBus.$off("toggleCaptions")
    },
    doubleClickHandler() {
      this.$bvModal.show("consumerSettings")
    },
    joinSession(id) {
      this.id = id
      this.clear()
      if (id == "local") {
        let uri = window.location.href.split("http://")[1]
        EventBus.$emit("setQRCodeURL", window.location.href)
        this.websocket = new wsConnection(this.quill,  "ws://" + uri + "conn/" + "local");
        return
      }
      EventBus.$emit("setQRCodeURL", window.location.href)
      this.websocket = new wsConnection(this.quill, this.$collab + "conn/" + id, this.password);

    },
    joinedEmptySession() {
      this.waiting = true
      this.quill.setText("Du är en ansluten till tolkningen, men den har inte börjat än.", "collab")
    },
    joinedRunningSession() {
      this.$toast.success("Distanstolkningen har startat")
      this.quill.setText("...", "collab")
      this.waiting = false
    },
    websocketConnected() {
      this.$toast.success("Du kopplades upp")
    },
    websocketFailed(err) {
      this.$toast.error("Anslutningen misslyckades")
      this.quill.setText("Tyvärr kunde du inte ansluta till distanstolkningen. Meddela din kontaktperson för uppdraget så får du hjälp.\nFelmeddelande: \n" + JSON.stringify(err, ["message", "argumentes", "type", "name"]), "collab")
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
        const msg = "Försöker ansluta igen ... (#" + tries + ")"
        this.$toast.info(msg, tries)

    },
    setSessionData(data) {
      this.websocket.sendSessionData(data)
    },
    enterPassword() {
      this.$bvModal.show("login")
    },
    toggleCaptions(val) {
      this.captions = val
    },
    sendChat(data) {
      this.websocket.sendChat(data)
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
    changeTextSize(inc) {
      this.settings.font.fontSize = Text.changeTextSize(
        inc,
        this.settings.font.fontSize
      );
      this.$store.commit("setFontSize", this.settings.font.fontSize.replace("px", ""))
    },
    changeColor() {
      let colors = Text.changeColor(this.settings.font.colorID);
      this.settings.font.backgroundColor = colors.background;
      this.settings.font.color = colors.foreground;
      this.settings.font.colorID = colors.colorID;
      this.$store.commit("setFontColorID", this.settings.font.colorID)
    },
    focus() {
      this.loadTextSettings();
      /*
      this.$nextTick(() => {
        this.quill.focus();
        this.quill.setSelection(this.quill.getText().length);
      })
      */
      setTimeout(() => {
        this.quill.focus();
        this.quill.setSelection(this.quill.getText().length);
      }, 5)
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
    loadTextSettings() {
      let settings = Text.loadTextSettings();
      this.settings.font = settings.font;
    },
    setupEditor() {
      this.loadTextSettings();
      Text.initText()

      const editorConfig = {
        debug: false,
        theme: "snow",
        readOnly: true,
        modules: {
          toolbar: null,
        },
      };

      Quill.register("formats/preview", Preview, true);
      //    Quill.register("modules/keyboard", keyboard, true);
      this.quill = new Quill(this.$refs.quillContainer, editorConfig);

      this.quill.on("text-change", (delta, oldDelta, source) => {
        if (this.websocket && source != "collab" && source != "init") {
          this.websocket.sendDelta(delta);
        }
      });
    },
    saveSettings() {
      let settings = this.$store.state.settings;
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
  display: none;
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
