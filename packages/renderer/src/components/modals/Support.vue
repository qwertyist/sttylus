<template>
  <b-modal
    id="support"
    size="lg"
    hideFooter
    hideHeader
    noFade
    hideBackdrop
    @keyup.esc="catchEscape"
    @show="showModal"
    @hide="closeModal"
    tabindex="1000"
  >
    <div class="d-flex justify-content-between align-items-center">
      <h4>Slå upp förkortning</h4>

      <b-badge>Fokusera med CTRL+3</b-badge>
    </div>
    <b-form autocomplete="off" @submit.prevent="lookup">
      <b-form-input
        ref="lookupInput"
        v-model="lookupPhrase"
        :autofocus="!sharingList"
        placeholder
      />
      <b-list-group>
        <b-list-group-item v-for="m in lookupResults" :key="'lookup_' + m.id">
          <span v-if="m.id >= 0">
            {{ m.match }}
          </span>
          <span v-else-if="m.id == -1">
            <b-badge variant="danger">{{ m.match }}</b-badge>
          </span>
          <span v-else>
            <b-badge variant="warning">{{ m.match }}</b-badge>
          </span>
        </b-list-group-item>
      </b-list-group>
    </b-form>

    <hr />

    <div class="d-flex justify-content-between align-items-center">
      <h4>Missade förkortningar</h4>
      <b-badge>Visa/dölj med CTRL+4</b-badge>
    </div>
    <div v-if="show.missedAbbs">
      <b-list-group
        style="overflow-y: auto; overflow-x: hidden; max-height: 40vh"
      >
        <b-list-group-item
          class="p-0 pl-2"
          v-for="abb in filteredMissedAbbs"
          :key="'missed_' + abb.word"
        >
          <b-row>
            <b-col cols="2">{{ abb.abb }}</b-col>
            <b-col cols="4"> {{ abb.word }}</b-col>
            <b-col cols="4"> {{ abb.counter }}</b-col>
            <b-col>
              <b-button
                size="sm"
                @click="dontRemindMissed(abb.abb)"
                variant="warning"
                >Påminn inte</b-button
              >
            </b-col>
          </b-row>
        </b-list-group-item>
      </b-list-group>
      <br />
      <b-button class="float-right" size="sm" @click="resetMissedAbbs()"
        >Nollställ</b-button
      >
    </div>
    <div v-else>
      <b-badge>...</b-badge>
    </div>
    <hr />

    <!--
        <div class="d-flex justify-content-between align-items-center">
            <h4>Ord som saknar förkortningar</h4>
            <b-badge>Visa/dölj med CTRL+5</b-badge>
        </div>
        <div v-if="show.suggestedAbbs">
            <b-list-group size="sm" style="max-height: 40vh; overflow-y: auto">
                <b-list-group-item
                    class="p-0 pl-2"
                    v-for="(s, index) in suggestions"
                    :key="'suggestion_' + index"
                >
                    {{ s }}
                    <span class="float-right">
                        <b-button size="sm" variant="primary" @click="addAbb(s)"
                            >Lägg till</b-button
                        >
                        <b-button size="sm" @click="ignoreSuggestion(s)"
                            >Ignorera</b-button
                        >
                    </span>
                </b-list-group-item>
            </b-list-group>
            <div class="float-right">
                <b-button
                    variant="danger"
                    size="sm"
                    @click="ignoreAllSuggestions"
                    >Ignorera alla</b-button
                >
            </div>
        </div>
        <div v-else>
            <b-badge>...</b-badge>
        </div>
        -->
    <b-form-datalist id="no-datalist"></b-form-datalist>
  </b-modal>
</template>

<script>
import api from '../../api/api'
import EventBus from '../../eventbus.js'

