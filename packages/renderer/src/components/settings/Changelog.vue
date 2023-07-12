<template>
    <!--<b-tabs>
    <b-tab title="Uppdateringslogg">
    -->
    <div>
        <b-row style="max-height: 88vh; overflow-y: scroll">
            <b-col cols="8">
                <b-card-group
                    v-for="d in changelog"
                    v-bind:key="'updates_on_' + d.date"
                >
                    <b-card :header="d.date | formatChangeLogDate">
                        <b-list-group
                            v-for="c in d.changes"
                            v-bind:key="'update_' + c.id"
                        >
                            <span v-if="c.type == 'version'">
                                <h3>{{ c.description }}</h3>
                            </span>
                            <b-list-group-item v-else>
                                <span v-if="c.type == 'fix'">
                                    <b>Buggfix: </b>
                                </span>
                                <span v-html="c.description"></span>
                                <span
                                    class="float-right"
                                    v-if="c.type != 'minor' && c.action != null"
                                >
                                    <small>Åtgärd: {{ c.action }}</small>
                                </span>
                            </b-list-group-item>
                        </b-list-group>
                    </b-card>
                </b-card-group>
                <hr />
                <h2>Uppdateringar äldre än {{ older.before | formatChangeLogDate }}</h2>
                <hr />
                      <b-list-group
                            v-for="c in older.changes"
                            v-bind:key="'update_' + c.id"
                        >
                        <b-list-group-item>
                                <span v-if="c.type == 'fix'">
                                    <b>Buggfix: </b>
                                </span>
                                <span v-html="c.description"></span>
                                <span
                                    class="float-right"
                                    v-if="c.type != 'minor' && c.action != null"
                                >
                                </span>
                            </b-list-group-item>
                        </b-list-group>
            </b-col>
            <b-col>

        <p class="lead">
            <br />
            Här listas uppdateringar, ändringar och buggfixar som gjorts i
            programmet eller distanstolkningstjänsten.
        </p>

        <b-alert show variant="warning" v-if="error && !desktop">
            Kunde inte hämta aktuell uppdateringslogg. Visar senast hämtade.
        </b-alert>
            </b-col>
        </b-row>
        </div>
        <!--</b-tab>
    <b-tab title="Version">
      <b-overlay :show="loading">
        <b-jumbotron>
          Du använder STTylus version:
          <code>{}</code>
          <br />Den senaste tillgängliga versionen av programmet är:
          <code>{}</code>
          <br />
        </b-jumbotron>
      </b-overlay>
    </b-tab>
  </b-tabs>
  -->
    </div>
</template>

<script>

import { changelog, older } from "../../changelog.json"
import axios from 'axios'
import EventBus from '../../eventbus'
export default {
    data() {
        return {
            loading: true,
            error: null,
            latest: new Date('2023, 07, 12'),
            changelog: [],
            older: [],
        }
    },
    computed: {
        desktop() {
            return this.$mode == "desktop";
        },
        version() {
            return {
                version: '0.0.0',
                frontend: '0.0.0',
                backend: '0.0.0',
                latest: {},
            }
        },
    },
    methods: {
      getChangelog() {
        axios.get("https://sttylus.se/changelog.json")
          .then(resp => {
            this.changelog = resp.data.changelog
            this.older = resp.data.older
          })
          .catch(err => {
            this.error = err
            this.changelog = changelog
            this.older = older;
            console.log(older);
            this.$forceUpdate();
          })
      }
    },
    mounted() {
        this.getChangelog()
        if (this.$store.state.lastLogin > this.latest) {
            console.log('number of changes:', this.changelog[0].changes.length)
            setTimeout(() => {
                EventBus.$emit('newChangeLog', this.changelog[0].changes.length)
            }, 125)
        }
    },
}
</script>
