<template>
    <div>
        <b-form v-if="show" @submit="onSubmit" @reset="onReset">
            <b-form-group
                id="licensekey"
                label="Licensnyckel"
                label-for="licensekey-input"
                description="För att använda sttylus behöver du en licensnyckel"
            >
                <b-form-input
                    v-model="form.licenseKey"
                    id="licensekey-input"
                    required
                    placeholder="Skriv in din licensnyckel"
                />
            </b-form-group>
            <b-button type="submit" variant="primary"> Skicka in </b-button>
            <b-button type="reset" variant="danger">
                Återställ formulär
            </b-button>
        </b-form>
    </div>
</template>

<script>
import axios from 'axios'
export default {
    name: 'License',
    data() {
        return {
            form: {
                licenseKey: '',
            },
            show: true,
        }
    },
    methods: {
        onSubmit(evt) {
            evt.preventDefault()
            axios
                .post(this.$backend + '/api/user/auth', {
                    license_key: this.form.licenseKey,
                    machine_id: this.$store.state.machineId,
                })
                .then((response) => {
                    this.$store.commit('setLicenseKey', this.form.licenseKey)
                    const user = {
                        id: response.data.id,
                        name: response.data.name,
                        company: response.data.company,
                        description: response.data.description,
                        email: response.data.email,
                        phone: response.data.phone,
                        role: response.data.role,
                    }
                    this.$store.commit('setUserData', user)
                })
                .catch((err) => {
                    console.log(err)
                })
        },
        onReset(evt) {
            evt.preventDefault()
            this.form.licenseKey = this.$store.state.licenseKey
            this.show = false
            this.$nextTick(() => {
                this.show = true
            })
        },
    },
    mounted() {
        this.form.licenseKey = this.$store.state.licenseKey
    },
}
</script>
