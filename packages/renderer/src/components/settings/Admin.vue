<template>
  <div>
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
        <template #cell(remove)="row">
          <!-- <b-button
              v-model="row.item.commend"
              type="submit"
              size="md"
              class="mt-0"
              variant="primary"
              @click="commendAbb(row.item, viewedList.id)"
            >
              <b-icon-arrow-up />
            </b-button>-->
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
      <h4>Förkortningslistor</h4>
      <b-table
        responsive
        fluid
        sticky-header
        small
        striped
        hover
        :items="lists"
        :fields="listFields"
        :sort-by="sortListBy"
        sort-desc
        @sort-changed="listSortChanged"
      >
        <template #cell(type)="row">
          <span v-if="row.item.type">Ämneslista</span><span v-else>Standardlista</span>
        </template><template #cell(download)="row">
          <b-button @click="downloadList(row.item.id, row.item.owner)">
            <b-icon icon="download" />
          </b-button>
        </template>
      </b-table>
    </b-container>

    <CreateUser />
  </div>
</template>

<script>
import api from "../../api/api.js";
import EventBus from "../../eventbus.js";

import CreateUser from "../modals/CreateUser.vue";
export default {
  name: "Admin",
  components: { CreateUser },
  data() {
    return {
      show: false,
      userFields: [
        { key: "email", label: "E-post" },
        { key: "role", label: "Roll" },
        { key: "updated", label: "Senaste uppdatering" },
        { key: "id", label: "ID" },
        { remove: { label: "Ta bort" } },
      ],
      users: [],
      listFields: [
        { key: "name", label: "Namn", sortable: true },
        { key: "type", label: "Typ" },
        { key: "owner", label: "Ägare", sortable: true },
        { key: "counter", label: "Antal förkortningar", sortable: true },
        { download: { label: "Ladda ner" } },
      ],
      lists: [],
      sortListBy: "counter",
      errors: [],
      sortBy: "email",
      sortDesc: false,
      editing: "",
      deleted: "",
      roles: [
        { value: "user", text: "Skrivtolk" },
        { value: "tester", text: "Testanvändare" },
        { value: "admin", text: "Administratörkonto" },
      ],
    };
  },
  methods: {
    getUsers(asc) {
      api
        .getUsers()
        .then((resp) => {
          this.users = resp.data;
          this.users.map((u) => {
            this.getLists(u);
          });
        })
        .catch((err) => {
          this.errors = err;
        });
    },
    createUser() {
      this.$bvModal.show("createUser");
    },
    deleteUser(user) {
      this.deleted = "";
      this.$bvModal
        .msgBoxConfirm(
          "Du är på väg att ta bort " + user.email + ". Är du säker?",
          {
            title: "Bekräfta borttagning",
            size: "md",
            buttonSize: "md",
            okVariant: "danger",
            okTitle: "Ta bort",
            cancelTitle: "Avbryt",
            footerClass: "p-2",
            hideHeaderClose: false,
            centered: true,
          }
        )
        .then((value) => {
          if (value == true) {
            api
              .deleteUser(user)
              .then((resp) => {
                this.deleted = resp.data;
                this.getUsers();
              })
              .catch((err) => {
                console.log("failed in api call", err);
              });
          }
        })
        .catch((err) => {
          console.log("failed in modal", err);
        });
    },
    changeRole(user) {
      this.editing = user.id;
    },
    onUpdateRole(user) {
      api.updateUser(user).then(() => {
        this.getUsers();
      });
    },
    getLists(user) {
      api.getUserListsByID(user.id).then((resp) => {
        if (resp.data != null) {
          resp.data.map((list) => {
            list.owner = user.name;
            this.lists.push(list);
          });
        }
      });
    },
    listSortChanged(sort) {
    },
    downloadList(id, owner) {
      let exportedList = { meta: null, abbs: null };
      api.getList(id).then((resp) => {
        exportedList.meta = resp.data;
        api.getAbbs(id).then((resp) => {
          exportedList.abbs = resp.data;
          var blob = new Blob([JSON.stringify(exportedList)], {
            type: "application/json",
          });
          let link = document.createElement("a");
          link.href = window.URL.createObjectURL(blob);
          link.download = owner + "_" + exportedList.meta.name + ".json";
          link.click();
        });
      });
    },
  },
  mounted() {
    this.getUsers();
    EventBus.$on("openAdminView", () => {
      console.log("open admin view");
      this.show = true;
    });
    EventBus.$on("createdUser", () => {
      this.getUsers(true);
      this.sortBy = "created";
      this.sortDesc = true;
    });
  },
  beforeDestroy() {
    EventBus.$off("createdUser");
    EventBus.$off("openAdminView");
  },
};
</script>

<style scoped>
</style>
