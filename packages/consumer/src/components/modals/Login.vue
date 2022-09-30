<template>
  <div>
    <b-modal
      ok-title="Anslut"
      ok-only
      size="lg"
      :title="'Koppla upp till distanstolkning: ' + id"
      id="login"
      @show="getSessionInfo"
      @hide="hide"
    >
      Skriv in ditt namn:
      <b-form @submit.prevent="hide">
        <b-input v-model="name" :placeholder="'Tolkanvändare'" />

        <div v-if="passwordProtected">
          <br />
          Bokningen är lösenordsskyddad:
          <b-input v-model="password" placeholder="Bokningens lösenord" />
          <b-badge variant="danger" v-if="failed"
            >Inloggningen misslyckades</b-badge
          >
        </div>
        <br />
        Genom att ansluta så godkänner jag lagring av cookies för att komma ihåg
        mitt användarnamn och mina textinställningar.<br />
        Textinställningarna nås genom att dubbelklicka i rutan.
      </b-form>
    </b-modal>
  </div>
</template>

<script>
import axios from "axios";
import EventBus from "../../eventbus";
export default {
  data() {
    return {
      passwordProtected: false,
      id: "",
      name: "",
      remember: true,
      password: "",
      failed: false,
    };
  },
  computed: {
    userCount() {
      return this.$store.state.session.users;
    },
  },
  mounted() {
    this.name = this.$cookie.get("name");
  },
  methods: {
    getSessionInfo() {
      this.id = this.$store.state.session.id;
      this.passwordProtected = this.$store.state.session.protected;
    },
    hide() {
      console.log("hiding");
      axios
        .post(this.$collabAPI + "/login", {
          id: this.$store.state.session.id,
          password: this.password,
        })
        .then((resp) => {
          if (resp.data == "OK") {
            if (this.name == "") {
              this.name = "Namnlös tolkanvändare";
            }
            this.$store.commit("setName", {
              name: this.name,
              remember: this.remember,
            });
            EventBus.$emit("join");
            this.$bvModal.hide("login");
            this.failed = false;
          }
        })
        .catch((err) => {
          console.log("Failed somehow?", err);
          this.failed = true;
          this.$bvModal.show("login");
        });
    },
  },
};
</script>