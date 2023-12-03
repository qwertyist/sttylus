<template>
    <div>
        <b-tabs v-model="tabIndex" content-class="subTab">
            <b-tab title="Välj/visa" class="mt2">
                <b-row>
                    <b-col cols="3">
                        <h3 class="mt-2">
                            Valda manuskript
                            <b-button
                                variant="danger"
                                class="float-right mt-0"
                                @click="deselectAllManuscripts"
                                >Avmarkera alla</b-button
                            >
                        </h3>
                        <hr />
                        <div>
                            <b-form autocomplete="off">
                                <b-list-group
                                    style="margin-bottom: 10px; max-height: calc(100vh - (42px*4); overflow-y: auto"
                                >
                                    <b-list-group-item
                                        v-for="(
                                            s, index
                                        ) in selectedManuscripts"
                                        v-bind:key="index"
                                    >
                                        <b-row>
                                            <b-col
                                                cols="7"
                                                class="align-self-center"
                                            >
                                                {{ s.title }}
                                            </b-col>
                                            <b-col class="align-self-center">
                                                <b-input-group>
                                                    <b-input
                                                        placeholder="Förkortning"
                                                        v-model="s.abb"
                                                        :formatter="
                                                            abbFormatter
                                                        "
                                                        @focus="
                                                            setCurrentAbb(
                                                                index,
                                                                s.abb
                                                            )
                                                        "
                                                        @blur="
                                                            setManuscriptAbb(
                                                                index,
                                                                s.abb
                                                            )
                                                        "
                                                    ></b-input>
                                                    <b-input-group-append>
                                                        <b-button
                                                            style="margin: 0px"
                                                            variant="danger"
                                                            @click="
                                                                deselectManuscript(
                                                                    index,
                                                                    ''
                                                                )
                                                            "
                                                            >x</b-button
                                                        >
                                                    </b-input-group-append>
                                                </b-input-group>
                                            </b-col>
                                        </b-row>
                                    </b-list-group-item>
                                </b-list-group>
                            </b-form>
                        </div>
                    </b-col>
                    <b-col cols="9">
                        <div>
                            <b-list-group
                                scrollable
                                style="
                                    max-height: calc(100vh - (50px * 3));
                                    overflow-y: auto;
                                "
                            >
                                <b-list-group-item
                                    v-for="m in manuscripts"
                                    v-bind:key="m.id"
                                >
                                    <b-row>
                                        <b-col>
                                            <h4>{{ m.title }}</h4>
                                        </b-col>
                                        <b-col>
                                            <b-form-group class="float-right">
                                                <b-button
                                                    class="float-right"
                                                    variant="primary"
                                                    @click="selectManuscript(m)"
                                                    >Välj</b-button
                                                >
                                                <b-button
                                                    class="float-right"
                                                    variant="warning"
                                                    @click="editManuscript(m)"
                                                    >Redigera</b-button
                                                >
                                                <b-button
                                                    class="float-right"
                                                    variant="danger"
                                                    @click="
                                                        openConfirmRemoveModal(
                                                            m
                                                        )
                                                    "
                                                    >Ta bort</b-button
                                                >
                                            </b-form-group>
                                        </b-col>
                                    </b-row>
                                    <!--<div class="preview">
                    <span class="preview" v-html="m.content"></span>
                  </div>-->
                                </b-list-group-item>
                            </b-list-group>
                        </div>
                    </b-col>
                    <b-col cols="2">
                        <b-list-group></b-list-group>
                    </b-col>
                </b-row>
                <b-modal
                    ref="confirmRemoveModal"
                    size="lg"
                    hide-footer
                    hide-backdrop
                    @hide="closeConfirmRemoveModal"
                    title="Bekräfta ta bort"
                >
                    <h4>{{ selectedForRemoval.title }}</h4>
                    <hr />

                    <div class="preview" style="height: 45vh !important">
                        <span
                            class="preview"
                            v-html="selectedForRemoval.preview"
                        ></span>
                    </div>
                    <hr />
                    <b-button
                        class="float-right"
                        variant="secondary"
                        @click="closeConfirmRemoveModal()"
                        >Avbryt</b-button
                    >

                    <b-button
                        class="float-right"
                        variant="danger"
                        @click="removeManuscript()"
                        >Ta bort</b-button
                    >
                </b-modal>
                <b-modal
                    ref="confirmCreateNewModal"
                    hide-footer
                    hide-header
                    hide-backdrop
                    @hide="closeConfirmCreateNewModal()"
                >
                    <p>
                        Vill du spara manuset du redigerar innan du skapar ett
                        nytt?
                    </p>

                    <b-button
                        class="float-right"
                        variant="primary"
                        @click="saveAndCreateNew()"
                        >Spara och skapa nytt</b-button
                    >
                    <b-button
                        class="float-right"
                        variant="danger"
                        @click="createNewWithoutSaving()"
                        >Skapa nytt</b-button
                    >
                </b-modal>
                <b-modal
                    ref="stopEditingModal"
                    hide-footer
                    hide-header
                    hide-backdrop
                    @hide="closeStopEditingModal()"
                >
                    <p>Vill du spara manuset du redigerar innan du stänger?</p>
                    <b-button
                        class="float-right"
                        variant="primary"
                        @click="saveAndClose()"
                        >Spara</b-button
                    >
                    <b-button
                        class="float-right"
                        variant="danger"
                        @click="dontSaveAndClose()"
                        >Stäng utan att spara</b-button
                    >
                </b-modal>
            </b-tab>
            <b-tab
                :title="createAction.text"
                :title-link-class="createAction.class"
                class="manuscriptTab"
            >
                <div class="manuscriptEditor">
                    <ManuscriptEditor v-bind:id="id" />
                </div>
                <br />
                <b-button @click="saveManuscript">Spara</b-button>
            </b-tab>
            <template
                v-if="id && id !== 'parsed' && id !== 'new'"
                v-slot:tabs-end
            >
                <b-nav-item
                    @mouseover="hoverClose(true)"
                    @mouseout="hoverClose(false)"
                    v-bind:style="editAction.class"
                    role="presentation"
                    @click="openStopEditingModal()"
                >
                    <span :style="editAction.class">
                        {{ editAction.text }}
                    </span>
                </b-nav-item>
                <b-nav-item
                    href="#"
                    role="presentation"
                    @click="openConfirmCreateNewModal()"
                    >Skapa nytt</b-nav-item
                >
            </template>
        </b-tabs>
        <SaveManuscript
          v-bind:manuscript="manuscript"
          v-bind:id="id"
          v-bind:close="close"
        />
    </div>
