<template>
  <div>
    <h2>Administrera förkortningslistor</h2>
    <div class="d-flex justify-content-between align-items-center">
      <h4>Verktygslåda</h4>
      <b-badge>Visa/dölj med ALT+1</b-badge>
    </div>
    <div v-if="show.toolbox">
      <b-button @click="$bvModal.show('batchImportModal')"
        >Importera till användare</b-button
      >
    </div>
    <div v-else>
      <b-badge>...</b-badge>
    </div>
    <div class="d-flex justify-content-between align-items-center">
      <h4>Alla listor</h4>
      <b-badge>Visa/dölj med ALT+2</b-badge>
    </div>
    <div v-if="show.lists">
      <b-table
        responsive
        fluid
        sticky-header
        small
        striped
        hover
        :items="lists"
        :fields="fields"
        :sort-by="sortBy"
        sort-desc
        @sort-changed="sortChanged"
      >
        <template #cell(type)="row">
          <span v-if="row.item.type">Ämneslista</span
          ><span v-else>Standardlista</span> </template
        ><template #cell(download)="row">
          <b-button @click="downloadList(row.item.id, row.item.owner)">
            <b-icon icon="download" />
          </b-button>
        </template>
      </b-table>
    </div>
    <div v-else>
      <b-badge>...</b-badge>
    </div>

    <b-modal
      id="batchImportModal"
      title="Importera till användare"
      size="lg"
      hide-footer
    >
      <div v-if="batchImport.upload">
        <b-form @submit.prevent="uploadBatchImport">
          <b-form-group label="Användare">
            <b-form-select v-model="batchImport.user" autofocus required>
              <option :value="null">Välj användare...</option>
              <option value="all">Auto</option>
              <option v-for="user in users" :key="user.id" :value="user">
                {{ user.name }}
              </option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="Ladda upp en eller flera listor">
            <b-form-file
              v-model="batchImport.files"
              multiple
              accept=".json"
              required
            />
          </b-form-group>
          <b-button type="submit">Ladda upp</b-button>
        </b-form>
      </div>
      <div v-else>
        <b-form
          @submit.prevent="confirmBatchImport"
          @reset.prevent="resetBatchImport"
        >
          <b-list-group
            v-for="pending in pendingLists"
            :key="pending.creator + '_' + pending.id"
          >
            <b-list-group-item>
              <div class="d-flex justify-content-between align-items-center">
                <span>{{ pending.creator_name }}</span
                ><span>{{ pending.name }}</span>
                <span v-if="!pending.clean">
                  <b-badge variant="success">{{ pending.new }}</b-badge>
                  <b-badge variant="danger">{{
                    pending.conflicts.length
                  }}</b-badge>
                  <b-button
                    variant="warning"
                    @click="confirmSingleImport(pending)"
                    >Importera</b-button
                  >
                </span>
                <span v-else>
                  <b-badge variant="success">{{ pending.new }}</b-badge>
                  <b-button @click="confirmSingleImport(pending)"
                    >Skapa</b-button
                  >
                </span>
              </div>
            </b-list-group-item>
          </b-list-group>
          <b-button type="reset">Gå tillbaka</b-button>
          <b-button type="submit">Importera allt</b-button>
        </b-form>
      </div>
    </b-modal>
  </div>
