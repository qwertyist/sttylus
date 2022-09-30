<template>
    <b-modal
        id="local"
        title="Lokal tolkning"
        size="lg"
        hide-footer
        hide-backdrop
        hide-header
        return-focus=".ql-editor"
        no-fade
        @show="openModal"
        @keydown.esc="closeModal"
        @hide="closeModal"
    >
        <h3>Lokal tolkning</h3>
        <b-overlay :show="remotesession" rounded="sm">
        <template #overlay>
            <div class="text-center">
                <b-icon icon="exclamation-circle-fill" font-scale="3" animation="cylon"></b-icon>
                <p class="lead">Du är redan ansluten till en distanstolkning.</p>
            </div>
        </template>
        <b-card
            header="Tolkning på lokalt nätverk"
            sub-title="Låt andra enheter ansluta till STTylus och ta del av texten"
        >
            <b-card-text v-if="ip">
                <div v-if="local">
                    Andra enheter på nätverket ansluter till:
                    <br />
                    <br />
                    <h4>http://{{ ip }}/</h4>
                    <br />
                    <b-button variant="danger" @click="stopLocal">Avsluta</b-button>
                </div>
                <div v-else>
                    <b-button @click="startLocal">Starta</b-button>
                </div>
            </b-card-text>

            <b-card-text v-else>
                <b>Du är inte ansluten till ett lokalt nätverk.</b>
                <br />
                <b-button @click="checkConnection" variant="info">Uppdatera</b-button>
            </b-card-text>
        </b-card>
        <b-card
            header="Tolkning på projektorduk eller extern skärm"
            sub-title="Visa texten utan distraktioner på en annan skärm"
        >
            <b-card-text>
                <div v-if="screens.length == 1">Du har bara en skärm ansluten</div>
                <div v-else>Du har en extern skärm eller projektor ansluten.</div>
                <div v-if="!presentation">
                    <b-button @click="startPresentation">Utvidga</b-button>
                </div>
                <div v-else>
                    {{ window }}
                    <b-button @click="stopPresentation">Stäng</b-button>
                </div>
            </b-card-text>
        </b-card>
        </b-overlay>
    </b-modal>
</template>
<script>
import EventBus from "../../eventbus.js"
import api from "../../api/api.js";
export default {
    data() {
        return {
            connected: false,
            local: false,
            ip: "",
            presentation: false,
            screens: [],
            remotesession: null,
            window: null,
        }
    },
    mounted() {
        nw.Screen.Init()
        this.updateConnectedScreens()
        this.checkConnection()
        this.addEventListeners()
    },
    beforeDestroy() {
        this.removeEventListeners()
    },
    methods: {
        addEventListeners() {
            EventBus.$on("localConnection", this.setLocalConnection)

            console.log("Listen to screen changes")
            nw.Screen.on("displayBoundsChanged", this.displayBoundsChanged)
            nw.Screen.on("displayAdded", this.displayAdded)
            nw.Screen.on("displayRemoved", this.displayRemoved)
        },
        removeEventListeners() {
            EventBus.$off("localConnection")
        },
        setLocalConnection(status) {
            console.log("Update local...")
            this.local = status
        },
        updateConnectedScreens() {

        this.screens = nw.Screen.screens
        },
        displayBoundsChanged(screen) {
            console.log("display bounds changed:", screen)
        },
        displayAdded(screen) {
            this.$toast.warning("Extern skärm eller projektor anlöts")
        this.screens = nw.Screen.screens
            console.log("display added:", screen)
        },
        displayRemoved(screen) {
            console.log("display removed:", screen)
            this.$toast.warning("Extern skärm eller projektor kopplades ur")
            this.presentation = false
        },
        openModal() {
            if (this.$store.state.session.connected) {
                console.log("Already connected to remote session")
                this.remotesession = true
            } else {
                console.log("Not connected to remote session")
                this.remotesession = false
            }

            this.$store.commit("setModalOpen", true)

            EventBus.$emit("modalOpened");
        },
        closeModal() {
            this.$store.commit("setModalOpen", false)
            EventBus.$emit("modalClosed");
        },
        checkConnection() {

            api.getLocalIP().then(resp => {
                this.ip = resp.data.ip
                this.connected = true
            })
                .catch(err => {
                    console.error("Couldn't get local ip", err)
                    this.connected = false
                })
        },
        startLocal() {
            this.local = true
            EventBus.$emit("connectLocal")
        },
        stopLocal() {
            this.local = false
            EventBus.$emit("disconnectLocal")
        },
        startPresentation() {
            nw.Window.open(window.location.href + "presentation", {}, (child) => {
                child.on("close", () => {
                    console.log("Do something on closing child window...")
                    EventBus.$emit("stopPresentation")
                    this.presentation = false
                    child.hide()
                    child.close(true)
                })
                child.on("loaded", () => {
                    this.presentation = true
                    console.log("child loaded...")
                    EventBus.$emit("startPresentation")
                    this.child = child
                })
            })

        },
        stopPresentation() {
            this.child.close()

            EventBus.$emit("stopPresentation")
        }
    }
}

</script>
