<template>
    <div class="manuscriptEditor quillWrapper">
        <div
            ref="quillContainer"
            :class="{ 'ql-container': true }"
            :style="manuscriptStyle"
            spellcheck="false"
        ></div>
    </div>
</template>
<script>
import api from '../api/api'
import EventBus from '../eventbus'
import Quill from 'quill'
import Delta from 'quill-delta'
import keyboard from './tabula/keyboard'
import Preview from './tabula/preview'
import AddAbbreviation from './modals/AddAbbreviation.vue'
import Text from './tabula/text.js'
import wsConnection from './tabula/websocket'
export default {
    components: {
        AddAbbreviation,
    },
    props: {
        nav: Boolean,
        id: String,
    },
    data() {
        return {
            quill: null,
            websocket: null,
            name: '',
            capitalize: true,
            version: 0,
            settings: {
                font: {
                    fontFamily: 'Arial',
                    fontSize: '45px !important',
                    backgroundColor: '#000000',
                    color: '#ffff00',
                    lineHeight: '1.2',
                    height: '100vh',
                    colorID: 2,
                    padding: 20 + 'px',
                    backgroundColor: '#000000',
                    customColors: {},
                },
                behaviour: {
                    capitalizeOnNewLine: true,
                },
            },
            margins: null,
            userColor: 'red',
        }
    },
    computed: {
        manuscriptStyle() {
            return {
                fontSize: 35 + 'px',
                maxWidth: 100 + '% !important',
                height: 68 + 'vh !important',
            }
        },
        wrapper() {
            return {
                padding: 15 + 'px',
            }
        },
    },
    watch: {
        id(newVal, oldVal) {
            if (newVal == 'new') {
                this.quill.setText('')
                return
            }
            console.log('ID changed to:', this.id)
            this.loadManuscript()
        },
    },
    mounted() {
        api.cacheAbbs()
            .then((resp) => {})
            .catch((err) => {
                console.log("couldn't create cache", err)
            })
        this.websocket = null
        this.initializeEditor()
        this.quill.version = 0
        /* this.quill.setText(
      "We have been subordinate to our limitations until now. The time has come to cast aside these bonds and to elevate our consciousness to a higher plane. It is time to become a part of all things.",
      "init"
    );*/
        this.addEventListeners()
    },
    beforeDestroy() {
        Text.saveTextSettings(this.settings)
    },
    methods: {
        addEventListeners() {
            EventBus.$on('closeManuscriptEditor', this.removeEventListeners)
            EventBus.$on('addAbbreviation', this.openAddModal)
            EventBus.$on('refocus', this.focus)
            EventBus.$on('clear', this.clear)
            EventBus.$on('newLine', this.newline)
            EventBus.$on('sizeChange', this.changeTextSize)
            EventBus.$on('colorChange', this.changeColor)
            EventBus.$on('getManuscript', this.getManuscript)
        },
        removeEventListeners() {
            EventBus.$off('closeManuscriptEditor')
            EventBus.$off('addAbbreviation')
            EventBus.$off('refocus')
            EventBus.$off('clear')
            EventBus.$off('newLine')
            EventBus.$off('sizeChange')
            EventBus.$off('colorChange')
            EventBus.$off('getManuscript', this.getManuscript)
        },
        loadManuscript() {
            api.getManuscript(this.id).then((resp) => {
                this.quill.setContents(JSON.parse(resp.data.content), 'api')
            })
        },
        getManuscript() {
            console.log('get manuscript content:', this.quill.getContents())
            EventBus.$emit('editedManuscript', this.quill.getContents())
        },
        changeTextSize(inc) {
            this.settings.font.fontSize = Text.changeTextSize(
                inc,
                this.settings.font.fontSize
            )
            this.$store.commit(
                'setFontSize',
                this.settings.font.fontSize.replace('px', '')
            )
        },
        changeColor() {
            let colors = Text.changeColor(this.settings.font.colorID)
            this.settings.font.backgroundColor = colors.background
            this.settings.font.color = colors.foreground
            this.settings.font.colorID = colors.colorID
            this.$store.commit('setFontColorID', this.settings.font.colorID)
        },
        openAddModal(phrase) {
            this.$store.commit('setSelectedWord', phrase)
            this.$bvModal.show('addAbb')
        },
        addAbb() {},
        focus() {
            Text.initText()
            let settings = Text.loadTextSettings()
            this.settings.font = settings.font

            this.$nextTick(() => {
                this.quill.focus()
                this.quill.setSelection(this.quill.getText().length)
            })
        },
        clear() {
            this.quill.setText('')
            window.scrollTo(0, 0)
            if (this.websocket) {
                this.websocket.sendClear()
                this.quill.version = 0
            }
        },
        newline(scroll) {
            window.scrollTo(0, 0)
        },
        initializeEditor() {
            this.setupEditor()
            this.$emit('ready', this.quill)
            this.$nextTick(() => {
                this.quill.focus()
                this.quill.setSelection(this.quill.getText().length)
            })
        },
        loadTextSettings() {
            let settings = Text.loadTextSettings()
            this.settings.font = settings.font
            this.settings.behaviour = settings.behaviour
        },
        setupEditor() {
            this.loadTextSettings()
            Text.initText()

            const editorConfig = {
                debug: false,
                theme: 'snow',
                modules: {
                    toolbar: [
                        'bold',
                        'italic',
                        'underline' /*, { "align": []}*/,
                    ],
                    keyboard: {
                        capitalizeOnNewLine:
                            this.settings.behaviour.capitalizeOnNewLine,
                        manuscriptEditor: true,
                    },
                },
            }

            Quill.register('modules/keyboard', keyboard, true)
            Quill.register('formats/preview', Preview, true)
            //    Quill.register("modules/keyboard", keyboard, true);
            this.quill = new Quill(this.$refs.quillContainer, editorConfig)

            this.quill.clipboard.addMatcher(
                Node.ELEMENT_NODE,
                (node, delta) => {
                    delta.ops = delta.ops.map((op) => {
                        if (typeof op.insert !== 'string') {
                            return { insert: '' }
                        }
                        return {
                            insert: op.insert,
                        }
                    })
                    return delta
                }
            )
        },
        saveSettings() {
            let settings = this.$store.state.settings
            console.log('save settings:', settings)
        },
    },
}
</script>
<style src="./tabula/manuscript/quill.scss" lang="scss"></style>
<style src="./tabula/manuscript/quill.snow.css"></style>
<style scoped>
.manuscriptEditor {
    max-height: 50% !important;
}
.manuscriptEditor .quillWrapper {
    height: 65vh;
    width: 65vw !important;
}
.ql-editor {
    outline: 1px solid red;
    width: 100% !important;
}

.ql-container {
    width: 65vw !important;
    height: 95% !important;
}
.navMargin {
    position: fixed;
    width: 100vw;
}
.preview {
    opacity: 0.6;
}
.manuscriptWrapper {
    height: 60vh;
}
</style>
