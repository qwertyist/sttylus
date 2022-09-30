<template>
  <div class="quillWrapper" :style="wrapper">
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
import api from "../api/api";
import EventBus from "../eventbus";
import Quill from "quill";
import Preview from "./tabula/preview";
import Text from "./tabula/text.js";
export default {
  data() {
    return {
      quill: null,
      parent: null,
      name: "",
      fullscreen: false,
      mouse: {
        mousedown: false,
        x: 0,
        y: 0,
        clicks: 0
      },
      timer: null,
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
      console.log("nav:", this.nav)
      return {
        backgroundColor: this.settings.font.backgroundColor,
        padding: 20 + "px",
        height: 100 + "vh",
      };
    },
  },
  mounted() {
    this.initializeEditor();
    this.quill.version = 0;
    /* this.quill.setText(
      "We have been subordinate to our limitations until now. The time has come to cast aside these bonds and to elevate our consciousness to a higher plane. It is time to become a part of all things.",
      "init"
    );*/
    window.addEventListener("message", m => {
      //console.log("message:", m)
      this.parentMessageHandler(JSON.parse(m.data))
    })
    this.addEventListeners();
    setTimeout(() => {
      this.quill.focus();
    }, 500)
  },
  beforeDestroy() {
    console.log("destroying textview");
    this.removeEventListeners();
    Text.saveTextSettings(this.settings);
  },
  methods: {
    parentMessageHandler(m) { 
      //console.log("parent message type:", m.type)
      switch(m.type) {
        case "init":
          this.quill.version = m.version
          this.quill.setContents(m.delta)
          break
        case "kill":
          nw.Window.get().Close(true)
          break
        case "delta":
          this.quill.updateContents(m.delta)
          this.quill.version = m.version
          this.quill.setSelection(this.quill.scroll.length())
          //console.log("update text")
          break
        case "clear":
          //console.log("clear text")
          break
        case "inc":
          this.changeTextSize(true)
          //console.log("increase text size")
          break
        case "dec":
          this.changeTextSize(false)
          //console.log("decrease text size")
          break
        case "color":
          this.changeColor()
          //console.log("change text color")
          break
        default:
          //console.log("misc message...")
      }
    },
    addEventListeners() {
      this.addMouseEventListeners()
    },
    removeEventListeners() {
      document.removeEventListener("mousedown")
    },
    changeTextSize(inc) {
      this.settings.font.fontSize = Text.changeTextSize(
        inc,
        this.settings.font.fontSize
      );
    },
    changeColor() {
      let colors = Text.changeColor(this.settings.font.colorID);
      this.settings.font.backgroundColor = colors.background;
      this.settings.font.color = colors.foreground;
      this.settings.font.colorID = colors.colorID;
    },
    addMouseEventListeners() {
      document.addEventListener("mousedown", (e) => {
        this.mouse.mousedown = true
        this.mouse.clicks++;
        this.mouse.x = e.clientX;
        this.mouse.y = e.clientY;
        if (this.mouse.clicks === 1) {
          this.timer = setTimeout(() => {
              this.$toast.info("Dubbelklicka för växla att till/från fullskärm")
              this.mouse.clicks = 0;
          }, 700)
        } else {
          clearTimeout(this.timer);
          this.mouse.clicks = 0;
          nw.Window.get().toggleFullscreen()
        }
      })

    },
    focus() {
      Text.initText()
      let settings = Text.loadTextSettings();
      this.settings.font = settings.font;

      /*
      this.$nextTick(() => {
        this.quill.focus();
        this.quill.setSelection(this.quill.getText().length);
      })
      */
      setTimeout(() => {
        this.quill.focus();
      }, 5)
    },
    clear() {
      window.scrollTo(0, 0);
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
      this.settings.behaviour = settings.behaviour;
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

      Preview.className = "prerender"
      Quill.register("formats/preview", Preview, true);
      //    Quill.register("modules/keyboard", keyboard, true);
      this.quill = new Quill(this.$refs.quillContainer, editorConfig);
    },
    onHello(msg) {
      alert("parent says hello", msg)
    },
    onDelta(delta) {
      //console.log("onDelta:", delta)
    },
    onSizeChange(inc) {
      console.log("onSizeChange:", inc)
    },
    onColorChange() {
      console.log("onColorChange")
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

.prerender {
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