<template>
  <div @dblclick="doubleClickHandler" class="quillWrapper" :style="wrapper">
    <slot name="toolbar"></slot>
    <div
      ref="quillContainer"
      :class="{ 'ql-container': true }"
      :style="settings.font"
      spellcheck="false"
    ></div>
  </div>
</template>
<script>
import EventBus from "../eventbus";
import Quill from "quill";
import Preview from "./tabula/preview";
import Text from "./tabula/text"
import wsConnection from "./tabula/websocket";
export default {
  props: {
    nav: Boolean,
  },
  data() {
    return {
      quill: null,
      websocket: null,
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
      console.log("nav:", this.nav)
      return {
        backgroundColor: this.settings.font.backgroundColor,
        padding: 20 + "px",
        height: 100 + "vh",
      };
    },
  },
  mounted() {
    this.websocket = null;
    this.initializeEditor();
    this.quill.version = 0;
    /* this.quill.setText(
      "We have been subordinate to our limitations until now. The time has come to cast aside these bonds and to elevate our consciousness to a higher plane. It is time to become a part of all things.",
      "init"
    );*/
    this.addEventListeners();
    setTimeout(() => {
      this.quill.focus();

      let uri = window.location.href.split("?");
      console.log("URI:", uri);
      if (uri.length == 2) {
        let vars = uri[1].split("&");
        let id;
        let tmp = "";
        vars.forEach(function (v) {
          tmp = v.split("=");
          if (tmp.length == 2) id = tmp[1];
        });
        if (id !== "") {
          console.log("ID:", id);
          this.$store.commit("setSessionID", id);
          this.joinSession(id)
          return
        }
        this.quill.setContents(JSON.parse('{"ops":[{"insert":"Det här är en testbokning för att se hur STTylus distanstolkningar fungerar på olika enheter. "},{"retain":1},{"insert":"\\nSka vi byta?\\nSka vi byta "},{"insert":{"preview":"grejer? "}},{"insert":"\\n"}]}'))
        this.quill.updateContents()

      } else {
        this.$store.commit("setSessionID", "local")
        this.joinSession("local")
      }
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
      var editor = document.getElementsByClassName("ql-editor")
      editor[0].addEventListener("scroll", (e) => {
        editor[0].scrollTop = editor[0].scrollHeight;
      }) 
      document.addEventListener('keypress', (e) => {
        if (e.key == "~") {
          e.preventDefault();
          console.log(this.quill.getContents())

        }
      });
      EventBus.$on("addAbbreviation", this.openAddModal);
      EventBus.$on("refocus", this.focus);
      EventBus.$on("clear", this.clear);
      EventBus.$on("newLine", this.newline);
      EventBus.$on("sizeChange", this.changeTextSize);
      EventBus.$on("colorChange", this.changeColor);
      EventBus.$on("joinRemoteSession", this.joinSession);
      EventBus.$on("createSession", this.createSession);
      EventBus.$on("setSessionData", this.setSessionData)
      EventBus.$on("leaveRemoteSession", this.leaveSession)


      EventBus.$on("websocketConnected", this.websocketConnected)
      EventBus.$on("websocketClosed", this.websocketClosed)
      EventBus.$on("websocketDropped", this.websocketDropped)
      EventBus.$on("websocketFailed", this.websocketFailed)
      EventBus.$on("websocketReconnecting", this.websocketReconnecting)
    },
    removeEventListeners() {
      EventBus.$off("addAbbreviation");
      EventBus.$off("refocus");
      EventBus.$off("clear");
      EventBus.$off("newLine");
      EventBus.$off("sizeChange");
      EventBus.$off("colorChange");
      EventBus.$off("joinRemoteSession");
      EventBus.$off("createSession");
      EventBus.$off("setSessionData")

      EventBus.$off("websocketConnected")
      EventBus.$off("websocketClosed")
      EventBus.$off("websocketDropped")
      EventBus.$off("websocketFailed")
      EventBus.$off("websocketReconnecting")
    },
    doubleClickHandler() {
      console.log("open settings");
      this.$bvModal.show("consumerSettings")
    },
    joinSession(id) {
      console.log("join Session with id:", id);
      this.clear()
      if (id == "local") {
        let uri = window.location.href.split("http://")[1]
        console.log("uri:", uri)

        this.websocket = new wsConnection(this.quill,  "ws://" + uri + "conn/" + "local");
        return
      }
      this.websocket = new wsConnection(this.quill, this.$collab + "conn/" + id);

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
    setSessionData(data) {
      console.log("textview got event, should send session data")
      this.websocket.sendSessionData(data)
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
    openAddModal(phrase) {
      console.log("add phrase", phrase);
      this.$store.commit("setSelectedWord", phrase);
      this.$bvModal.show("addAbb");
    },
    addAbb() {

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
      console.log(this.settings.font)
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