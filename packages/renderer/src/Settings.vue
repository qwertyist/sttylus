<template>
  <div id="app">
    <UserSettings />
  </div>
</template>

<script>
import Navigation from "./components/Navigation.vue";
import UserSettings from "./components/UserSettings.vue";

export default {
  name: "Settings",
  components: {
    Navigation,
    UserSettings,
  },
  data() { 
    return {
      toast: false
    }
  },
  mounted() {
    window.addEventListener("keydown", this.hotKeys);
  },
  beforeDestroy() {
    //this.$store.commit("storeSettings")
    window.removeEventListener("keydown", this.hotKeys);
  },
  methods: {
    hotKeys(e) {
      if (e.key == "F1") {
        e.preventDefault();
        this.$bvModal.show("support");
      }
      if (e.key == "F2") {
        e.preventDefault();
        let currentLists = this.$store.state.settings.selectedLists;
        let targetList = this.$store.state.targetList.id;
        if (
          [currentLists.standard]
            .concat(currentLists.addon)
            .indexOf(targetList) == -1
        ) {
          if (!this.toast) {
            this.toast = true
            this.$toast.warning(
              "Du kan bara lägga till förkortningar i listor du valt.", {
                duration: 3000,
                onDismiss: () => { this.toast = false },
            });
          }
        } else {
          this.$bvModal.show("addAbb");
        }
      }
      if (e.key == "F3") {
        e.preventDefault();
      }
     /* if (e.key == "F5") {
        e.preventDefault();
        this.$router.push("/");
      }*/
      if (e.key == "F6") {
        e.preventDefault();
      }
      if (e.key == "F7") {
        e.preventDefault();
      }
      if (e.key == "F8") {
        e.preventDefault();
      }
      if (e.key == "F12") {
        e.preventDefault();
      }
    },
  },
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
