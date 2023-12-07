<template>
    <b-modal
        @show="onShow()"
        id="editList"
        title="Redigera förkortningslista"
        hide-footer
        no-fade
        ref="editListModal"
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
                        Bekräfta
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
    props: ['list'],
    data() {
        return {
            form: {
                id: null,
                name: null,
                standard: false,
            },
        }
    },
    methods: {
        onShow() {
            this.form = { ...this.list }
            this.form.standard = !Boolean(this.list.type)
        },
        onSubmit(evt) {
            evt.preventDefault()
            const data = this.form
            ;(data.type = this.form.standard ? 0 : 1),
                api
                    .updateList(data)
                    .then((resp) => {
                        this.form.name = ''
                        this.form.standard = false
                        EventBus.$emit('updatedList', data)
                        this.closeModal()
                    })
                    .catch((err) => {
                        console.log('Failed updating list:', err)
                    })
        },
        closeModal() {
            this.$bvModal.hide('editList')
        },
    },
}
</script>
