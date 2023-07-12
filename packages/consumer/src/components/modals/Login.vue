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
      <b-form @submit.prevent="hide">
        Tolkningen är lösenordsskyddad: <br />
          <b-input v-model="password" placeholder="Bokningens lösenord" />
          <b-badge variant="danger" v-if="failed"
            >Inloggningen misslyckades</b-badge
          >
        <br />
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
      id: "",
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
    },
    hide() {
      console.log("hiding");
      this.$router.push(
        { path: "/" + this.id + "/" + this.password,
        }
      ).then(() => {
        this.$router.go()
      })
    },
  },
};
</script>
