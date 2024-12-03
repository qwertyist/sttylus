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
          <b-icon
            icon="exclamation-circle-fill"
            font-scale="3"
            animation="cylon"
          ></b-icon>
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
                  <b-badge variant="warning"
                    >Du är ansluten till en distanstolkning.</b-badge
                  >
                </div>
              </b-card-header>
              <b-card-body
                style="max-height: 40vh; overflow-y: auto; width: 100%"
              >
                <b-list-group
                  size="sm"
                  v-for="sess in sessions"
                  v-bind:key="sess.id"
                >
                  <b-list-group-item v-if="!inSession || sess.id == session.id">
                    <b-row>
                      <b-col>
                        <b>
                          {{ sess.name }}
                        </b>
                        <br />
                        <b-button
                          v-if="ownSession(sess) && !sess.public && !inSession"
                          @click="removeRemoteSession(sess.id)"
                          size="sm"
                          variant="danger"
                          >Ta bort</b-button
                        >
                        <b-button
                          v-if="!ownSession(sess) && !sess.public && !inSession"
                          @click="leaveRemoteSession(sess.id)"
                          size="sm"
                          variant="warning"
                          >Lämna/dölj</b-button
                        >
                        <!--
                                                <b-button
                                                    @click.prevent
                                                    :id="'details' + sess.id"
                                                    variant="info"
                                                    >Visa detaljer</b-button
                                                >
                                                <b-popover
                                                    :target="
                                                        'details' + sess.id
                                                    "
                                                    placement="left"
                                                >
                                                    <template #title
                                                        >Bokningsdetaljer:
                                                        {{ sess.id }}</template
                                                    >
                                                    <b-row>
                                                        <b-col cols="4"
                                                            >Från:</b-col
                                                        >
                                                        <b-col>{{
                                                            sess.from
                                                                | formatDate
                                                        }}</b-col>
                                                    </b-row>
                                                    <b-row>
                                                        <b-col cols="4"
                                                            >Till:</b-col
                                                        >
                                                        <b-col>{{
                                                            sess.to | formatDate
                                                        }}</b-col>
                                                    </b-row>
                                                    <b-row>
                                                        <b-col cols="4"></b-col>
                                                        <b-col>
                                                            <span
                                                                v-if="
                                                                    sess.recurring
                                                                "
                                                                >Återkommande</span
                                                            >
                                                        </b-col>
                                                    </b-row>
                                                    <b-row>
                                                        <b-col cols="4"></b-col>
                                                        <b-col></b-col>
                                                    </b-row>

                                                    <b-row>
                                                        <b-col cols="4"
                                                            >Typ:</b-col
                                                        >
                                                        <b-col>{{
                                                            getSessionType(
                                                                sess.type
                                                            )
                                                        }}</b-col>
                                                    </b-row>

                                                    <b-row>
                                                        <b-col cols="4"
                                                            >Beställare:</b-col
                                                        >
                                                        <b-col>{{
                                                            sess.ref
                                                        }}</b-col>
                                                    </b-row>
                                                    <hr />
                                                    Övrig information:
                                                    <br />
                                                    <span
                                                        v-html="
                                                            sess.description
                                                        "
                                                    />
                                                    <hr />
                                                    <br />
                                                    <br />
                                                </b-popover>
                              -->
                      </b-col>
                      <b-col>
                        <b-button
                          @click="toClipboard('user', sess)"
                          variant="secondary"
                          >Kopiera länk</b-button
                        >
                        <b-button
                          @click="toggleQRCode(sess)"
                          variant="secondary"
                          >Visa QR-kod</b-button
                        >
                        <b-button
                          v-if="!inSession"
                          variant="primary"
                          @click="joinRemoteSession(sess.id)"
                          >Anslut</b-button
                        >
                      </b-col>
                    </b-row>
                  </b-list-group-item>
                </b-list-group>
              </b-card-body>
            </b-card>
            <b-card
              title="Anslut till kollegas distanstolkning"
              sub-title="Skriv in bokningens ID-nummer eller klistra in länken till skrivtolkningen."
            >
              <b-card-text>
                <br />
                <b-form inline class="float-right" @submit.prevent>
                  <b-form-input
                    v-model="session.id"
                    autofocus
                    placeholder="Ange sessionens ID-nummer"
                    :disabled="inSession"
                    @keydown.enter="joinRemoteSession(session.id)"
                  />
                  <b-button
                    v-if="!inSession"
                    variant="primary"
                    @click="joinRemoteSession(session.id)"
                    >Anslut</b-button
                  >
                  <b-button v-else variant="danger" @click="closeRemoteSession"
                    >Koppla ner</b-button
                  >
                </b-form>
              </b-card-text>
            </b-card>
          </b-tab>
          <b-tab title="Skapa">
            <b-card
              title="Skapa distanstolkning"
              sub-title="Låt andra ta del av texten via internet"
            >
              <b-card-text>
                <b-form @submit.prevent="createRemoteSession">
                  <div v-show="inSession == false">
                    <b-button type="submit">Skapa</b-button>
                  </div>
                </b-form>
                <div v-show="inSession == true">
                  <b-card
                    title="Sessionsuppgifter"
                    sub-title="Dessa uppgifter delar du med dig om du vill låta användare eller kollegor ansluta till din tolkning"
                  >
                    <b-card-text>
                      <b-form>
                        <b-row>
                          <b-col cols="3"> Tolkanvändarvy </b-col>
                          <b-col>
                            <b-form-input
                              readonly
                              :value="
                                'https://sttylus.se/visa/#/' +
                                session.id +
                                '/' +
                                session.password
                              "
                              @click="toClipboard('view', session)"
                            />
                          </b-col>
                        </b-row>

                        <b-row>
                          <b-col cols="3"> Sessions-ID </b-col>
                          <b-col>
                            <b-form-input
                              readonly
                              :value="session.id"
                              @click="toClipboard('id', session)"
                            />
                          </b-col>
                        </b-row>
                        <b-row>
                          <b-col cols="3">Lösenord</b-col>
                          <b-col>
                            <b-form inline @submit.prevent="setSessionPassword">
                              <b-form-input
                                placeholder="Välj ett lösenord..."
                                v-model="session.password"
                              />
                              <b-button
                                v-if="passwordChanged"
                                @click="setSessionPassword"
                                >Ändra</b-button
                              >
                            </b-form>
                          </b-col>
                        </b-row>
                        <!--<b-form-input
                                                    :value="
                                                            'http://localhost:3345/#/' + session.id + '/' + session.password
                                                    "
                                                />#
                                                -->
                      </b-form>
                    </b-card-text>
                  </b-card>
                  <br />
                  <b-button @click="closeRemoteSession"
                    >Koppla ner från session</b-button
                  >
                </div>
              </b-card-text>
            </b-card>
            <b-card
              v-show="inSession == true"
              title="Skrivtolka till andra plattformar"
              subTitle="Välj att få skrivtolkningen som closed captions till t.ex. videomöten eller livestreams."
            >
              <b-card-text>
                <b-form @submit.prevent="bindAPIToken">
                  <b-row>
                    <b-col cols="3">
                      <b-form-select
                        v-model="thirdpartyservice"
                        :options="thirdpartyOptions"
                      />
                    </b-col>
                    <b-col cols="6" v-if="thirdpartyservice == 'captions'">
                      <b-form-input
                        :value="'https://sttylus.se/text/' + session.id"
                        @click="toClipboard('captions', session)"
                      />
                      <b-form-input
                        :value="'https://sttylus.se/GETlivecap/' + session.id"
                        @click="toClipboard('glc', session)"
                      />
                    </b-col>
                    <b-col cols="6" v-if="thirdpartyservice == 'zoom'">
                      <b-form-input
                        :disabled="connected3rdparty"
                        placeholder="Mötets API-token"
                        v-model="session.token"
                      />
                    </b-col>
                    <b-col v-if="thirdpartyservice == 'zoom'">
                      <b-button @click="bindAPIToken">Anslut</b-button>
                    </b-col>

                    <b-form
                      v-if="breakout"
                      inline
                      class="float-right"
                      @submit.prevent="bindAPIToken"
                    >
                      <b-form-input
                        :disabled="true"
                        placeholder="Breakout room ej implementerat"
                        v-model="session.breakout"
                      />
                      <br />
                      <b-button
                        v-if="!connectedBreakout"
                        @click="bindBreakoutRoom"
                        >Byt rum</b-button
                      >
                      <b-button
                        v-if="connectedBreakout"
                        @click="leaveBreakoutRoom"
                        >Lämna rum</b-button
                      >
                    </b-form>
                    <div
                      v-if="thirdPartyError && !session.token"
                      class="text-danger float-right"
                    >
                      <small>Felmeddelande: {{ thirdPartyError }}</small>
                    </div>
                  </b-row>
                </b-form>
              </b-card-text>
            </b-card>
          </b-tab>
        </b-tabs>
      </b-overlay>
    </b-overlay>
  </b-modal>
