<template>
    <b-modal
        id="colorPicker"
        size="md"
        hideBackdrop
        noFade
        ref="colorPicker"
        title="Välj färger till eget färgschema"
    >
        <div
            :style="{
                padding: 15 + 'px',
                color: picker.textColor,
                backgroundColor: picker.bgColor,
            }"
            class="text-wrap"
        >
            <b>Lorem</b> <i>ipsum</i> dolor <u>sit</u> amet, consectetur
            adipisicing elit, sed do eiusmod tempor incididunt ut labore et
            dolore magna aliqua.
        </div>
        <br />
        <b-row>
            <b-col cols="5"></b-col>
            <b-col cols="3"><div class="float-right">Textfärg</div> </b-col>
            <b-col>
                <v-swatches
                    swatches="text-advanced"
                    v-model="picker.textColor"
                ></v-swatches>
            </b-col>
        </b-row>
        <b-row>
            <b-col cols="5"></b-col>
            <b-col cols="3"><div class="float-right">Bakgrundsfärg</div></b-col>
            <b-col>
                <v-swatches
                    swatches="text-advanced"
                    v-model="picker.bgColor"
                ></v-swatches> </b-col
            ><b-col></b-col
        ></b-row>
        <template #modal-footer>
            <b-button variant="secondary" size="sm" @click="closeModal"
                >Avbryt</b-button
            >
            <b-button variant="primary" size="sm" @click="saveColor"
                >Spara</b-button
            >
        </template>
    </b-modal>
</template>
<script>
import VSwatches from 'vue-swatches'
import 'vue-swatches/dist/vue-swatches.css'
import EventBus from '../../eventbus'
export default {
    components: { VSwatches },
    data() {
        return {
            picker: { textColor: 'white', bgColor: 'black' },
        }
    },
    mounted() {
        let custom = this.$store.state.settings.font.customColors || {
            foreground: 'white',
            background: 'black',
        }
        this.picker = {
            textColor: custom.foreground,
            bgColor: custom.background,
        }
    },
    methods: {
        saveColor() {
            let custom = {
                foreground: this.picker.textColor,
                background: this.picker.bgColor,
            }
            this.$store.commit('setCustomColors', custom)
            EventBus.$emit('savedColors', custom)
            this.$bvModal.hide('colorPicker')
        },
        closeModal() {
            this.$bvModal.hide('colorPicker')
        },
    },
}
</script>
<style>
.vue-swatches__trigger {
    border: 2px #000 solid;
}
</style>
