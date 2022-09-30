<template>
  <div>
    <Navigation :view="'login'" />
    <b-overlay :show="loading">
      <b-jumbotron>
        <div v-if="local && users">
          <h1>Välj användare</h1>
          <b-list-group v-for="user in users" v-bind:key="user.id"><b-list-group-item @dblclick="loginAs(user)">{{ user.name }}</b-button></b-list-group-item></b-list-group>
        </div>
        <hr />
        <div v-if="local"><h1>Logga in och ladda ner användare lokalt</h1></div>
        <h1 v-else>Logga in</h1>
        <div v-if="step == 0">
          <b-form @submit.prevent="next">
            <b-form-group label="E-post">
              <b-form-input type="email" v-model="form.email" />
            </b-form-group>

            <b-button variant="info" type="submit">Nästa</b-button>
            <span v-if="error" class="text-danger"> {{ error }}</span>
          </b-form>
        </div>
        <div v-if="step == 1">
          <div
            v-if="passwordprompt == 'sync'"
          >STTylus kommer att synka dina lokala filer med användardata som finns på webben.</div>
          <div v-if="passwordprompt == 'login' || passwordprompt == 'sync'">
            <ValidationObserver v-slot="{ handleSubmit }">
              <b-form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider rules="required" v-slot="{ errors, valid }">
                  <b-form-group label="Lösenord">
                    <b-form-input type="password" :state="valid" v-model="form.password" autofocus />
                    <span>{{ errors[0] }}</span>
                  </b-form-group>
                </ValidationProvider>

                <b-button variant="info" type="submit">Logga in</b-button>
              </b-form>
            </ValidationObserver>
            {{ syncResponse }}
          </div>
          <div v-if="passwordprompt == 'register'">
            I och med att detta är din första inloggning så väljer du ditt
            lösenord nedan.
          </div>
          <div v-if="passwordprompt == 'reset'">
            Lösenordet för ditt konto har nollställts. Skapa ett nytt lösenord
            nedan.
          </div>
          <div v-if="passwordprompt == 'register' || passwordprompt == 'reset'">
            <ValidationObserver v-slot="{ handleSubmit }">
              <b-form @submit.prevent="handleSubmit(onSubmit)">
                <ValidationProvider
                  rules="required|confirmed:confirm"
                  vid="password"
                  v-slot="{ errors, valid }"
                >
                  <b-form-group label="Lösenord">
                    <b-form-input type="password" :state="valid" v-model="form.password" autofocus />
                    <span>{{ errors[0] }}</span>
                  </b-form-group>
                </ValidationProvider>

                <ValidationProvider
                  vid="confirm"
                  rules="required|confirmed:password"
                  v-slot="{ errors, valid }"
                >
                  <b-form-group label="Bekräfta lösenord">
                    <b-form-input :state="valid" type="password" v-model="form.confirm" />
                    <span>{{ errors[0] }}</span>
                  </b-form-group>
                </ValidationProvider>
                <span v-if="passwordprompt == 'register'">
                  <b-button variant="info" type="submit">Registrera</b-button>
                </span>
                <span v-if="passwordprompt == 'reset'">
                  <b-button variant="info" type="submit">Skapa nytt lösenord</b-button>
                </span>
              </b-form>
            </ValidationObserver>
          </div>
        </div>
      </b-jumbotron>
      <template #overlay>
        <div class="text-center">
          <div v-if="passwordprompt == 'register'">
            <b-spinner label="Skapar användare..." />
            <p>Vänta medan din användare skapas...</p>
          </div>
          <div v-else>
            <b-spinner label="Synkroniserar användare..." />
            <p>Vänta medan din användare laddas ner...</p>
          </div>
        </div>
      </template>
    </b-overlay>
    <!--<b-jumbotron>
      <h2>Testa STTylus</h2>
      Om du inte har ett registrerat konto så kan du köra ett begränsat
      testkonto.<br />
      Är du intresserad av att registrera ett konto så kontaktar du
      info@sttylus.se.<br />
      <b-button variant="info" type="submit" @click.prevent="onLoginTester"
        >Logga in med testkonto</b-button
      >
    </b-jumbotron>-->
  </div>
</template>
<script>
import api from "./api/api.js";
import Navigation from "./components/Navigation.vue";