export default {
  props: ['view'],
  data() {
    return {
      collab: false,
      selectedBaseList: null,
      selectedLists: [],
      form: {
        sharedAbb: '',
        sharedWord: '',
        sharedList: '',
      },
      sharingList: false,
      sharedAbbs: null,
      show: {
        sharedAbbs: false,
        missedAbbs: false,
        suggestedAbbs: false,
      },
      connected: false,
      letters: ['Alla'],
      currentLetter: 'Alla',
      lastLetter: '',
      lookupPhrase: '',
      lookupResults: [],
      filteredMissedAbbs: [],
      suggestions: [],
    }
  },
  computed: {
    sharedAbbsByLetter() {
      if (this.currentLetter == 'Alla') {
        return this.sharedAbbs
      } else if (this.currentLetter == '#') {
        return this.sharedAbbs.filter((a) => {
          return a.abb[0].toLowerCase() == a.abb[0].toUpperCase()
        })
      } else {
        return this.sharedAbbs.filter((a) => {
          return a.abb.startsWith(this.currentLetter)
        })
      }
    },
  },
  methods: {
    catchEscape() {
      this.closeModal()
      this.$bvModal.hide('support')
    },
    showModal() {
      this.$store.commit('setModalOpen', true)
      EventBus.$on('sharedAbbEvent', this.getSharedAbbs)
      EventBus.$emit('modalOpened')
      document.addEventListener('keydown', this.preventDefaults)
      this.getSelectedLists()
      this.getSharedAbbs()
      this.filterMissedAbbs()
      this.getSuggestions()
      this.connected = this.$store.state.session.connected
      if (this.connected) {
        this.collab = true
        this.sharedList = true
      }
    },
    filterMissedAbbs() {
      let missedAbbs = [...this.$store.state.missedAbbs.keys()]
      this.filteredMissedAbbs = []
      missedAbbs.forEach((a) => {
        console.log(a)
        this.filteredMissedAbbs.push(this.$store.state.missedAbbs.get(a))
      })
    },
    closeModal() {
      this.lookupPhrase = ''
      this.lookupResults = []
      this.$store.commit('setModalOpen', false)
      EventBus.$off('sharedAbbEvent')
      this.$bvModal.hide('addAbb')
      document.removeEventListener('keydown', this.preventDefaults)
      this.$store.commit('setModalOpen', true)
      EventBus.$emit('modalClosed')
      EventBus.$emit('abbModalClosed')
      EventBus.$emit('closeNav')
    },
    onSubmitSharedList() {
      if (this.form.sharedList != '' && this.selectedBaseList == null) {
        api
          .joinSharedList(this.form.sharedList.trim())
          .then((resp) => {
            this.$store.commit('setSharedList', { id: resp.data })
            if (resp.data != this.form.sharedList) {
              this.$toast.info(
                'List-ID:t existerar inte, skapar en ny delad lista'
              )
              this.form.sharedList = resp.data
            } else {
              if (this.selectedBaseList) {
                this.$toast.warning(
                  'Bara skaparen av en delad lista kan dela en befintlig lista'
                )
                this.form.sharedList = ''
                return
              }
              this.$toast.info('Går med i delad lista')
              this.form.sharedList = resp.data
            }
            this.$nextTick(() => {
              this.sharingList = true
            })
            this.getSharedAbbs()
          })
          .catch((err) => {
            this.$toast.error('Kunde inte gå med i delad lista:', err)
            this.form.sharedList = ''
            this.sharingList = false
          })
        this.sharingList = true
      } else {
        api
          .initSharedList(this.selectedBaseList)
          .then((resp) => {
            this.$toast.info('Skapade delad lista')
            this.form.sharedList = resp.data
            this.$nextTick(() => {
              this.getSharedAbbs()
              this.sharingList = true
              this.$refs.sharedAbbInput.focus()
            })
            this.$store.commit('setSharedList', {
              id: resp.data,
              base: this.selectedBaseList,
            })
          })
          .catch((err) => {
            console.log(err)
            this.$toast.error('Kunde inte skapa delad lista:', err)
            this.sharingList = ''
            this.$store.commit('setSharedList', {})
          })
      }
    },
    onLeaveSharedList(disconnected = false) {
      console.log('this.selectedBaseList:', this.selectedBaseList)
      if (this.selectedBaseList) {
        this.$store.commit('setSharedAbbs', [])
        this.sharingList = false
        this.form.sharedList = ''
        this.sharedAbbs = null
        this.$store.commit('setSharedList', {})
        api
          .leaveSharedList()
          .then(() => {
            this.selectedBaseList = null
          })
          .catch((err) => {
            console.error("Couldn't leave shared list:", err)
          })
        return
      }
      if (this.sharedAbbs != null && this.sharedAbbs.length > 0) {
        this.$bvModal
          .msgBoxConfirm(
            'Vill du spara förkortningarna innan du lämnar den delade listan?',
            {
              okTitle: 'Spara',
              cancelTitle: 'Släng',
            }
          )
          .then((value) => {
            console.log('leave msgbox:', value)
            if (value == null) {
              return
            }
            if (value) {
              this.onSaveSharedAbbs()
            } else {
              this.$store.commit('setSharedAbbs', [])
              this.sharingList = false
              this.form.sharedList = ''
              this.sharedAbbs = null
              api.leaveSharedList().then((resp) => {
                this.$store.commit('setSharedList', {})
              })
            }
          })
      } else {
        this.$store.commit('setSharedList', { id: '', base: '' })
        this.$store.commit('setSharedAbbs', [])
        this.form.sharedList = ''
        this.sharedAbbs = []
        this.sharingList = false
        this.form.sharedList = ''
      }
    },
    onSubmitSharedAbb() {
      let abb = this.form.sharedAbb.trim().replaceAll(' ', '').toLowerCase()
      let word = this.form.sharedWord.trim()
      let obj = {
        base: this.selectedBaseList,
        abb: '',
        word: '',
        create: false,
        override: false,
        delete: false,
        me: true,
      }
      this.show.sharedAbbs = true
      this.lookupPhrase = ''
      this.lookupResults = []
      if (abb != '') {
        if (abb == word) {
          api
            .deleteSharedAbb(this.form.sharedList, abb)
            .then((resp) => {
              this.form.sharedAbb = ''
              this.form.sharedWord = ''
              this.getSharedAbbs()
              obj = {
                ...obj,
                delete: true,
                abb: abb,
                me: true,
              }
              EventBus.$emit('sharedAbb', obj)
              this.$nextTick(() => {
                this.$refs.sharedAbbInput.focus()
              })

              if (obj.base) {
                api
                  .deleteAbb(obj.base, obj)
                  .then((resp) => {
                    console.log('shared abb also deleted user abb:', resp.data)
                  })
                  .catch((err) => {
                    console.error('shared abb couldnt delete user abb', err)
                  })
              }
            })
            .catch((err) => {
              this.$toast.error('Ta bort: Den delade listan existerar inte')
              console.error('couldnt delete shared abb', err)
              this.form.sharedList = ''
            })
        } else {
          obj = {
            ...obj,
            create: true,
            abb: abb,
            word: word,
            me: true,
          }
          if (word == '') {
            obj.word = abb
            obj.create = false
            obj.override = true
          }
          api
            .createSharedAbb(this.form.sharedList, obj)
            .then(() => {
              this.form.sharedAbb = ''
              this.form.sharedWord = ''
              this.getSharedAbbs()
              this.ignoreSuggestion(word)
              this.$nextTick(() => {
                this.$refs.sharedAbbInput.focus()
              })
              EventBus.$emit('sharedAbb', obj)
              if (obj.base) {
                if (obj.word == abb) {
                  return
                }
                api
                  .createAbb(obj.base, obj)
                  .then((resp) => {
                    console.log('shared abb also created user abb:', resp.data)
                    api.cacheAbbs().then(() => {})
                  })
                  .catch((err) => {
                    console.error('shared abb couldnt create user abb', err)
                  })
              }
            })
            .catch((err) => {
              this.$toast.error('Den delade listan existerar inte')
              this.form.sharedList = ''
            })
        }
      }
    },
    onRemoveSharedAbb(abb) {
      api.deleteSharedAbb(this.form.sharedList, abb).then((resp) => {
        this.getSharedAbbs()
        let obj = {
          base: this.selectedBaseList,
          delete: true,
          abb: abb,
          me: true,
        }
        EventBus.$emit('sharedAbb', obj)
        if (obj.base) {
          api
            .deleteAbb(obj.base, obj)
            .then((resp) => {
              console.log('shared abb also deleted user abb:', resp.data)
            })
            .catch((err) => {
              console.error('shared abb couldnt delete user abb', err)
            })
        }
      })
    },
    setCurrentLetter(letter) {
      this.currentLetter = letter
    },

    getSelectedLists() {
      let selected = this.$store.state.settings.selectedLists
      if (selected) {
        console.log('Found selected lists')
        api
          .getLists([selected.standard].concat(selected.addon))
          .then((resp) => {
            this.selectedLists = resp.data.filter((l) => {
              return l.type == 1
            })
          })
          .catch((err) => {
            console.error('support.getselectedLists failed:', err)
          })
      }
    },
    getSharedAbbs() {
      if (this.form.sharedList == '') {
        return
      }
      api.getSharedAbbs(this.form.sharedList).then((resp) => {
        if (resp.data == null || resp.status == '204') {
          this.sharedAbbs = []
          return
        }
        let sorted = resp.data.sort((a, b) => {
          if (a.abb < b.abb) {
            return -1
          }
          if (a.abb > b.abb) {
            return 1
          }
        })
        this.sharedAbbs = sorted

        this.show.sharedAbbs = true
      })
    },
    onSaveSharedAbbs() {
      this.$store.commit('setSharedAbbs', this.sharedAbbs)
      EventBus.$emit('openSettings', 'import')
      EventBus.$emit('importSharedAbbs')
      this.$bvModal.hide('support')
    },
    created() {},
    mounted() {
      this.letters = this.letters.concat(
        '#abcdefghijklmnopqrstuvwxyzåäö'.split('')
      )
      let sharedListId = this.$store.state.sharedList.id
      if (sharedListId != '') {
        api
          .joinSharedList(sharedListId)
          .then((resp) => {
            this.form.sharedList = resp.data
            this.$nextTick(() => {
              this.sharingList = true
            })
            this.getSharedAbbs()
            this.$store.commit('setSharedList', { id: resp.data })
          })
          .catch((err) => {
            this.form.sharedList = ''
            this.$nextTick(() => {
              this.sharingList = false
            })
          })
      } else {
        this.form.sharedList = ''
        this.sharingList = false
      }

      EventBus.$on('lookupPhrase', (phrase) => {
        this.lookupPhrase = phrase.toLowerCase()
        this.lookup()
      })

      EventBus.$on('createdAbb', (word) => {
        this.ignoreSuggestion(word.word)
      })
    },
    beforeDestroy() {
      EventBus.$off('createdAbb')
      EventBus.$off('lookupPhrase')
    },
    lookup() {
      this.show.missedAbbs = false
      this.show.sharedAbbs = false
      this.show.suggestedAbbs = false
      this.lookupPhrase = this.lookupPhrase.replace(/\s+/g, ' ').trim()
      if (this.lookupPhrase !== '') {
        api
          .lookup(this.lookupPhrase)
          .then((resp) => {
            if (resp.status == '204' || resp.data == {}) {
              this.lookupResults = [
                { id: -1, match: 'Hittade inga förkortningar' },
              ]
              return
            }
            var results = []
            var i = 0
            for (const [list, matches] of Object.entries(resp.data)) {
              if (matches.length == 0) {
                this.lookupResults = [
                  {
                    id: -1,
                    match: 'Hittade inga förkortningar',
                  },
                ]
                return
              }
              if (matches.length > 50) {
                this.lookupResults = [{ id: -500, match: 'För många träffar' }]
                return
              }
              matches.forEach((match) => {
                results.push({
                  id: i,
                  match: `${match} (${list})`,
                })
                i += 1
              })
            }
            this.lookupResults = results
            console.log('Lookup results:', this.lookupResults)
          })
          .catch((err) => {
            this.$toast.error('kunde inte slå upp fras', err)
            console.error(err)
            if (err.response) {
              console.log(err.response)
            }
          })
      } else {
        this.lookupResults = []
      }
    },
    stopTracking(abb) {},
    ignoreAllMissedAbbs() {},
    ignoreAllSuggestions() {
      api.ignoreAllSuggestions().then(() => {
        this.getSuggestions()
      })
    },
    ignoreSuggestion(s) {
      api.ignoreSuggestion(s).then((resp) => {
        this.getSuggestions()
        if (resp.data == 'OK') {
          this.show.suggestedAbbs = true
          this.show.sharedAbbs = false
        } else {
          this.show.suggestedAbbs = false
          this.show.sharedAbbs = true
        }
      })
    },
    getSuggestions() {
      api.getSuggestions().then((response) => {
        this.suggestions = response.data.suggestions
      })
    },

    addAbb(word) {
      this.$emit('addAbb', word)
    },
    resetMissedAbbs() {
      this.$store.commit('createMissedAbbsMap')
      this.filterMissedAbbs()
    },
    hideMissed(abb) {
      this.$store.commit('forgetMissedAbb', abb)
      this.filterMissedAbbs()
    },
    dontRemindMissed(abb) {
      api.dontRemind(abb).then((resp) => {
        this.hideMissed(abb)
      })
    },
    preventDefaults(e) {
      if (e.ctrlKey) {
        if (e.key == '1') {
          if (this.collab) {
            e.preventDefault()
            this.$refs.sharedAbbInput.focus()
          }
          return
        }
        if (e.key == '2') {
          e.preventDefault()
          this.show.sharedAbbs = !this.show.sharedAbbs
        }
        if (e.key == '3') {
          e.preventDefault()
          this.show.suggestedAbbs = false
          this.$refs['lookupInput'].focus()
        }
        if (e.key == '4') {
          e.preventDefault()
          this.show.missedAbbs = !this.show.missedAbbs
        }
        if (e.key == '5') {
          e.preventDefault()
          this.lookupPhrase = ''
          this.lookupResults = []
          this.show.suggestedAbbs = !this.show.suggestedAbbs
          this.show.sharedAbbs = false
        }
        console.log('e.key', e.key, 'length:', e.key.length)
        if (e.key.length == 1) {
          let match = e.key.match(/\p{Letter}+/gu)
          if (e.key == '0' || match) {
            let letter = ''
            if (e.key == '0') {
              letter = '#'
            } else {
              letter = match[0]
            }
            e.preventDefault()
            console.log('Switch to show letter:', letter)
            console.log('Current letter:', this.currentLetter)
            console.log('Last letter:', this.lastLetter)
            if (this.lastLetter == letter) {
              letter = 'Alla'
              this.currentLetter = letter
            } else {
              this.currentLetter = letter
            }
            this.lastLetter = letter
          }
        }
      }
      if (e.key == 'F1') {
        e.preventDefault()
      }
    },
  },
}
</script>

<style scoped>
.list-group {
  max-height: 40vh !important;
  overflow-y: scroll;
  -webkit-overflow-scrolling: touch;
}
</style>
