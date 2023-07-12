<template>
    <div>
        <b-modal
            id="removeList"
            scrollable
            :title="listAction"
            hide-footer
            ref="remove-list-modal"
        >
            <template v-if="list.type < 3" />
            <template v-else />
            <b-form @submit="onSubmit" @reset="onReset">
                <h3>{{ list.name }}</h3>
                <b-badge pill variant="info">
                    {{ list.type ? 'Ämneslista' : 'Standardlista' }}
                </b-badge>
                med {{ list.counter }} förkortningar
                <hr />
                <div
                    class="abbsForRemoval"
                    style="
                        height: 60vh;
                        overflow: auto;
                        word-wrap: break-word;
                        white-space: normal;
                    "
                >
                    <b-list-group
                        v-for="(abb, idx) in abbs"
                        :key="'remove-abb-list-' + idx"
                    >
                        <b-list-group-item>
                            {{ abb.abb }} - {{ abb.word }}
                        </b-list-group-item>
                    </b-list-group>
                </div>
                <b-button variant="primary" @click="closeModal">
                    Avbryt
                </b-button>
                <template v-if="list.type < 3">
                    <b-button type="submit" variant="danger">
                        Ta bort
                    </b-button>
                </template>
                <template v-else>
                    <b-button type="submit" variant="danger">
                        Dölj lista
                    </b-button>
                </template>
            </b-form>
        </b-modal>
    </div>
</template>

<script>
import EventBus from '../../eventbus.js'
import api from '../../api/api.js'

export default {
    name: 'RemoveList',
    props: ['list', 'abbs'],
    data() {
        return {}
    },
    computed: {
        listAction() {
            if (this.list.type < 3) {
                return 'Ta bort lista'
            } else {
                return 'Dölj lista'
            }
        },
    },
    methods: {
        closeModal() {},
        onSubmit(evt) {
            evt.preventDefault()
            if (this.list.type < 3) {
                api.deleteList(this.list.id)
                    .then((response) => {
                        EventBus.$emit('removedList', this.list)
                        this.$bvModal.hide('removeList')
                    })
                    .catch((err) => {
                        alert(err)
                    })
            } else {
                console.log('Göm lista')
            }
        },
        onReset(evt) {},
    },
    mounted() {},
}
</script>
