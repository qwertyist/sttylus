<template>
    <div>
        <b-modal id="resolve" @show="getConflicts">
            <h3>Bekräfta importering av lista</h3>
            {{ this.conflicts.length }}
            <b-table small :fields="fields" :items="conflicts" />
        </b-modal>
    </div>
</template>

<script>
export default {
    data() {
        return {
            conflicts: [],
            fields: [
                {
                    abb: {
                        label: 'Förkortning',
                        sortable: true,
                        tdClass: 'mt-2',
                    },
                },
                { new: { label: 'Text' } },
                { old: { label: 'Skriver över:' } },
            ],
        }
    },
    mounted() {},
    methods: {
        getConflicts() {
            this.conflicts = this.$store.state.conflictingAbbs

            console.log('Current conflicts are:', this.conflicts)
        },
    },
}
</script>

<style scoped></style>
