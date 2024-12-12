<template>
  <b-modal
    id="addAbb"
    size="lg"
    :modal-class="addModalClass"
    hide-footer
    hide-header
    hide-backdrop
    no-fade
    @show="showModal"
    @hide="closeModal"
    ref="addabbmodal"
  >
    <b-row>
      <b-form
        inline
        @keydown.enter="onSubmit"
        @keydown.113="nextList"
        autocomplete="off"
      >
        <b-form-input
          v-model="form.abb"
          id="abb"
          size="lg"
          class="mb-2 mb-sm-0"
          placeholder="Förkortning"
          autofocus
          ref="abbInput"
          :state="exists"
          @input="checkExists"
        />

        <b-form-datalist id="abb" />
        <b-form-input
          v-model="form.word"
          id="word"
          size="lg"
          class="mb-2 mr-sm-2 mb-sm-0"
          placeholder="Text"
          ref="wordInput"
        />
        <b-form-select
          v-model="form.selected"
          size="lg"
          style="width: 285px"
          ref="listSelect"
          disabled
        >
          <option
            :style="listColorID"
            v-for="list in form.lists"
            :value="list.id"
            :key="list.id"
          >
            {{ list.name }}
          </option>
        </b-form-select>
        <b-button type="submit" style="display: none" />
      </b-form>
      <div style="position: relative">
        <span v-if="existing" class="text-danger">
          <small>existerar som: {{ existing }}</small>
        </span>
        <span
          v-if="form.lists.length > 1"
          class="float-right"
          style="position: relative; right: 10px"
        >
          <small>
            Växla lista med
            <kbd>F2</kbd>
          </small>
        </span>
      </div>
    </b-row>
  </b-modal>
</template>
<script>
import api from '../../api/api.js'
import db from '../../store/db.js'
import EventBus from '../../eventbus.js'

export default {
  name: 'AddAbbreviation',
  props: ['query'],
  data() {
    return {
      addModalClass: ['addModalClass'],
      currentLists: {},
      form: {
        abb: '',
        word: '',
        lists: [],
        selected: '',
        index: 0,
      },
      exists: null,
      existing: '',
      secondary: false,
    }
  },
  computed: {
    listColorID() {
      return { color: 'red' }
    },
    sharedList() {
      return this.$store.state.sharedList
    },
  },
  methods: {
    catchEscape() {
      this.$toast.info('haha')
    },
    handleKeydown(e) {
      console.log(e)
    },
    checkExists() {
      if (this.form.abb !== '') {
        if (this.form.selected) {
          api.getAbb(this.form.selected, this.form.abb.trim()).then((resp) => {
            if (resp.data) {
              this.exists = false
              this.existing = resp.data.word
            } else {
              this.exists = null
              this.existing = ''
            }
          })
        }
      }
    },
    handleHotkeys(e) {
      if (e.ctrlKey && this.view != 'tabula') {
        if (['1', '2', '3', '4', '5'].indexOf(e.key) != -1) {
          EventBus.$emit('changeStandardList', e.key)
          e.preventDefault()
        }
      }
    },
    showModal() {
      document.addEventListener('keyup', (e) => {
        if (e.key == 'Escape') {
          this.$bvModal.hide('support')
        }
      })
      this.$store.commit('setModalOpen', true)
      EventBus.$emit('modalOpened')
      this.updateLists()
      this.form.word = this.$store.state.selectedWord
      this.form.selected = this.$store.state.targetList.id
      this.index = this.$store.state.targetList.index
    },
    closeModal() {
      document.removeEventListener('keydown', (e) => this.handleHotkeys(e))
      document.removeEventListener('keyup', (e) => {
        if (e.key == 'Escape') {
          this.$bvModal.hide('support')
        }
      })
      this.$store.commit('setModalOpen', false)
      this.$bvModal.hide('addAbb')
      EventBus.$emit('closeNav')
      EventBus.$emit('abbModalClosed')

      this.form.abb = ''
      this.form.word = ''
      this.exists = null
      this.existing = ''
    },
    updateLists() {
      this.currentLists = this.$store.state.settings.selectedLists
      api
        .getLists([this.currentLists.standard].concat(this.currentLists.addon))
        .then((resp) => {
          if (resp.data == null || resp.data == 'null') {
            console.log('addAbb updateLists no result', resp.data)
            this.form.lists = []
            return
          }
          this.form.lists = resp.data
        })
        .catch((err) => {
          console.error('Couldnt get addAbb lists:', err)
        })
    },
    onSubmit(e) {
      e.preventDefault()
      let targetListId = null
      if (this.form.abb === '' || this.form.word === '') {
        this.$store.commit('setSelectedWord', '')
        this.$bvModal.hide('addAbb')
        return
      }

      if (this.form.selected.length === 0) {
        targetListId = this.$store.state.selectedLists.standard.id
      } else {
        targetListId = this.form.selected
      }
      let data = {
        abb: this.form.abb.trim(),
        word: this.form.word.trim(),
        creator: this.$store.state.userData.id,
        targetListId,
      }

      this.$bvModal.hide('addAbb')

      var deleted = false
      if (data.abb === data.word) {
        db.deleteAbb(data.abb, targetListId)
        deleted = true

      }
      api
        .createAbb(targetListId, {
          abb: data.abb,
          word: data.word,
          creator: data.creator,
        })
        .then((resp) => {
          console.log('hello')
          EventBus.$emit('createdAbb', data)
          if (!deleted) db.addAbb(resp.data, targetListId)
          this.$store.commit('setSelectedWord', '')
          this.$bvModal.hide('addAbb')
          if (e.shiftKey) {
            this.$store.commit('setSelectedWord', data.word)
            this.$bvModal.show('addAbb')
          }
        })
        .catch((err) => {
          console.log('Failed creating abbreviation:', err)
          this.$bvModal.hide('addAbb')
        })
    },
    nextList() {
      if (this.form.lists.length < 2) return
      if (this.index < this.form.lists.length - 1) {
        this.index++
        this.form.selected = this.form.lists[this.index].id
      } else {
        this.index = 0
        this.form.selected = this.form.lists[this.index].id
      }
      let targetList = {
        index: this.index,
        id: this.form.selected,
      }
      this.$store.commit('setTargetList', targetList)
      this.checkExists()
    },
  },

  mounted() {
    EventBus.$on('updatedSelectedLists', this.updateLists)
    EventBus.$on('showTextView', () => {
      this.$bvModal.hide('addAbb')
      this.$bvModal.hide('support')
    })
    this.currentLists = this.$store.state.settings.selectedLists
    api
      .getLists([this.currentLists.standard].concat(this.currentLists.addon))
      .then((resp) => {
        if (resp.data == null) {
          this.form.lists = []
          return
        }
        this.form.lists = resp.data
      })
      .catch((err) => {
        console.error('addabb failed to get lists after mounted:', err)
      })
  },
  beforeDestroy() {
    EventBus.$off('showTextView')
    EventBus.$off('updatedSelectedLists')
  },
}
</script>

<style scope>
.addModalClass > div {
  position: absolute;
  height: 80px;

  margin-left: -470px;
  background: #063;
  bottom: 25px;
  left: 50%;
  border-radius: none;
}

.addModalClass > .modal-dialog > .modal-content {
  border-radius: none !important;
  background-color: white !important;
  width: 940px !important;
  height: 120px;
  padding: 1em;
}

.form-control.is-invalid,
.was-validated .form-control:invalid {
  padding-right: 1rem !important;
}
</style>
