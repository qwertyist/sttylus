<template>
    <div>
        <b-navbar
            aria-controls="disabled"
            toggleable="lg"
            type="dark"
            class="nav no-key-nav"
        >
            <b-navbar-brand @click="openTextView" to="#"
                >STTylus</b-navbar-brand
            >
            <b-navbar-brand v-if="user_name">
                <small>Inloggad som: {{ user_name }}</small>
            </b-navbar-brand>

            <b-navbar-toggle target="nav-collapse" />

            <b-collapse id="nav-collapse" is-nav>
                <!-- <b-navbar-nav>
          <b-nav-item to="/news">Nyheter</b-nav-item>
        </b-navbar-nav>
        -->
                <!-- Right aligned nav items -->

                <b-navbar-nav class="ml-auto nav-fill">
                    <b-nav-item
                        v-show="notLogin"
                        no-key-nav
                        @click="addAbb('')"
                    >
                        <b>Skapa förkortning</b>
                    </b-nav-item>
                    <b-nav-item v-show="notLogin" @click="showSupport"
                        >Stödtolkning</b-nav-item
                    >
                    <b-nav-item
                        v-if="!desktop"
                        v-show="textViewOnly"
                        @click="showRemoteSettings"
                        >Distanstolkning</b-nav-item
                    >
                    <!--<b-nav-item v-show="textViewOnly && inSession" i to="#" @click="showRemotePane">
            <b-badge>
              Visa {{ connectedClients }}
              <span v-if="connectedClients > 1">anslutna</span>
              <span v-else>ansluten</span>
            </b-badge>
          </b-nav-item>
          -->
                    <b-nav-item
                        v-if="desktop"
                        v-show="notLogin"
                        @click="showLocalSettings"
                        >Lokal tolkning</b-nav-item
                    >
                    <b-nav-item v-show="textViewOnly" @click="openSettings"
                        >Inställningar</b-nav-item
                    >
                </b-navbar-nav>
            </b-collapse>
        </b-navbar>
        <Support @addAbb="addAbb" :view="view" />
        <RemoteInterpretation />
        <LocalInterpretation />
    </div>
</template>

<script>
import Support from './modals/Support.vue'
import RemoteInterpretation from './modals/RemoteInterpretation.vue'
import LocalInterpretation from './modals/LocalInterpretation.vue'
import EventBus from '../eventbus.js'

export default {
    name: 'Navigation',
    components: { Support, RemoteInterpretation, LocalInterpretation },
    props: ['view'],
    data() {
        return {
            addModalClass: ['addModalClass'],
            form: {
                abb: '',
                word: '',
                lists: [],
                selected: '',
            },
            connectedClients: 0,
            query: '',
            showRemote: false,
            secondary: false,
            toast: false,
        }
    },
    computed: {
        version() {
            return this.$version
        },
        inSession() {
            return this.$store.state.session.connected
        },
        lastUpdate() {
            return this.$lastUpdate
        },
        notLogin() {
            if (this.view !== 'login') {
                return true
            }
            return false
        },
        textViewOnly() {
            if (this.view !== 'tabula') {
                return false
            } else {
                return true
            }
        },
        user_name() {
            return this.$store.state.userData.name
        },
        desktop() {
            return this.$mode == 'desktop'
        },
        tester() {
            if (this.$store.state.userData.role == 'tester') {
                return true
            }
            return false
        },
    },
    methods: {
        connected(n) {
            this.connectedClients = n
        },
        connects() {
            this.connectedClients++
        },
        disconnects() {
            this.connectedClients--
        },
        addEventListeners() {
            EventBus.$on('clients_connected', this.connected)
            EventBus.$on('client_disconnects', this.disconnects)
            EventBus.$on('client_connects', this.connects)
        },
        removeEventListeners() {
            EventBus.$off('clients_connected', this.connected)
            EventBus.$off('client_disconnects', this.disconnects)
            EventBus.$off('client_connects', this.connects)
        },
        openTextView() {
            EventBus.$emit('openTextView')
        },
        openSettings() {
            EventBus.$emit('chatHidden')
            EventBus.$emit('openSettings')
        },
        showRemotePane() {
            EventBus.$emit('toggleChat')
        },
        watchForUnsavedChanges() {
            this.$router.push('/')
        },
        addAbb(word) {
            this.$store.commit('setSelectedWord', word)
            this.$bvModal.show('addAbb')
        },
        showSupport() {
            this.$bvModal.show('support')
        },
        showRemoteSettings() {
            this.$bvModal.show('remote')
        },
        showLocalSettings() {
            this.$bvModal.show('local')
        },
    },
    mounted() {
        this.toast = false
        this.addEventListeners()
    },
    beforeDestroy() {
        this.removeEventListeners()
    },
}
</script>

<style lang="scss">
.nav {
    background-color: #393e41;
    z-index: 1000 !important;
}
.dropdown-menu {
    width: 300px !important;
    height: 14em !important;
    padding: 1em;
}
</style>
