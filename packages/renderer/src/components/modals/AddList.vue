<template>
    <b-modal
        id="addList"
        title="Lägg till förkortningslista"
        hide-footer
        no-fade
        ref="addlistmodal"
    >
        <b-form @submit="onSubmit" autocomplete="off">
            <b-form-group
                id="list-name"
                label-for="list-name-input"
                description="Förkortningslistans namn"
            >
                <b-form-input
                    v-model="form.name"
                    id="list-name-input"
                    required
                    autofocus
                    :max-length="20"
                    placeholder="Ge förkortningslistan ett namn..."
                />
            </b-form-group>
            <b-row>
                <b-col cols="7">
                    Ämneslista
                    <b-form-checkbox
                        v-model="form.standard"
                        inline
                        name="standard-switch"
                        switch
                    >
                        Standardlista
                    </b-form-checkbox>
                </b-col>
                <b-col>
                    <b-button variant="danger" @click="closeModal">
                        Avbryt
                    </b-button>
                    <b-button type="submit" variant="primary">
                        Lägg till
                    </b-button>
                </b-col>
            </b-row>
        </b-form>
    </b-modal>
</template>
<script>
import api from '../../api/api.js'
import EventBus from '../../eventbus.js'
export default {
    data() {
        return {
            form: {
                name: null,
                standard: false,
            },
        }
    },
    methods: {
        onSubmit(evt) {
            evt.preventDefault()
            const data = {
                type: this.form.standard ? 0 : 1,
                name: this.form.name,
            }
            api.createList(data)
                .then((resp) => {
                    this.form.name = ''
                    this.form.standard = false
                    this.$store.commit('subscribeList', resp.data.id)
                    EventBus.$emit('createdList')
                    this.closeModal()
                })
                .catch((err) => {
                    console.log('Failed creating list:', err)
                })
        },
        closeModal() {
            this.form.abb = ''
            this.form.word = ''
            this.$bvModal.hide('addList')
        },
    },
}
</script>