export default {
  name: "Login",
  components: {
    Navigation,
  },
  data() {
    return {
      users: [],
      step: 0,
      firstLogin: false,
      passwordprompt: "",
      loading: false,
      form: {
        email: "",
        password: "",
        confirm: "",
        last_sync: new Date(0),
      },
      syncResponse: "",
      error: "",
    };
  },
  computed: {
    local() {
      console.log(this.$mode)
      return this.$mode == "desktop" ? true : false
    }
  },
  mounted() {
    if (this.local) {
      console.log("Local login")
      api.getUsers().then(resp => {
        console.log(resp.data)
        this.users = resp.data
      })
      this.form.last_sync = this.$store.state.lastSync || new Date(0)
    }
  },
  methods: {
    next() {
      api
        .isRegistered(this.form.email)
        .then((resp) => {
          console.log(
            "check if " + this.form.email + " is registered:",
            resp.data
          );
          if (resp.data == "login") {
            this.passwordprompt = "login";
          } else if (resp.data == "register") {
            this.passwordprompt = "register";
          } else if (resp.data == "reset") {
            this.passwordprompt = "reset";
          } else if (resp.data == "sync") {
            this.passwordprompt = "sync"
          } else {
            this.$toast.error(
              "Något gick fel vid inloggningen, kontakta supporten"
            );
            return;
          }
          this.$nextTick(() => {
            this.step++;
          });
        })
        .catch((err) => {
        if (err.response) {
          if (err.response.status == 502) {
            if (this.users) {
              this.$bvModal.msgBoxOk("Du är inte ansluten till internet. Logga in på en befintlig lokal användare.", {
                title: "Internetanslutning saknas"
              })
              this.error = "Du är inte ansluten till internet";
              return
            }
            this.$bvModal.msgBoxOk("Du är inte ansluten till internet och har heller ingen användare lokalt. Kontakta supporten för att få hjälp.", 
            {title: "Användarkonto och internetanslutning saknas"})
          }
          this.error = "Inget användarkonto med den e-postadressen";
        }});
    },
    onSubmit() {
      switch (this.passwordprompt) {
        case "login":
          this.onLogin();
          break;
        case "sync":
          this.onSync();
          break;
        case "register":
          this.onNewPassword();
          break;
        case "reset":
          this.onNewPassword();
          break;
        default:
          this.$toast.warning("Fel vid inloggning");
      }
    },
    onLogin() {
      const data = JSON.stringify(this.form);
      this.$store
        .dispatch("AUTH_REQUEST", {
          email: this.form.email,
          password: this.form.password,
        })
        .then((resp) => {
          setTimeout(() => {
            this.$router.push("/");
          }, 250);
        })
        .catch((err) => {
          this.validPassword = false;
        });
    },
    onNewPassword() {
      console.log(this.errors);
      if (this.form.confirm != this.form.password) {
        return;
      }
      const data = JSON.stringify(this.form);
      if (this.passwordprompt == "register") {
        api.register(data).then((resp) => {
          //console.log("First time registration for user");
          this.$store
            .dispatch("AUTH_REQUEST", {
              email: this.form.email,
              password: this.form.password,
            })
            .then(() => {
              this.loading = true;
              api
                .createStandard(this.$store.state.userData.id)
                .then((resp) => {
                  this.$store.commit("setSelectedStandard", resp.data);
                  this.loading = false;
                  this.$router.push("/");
                })
                .catch((err) => {
                  console.log("create standard failed", err);
                });
            })
            .catch((err) => {
              this.validPassword = false;
            });
        });
      } else if (this.passwordprompt == "reset") {
        console.log("Proceed to reset password");
        api.resetPassword(this.form.email, this.form.password).then((resp) => {
          console.log(resp.data);
          this.$store
            .dispatch("AUTH_REQUEST", {
              email: this.form.email,
              password: this.form.password,
            })
            .then(() => {
              this.$router.push("/");
            });
        });
      }
    },
    loginAs(user) {
      console.log("try to log in as:", user)
      this.form = {}
      this.form.id = user.id
      this.form.email = user.email
      this.onSync()
    },
    onSync() {
      this.loading = true;
      console.log("email:", this.form.email)
      if(this.form.email != "info@sttylus.se") {
        api.syncUser(this.form).then((resp) => {
          console.log("sync user response:", resp.data);
          if (resp.status == 204) {
            console.log("No new content")
            resp.data = {
              email: this.form.email,
              id: this.form.id,
            }
          } else {
            this.$toast.info("Synkroniserar listor")
          }
          this.$store.commit("setLastSync", resp.data.last_sync)
        })
      }
      this.$store
        .dispatch("AUTH_REQUEST", {
          email: this.form.email,
          id: this.form.id,
          local: true,
      })
      .then((resp) => {
        setTimeout(() => {
          this.$router.push("/");
        }, 250);
      })
      .catch((err) => {
        console.log("api syncUser failed:", err)
        this.validPassword = false;
        this.$toast.error("Kunde inte synkronisera listor. Kontakta admin", err)
      });
    }
  },
};
</script>
<style>
</style>
