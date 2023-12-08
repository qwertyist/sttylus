<template>
  <div>
    <b-modal
      id="consumerSettings"
      size="lg"
      centered
      scrollable
      noFade
      hideBackdrop
      title="Inställningar"
      hideHeader
      hideFooter
      hideHeaderClose
      returnFocus="body"
      ref="consumer-settings-modal"
      @show="onOpen"
      @hide="onClose"
    >
    <b-tabs lazy>
      <b-tab
        style="height: 50vh"
        title="Chat"
      >
      <ConsumerChat />
      </b-tab>
      <b-tab
        title="Textinställningar"
        style="height: 50vh"
        @click="getFontSettings"
      >
      <br />
        <div
          :style="{
            display: 'block',
            height: (orientation == 'LANDSCAPE' ? 25 : 37 ) + 'vh',
            width: '100%',
            overflow: 'hidden',
            fontSize: fontSettings.size + 'px',
            fontFamily: fontSettings.family,
            lineHeight: fontSettings.lineHeight,
            backgroundColor: fontSettings.background,
            color: fontSettings.foreground,
          }"
          v-html="example"
        ></div>
      <b-form inline class="mb-2">
        <div v-show="!isMobile">
          <b-button size="md" @click="decreaseTextSize">
            <b-icon icon="dash-square-fill" aria-hidden="true"></b-icon>
            Mindre
          </b-button>
          <b-button size="md" @click="increaseTextSize">
            <b-icon icon="plus-square-fill" aria-hidden="true"></b-icon>
            Större
          </b-button>
          <b-button size="md" @click="changeColor">Byt färger</b-button>
        </div>
        <b-input-group prepend="Typsnitt">
          <b-form-select
            @change="changeFontFamily"
            prepend="Typsnitt"
            :options="form.family.options"
            v-model="fontSettings.family"
            style="width: 180px"
          />
        </b-input-group>
        <b-input-group prepend="Radavstånd">
          <b-form-select
            @change="changeLineHeight"
            prepend="Typsnitt"
            :options="form.lineHeights"
            v-model="fontSettings.lineHeight"
            style="width: 80px"
          />
        </b-input-group>
      </b-form>
        <template v-if="isMobile">
        <div class="d-flex justify-content-between align-items-center">
            <b-button size="md" @click="decreaseTextSize">
              <b-icon icon="dash-square-fill" aria-hidden="true"></b-icon>
              Mindre
            </b-button>
            <b-button size="md" @click="increaseTextSize">
              <b-icon icon="plus-square-fill" aria-hidden="true"></b-icon>
              Större
            </b-button>
            <b-button size="md" @click="changeColor">Byt färger</b-button>
        </div>
        <br />
      </template>
      </b-tab>
      <b-tab
        title="Övrigt"
        @click="getMiscSettings()"
      >
          <b-form
            style="height: 50vh"
            @submit.prevent=""
          >
          Namn som visas i distanstolkningar:
          <b-row>
            <b-col cols=4>
              <b-form-input v-model="name" />
            </b-col>
            <b-col>
              <b-button type="submit" @click="updateName()">Spara</b-button>
            </b-col>
          </b-row>
        <hr />
        <b-row>
          <b-col cols=4>
          Dela på distanstolkning
          </b-col>
          <b-col>
            <b-button @click="toggleQRCode()" size="md" >Visa QR-kod</b-button>
          </b-col>
        </b-row>
          </b-form>
      </b-tab>
    </b-tabs>
    </b-modal>
  </div>

</template>

<script>
import EventBus from "../../eventbus.ts";
import ConsumerChat from "../ConsumerChat.vue";
export default {
  name: "ConsumerSettings",
  components: { ConsumerChat },
  props: { visible: Boolean },
  data() {
    return {
      modalOpen: false,
      name: "",
      orientation: "PORTRAIT",
      example_full:
        "<em>Såhär</em> ser <b>texten</b> ut...<br /> Dina inställningar lagras till nästa distanstolkning.",
      example_short: "<em>Såhär</em> ser <b>texten</b><br/><u>ut...</u>",
      fontSettings: {
        family: "Arial",
        size: 32,
        foreground: "#ffff00",
        background: "#000000",
        selected: 3,
        lineHeight: 1.25,
      },
      form: {
        family: {
          options: ["Times New Roman", "Arial", "Roboto Mono", "Verdana"],
        },
        lineHeights: [1.0, 1.25, 1.5, 1.75, 2],
        colors: [
          {
            foreground: "#ffffff",
            background: "#000000",
          },
          {
            foreground: "#000000",
            background: "#ffffff",
          },
          {
            foreground: "#ffff00",
            background: "#000000",
          },
          {
            foreground: "#000000",
            background: "#ffff00",
          },
          {
            foreground: "#ffff00",
            background: "#0000ff",
          },
          {
            foreground: "#0000ff",
            background: "#ffff00",
          },
        ],
      },
    };
  },
  computed: {
    isMobile() {
      return this.$mobile;
    },
    example() {
      if(this.isMobile) {
        return this.example_short
      } else  {
        return this.example_full
      }

    }
  },
  mounted() {
    this.getFontSettings()
    this.getMiscSettings()
  },
  methods: {
    getFontSettings() {
    if (localStorage.getItem("fontSettings")) {
      console.log("got stored font settings");
      this.fontSettings = JSON.parse(localStorage.getItem("fontSettings"));
    } else {
      console.log("no stored font settings, storing defaults");
      const data = JSON.stringify(this.fontSettings);
      localStorage.setItem("fontSettings", data);
    }

    },
    getMiscSettings() {
      this.name = this.$store.state.name
    },
    toggleQRCode() {
      this.$bvModal.hide("consumerSettings")
      EventBus.$emit("toggleQRCode");
    },
    updateName() {
      if(this.name != "") {
        this.$store.commit("setName", this.name)
      }
    },
    increaseTextSize() {
      this.fontSettings.size += 4;
      this.updateSettings();
    },
    decreaseTextSize() {
      this.fontSettings.size -= 4;
      this.updateSettings();
    },
    changeColor() {
      if (this.fontSettings.selected != null) {
        this.fontSettings.selected++;
        if (this.fontSettings.selected > this.form.colors.length - 1) {
          this.fontSettings.selected = 0;
        }
      } else {
        this.fontSettings.selected = 0;
      }
      this.fontSettings.background =
        this.form.colors[this.fontSettings.selected].background;
      this.fontSettings.foreground =
        this.form.colors[this.fontSettings.selected].foreground;
      this.updateSettings();
    },
    changeLineHeight() {
      console.log("change lineheight");
      this.updateSettings();
    },
    changeFontFamily() {
      console.log("change font family");
      this.updateSettings();
    },
    updateSettings() {
      let data = JSON.stringify(this.fontSettings);
      localStorage.setItem("fontSettings", data);
    },
    onOpen() {
      EventBus.$emit("modalOpen", true)
      if(window.screen.orientation.angle == 90) {
        this.orientation = "LANDSCAPE"
      } else {
        this.orientation = "PORTRAIT"
      }
      console.log(this.orientation)
      let stored = JSON.parse(localStorage.getItem("fontSettings"));
      chttps://sttylus.se/visa/#/38228437/spännrambjörntrådonsole.log("on open:", stored);
      this.fontSettings = stored;
    },
    onClose() {
      EventBus.$emit("modalOpen", false)
      this.updateSettings();
      EventBus.$emit("refocus", "");
    },
  },
};
</script>

<style scoped>
</style>