</template>
<script>
import api from '../../api/api.js'
import Vue from 'vue'
import EventBus from '../../eventbus.js'
import ManuscriptEditor from '../ManuscriptEditor.vue'
import SaveManuscript from '../modals/SaveManuscript.vue'
import Delta from 'quill-delta'
export default {
    name: 'AddAbbreviation',
    components: {
        ManuscriptEditor,
        SaveManuscript,
    },
    created() {
        this.actions = [
            { text: 'Redigerar', class: '' },
            { text: 'Stäng', class: 'background-color: #F00; color: #FFF;' },
            { text: 'Skapa', class: '' },
            { text: 'Fortsätt redigera', class: '' },
            { text: '...', class: '' },
        ]
    },
    data() {
        return {
            tabIndex: 0,
            selectedManuscripts: [],
            manuscripts: [],
            manuscript: {},
            currentAbb: '',
            i: 1,
            freeIndexes: [],
            id: 'new',
            title: '',
            close: true,
            content: '',
            modalOpen: false,
            createAction: {
                text: 'Skapa',
                class: '',
            },
            editAction: {
                text: '',
                class: '',
            },
            selectedForRemoval: {
                title: '',
                content: '',
            },
            changed: false,
            saved: false,
        }
    },
    mounted() {
        this.getManuscripts()
        this.addEventListener()
        this.selectedManuscripts =
            this.$store.state.settings.selectedManuscripts || []
        this.selectedManuscripts = this.selectedManuscripts.sort()
        this.i = this.selectedManuscripts.length + 1

        if (this.i == 0) {
            this.i = 1
        }
    },
    beforeDestroy() {
        this.removeEventListener()
    },
    watch: {
        selectedManuscripts: {
            deep: true,
            handler: function (change) {},
        },
        tabIndex: {
            handler: function (change) {
                if (change == 0) {
                    EventBus.$emit('leaveManuscriptEditor')
                    if (this.id == '' || this.id == 'new') {
                        this.createAction = { ...this.actions[2] }
                    } else {
                        this.createAction = { ...this.actions[3] }
                        this.createAction.text += ': ' + this.title
                        this.editAction = { ...this.actions[1] }
                    }
                } else {
                    if (this.id == '' || this.id == 'new') {
                        this.createAction = { ...this.actions[2] }
                    } else {
                        this.createAction = { ...this.actions[4] }
                        this.editAction = { ...this.actions[0] }
                        this.editAction.text += ': ' + this.title
                    }
                }
            },
        },
    },
    methods: {
        addEventListener() {
            EventBus.$on('savedManuscript', this.savedManuscript)
            EventBus.$on('manuscriptChanged', this.manupscriptChanged)
            EventBus.$on('manuscriptNotSaved', this.hideModal)
        },
        removeEventListener() {
            EventBus.$off('savedManuscript')
            EventBus.$off('manuscriptChanged')
            EventBus.$off('manuscriptNotSaved')
        },
        saveManuscript() {
            console.log('Should save manuscript')
            this.$bvModal.show('savemanuscript')
            this.close = true
        },
        manupscriptChanged(value) {
            this.changed = value
        },
        getManuscripts() {
            api.getManuscripts({
                user_id: this.$store.state.userData.id,
            }).then((resp) => {
                if (!resp.data) {
                    this.manuscripts = []
                } else {
                    this.manuscripts = resp.data.sort((a, b) =>
                        a.title > b.title ? 1 : -1
                    )
                }
            })
        },
        createNew() {
            openConfirmCreateNewModal()
        },
        selectManuscript(m) {
            console.log('selected manuscripts:', this.selectedManuscripts)
            let index = this.selectedManuscripts.findIndex(
                (manuscript) => manuscript.id == m.id
            )
            console.log(index)
            if (index !== -1) {
                this.updateStore()
                return
            }
            if (this.freeIndexes.length > 0) {
                m.abb = this.freeIndexes.shift()
                this.selectedManuscripts.push(m)
            } else {
                m.abb = 'text' + this.i
                this.selectedManuscripts.push(m)
                this.i++
            }
            this.updateStore()
        },
        deselectManuscript(m) {
            console.log('deselect:', m)
            console.log('selectedManuscripts', this.selectedManuscripts)
            if (this.selectedManuscripts.length > 0) {
                let abb = this.selectedManuscripts[m].abb
                if (/text\d+/.test(abb)) {
                    console.log(
                        'auto generated abb [',
                        abb,
                        '], proceed with caution'
                    )
                    this.freeIndexes.push(abb)
                }
                this.selectedManuscripts.splice(m, 1)

                this.updateStore()
            }
        },
        deselectAllManuscripts() {
            this.selectedManuscripts = []
            this.i = 1
            this.freeIndexes = []
            this.updateStore()
        },
        abbFormatter(value) {
            let index = this.selectedManuscripts
                .filter((manuscript, index) => {
                    return index !== this.currentAbb[0]
                })
                .findIndex((manuscript) => manuscript.abb === value)
            if (index !== -1) {
                this.selectedManuscripts[this.currentAbb[0]].abb =
                    this.currentAbb[1]
            } else {
                this.selectedManuscripts[this.currentAbb[0]].abb = value
            }
            this.updateStore()
            return value.replace(/\s/g, '')
        },
        setManuscriptAbb(i, abb) {
            this.updateStore()
        },
        setCurrentAbb(i, abb) {
            this.currentAbb = [i, abb]
            this.updateStore()
        },
        updateStore() {
            console.log(
                'updateStore, view selectedManuscripts',
                this.selectedManuscripts
            )
            const manuscriptAbbs = this.selectedManuscripts.map((m) => {
                return { abb: m.abb, id: m.id, title: m.title }
            })
            const data = JSON.stringify(manuscriptAbbs)
            this.$store.commit('setSelectedManuscripts', manuscriptAbbs)
        },
        editManuscript(m) {
            this.editAction = { ...this.actions[0] }
            this.createAction = { ...this.actions[4] }
            this.id = m.id
            this.title = m.title
            this.manuscript = m
            this.tabIndex = 1
        },
        editParsedDoc(doc) {
            this.id = 'parsed'

            this.$nextTick(() => {
                this.tabIndex = 1
            })
        },
        hideModal() {
            console.log('should hide modal')
            this.$bvModal.hide('savemanuscript')
        },
        savedManuscript({ id, close }) {
            console.log('eventbus got savedManuscript:', id, 'close:', close)
            this.id = id
            this.title = ''
            this.getManuscripts()
            this.$bvModal.hide('savemanuscript')
            this.manuscript = {}
                this.$nextTick(function () {
                    // DOM is now updated
                    // `this` is bound to the current instance
                    this.$forceUpdate()
                })
            if (close) {
                EventBus.$emit("clear")
                console.log('Stäng och gå tillbaka till listan över manuskript')
                this.tabIndex = 0
                this.id = 'new'
                this.createAction = this.actions[2]
            } else {
                this.tabIndex = 1
                this.createAction = this.actions[2]
                console.log('Nytt manuskript')
            }
        },
        removeManuscript() {
            api.deleteManuscript(this.selectedForRemoval.id).then(
                (response) => {
                    let index = this.selectedManuscripts.findIndex(
                        (manuscript) =>
                            manuscript.id === this.selectedForRemoval.id
                    )
                    if (index > -1) {
                        this.deselectManuscript(index)
                    }
                    this.getManuscripts()
                    this.selectedForRemoval = {}
                    this.$refs.confirmRemoveModal.hide()
                }
            )
        },
        hoverClose(on) {
            if (this.tabIndex == 0 || this.modalOpen) {
                return
            }
            if (on) {
                this.editAction = { ...this.actions[1] }
                this.editAction.text += ': ' + this.title
            } else {
                this.editAction = { ...this.actions[0] }
                this.editAction.text += ': ' + this.title
            }
        },
        openStopEditingModal() {
            this.modalOpen = true
            this.editAction = { ...this.actions[1] }
            if (this.tabIndex > 0) {
                this.editAction.text += ': ' + this.title
            }
            //EventBus.$emit("stopEditingManuscript");
            this.$refs.stopEditingModal.show()
        },
        closeStopEditingModal() {
            console.log('stängs')
            if (this.tabIndex > 0) {
                this.editAction = { ...this.actions[0] }
                this.editAction.text += ': ' + this.title
            }
            this.modalOpen = false
            this.$refs.stopEditingModal.hide()
        },
        openConfirmRemoveModal(m) {
            let content = new Delta(JSON.parse(m.content))
            console.log(content)
            this.selectedForRemoval = m
            let preview = []
            content.eachLine((line, attributes) => {
                line.ops.map((op) => {
                    if (op.insert) {
                        preview.push(op.insert)
                    }
                })
                preview.push('\n')
            })
            console.log(preview)
            this.selectedForRemoval.preview = preview.join('\n')
            this.$refs.confirmRemoveModal.show()
        },
        closeConfirmCreateNewModal() {},
        closeConfirmRemoveModal() {
            this.selectedForRemoval = {}
            this.$refs.confirmRemoveModal.hide()
        },
        openConfirmCreateNewModal() {
            this.$refs.confirmCreateNewModal.show()
        },
        createNewWithoutSaving() {
            this.id = 'new'
            this.title = ''
            this.manuscript = {}
            this.createAction = { ...this.actions[2] }
            this.$refs.confirmCreateNewModal.hide()
        },
        saveAndCreateNew() {
            this.close = false
            this.$refs.confirmCreateNewModal.hide()
            this.$bvModal.show('savemanuscript')
            EventBus.$emit('saveManuscript', false)
        },
        dontSaveAndClose() {
            this.id = 'new'
            this.title = ''
            this.tabIndex = 0
            this.manuscript = {}
            this.$nextTick(() => {
                if (this.tabIndex > 0) {
                    this.editAction = { ...this.actions[1] }
                    this.editAction.text += ': ' + this.title
                } else {
                    this.createAction = { ...this.actions[2] }
                }
            })
            this.$refs.stopEditingModal.hide()
        },
        saveAndClose() {
            this.$refs.stopEditingModal.hide()
            this.$bvModal.show('savemanuscript')
            EventBus.$emit('saveManuscript', true)
        },
    },
}
</script>

<style scoped>
.manuscriptEditor {
    width: 100%;
    max-height: 98% !important;
    height: 98% !important;
}
p {
    max-width: 100% !important;
}
.preview {
    white-space: pre-wrap !important;

    overflow-y: scroll !important;
    height: 10em !important;
    word-wrap: break-word !important;
    width: 100%;
}
.close {
    background: rgb(247, 119, 119);
}
</style>
