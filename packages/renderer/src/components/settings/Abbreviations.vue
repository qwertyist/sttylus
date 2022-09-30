<template>
  <div>
    <b-row>
      <b-col
        cols="4"
        style="height: 85vh !important; overflow-y: auto"
      >
        <h3 class="float-left">
          Dina listor
        </h3>
        <div v-if="!tester">
          <b-button
            size="sm"
            class="float-right"
            variant="primary"
            @click="addList"
          >
            Lägg till lista
          </b-button>
        </div>
        <b-table
          fixed
          @row-clicked="
            (item) => {
              viewList(item);
            }
          "
          hover
          :items="standardLists"
          :fields="standardlistFields"
        >
          <template #cell(name)="data">
            <span
              :class="[
                viewedList.id == data.item.id ? 'active_row' : '',
                'onHover',
              ]"
            >{{ data.item.name }}</span>
          </template>
          <template #cell(counter)="data">
            <span class="pl-3">{{ data.item.counter }}</span>
          </template>
          <template #cell(checkbox)="data">
            <span v-if="data.item.temp"><em class="float-right">Alltid vald</em></span><span v-else><b-form-radio
              class="float-right"
              v-model="selectedStandard"
              :value="data.item.id"
            /></span>
          </template>
        </b-table>
        <b-table
          fixed
          @row-clicked="
            (item) => {
              viewList(item);
            }
          "
          hover
          :items="addonLists"
          :fields="addonlistFields"
        >
          <template #cell(name)="data">
            <span
              :class="[
                viewedList.id == data.item.id ? 'active_row' : '',
                'onHover',
              ]"
            >{{ data.item.name }}</span>
          </template>
          <template #cell(counter)="data">
            <span class="pl-3">{{ data.item.counter }}</span>
          </template>
          <template #cell(checkbox)="data">
            <b-form-checkbox
              class="float-right"
              v-model="selectedAddons"
              :value="data.item.id"
              @change="toggleSelectAddon($event)"
            />
          </template>
        </b-table>
      </b-col>
      <b-col>
        <!--   <b-tabs v-model="tabIndex" pill contentClass="mt-3" @input="changeTab">
          <b-tab v-for="l in lists" :title="l.name" :key="'dyn-list' + l.id"> -->
        <b-row>
          <b-col>
            <div class="float-left">
              <h2 class="mb-0">
                {{ viewedList.name }}
              </h2>
              <div v-if="!tester">
                <b-button
                  type="submit"
                  size="sm"
                  variant="primary"
                  class="m1-0 mb-1 text-right"
                  @click="editList(viewedList.id)"
                >
                  Byt namn/typ
                </b-button>
                <b-button
                  v-if="!viewedList.temp"
                  type="submit"
                  size="sm"
                  variant="danger"
                  class="ml-0 mb-1 text-right"
                  @click="removeList(viewedList.id)"
                >
                  Ta bort lista
                  <b-icon-trash />
                </b-button>
              </div>
            </div>
            <div class="float-right">
              <b-form-input
                size="sm"
                v-model="filter"
                placeholder="Sök på förkortning eller fras..."
                filter-debounce="200"
              />
              <b-form-input
                v-model="currentPage"
                style="max-width: 55px"
                class="float-left"
                size="sm"
              /><b-pagination
                size="sm"
                class="float-right"
                v-model="currentPage"
                :total-rows="abbs.length"
                :per-page="perPage"
                aria-controls="abbs"
              />
            </div>
          </b-col>
        </b-row>
        <div
          style="height: 75vh !important; overflow-y: auto;"
          @scroll="onScrollAbbs"
        >
          <b-table
            id="abbs"
            striped
            hover
            small
            responsive
            :items="abbs"
            :fields="abbFields"
            :filter="filter"
            :filter-function="searchAbb"
            :per-page="perPage"
            :current-page="currentPage"
            :sort-by="sortBy"
            sort-desc
            @sort-changed="sortChanged"
          >
            <template #cell(abb)="row">
              <div>
                {{ row.item.abb }}
              </div>
            </template>
            <template #cell(word)="row">
              <template v-if="editing != row.item.id">
                <div
                  @click="editAbb(row.item)"
                  style="
                    word-wrap: break-word;
                    min-width: 220px;
                    white-space: normal;
                  "
                >
                  {{ row.item.word }}
                </div>
              </template>
              <b-form-input
                v-else
                v-model="row.item.word"
                @change="onUpdateAbb(row.item, viewedList.id)"
              />
              <!--     <div @click="editAbb(row.item)" v-else style="word-wrap: break-word;min-width: 240px;max-width: 240px;white-space:normal;">{{ row.item.word }}</div>
        -->
            </template>
            <template #cell(updated)="row">
              <div
                v-if="viewedList.id != 'temp'"
                class=""
                style="min-width: 60px"
              >
                {{ row.item.updated | formatDate }}
              </div>
            </template>
            <template
              v-if="viewedList.type < 2"
              #cell(remove)="row"
            >
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
                @click="removeAbb(row.item, viewedList.id)"
              >
                <b-icon-trash />
              </b-button>
            </template>
          </b-table>
        </div>
        <!--        <b-row>
          <b-col class="text-right"> </b-col>
        </b-row>
        -->
      </b-col>
    </b-row>
    <AddList />
    <RemoveList
      :list="viewedList"
      :abbs="abbs"
    />
    <EditList :list="viewedList" />
  </div>