</template>
<script>
import api from '../../../api/api.js'
import axios from 'axios'
export default {
  name: 'AdminListView',
  props: ['users'],
  data() {
    return {
      show: {
        toolbox: false,
        lists: false,
      },
      batchImport: {
        upload: true,
        creators: null,
        user: 'all',
        files: [],
        lists: null,
        conflicts: [],
      },
      pendingLists: [],
      lists: [],
      fields: [
        { key: 'name', label: 'Namn', sortable: true },
        { key: 'type', label: 'Typ' },
        { key: 'owner', label: 'Ägare', sortable: true },
        { key: 'counter', label: 'Antal', sortable: true },
        { download: { label: 'Ladda ner' } },
      ],
      sortBy: 'email',
    }
  },
  mounted() {
    this.batchImport.creators = new Map()
    document.addEventListener('keydown', this.preventDefaults)
    this.users = this.users.sort((a, b) => {
      if (a.name < b.name) {
        return -1
      }
      if (a.name > b.name) {
        return 1
      }
      return 0
    })
    this.users.forEach((u) => {
      this.getLists(u)
    })
  },
  methods: {
    logger() {},
    sortChanged() {},
    preventDefaults(e) {
      if (e.altKey) {
        if (e.key == '1') {
          e.preventDefault()
          this.show.toolbox = !this.show.toolbox
          return
        }
        if (e.key == '2') {
          e.preventDefault()
          this.show.lists = !this.show.lists
        }
      }
    },
    showImportToUserModal() {
      this.$bvModal.show('batchImportModal')
    },
    resetBatchImport() {
      this.batchImport.upload = true
      this.pendingLists = []
    },
    uploadBatchImport() {
      this.pendingLists = []
      this.batchImport.lists = new Map()
      if (this.batchImport.files.length <= 0) return
      if (this.batchImport.user == null) return
      this.batchImport.upload = false
      this.batchImport.files.forEach((file) => {
        const fr = new FileReader()
        fr.onload = (e) => {
          let pending = JSON.parse(e.target.result)
          pending.filename = file.name
          this.batchImport.creators.set(
            pending.meta.creator,
            this.users.find((u) => u.id == pending.meta.creator)
          )
          if (pending.abbs) {
            if (pending.abbs.length > 0) {
              api
                .getList(pending.meta.id)
                .then((resp) => {
                  if (resp.status != 200) {
                    const pendingListInfo = {
                      creator: pending.meta.creator,
                      creator_name: this.getCreator(
                        pending.meta.creator,
                        pending.filename
                      ),
                      name: pending.meta.name,
                      filename: pending.filename,
                      id: pending.meta.id,
                      clean: true,
                      new: pending.meta.counter,
                      conflicts: [],
                    }
                    this.batchImport.lists.set(pending.meta.id, pending)
                    this.pendingLists.push(pendingListInfo)
                  } else if (resp.data != null) {
                    this.detectConflicts(resp.data, pending)
                  }
                })
                .catch((err) => {
                  console.error("couldn't get current list", err)
                })
            }
          }
        }
        fr.readAsText(file)
      })
    },
    detectConflicts(curr, pending) {
      let currAbbs = new Map()
      let conflicts = []
      api
        .getAbbs(curr.id)
        .then((resp) => {
          resp.data.map((abb) => {
            currAbbs.set(abb.abb, abb)
          })
          let newAbbs = pending.abbs.length - resp.data.length
          pending.abbs.map((newAbb) => {
            const currAbb = currAbbs.get(newAbb.abb)
            if (currAbb) {
              if (currAbb.word != newAbb.word) {
                conflicts.push({
                  abb: newAbb.abb,
                  old: currAbb.word,
                  new: newAbb.word,
                })
              }
            }
          })

          if (pending.abbs.length != resp.data.length && conflicts.length > 0) {
            let creator_name = this.getCreator(
              pending.meta.creator,
              pending.filename
            )
            const pendingListInfo = {
              name: curr.name,
              id: curr.id,
              creator: curr.creator,
              creator_name: creator_name,
              filename: pending.filename,
              clean: false,
              conflicts: conflicts,
              new: newAbbs,
            }
            this.batchImport.lists.set(curr.id, pending)
            this.pendingLists.push(pendingListInfo)
          }
        })
        .catch((err) => {
          console.error(
            'detectConflicts failed after getting current abbs',
            err
          )
        })
    },
    getCreator(c_id, filename) {
      let match = this.users.find((u) => u.id == c_id)
      if (!match) {
        let name = filename.split('_')[0]
        match = this.users.find((u) => u.name == name)
        if (!match) {
          console.error('filename belongs to wrong id:', c_id, name)
          return 'id mismatch'
        }
      }
      return match.name
    },
    confirmSingleImport(pendingList) {
      const list = this.batchImport.lists.get(pendingList.id)
      if (pendingList.creator_name == 'id mismatch') return
      if (pendingList.clean) {
        list.meta.creator = pendingList.creator
        api
          .createList(list.meta)
          .then((resp) => {
            let created = resp.data
            api
              .importList(created.id, list.abbs)
              .then((resp) => {
                let i = this.pendingLists.findIndex((l) => l.id == list.meta.id)
                this.pendingLists.splice(i, 1)
                this.batchImport.lists.delete(list.meta.id)
              })
              .catch((err) => {
                console.error('couldnt import to created list:', err)
              })
          })
          .catch((err) => {
            console.error('couldnt create list before importing:', err)
          })
        return
      } else {
        api
          .importList(list.meta.id, list.abbs)
          .then((resp) => {
            let i = this.pendingLists.findIndex((l) => l.id == list.meta.id)
            this.pendingLists.splice(i, 1)
            this.batchImport.lists.delete(list.meta.id)
          })
          .catch((err) => {
            console.error("couldn't import list...", err)
          })
      }
    },
    confirmBatchImport() {
      this.pendingList.map((pending) => {
        this.confirmSingleImport(pending)
      })
    },
    downloadList(id, owner) {
      let exportedList = { meta: null, abbs: null }
      api.getList(id).then((resp) => {
        exportedList.meta = resp.data
        api.getAbbs(id).then((resp) => {
          exportedList.abbs = resp.data
          var blob = new Blob([JSON.stringify(exportedList)], {
            type: 'application/json',
          })
          let link = document.createElement('a')
          link.href = window.URL.createObjectURL(blob)
          link.download = owner + '_' + exportedList.meta.name + '.json'
          link.click()
        })
      })
    },
    getLists(user) {
      axios
        .get('https://sttylus.se/api2/abbs/lists', {
          headers: { 'X-Id-Token': user.id },
        })
        .then((resp) => {
          if (resp.data != null) {
            resp.data.forEach((list) => {
              list.owner = user.name
              this.lists.push(list)
            })
          }
        })
        .catch((err) => {
          this.$toast.warning('kunde inte hämta användarlistor:', err)
          console.error('getLists failed:', err)
        })
    },
  },
}
</script>
