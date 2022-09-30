<template>
  <b-container
    fluid
    class="no-gutters navMargin"
  >
    <b-row class="d-block">
      <b-tabs
        v-model="tabIndex"
        class="navbar-dark"
        content-class=""
      >
        <b-tab
          title="Förkortningar"
          :title-item-class="'tab-title-class'"
          active
          href="#abbs"
          class="m-3"
        >
          <Abbreviations
            @commend="increment()"
            @decommend="decrement()"
          />
        </b-tab>
       <b-tab
          lazy
          title="Manuskript"
          :title-item-class="'tab-title-class'"
          href="#manuscript"
          class="manuscriptTab"
        >
          <Manuscript />
        </b-tab>
        <b-tab
          lazy
          title="Textinställningar"
          :title-items-class="'tab-title-class'"
          href="#text"
          class="m-3"
        >
          <Font />
        </b-tab>
        <div v-if="!tester">
          <b-tab
            lazy
            title="Importera/Exportera"
            :title-items-class="'tab-title-class'"
            href="#text"
            class="m-3"
          >
            <ImportExport />
          </b-tab>
        </div>
        <b-tab
          title="Snabbkommandon"
          :title-items-class="'tab-title-class'"
          href="#hotkey"
          class="m-3"
        >
          <Hotkeys />
        </b-tab>
        <b-tab
          lazy
          v-if="admin"
          title="Admin-vy"
          :title-item-class="'tab-title-class ml-auto'"
          href="#admin"
          class="m-3"
        >
          <Admin />
        </b-tab>
        <b-tab
          :title-item-class="'tab-title-class ml-auto'"
          href="#changelog"
          class="informationTab"
        >
          <template #title>
            Programinformation
            <b-badge v-if="changes > 0">
              {{ changes }}
            </b-badge>
          </template>
          <Changelog />
        </b-tab>
        <b-tab
          v-if="!tester"
          title="Återställ inställningar"
          :title-item-class="'tab-title-class'"
          @click="resetSettings"
        />

        <b-tab
          v-if="webapp && !tester"
          title="Logga ut"
          :title-item-class="'tab-title-class'"
          @click="logout"
        />
        <b-tab
          v-if="webapp && tester"
          title="Logga ut"
          :title-item-class="'tab-title-class ml-auto'"
          @click="logout"
        />

        <b-tab
          v-if="!webapp"
          title="Logga ut"
          :title-item-class="'tab-title-class'"
          @click="logout"
        />

        <!--<b-tab title="Användaruppgifter" :titleItemsClass="'tab-title-class'" href="#user">
          <h2>Användaruppgifter</h2>
          <User />
        </b-tab>
        <b-tab title="Licens" href="#license">
          <License />
        </b-tab>
        -->
      </b-tabs>
    </b-row>
  </b-container>
</template>

<script>
import Abbreviations from "./settings/Abbreviations.vue";
import Manuscript from "./settings/Manuscript.vue";
import Admin from "./settings/Admin.vue";
import Changelog from "./settings/Changelog.vue";
import Font from "./settings/Font.vue";
import Hotkeys from "./settings/Hotkeys.vue";
import ImportExport from "./settings/ImportExport.vue";
import AddAbbreviation from "./modals/AddAbbreviation.vue";
import EventBus from "../eventbus.js";