</template>
<script>
import EventBus from "../../eventbus.js";
import AddList from "../modals/AddList.vue";
import EditList from "../modals/EditList.vue";
import RemoveList from "../modals/RemoveList.vue";

import api from "../../api/api.js";
export default {
  name: "Abbreviations",
  components: { AddList, RemoveList, EditList },
  data() {
    return {
      listform: {
        name: null,
        standard: false,
      },
      abbform: {},
      editing: false,
      tabIndex: 0,
      currentPage: 1,
      perPage: 20,
      lists: [],
      standardLists: [],
      addonLists: [],
      addons: [],
      viewedList: "",
      viewedLists: [],
      filter: "",
      selectedAddons: [],
      templistFields: [
        { key: "name", label: "Slasklista" },
        { key: "counter", label: "Förkortningar" },

        { key: "save", label: "Spara" },
      ],
      standardlistFields: [
        { key: "name", label: "Standardlistor" },
        { key: "counter", label: "Förkortningar" },

        { key: "checkbox", label: "Välj en" },
      ],
      addonlistFields: [
        { key: "name", label: "Ämneslistor" },
        { key: "counter", label: "Förkortningar" },
        { key: "checkbox", label: "Välj en/flera", tdClass: "align-right" },
      ],
      abbFields: [
        { abb: { label: "Förkortning", sortable: true, tdClass: "mt-2" } },
        {
          word: {
            label: "Text (Tryck på frasen för att redigera)",
            sortable: true,
            tdClass: "mt-2 word",
          },
        },
        { updated: { label: "Senast ändrad", sortable: true } },

        { remove: { label: "Ta bort" } },
      ],
      sortBy: "updated",
      abbs: [],
      nominations: [],
      tempList: { name: "Temporär lista", temp: true, id: "temp" },
      tempAbbs: [],
    };
  },
  computed: {
    tester() {
      if (this.$store.state.userData.role == "tester") {
        return true;
      }
      return false;
    },
    selectedStandard: {
      get() {
        return this.$store.state.settings.selectedLists.standard;
      },
      set(list) {
        this.$store.commit("setSelectedStandard", list);
        return list;
      },
    },
  },
  methods: {
    toggleSelectAddon(listIDs) {
      console.log("set selected addon in view:", this.selectedAddons.length)
      this.selectedAddons = listIDs
      this.$store.commit("setSelectedAddons", this.selectedAddons);
    },
    listAbbCount(list) {
      return 0;
    },
    addList() {
      this.$bvModal.show("addList");
    },
    editList() {
      this.$bvModal.show("editList");
    },
    viewList(list) {
      console.log("view list:", list)
      this.perPage = 25;
      if (!list) {
        this.viewedList = { name: "Ingen lista vald", id: undefined };
      } else {
        if (list.name == "") {
          console.log("hitta lista via id")
        }
        this.viewedList = list;
        if (
          [this.selectedStandard]
            .concat(this.selectedAddons)
            .indexOf(list.id) !== -1
        ) {
          let targetList = {
            id: list.id,
            index: [this.selectedStandard]
              .concat(this.selectedAddons)
              .indexOf(list.id),
          };
          this.$store.commit("setTargetList", targetList);
        } else {
          this.$toast.warning(
            "Kom ihåg att välja listan om du vill lägga till förkortningar i den.", { duration: 2000 }
          );
        }
        if (this.viewedList.id !== undefined) {
          this.getAbbs(list.id);
        }
      }
      //console.log("this.viewedList:", this.viewedList);
    },
    sortChanged(value) {
      this.sortBy = value.sortBy;
    },
    removeList() {
      this.$bvModal.show("removeList");
    },
    searchAbb(data, query) {
      if (query.length === 0) {
        return true;
      }
      return (
        data.abb.toLowerCase().includes(query.toLowerCase()) ||
        data.word.toLowerCase().includes(query.toLowerCase())
      );
    },
    afterRemovedList(list) {
      this.getLists(true);
      if (this.viewedList.id == list.id) {
        this.viewedList = "temp";
      }
      this.$store.commit("unsubscribeList", list.id);
      this.$store.commit("unsetSelectedAddon", list);
      this.flushSelectedLists(list);
    },
    afterUpdatedList(list) {
      this.flushSelectedLists(list);
    },
    flushSelectedLists(list) {
      this.$store.commit("unsetSelectedAddon", list);
      this.toggleSelectAddon();
      this.getLists();
    },
    editAbb(abb) {
      this.editing = abb.id;
    },
    openRemoveListModal(listID) {
      api
        .getList(listID)
        .then((list) => {
          api
            .getAbbs(listID)
            .then((abbs) => {
              let n = 0;
              if (abbs.data !== null) {
                n = abbs.data.length;
              }
              const action =
                list.data.type < 3 ? "Ta bort lista" : "Dölj lista";
              this.removelist = {
                action: action,
                id: list.data.id,
                name: list.data.name,
                type: list.data.type,
                abbs: abbs.data,
                n: n,
              };
            })
            .catch((err) => {
              alert(err);
            });
          this.$refs["remove-list-modal"].show();
        })
        .catch((err) => {
          alert(err);
        });
    },
    hideRemoveListModal() {
      this.$refs["remove-list-modal"].hide();
    },
    onSubmitRemove(evt) {
      evt.preventDefault();
      if (this.removelist.type < 3) {
        api
          .deleteList(this.removeList.id)
          .then((response) => {
            this.$refs["remove-list-modal"].hide();
            this.getLists(true);
            this.removelist = {
              name: null,
              type: null,
              n: 0,
              abbs: [],
            };
          })
          .catch((err) => {
            console.log("something went wrong when removing list:", err);
            alert(err);
          });
      } else {
        console.log("Göm lista");
      }
    },
    onResetRemove() {},
    onUpdateAbb(abb, listId) {
      var currentPage = this.currentPage;
      abb.word = abb.word.trim();
      api
        .updateAbb(listId, abb)
        .then(() => {
          if (this.sortBy !== "updated") {
            this.getAbbs(listId);
            this.currentPage = currentPage;
          }
        })
        .catch((err) => {
          console.log("failed updating abb", err);
        });
    },
    removeAbb(abb, listID) {
      api.deleteAbb(listID, abb).then((response) => {
        this.getAbbs(listID);
        this.getLists(false);
      });
    },
    commendAbb(commended) {
      let nominations = this.$store.state.stagedNominations;
      if (nominations.indexOf(commended) == -1) {
        nominations.push(commended);
        this.$emit("commend");
      } else {
        nominations = nominations.filter(function (abb) {
          return abb !== commended;
        });
        this.$emit("decommend");
      }
      this.$store.commit("setStagedNominations", nominations);
    },
    createAndSelectUserLists() {
      /*
      api
        .createList({
          name: "Standardlista",
          type: 0,
        })
        .then((resp) => {
          this.$store.commit("subscribeList", resp.data.id);
        });
        */
    },
    getLists(view) {
      api.getUserLists().then((response) => {
        if (response.data == null) {
          return;
        }
        if (response.data !== null) {
          this.lists = response.data;
          this.standardLists = response.data.filter((l) => {
            if (l.type == 0) {
              return l;
            }
            if (l.type == 2) {
              l.name = "🌎 " + l.name;
              return l;
            }
          });

          this.addonLists = response.data.filter((l) => {
            if (l.type == 1) {
              return l;
            }
            if (l.type == 3) {
              l.name = "🌎 " + l.name;
              return;
            }
          });
        }

        if (this.addonLists.length == 0 && this.standardLists.length == 0) {
          this.createAndSelectUserLists();
        }
        if (view) {
          this.viewList(this.standardLists[0]);
        }
      });
    },
    getAbbs(listID) {
      api
        .getAbbs(listID)
        .then((response) => {
          if (response.data === null) {
            this.abbs = [];
          } else {
            this.abbs = response.data.sort((a, b) => {
              return new Date(b.updated) - new Date(a.updated);
            });
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    onScrollAbbs({target: { scrollTop, clientHeight, scrollHeight}}) {
      if(scrollTop + clientHeight >= scrollHeight) {
        this.perPage += 10;
      }
    },
  },
  mounted() {
    this.selectedAddons = this.$store.state.settings.selectedLists.addon;

    if (!Array.isArray(this.selectedAddons)) {
      this.selectedAddons = [];
      this.$store.commit("setSelectedAddons", this.selectedAddons);
    }

    window.addEventListener("scroll", this.onScrollAbbs);
    EventBus.$on("createdAbb", (abb) => {
      if(abb.targetListId == this.viewedList.id) {
        this.getAbbs(abb.targetListId);
      } 
    });
    EventBus.$on("createdList", (list) => {
      this.getLists(false);
      this.viewList(list);
    });
    EventBus.$on("removedList", this.afterRemovedList);
    EventBus.$on("updatedList", this.afterUpdatedList);
    this.getLists(true);
  },
  beforeDestroy() {
    window.removeEventListener("scroll", this.onScrollAbb);
    EventBus.$off("createdAbb");
    EventBus.$off("createdList");
    EventBus.$off("removedList");
    EventBus.$off("updatedList");
  },
};
</script>
<style scoped>
.active_row {
  background-color: #b8ffb8;
  padding: 0.7em;
  padding-left: 0em;
}
.word {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}
td {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 1px;
}
.abbsForRemoval {
  max-height: 450px;
  overflow-y: scroll;
}
</style>
