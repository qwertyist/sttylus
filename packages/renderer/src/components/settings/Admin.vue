<template>
  <div>
    <template v-if="usersLoaded">
      <AdminListView lazy :users="users" />
    </template>
    <b-container fluid>
      <b-row>
        <b-col>
          <h4>Användare</h4>
        </b-col>
        <b-col>
          <b-button
            class="float-right"
            size="sm"
            variant="primary"
            @click="createUser"
          >
            Lägg till användare
          </b-button>
        </b-col>
      </b-row>
      <b-table
        fluid
        sticky-header
        striped
        hover
        :items="users"
        :fields="userFields"
      >
        <template #cell(role)="row">
          <template v-if="editing != row.item.id">
            <div
              @click="changeRole(row.item)"
              style="
                word-wrap: break-word;
                min-width: 240px;
                max-width: 240px;
                white-space: normal;
              "
            >
              {{ row.item.role }}
            </div>
          </template>
          <b-form-select
            v-else
            v-model="row.item.role"
            :options="roles"
            @change="onUpdateRole(row.item)"
          />
        </template>
        <template #cell(updated)="row">
          {{ row.item.updated | formatDate }}
        </template>
        <template #cell(reset)="row">
          <b-button
            v-model="row.item.reset"
            type="submit"
            size="sm"
            variant="warning"
            class="mt-0"
            @click="resetUserPassword(row.item)"
          >
            <b-icon-key />
          </b-button>
        </template>
        <template #cell(remove)="row">
          <b-button
            v-model="row.item.remove"
            type="submit"
            size="sm"
            variant="danger"
            class="mt-0"
            @click="deleteUser(row.item)"
          >
            <b-icon-trash />
          </b-button>
        </template>
      </b-table>
    </b-container>
    <CreateUser />
  </div>
</template>

<script>
import api from '../../api/api.js'
import EventBus from '../../eventbus.js'

import CreateUser from '../modals/CreateUser.vue'
import AdminListView from './admin/ListView.vue'

export default {
  name: 'Admin',
  components: { CreateUser, AdminListView },
  data() {
    return {
      show: false,
      userFields: [
        { key: 'email', label: 'E-post' },
        { key: 'role', label: 'Roll' },
        { key: 'updated', label: 'Senaste uppdatering' },
        { key: 'id', label: 'ID' },
        { reset: { label: 'Återställ lösenord' } },
        { remove: { label: 'Ta bort' } },
      ],
      users: [],
      errors: [],
      sortBy: 'email',
      sortDesc: false,
      editing: '',
      deleted: '',
      roles: [
        { value: 'user', text: 'Skrivtolk' },
        { value: 'tester', text: 'Testanvändare' },
        { value: 'admin', text: 'Administratörkonto' },
      ],
    }
  },
  computed: {
    usersLoaded() {
      return this.users.length > 0
    },
  },
  mounted() {
    this.getUsers()
    EventBus.$on('openAdminView', () => {
      console.log('open admin view')
      this.show = true
    })
    EventBus.$on('createdUser', () => {
      this.getUsers(true)
      this.sortBy = 'created'
      this.sortDesc = true
    })
  },
  beforeDestroy() {
    EventBus.$off('createdUser')
    EventBus.$off('openAdminView')
  },
  methods: {
    getUsers() {
      api
        .getUsers()
        .then((resp) => {
          this.users = resp.data
        })
        .catch((err) => {
          this.$toast.warning('kunde inte hämta användarna:' + err)
        })
    },
    createUser() {
      this.$bvModal.show('createUser')
    },
    resetUserPassword(user) {
      this.$bvModal
        .msgBoxConfirm(
          'Vill du nollställa lösenordet för användaren ' + user.email + '?',
          {
            title: 'Nollställ lösenord',
            size: 'md',
            buttonSize: 'md',
            okVariant: 'warning',
            okTitle: 'Nollställ',
            cancelTitle: 'Avbryt',
            footerClass: 'p-2',
            hideHeaderClose: false,
            centered: true,
          }
        )
        .then((value) => {
          if (value == true) {
            api
              .resetUserPassword(user)
              .then(() => {
                this.$toast.info('Lösenordet nollställdes')
                this.getUsers()
              })
              .catch((err) => {
                this.$toast.warning(
                  'kunde inte nollställa användarens lösenord:',
                  err
                )
                console.log('resetUserPassword failed: ', err)
              })
          }
        })
        .catch((err) => {
          console.log('failed in modal', err)
        })
    },
    deleteUser(user) {
      this.deleted = ''
      this.$bvModal
        .msgBoxConfirm(
          'Du är på väg att ta bort ' + user.email + '. Är du säker?',
          {
            title: 'Bekräfta borttagning',
            size: 'md',
            buttonSize: 'md',
            okVariant: 'danger',
            okTitle: 'Ta bort',
            cancelTitle: 'Avbryt',
            footerClass: 'p-2',
            hideHeaderClose: false,
            centered: true,
          }
        )
        .then((value) => {
          if (value == true) {
            api
              .deleteUser(user)
              .then((resp) => {
                this.deleted = resp.data
                this.getUsers()
              })
              .catch((err) => {
                this.$toast.warning('kunde inte ta bort användaren:', err)
                console.log('deleteUser failed api call', err)
              })
          }
        })
        .catch((err) => {
          console.log('failed in modal', err)
        })
    },
    changeRole(user) {
      this.editing = user.id
    },
    onUpdateRole(user) {
      api.updateUser(user).then(() => {
        this.getUsers()
      })
    },
  },
}
</script>

<style scoped></style>