export default {
  name: "Settings",
  components: {
    Abbreviations,
    Manuscript,
    Font,
    Hotkeys,
    ImportExport,
    Admin,
    Changelog,
    AddAbbreviation,
  },
  data() {
    return {
      tabIndex: 0,
      tabs: [
        "#abbs",
        "#manuscripts",
        "#text",
        "#import",
        "#admin",
        "#changelog",
      ],
      changes: 0,
      nominationCounter: 0,
    };
  },
  computed: {
    webapp() {
      return this.$store.state.mode == "webapp";
    },
    tester() {
      if (this.$store.state.userData.role == "tester") {
        return true;
      }
      return false;
    },
    admin() {
      if (this.$store.state.userData.role == "admin") {
        return true;
      }
      return false;
    },
  },
  methods: {
    addEventListeners() {
      EventBus.$on("newChangeLog", this.alertNewChanges);
      EventBus.$on("openSettings", () => { this.tabIndex = 0; });
      EventBus.$on("importSharedAbbs", () => {
        this.tabIndex = 3;
      });
      EventBus.$on("unsavedManuscriptChanges", this.unsavedManuscriptChanges);
      EventBus.$on("createdList", () => {
        this.tabIndex = 0;
      });
    },
    removeEventListeners() {
      EventBus.$off("importSharedAbbs");
      EventBus.$off("newChangeLog");
      EventBus.$off("openSettings");
      EventBus.$off("unsavedManuscriptChanges");
      EventBus.$off("createdList");
    },
    increment() {
      this.nominationCounter += 1;
    },
    decrement() {
      this.nominationCounter -= 1;
    },
    reset() {
      this.nominationCounter = 0;
    },
    unsavedManuscriptChanges(changes) {
      if (changes) {
        this.$toast.warning("Manuskriptet du redigerar har osparade ändringar");
      }
    },
    alertNewChanges(changes) {
      this.changes = changes;

      /*
      this.$toast.info(
        "Nya uppdateringar finns. Läs mer i uppdateringsloggen."
      );
      */
    },
    resetSettings(e) {
      const messageNodes = this.$createElement("div", {
        domProps: {
          innerHTML:
            "Om programmet ser konstigt ut, så testa att återställa webbläsarens zoom-inställning med <kbd>CTRL+0</kbd>.<br/>Programmet fungerar bäst i fullskärm, som du aktiverar med <kbd>F11</kbd>.",
        },
      });
      this.$bvModal
        .msgBoxOk([messageNodes], {
          title: "Återställ inställningar",
          size: "lg",
          buttonSize: "sm",
          okTitle: "OK",
          cancelTitle: "Återgå",
          footerClass: "p-2",
          hideHeaderClose: true,
          centered: true,
        })
        .then(() => {
          this.$store.commit("reset");
          this.$nextTick(() => {
            this.tabIndex = 0;
          });
        });
    },
    logout(e) {
      this.$bvModal
        .msgBoxConfirm("Är du säker på att du vill logga ut?", {
          title: "Bekräfta utloggning",
          size: "lg",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Logga ut",
          cancelTitle: "Återgå",
          footerClass: "p-2",
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            //  api.setIdToken("");
            this.$store.commit("logout");
          } else {
            this.tabIndex = 0;
          }
        })
        .catch((err) => {
          // An error occurred
        });
    },
    close(e) {
      e.preventDefault();
      this.$bvModal
        .msgBoxConfirm("Är du säker på att du vill avsluta programmet?", {
          title: "Bekräfta avsluta STTylus",
          size: "lg",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Avsluta",
          cancelTitle: "Återgå",
          footerClass: "p-2",
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            var win = nw.Window.get();
            this.$store.commit("beforeClose");
            win.close();
          }
          this.tabIndex = 0;
        })
        .catch((err) => {
          // An error occurred
        });
    },
  },
  watch: {},
  mounted() {
    this.addEventListeners();
    //this.tabIndex = this.tabs.findIndex((tab) => tab === this.$route.hash);
    this.nominationCounter = this.$store.state.stagedNominations.length;
    this.$nextTick(() => {
      this.tabIndex = this.tabs.findIndex((tab) => tab === this.$route.hash);
    });
  },
  beforeDestroy() {
    this.removeEventListeners();
  },
};
</script>

<style>
.nav-tabs .nav-link {
  color: rgba(255, 255, 255, 0.5);
}

.nav-tabs .nav-link:hover {
  color: rgba(255, 255, 255, 0.75);
}

.nav.tabs .nav-link.active {
  color: rgba(73, 80, 87, 0.75);
}

.nav-tabs .nav-link.active:hover {
  color: rgba(73, 80, 87, 0.35);
}

.tab-content > active {
  max-height: 85vh;
  overflow-y: auto;
  overflow-x: hidden;
  margin: 0px !important;
}
.manuscriptTab {
  height: 75vh !important;
  margin: 0px !important;
}

.informationTab {
  height: 75vh !important;
  margin: 0px !important;
}
</style>
