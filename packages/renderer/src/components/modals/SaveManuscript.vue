<template>
    <div>
        <b-modal
            title="Spara manuskript"
            size="lg"
            hide-footer
            @show="show"
            @hide="hide"
            ref="savemanuscript"
            id="savemanuscript"
        >
            <form @submit.prevent="onSubmit" inline autocomplete="off">
                <b-form-input
                    required
                    v-model="form.title"
                    autofocud
                    placeholder="Titel..."
                    type="text"
                />
                <b-form-input v-model="form.id" hidden />
                <div class="float-right">
                    <b-button variant="danger" type="reset" @click="hide"
                        >Avbryt</b-button
                    >
                    <b-button variant="primary" type="submit">Spara</b-button>
                </div>
                <div v-if="id != 'new'">
                    <b-row>
                        <b-col cols="1">
                            <b-form-checkbox v-model="createNew" />
                        </b-col>
                        <b-col cols="7">Spara nytt manuskript</b-col>
                    </b-row>
                </div>
            </form>
        </b-modal>
    </div>
</template>
<script>
import api from '../../api/api.js'
import Delta from 'quill-delta'
import EventBus from '../../eventbus.js'
export default {
    name: 'SaveManuscript',
    props: {
        id: String,
        close: Boolean,
        manuscript: Object,
    },
    data() {
        return {
            ro: false,
            doc: {},
            createNew: false,
            form: {
                title: '',
                preview: '',
                content: {},
                id: '',
            },
        }
    },
    watch: {
        createNew() {
            this.$nextTick(() => {
                if (this.createNew) {
                    this.form.title = ''
                    this.form.id = 'new'
                } else {
                    this.form.id = this.manuscript.id
                    this.form.title = this.manuscript.title
                }
            })
        },
    },
    computed: {
        validForm() {
            if (this.form.title.length == 0) {
                return null
            }
            if (this.form.title.trim() == '') {
                return false
            } else {
                return true
            }
        },
    },
    methods: {
        show({ id, prevOiew, body, close }) {
            EventBus.$emit('getManuscript')
            this.createNew = false
            if (Object.keys(this.manuscript).length != 0) {
                this.form = { ...this.manuscript }
            } else {
                this.form = {
                    id: 'new',
                    title: '',
                }
            }
        },
        hide() {
            EventBus.$emit('manuscriptNotSaved')
        },
        onSubmit(evt) {
            this.form.content = JSON.stringify(this.doc)
            console.log('create new manuscript:', this.createNew)
            console.log('id:', this.form.id)
            if (!this.createNew && this.form.id != 'new') {
                api.updateManuscript(this.form)
                    .then((resp) => {
                        EventBus.$emit('savedManuscript', {
                            id: resp.data.id,
                            close: this.close,
                        })
                    })
                    .catch((err) => {
                        console.log('failed updating manuscript', err)
                    })
            } else {
                api.createManuscript(this.form)
                    .then((resp) => {
                        EventBus.$emit('savedManuscript', {
                            id: resp.data.id,
                            close: this.close,
                        })
                    })
                    .catch((err) =>
                        console.log('create manuscript failed:', err)
                    )
            }
        },
    },
    beforeDestroy() {
        EventBus.$off('editedManuscript')
    },
    mounted() {
        EventBus.$on('editedManuscript', (doc) => {
            this.doc = doc
        })
        this.form.title = ''
    },
}
</script>
