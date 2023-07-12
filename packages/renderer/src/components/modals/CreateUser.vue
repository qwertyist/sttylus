<template>
    <b-modal
        id="createUser"
        title="Skapa ny användare"
        hide-footer
        no-fade
        ref="createusermodal"
    >
        <b-form @submit="onSubmit" autocomplete="off">
            <b-form-group
                id="user-email"
                label-for="user-email-input"
                description="Ange skrivtolkens epostadress. Denna används för inloggning."
            >
                <b-form-input
                    type="email"
                    v-model="form.email"
                    :state="validateEmail"
                    id="user-email-input"
                    required
                />
            </b-form-group>
            <b-form-group
                id="user-name"
                label-for="user-name-input"
                description="Ange skrivtolkens namn."
            >
                <b-form-input
                    type="text"
                    v-model="form.name"
                    id="user-name-input"
                    required
                />
            </b-form-group>
            <b-form-group
                id="user-role"
                label-for="user-role-select"
                description="Bestäm användarens behörighetsnivå (roll)"
            >
                <b-form-select :options="roles" v-model="form.role" />
            </b-form-group>
            <b-row>
                <b-col cols="7" />
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
                email: '',
                name: '',
                role: 'user',
            },
            validateEmail: null,
            roles: [
                { value: 'user', text: 'Skrivtolk' },
                { value: 'tester', text: 'Testanvändare' },
                { value: 'admin', text: 'Administratörkonto' },
            ],
        }
    },
    methods: {
        onSubmit(evt) {
            evt.preventDefault()
            api.createUser(this.form)
                .then(() => {
                    this.$toast.info('Användaren skapades')
                    EventBus.$emit('createdUser')
                    this.closeModal()
                })
                .catch((err) => {
                    if (err.response.status == 409) {
                        this.$toast.error('E-postadressen är redan registrerad')
                        this.validateEmail = false
                    }
                })
        },
        closeModal() {
            this.form.email = ''
            this.form.name = ''
            this.form.role = 'user'
            this.valiadteEmail = null
            this.$bvModal.hide('createUser')
        },
    },
}
</script>
