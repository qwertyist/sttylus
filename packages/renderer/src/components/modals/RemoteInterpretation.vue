<template>
  <b-modal
    :key="id"
    id="remote"
    size="lg"
    hideFooter
    hideHeader
    hideBackdrop
    noFade
    v-show="show"
    return-focus=".ql-editor"
    @show="showModal"
    @keydown.esc="closeModal"
    @hide="closeModal"
  >
    <b-overlay :show="localsession" rounded="sm">
    <template #overlay>
      <div class="text-center">
        <b-icon icon="exclamation-circle-fill" font-scale="3" animation="cylon"></b-icon>
        <p class="lead">Du är redan ansluten till en lokal tolkning</p>
      </div>
    </template>
    <b-overlay :show="qr.show" @click="toggleQRCode" rounded="sm">
    <template #overlay>
      <div class="text-center">
        <qrcode-vue :value="qr.value" :size="400" level="H" />
      </div>
    </template>
    <b-row>
      <b-col>
        <h3>Distanstolkning</h3>
      </b-col>
      <b-col>
        <!--
        <div class="float-right">
          <b-badge v-if="connected === 'online'" variant="success">Du är ansluten till internet</b-badge>
          <b-badge
            v-else-if="connected === 'local'"
            variant="warning"
          >Du är ansluten till ett lokalt nätverk</b-badge>
          <b-badge
            v-else-if="connected === 'offline'"
            variant="danger"
          >Din dator är ej ansluten till ett nätverk/internet</b-badge>
          <span v-else>
            Testar din internetanslutning...
            <b-spinner small label="Loading..." />
          </span>
        </div>
        -->
      </b-col>
    </b-row>
    <b-tabs :tabindex="tabindex">
      <b-tab title="Anslut">
        <b-card>
          <b-card-header>
            <h5>Dina bokade tolkningar</h5>
            <div class="float-right" v-if="inSession">
              <b-badge variant="warning">Du är redan ansluten.</b-badge>
            </div>
          </b-card-header>
          <b-card-body style="max-height: 40vh; overflow-y: auto; width: 100%">
            <b-list-group size="sm"  v-for="sess in sessions" v-bind:key="sess.id">
              <b-list-group-item>
                <b-row>
                  <b-col>
                    <b>
                      {{ sess.name }}
                      <div v-if="sess.ref == you">(Din tolkning)</div>
                    </b>
                    <b-button @click.prevent :id="'details' + sess.id" variant="info">Visa detaljer</b-button>
                    <b-popover :target="'details' + sess.id" placement="left">
                      <template #title>Bokningsdetaljer: {{ sess.id }}</template>
                      <b-row>
                        <b-col cols="4">Från:</b-col>
                        <b-col>{{ sess.from | formatDate }}</b-col>
                      </b-row>
                      <b-row>
                        <b-col cols="4">Till:</b-col>
                        <b-col>{{ sess.to | formatDate }}</b-col>
                      </b-row>
                      <b-row>
                        <b-col cols="4"></b-col>
                        <b-col>
                          <span v-if="sess.recurring">Återkommande</span>
                        </b-col>
                      </b-row>
                      <b-row>
                        <b-col cols="4"></b-col>
                        <b-col></b-col>
                      </b-row>

                      <b-row>
                        <b-col cols="4">Typ:</b-col>
                        <b-col>{{ getSessionType(sess.type) }}</b-col>
                      </b-row>

                      <b-row>
                        <b-col cols="4">Beställare:</b-col>
                        <b-col>{{ sess.ref }}</b-col>
                      </b-row>
                      <hr />Övrig information:
                      <br />
                      <span v-html="sess.description" />
                      <hr />
                      <b-button
                        v-if="!sess.public && !inSession"
                        @click="removeRemoteSession(sess.id)"
                        class="float-right"
                        size="sm"
                        variant="danger"
                      >Ta bort</b-button>
                      <br />
                      <br />
                    </b-popover>
                  </b-col>
                  <b-col>
                    <b-button @click="copyLink(sess.id)" variant="secondary">Kopiera länk</b-button>
                    <b-button @click="toggleQRCode(sess.id)" variant="secondary">Visa QR-kod</b-button>
                    <b-button
                      v-if="!inSession"
                      variant="primary"
                      @click="joinRemoteSession(sess.id)"
                    >Anslut</b-button>
                  </b-col>
                </b-row>
              </b-list-group-item>
            </b-list-group>
          </b-card-body>
        </b-card>
        <b-card title="Anslut till kollegas distanstolkning">
          <b-card-text>
            <div class="float-right" v-if="inSession">
              <b-badge variant="warning">Du är redan ansluten.</b-badge>
            </div>
            <br />
            <b-form inline class="float-right" @submit.prevent>
              <b-form-input
                v-model="session.id"
                placeholder="Ange sessionens sexsiffriga ID-nummer"
                @keydown.enter="joinRemoteSession(session.id)"
                :disabled="inSession"
              />
              <b-button
                v-if="!inSession"
                @click="joinRemoteSession(session.id)"
                variant="primary"
              >Anslut</b-button>
              <b-button v-else @click="closeRemoteSession" variant="danger">Koppla ner</b-button>
            </b-form>
          </b-card-text>
        </b-card>
      </b-tab>
      <b-tab title="Skapa">
        <b-card title="Skapa distanstolkning" subTitle="Låt andra ta del av texten via internet">
          <b-card-text>
            <b-form @submit.prevent="createRemoteSession">
              <div v-show="inSession == false" class="float-right">
                <b-form-select v-model="broadcast" :options="broadcastOptions" />
                <b-button type="submit">Skapa</b-button>
              </div>
            </b-form>
            <div v-show="inSession == true">
              <b-card
                title="Sessionsuppgifter"
                subTitle="Dessa uppgifter delar du med dig om du vill låta användare eller kollegor ansluta till din tolkning"
              >
                <b-card-text>
                  <b-form inline class="float-right">
                    <b-form-input :value="'https://sttylus.se/view2.html?id=' + session.id" />#
                    <h3>{{ session.id }}</h3>
                  </b-form>

                  <b-button @click="closeRemoteSession">Koppla ner från session</b-button>
                </b-card-text>
              </b-card>
            </div>
          </b-card-text>
        </b-card>
        <!--<b-card v-show="inSession == true" title="Lösenordsskydda distanstolkningen">
          <b-card-text>
            <b-form inline class="float-right" @submit.prevent="setSessionPassword">
              <b-form-input placeholder="Välj ett lösenord..." v-model="session.password" />
              <br />
              <b-button @click="setSessionPassword">OK</b-button>
            </b-form>
          </b-card-text>
        </b-card>
        <b-card
          v-show="inSession == true"
          title="Koppla distanstolkning till videomöte"
          subTitle="Skicka texten till en videomötestjänst med stöd för undertexter"
        >
          <b-card-text>
            <b-form inline class="float-right" @submit.prevent="bindAPIToken">
              <b-form-select v-model="thirdpartyservice" :options="thirdpartyOptions" />
              <b-form-input
                :disabled="connected3rdparty"
                placeholder="Mötets API-token"
                v-model="session.token"
              />

              <br />
              <b-button v-if="!breakout" @click="bindAPIToken">Anslut</b-button>
            </b-form>

            <br />
            <br />
            <!--
            <b-form v-if="breakout" inline class="float-right" @submit.prevent="bindAPIToken">
              <b-form-input
                :disabled="true"
                placeholder="Breakout room ej implementerat"
                v-model="session.breakout"
              />
              <br />
              <b-button v-if="!connectedBreakout" @click="bindBreakoutRoom">Byt rum</b-button>
              <b-button v-if="connectedBreakout" @click="leaveBreakoutRoom">Lämna rum</b-button>
            </b-form>
            <div v-if="thirdPartyError && !session.token" class="text-danger float-right">
              <small>Felmeddelande: {{ thirdPartyError }}</small>
            </div>
          </b-card-text>
        </b-card>
        -->
      </b-tab>
    </b-tabs>
  </b-overlay>
  </b-overlay>
  </b-modal>