</template>
<script>
import QrcodeVue from 'qrcode.vue'

import EventBus from '../../eventbus.js'
import axios from 'axios'
import eventbus from '../../eventbus.js'

import { tools } from '../../../assets/passwordlist.json'

export default {
  components: {
    QrcodeVue,
  },
  data() {
    return {
      localsession: null,
      session: {
        id: '',
        password: '',
        token: '',
        breakout: '',
      },
      qr: {
        show: false,
        value: '',
      },
      id: 0,
      breakout: false,
      broadcast: 'remote',
      broadcastOptions: [
        { value: 'remote', text: 'Distanstolkning', disabled: false },
        { value: 'zoom', text: 'Zoom' },
      ],
      localIP: '127.0.0.1',
      thirdPartyError: '',
      thirdpartyservice: null,
      thirdpartyOptions: [
        { value: null, text: 'Ej valt...', disabled: true },
        { value: 'zoom', text: 'Zoom', disabled: false },
       // { value: 'captions', text: 'Livetextning', disabled: true },
      ],
      password: '',
      ZoomUser: false,
      APIToken: '',
      breakoutAPIToken: '',
      connected3rdparty: false,
      connectedBreakout: false,
      connected: false,
      show: true,
      tabindex: 0,
      sessions: [],
    }
  },
  methods: {
    addEventListeners() {
      EventBus.$on('passwordMessage', (msg) => {
        if (msg == undefined) {
          return
        }
        if (msg instanceof Object) {
          return
        }
        if (msg == 'ok') {
          this.$toast.success('Lösenordet byttes')
          this.password = this.session.password
        } else if (msg == 'whitespace') {
          this.$toast.warning('Lösenordet får inte innehålla blanksteg')
          this.session.password = this.password
        } else {
          this.session.password = msg
          this.password = msg
        }
      })
      EventBus.$on('networkStatusUpdate', (status) => {
        this.connected = status
      })
      EventBus.$on('zoomConnected', (success) => {
        if (success) {
          this.$toast.success('Anslöt till Zoom')
          this.connected3rdparty = true
        } else {
          this.$toast.warning(
            'Kunde inte ansluta till Zoom. Kontrollera API-token.'
          )
          this.connected3rdparty = false
        }
      })
      EventBus.$on('sessionUpdate', this.updateSessionInfo)
    },
    removeEventListeners() {
      EventBus.$off('sessionPasswordUpdated')
      EventBus.$off('passwordMessage')
      EventBus.$off('networkStatsUpdate', (status) => {
        this.connected = status
      })
      EventBus.$off('sessionUpdate', this.updateSessionInfo)
    },
    toClipboard(target, sess) {
      let copyText
      let successText

      if (Object.keys(sess).length === 0) {
        this.$toast.error('Fel: Inget att kopiera')
        navigator.clipboard.writeText('')
        return
      }

      if (sess.password == undefined) {
        sess.password = ''
      }
      if (target == 'captions') {
        copyText = 'https://sttylus.se/text/' + sess.id
        successText = 'Livetextningens webbadress kopierades'
      } else if (target == 'glc') {
        copyText = 'https://sttylus.se/GETlivecap/' + sess.id
        successText = 'Livetextningens GETlivecap-feed kopierades'
      } else if (target == 'id') {
        copyText = sess.id
        successText = 'Distanstolkningens sessions-ID kopierades'
      } else {
        // "target" == user
        copyText = 'https://sttylus.se/visa/#/' + sess.id + '/' + sess.password
        successText = 'Distanstolkningens webbadress kopierades'
      }
      navigator.clipboard.writeText(copyText)
      this.$toast.success(successText)
    },
    toggleQRCode(sess) {
      this.qr.value =
        'https://sttylus.se/visa/#/' + sess.id + '/' + sess.password
      this.qr.show = !this.qr.show
    },
    getSessionType(type) {
      if (type == 0) {
        return 'Distanstolkning'
      } else if (type == 1) {
        return 'Zoom CC-tolkning'
      } else if (type == 2) {
        return 'Konferenstolkning (GETlivecap)'
      }
    },
    ownSession(sess) {
      const userID = this.$store.state.userData.id
      if (sess.ref == userID) {
        return true
      }
      return false
    },
    getRecurring(recurring) {
      if (recurring) {
        return 'ja'
      } else {
        return 'nej'
      }
    },
    getSessions() {
      axios
        .get(this.$collabAPI + 'sessions/' + this.$store.state.userData.id)
        .then((resp) => {
          this.sessions = resp.data
        })
        .catch((err) => {
          this.$toast.warning('Fick inte kontakt med servern', err)
        })
    },
    showModal() {
      if (this.desktop && this.$store.state.local.connected) {
        console.log('Already connected to local session')
        this.localsession = true
      } else {
        console.log('Not connected to local session')
        this.localsession = false
      }

      this.$store.commit('setModalOpen', true)
      this.getSessions()
      EventBus.$emit('modalOpened')
      this.updateSessionInfo()
    },
    closeModal() {
      this.$store.commit('setModalOpen', false)
      EventBus.$emit('modalClosed')
    },
    getLocalIP() {
      api.getLocalIP().then((resp) => {
        this.localIP = resp.data
      })
    },

    createRemoteSession() {
      switch (this.broadcast) {
        case 'remote':
          this.generatePassword()
          axios
            .post(this.$collabAPI + 'session', {
              ref: this.$store.state.userData.id,
              password: this.session.password,
            })
            .then((resp) => {
              this.session.id = resp.data.id
              this.getSessions()
              EventBus.$emit(
                'joinRemoteSession',
                this.session.id,
                this.session.password
              )
            })
            .catch((err) => {
              console.log('create session failed:', err)
            })
          EventBus.$emit('createSession', 'online')
          break
        case 'offline':
          EventBus.$emit('createOfflineSession')
      }
      this.$store.commit('setSessionConnected', true)
    },
    closeRemoteSession(end = false) {
      this.$bvModal
        .msgBoxConfirm('Är du säker på att du vill koppla ner?', {
          title: 'Bekräfta nedstängning av session',
          size: 'md',
          buttonSize: 'sm',
          okVariant: 'danger',
          okTitle: 'Koppla ner',
          cancelTitle: 'Återgå till session',
          footerClass: 'p-2',
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            this.session.password = ''
            this.session.token = ''
            this.session.breakout = ''
            this.session.id = ''
            this.connected3rdparty = false
            EventBus.$emit('leaveRemoteSession', end)
            this.$store.commit('setSessionConnected', false)
          } else {
          }
        })
        .catch((err) => {
          // An error occurred
        })
    },
    joinRemoteSession() {
      if (this.session.id) {
        let tmp = this.session.id.match(/([0-9]{8})/g)
        if (tmp != "") {
          this.session.id = tmp[0]
        }
        axios
          .get(this.$collabAPI + 'session/' + this.session.id)
          .then((resp) => {
            if (resp.status == 204) {
              this.$toast.warning(
                'Ingen distanstolkning med det ID-numret. Anslut igen för att skapa.'
              )
              this.session.id = ''
              this.inSession = false
            } else {
              axios
                .get(
                  this.$collabAPI +
                    'join/' +
                    this.$store.state.userData.id +
                    '/' +
                    this.session.id
                )
                .then(() => {
                  this.getSessions()
                })
                .catch((err) => {
                  this.$toast.warning(
                    'Kunde inte knyta din användare till bokningen.',
                    err
                  )
                  console.error('join/', err)
                })
              EventBus.$emit('joinRemoteSession', this.session.id)
              this.inSession = true
            }
          })
          .catch((err) => {
            this.$toast.warning('Något gick fel vid anslutningen.')
            console.log('Failed joining session:', err)
          })
        this.updateSessionInfo()
      }
    },
    removeRemoteSession(id) {
      this.$bvModal
        .msgBoxConfirm('Är du säker på att du vill ta bort bokningen?', {
          title: 'Bekräfta borttagning av distanstolkning',
          size: 'md',
          buttonSize: 'sm',
          okVariant: 'danger',
          okTitle: 'Ta bort',
          cancelTitle: 'Avbryt',
          footerClass: 'p-2',
          hideHeaderClose: true,
          centered: true,
        })
        .then((value) => {
          if (value) {
            axios
              .delete(this.$collabAPI + 'session/' + id)
              .then((resp) => {
                this.getSessions()
              })
              .catch((err) => {
                console.error("couldn't remove session:", err)
              })
          }
        })
        .catch((err) => {
          // An error occurred
        })
    },
    leaveRemoteSession(id) {
      axios
        .get(
          this.$collabAPI + 'leave/' + this.$store.state.userData.id + '/' + id
        )
        .then(() => {
          this.getSessions()
        })
        .catch((err) => {
          this.$toast.warning(
            'Kunde inte koppla användare från bokning. Kontakta administratören',
            err
          )
        })
    },
    disconnectRemoteSession() {
      this.thirdpartyservice = null
      EventBus.$emit('leaveRemoteSession')
    },
    setSessionPassword() {
      console.log('submitted password')
      this.$toast.info('Byter lösenord...')
      EventBus.$emit('setSessionPassword', this.session.password)
    },
    bindAPIToken() {
      if (this.session.token.length < 35) {
        this.$toast.warning('Inte en giltig API-token')
        return
      }
      EventBus.$emit('sendSessionData', { token: this.session.token })
    },
    bindBreakoutRoom() {},
    leaveBreakoutRoom() {},
    updateSessionInfo() {
      EventBus.$emit('checkConnection')
      if (this.session.id) {
        console.log('session id:', this.session.id)
      }
    },
    generatePassword() {
      const password = tools
        .sort(() => 0.5 - Math.random())
        .slice(0, 2)
        .join('')
      this.password = password
      this.session.password = password
    },
  },
  computed: {
    desktop() {
      return this.$mode == 'desktop'
    },
    you() {
      return this.$store.getters.getUserId
    },
    inSession: {
      get() {
        return this.$store.state.session.connected
      },
      set(value) {
        this.$store.commit('setSessionConnected', value)
      },
    },
    passwordChanged() {
      return this.password != this.session.password
    },
    passwordMessage(msg) {},
  },
  mounted() {
    this.addEventListeners()
    this.connected3rdparty = false
    if (import.meta.env.VUE_APP_STTYLUS_MODE == 'webapp') {
      this.broadcastOptions.splice(1)
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
  },
  beforeDestroy() {
    this.removeEventListeners()
  },
}
</script>
<style scoped>
.nav {
  background-color: #fff;
}
</style>