</template>
<script>
import QrcodeVue from 'qrcode.vue'

import EventBus from "../../eventbus.js";
import axios from "axios";
import eventbus from '../../eventbus.js';

export default {
  components: {
    QrcodeVue,
  },
  data() {
    return {
      localsession: null,
      session: {
        id: "",
        password: "",
        token: "",
        breakout: "",
      },
      qr: {
        show: false,
        value: ""
      },
      id: 0,
      breakout: false,
      broadcast: "remote",
      broadcastOptions: [{ value: "remote", text: "Distanstolkning", disabled: false },{ value: "zoom", text:  "Zoom"}],
      localIP: "127.0.0.1",
      thirdPartyError: "",
      thirdpartyservice: "zoom",
      thirdpartyOptions: [{ value: "zoom", text: "Zoom", disabled: false }],
      password: "",
      ZoomUser: false,
      APIToken: "",
      breakoutAPIToken: "",
      connected3rdparty: false,
      connectedBreakout: false,
      connected: false,
      show: true,
      tabindex: 0,
      sessions: [],
    };
  },
  methods: {
    copyLink(id) {
      var copyText = "https://sttylus.se/view2.html?id=" + id
      navigator.clipboard.writeText(copyText);
    },
    toggleQRCode(id) {
      this.qr.value = "https://sttylus.se/view2.html?id=" + id
      this.qr.show = !this.qr.show;
    },
    getSessionType(type) {
      if (type == 0) {
        return "Distanstolkning";
      } else if (type == 1) {
        return "Zoom CC-tolkning";
      } else if (type == 2) {
        return "Konferenstolkning (GETlivecap)";
      }
    },
    getRecurring(recurring) {
      if (recurring) {
        return "ja";
      } else {
        return "nej";
      }
    },
    getSessions() {
      axios.get("https://sttylus.se/ws/" + "sessions/" + this.$store.state.userData.id).then((resp) => {
        this.sessions = resp.data;
      }).catch(err => {
        this.$toast.warning("Fick inte kontakt med servern", err)
      });
    },
    showModal() {
      if (this.desktop && this.$store.state.local.connected) {
          console.log("Already connected to local session")
          this.localsession = true
      } else {
          console.log("Not connected to local session")
          this.localsession = false
      }

      this.$store.commit("setModalOpen", true)
      this.getSessions()
      EventBus.$emit("modalOpened");
      this.updateSessionInfo();
    },
    closeModal() {
      this.$store.commit("setModalOpen", false)
      EventBus.$emit("modalClosed");
    },
    getLocalIP() {
      api.getLocalIP().then((resp) => {
        this.localIP = resp.data;
      });
    },

    createRemoteSession() {
      switch (this.broadcast) {
        case "remote":
          console.log("Skapa distanstolkning...")
          axios.post("https://sttylus.se/ws/session",
            { ref: this.$store.state.userData.id }
          )
            .then(resp => {
              console.log("create session", resp.data)
              this.session.id = resp.data.id;
              this.getSessions();
              EventBus.$emit("joinRemoteSession", this.session.id);
            }).catch(err => {
              console.log("create session failed:", err)
            })
          EventBus.$emit("createRemoteSession", "online");
          break;
        case "offline":
          EventBus.$emit("createOfflineSession")
      }
      this.$store.commit("setSessionConnected", true);
    },
    closeRemoteSession(end = false) {
      this.$bvModal
        .msgBoxConfirm("Är du säker på att du vill koppla ner?", {
          title: "Bekräfta nedstängning av session",
          size: "md",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Koppla ner",
          cancelTitle: "Återgå till session",
          footerClass: "p-2",
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            this.session.password = "";
            this.session.token = "";
            this.session.breakout = "";
            this.session.id = "";
            this.connected3rdparty = false
            console.log("Nollar API Tokens");
            EventBus.$emit("leaveRemoteSession", end);
            this.$store.commit("setSessionConnected", false);
          } else {
          }
        })
        .catch((err) => {
          // An error occurred
        });
    },
    joinRemoteSession(id) {
      console.log("JoinRemoteSession", id);
      if (id) {
        this.session.id = id;
      }
      if (this.session.id) {
        axios
          .get("https://sttylus.se/ws/" + "session/" + this.session.id)
          .then((resp) => {
            if (resp.status == 204) {
              this.$toast.warning("Ingen distanstolkning med det ID-numret");
              this.session.id = "";
              this.inSession = false;
            } else {
              console.log("join a session")
              EventBus.$emit("joinRemoteSession", this.session.id);
              this.inSession = true;
            }
          })
          .catch((err) => {
            this.$toast.warning("Något gick fel vid anslutningen.");
            console.log("Failed joining session:", err)
          });
        this.updateSessionInfo();
      }
    },
    removeRemoteSession(id) {

      this.$bvModal
        .msgBoxConfirm("Är du säker på att du vill ta bort bokningen?", {
          title: "Bekräfta borttagning av distanstolkning",
          size: "md",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Ta bort",
          cancelTitle: "Avbryt",
          footerClass: "p-2",
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            axios.delete(this.$collabAPI + "session/" + id)
              .then(resp => {
                this.getSessions()
              })
              .catch(err => {
                console.error("couldn't remove session:", err)
              })
          } 
        })
        .catch((err) => {
          // An error occurred
        });
    },
    disconnectRemoteSession() {
      EventBus.$emit("leaveRemoteSession");
    },
    setSessionPassword() {
      axios
        .put("https://sttylus.se/ws/" + "session", {
          id: this.session.id,
          password: this.session.password,
          breakout: this.session.breakout,
          token: this.session.token,
        })
        .then((resp) => {
          this.session.password = resp.data.password;
          this.$store.commit("setSessionPassword", resp.data.password);
          this.$toast.success("Lösenordet uppdaterades");
        })
        .catch((err) => {
          this.session.password=""
          console.log(err);
        });
    },
    bindAPIToken() {
      axios
        .put("https://sttylus.se/ws/session",{
          id: this.session.id,
          token: this.session.token,
          breakout: "",
          password: this.session.password,
        })
        .then((resp) => {
          console.log("Zoom resp:", resp);
          this.connected3rdparty = true;
          this.breakout = true;
          EventBus.$emit("sendSessionData", this.session)
          EventBus.$emit("join3rdParty", this.session.id);
        })
        .catch((err) => {
          this.thirdPartyError = "Kunde inte ansluta";
          this.APIToken = "";
        });
    },
    bindBreakoutRoom() {
      axios
        .post(this.$collabAPI + "/bind", {
          id: this.session.id,
          token: this.session.token,
          breakout: this.session.breakout,
        })
        .then((resp) => {
          console.log("Zoom resp:", resp);
          this.connectedBreakout = true;
        })
        .catch((err) => {
          this.thirdPartyError = "Kunde inte ansluta";
          this.breakoutAPIToken = "";
        });
    },
    leaveBreakoutRoom() {
      axios
        .post(this.$collabAPI + "/bind", { id: this.session.id, token: this.session.token, password: this.session.password, breakout: "" })
        .then((resp) => {
          console.log("Zoom resp: ", resp);
          this.session.token = "";
          this.connectedBreakout = false;
        });
    },
    updateSessionInfo() {
      EventBus.$emit("checkConnection");
      if (this.session.id) {
        console.log("session id:", this.session.id)
        axios
          .post(this.$collabAPI + "/info", { id: this.session.id })
          .then((resp) => {
            if (resp.data.ZoomAPIToken !== "") {
              this.connected3rdparty = true;
              this.breakout = true;

              this.session.token = resp.data.ZoomAPIToken;
              this.session.breakout = resp.data.ZoomBreakoutAPIToken;
              EventBus.$emit("join3rdParty", this.session.id);
            }
            this.session.password = resp.data.password;
          })
          .catch((err) => {
            console.log("/info failed:", err);
          });
      }
    },
  },
  computed: {
    desktop() {
      return this.$mode == "desktop"
    },
    you() {
      return this.$store.getters.getUserId
    },
    inSession: {
      get() {
        return this.$store.state.session.connected;
      },
      set(value) {
        this.$store.commit("setSessionConnected", value);
      },
    },
  },
  mounted() {
    EventBus.$on("networkStatusUpdate", (status) => {
      this.connected = status;
    });
    this.connected3rdparty=false
    if (import.meta.env.VUE_APP_STTYLUS_MODE == "webapp") {
      this.broadcastOptions.splice(1);
    }
    /*
    axios
      .post(this.$collabAPI + "auth", {
        licenseKey: this.$store.state.userData.id,
      })
      .then((resp) => {
        if (resp.data == "OK") {
          this.ZoomUser = true;
        } else {
          this.ZoomUser = false;
        }
      })
      .catch((err) => {
        this.ZoomUser = false;
      });
      */
    EventBus.$on("sessionUpdate", this.updateSessionInfo);
  },
  beforeDestroy() {
    EventBus.$off("networkStatsUpdate", (status) => {
      this.connected = status;
    });
    EventBus.$off("sessionUpdate", this.updateSessionInfo);
  },
};
</script>
<style scoped>
.nav {
  background-color: #fff;
}
</style>
