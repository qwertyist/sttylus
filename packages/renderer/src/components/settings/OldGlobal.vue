<template>
  <div>
    <b-tabs v-model="tabIndex">
      <b-tab title="Gemensamma listor">
        <b-tabs
          v-model="listTabIndex"
          pill
          content-class="mt-3"
          @input="changeListTab"
        >
          <b-tab
            v-for="l in globalLists"
            :title="l.name"
            :key="'dyn-list' + l.id"
          >
            <b-table
              striped
              hover
              :items="globalAbbs"
              :fields="abbFields"
            >
              <template #cell(abb)="row">
                <b-form-input
                  v-model="row.item.abb"
                  disabled
                />
              </template>
              <template #cell(word)="row">
                <b-form-input
                  v-model="row.item.word"
                  disabled
                />
              </template>
            </b-table>
          </b-tab>
        </b-tabs>
      </b-tab>
      <b-tab
        v-if="stagedNominations.length > 0"
        title="Nominera förkortningar"
      >
        <b-row>
          <b-col cols="4">
            <b-form
              @submit.prevent="onSubmitCommend"
              @reset.prevent="onResetCommend"
            >
              <!--<b-form-group
                id="commend-level"
                label="Nivå"
                label-for="commend-level-select"
                description="Välj om förkortningarna ska göras globala eller delas till en tolk eller en grupp tolkar"
              >
                <b-form-select
                  id="commend-level-select"
                  v-model="commendForm.selectedLevel"
                  :options="commendLevels"
                  class="mb-3"
                >
                  <template v-slot:first>
                    <b-form-select-option :value="null" disabled>-- Välj nivå --</b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>-->
              <b-form-group
                id="commend-target-list"
                label="Förkortningslista"
                label-for="commend-target-list-select"
                description="Välj vilken förkortningslista du vill nominera förkortningarna till"
              >
                <b-form-select
                  v-model="commendForm.selectedTargetList"
                  id="commend-target-list-select"
                  :options="targetLists"
                  required
                  class="mb-3"
                >
                  <template #first>
                    <b-form-select-option
                      :value="null"
                      disabled
                    >
                      -- Välj förkortningslista --
                    </b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>
              <b-form-group
                id="commend-comment"
                label="Kommentar/Beskrivning"
                label-for="commend-title-input"
                description="Frivilligt: Motivera din nominering"
              >
                <b-form-input
                  v-model="commendForm.comment"
                  id="commend-comment-input"
                  required
                  placeholder="Kommentar"
                />
              </b-form-group>
              <b-button
                type="reset"
                variant="danger"
              >
                Avbryt och återställ
              </b-button>
              <b-button
                type="submit"
                variant="primary"
              >
                Nominera
              </b-button>
            </b-form>
          </b-col>
          <b-col>
            <b-table
              striped
              hover
              :items="stagedNominations"
              :fields="stagedNominationFields"
            >
              <template #cell(abb)="row">
                <b-form-input v-model="row.item.abb" />
              </template>
              <template #cell(word)="row">
                <b-form-input v-model="row.item.word" />
              </template>
            </b-table>
          </b-col>
        </b-row>
      </b-tab>
      <b-tab>
        <template #title>
          Nomineringar
          <b-badge>{{ commendRequests.length }}</b-badge>
        </template>
        <b-row>
          <b-col>
            <b-list-group
              v-for="c in commendRequests"
              :key="c.id"
            >
              <b-list-group-item>
                <b-card
                  :title="c.comments[0].message"
                  :sub-title="'till Svensk standardlista'"
                  class="mb-2"
                >
                  <!-- <b-card-header v-if="c.type < 3">global</b-card-header>
                  <b-card-header v-else>ej global</b-card-header>-->
                  <b-card-body>
                    <b-card-text>
                      <div class="abbsInCommend">
                        <b-list-group>
                          <b-list-group-item
                            v-for="abb in c.abbs"
                            class="py-1 d-flex justify-content-between align-items-center"
                            @change="toggleKeepAbb(abb)"
                            :key="abb.id"
                          >
                            <!-- <b-form-checkbox
                              inline
                              v-model="abb.selected"
                            >-->
                            <div>{{ abb.abb }} - {{ abb.word }}</div>
                            <div>
                              {{ abb.comment ? [abb.comment] : "kommentar" }}
                              <small>{{
                                new Date(abb.updated).toLocaleDateString(
                                  undefined,
                                  {
                                    day: "numeric",
                                    month: "long",
                                    year: "numeric",
                                  }
                                )
                              }}</small>
                            </div>
                          </b-list-group-item>
                        </b-list-group>
                      </div>

                      <b-list-group
                        v-for="(comment, i) in c.comments.slice(1)"
                        :key="comment.id + '_' + i"
                      >
                        <b-list-group-item>
                          {{
                            comment.message
                          }}
                        </b-list-group-item>
                      </b-list-group>
                    </b-card-text>
                  </b-card-body>

                  <b-card-footer>
                    <div class="float-right">
                      <b-form inline>
                        <b-form-group
                          label="Kommentera nominering"
                          :label-sr-only="true"
                        >
                          <b-form-input placeholder="Lämna en kommentar" />
                          <b-button variant="secondary">
                            Nej
                            <b-icon-arrow-down-square class="commendButtons" />
                          </b-button>
                          <b-button variant="primary">
                            Ja
                            <b-icon-arrow-up-square class="commendButtons" />
                          </b-button>
                          <b-button
                            variant="primary"
                            @click="showCommendModal(c, 'view')"
                          >
                            Visa mer
                            <b-icon-box-arrow-up-right class="commendButtons" />
                          </b-button>
                        </b-form-group>
                      </b-form>

                      <div
                        v-if="admin"
                        class="float-right"
                      >
                        <b-icon-check
                          variant="success"
                          class="commendButtons"
                          @click="showCommendModal(c, 'commend')"
                        />
                        <b-icon-trash-fill
                          variant="danger"
                          class="commendButtons"
                          @click="showCommendModal(c, 'deny')"
                        />
                      </div>
                    </div>
                  </b-card-footer>
                </b-card>
              </b-list-group-item>
            </b-list-group>
          </b-col>
        </b-row>
      </b-tab>
    </b-tabs>

    <b-modal
      id="commend-modal"
      title="Visa nominering"
      hide-footer
      ref="commend-modal"
    >
      <p>
        <b>Förkortningslista:</b>
        <span v-show="selectedNomination.list.name == 'dave'"> Svensk standardlista </span>
      </p>
      <hr>
      <p>
        <b>Förkortningar</b>
      </p>
      <p>
        <b-list-group style="max-height: 400px; overflow-y: scroll;">
          <b-list-group-item
            v-for="abb in selectedNomination.abbs"
            class="d-flex justify-content-between align-items-center"
            :key="abb.id"
          >
            {{ abb.abb }} - {{ abb.word }}
          <!--  <b-badge v-b-popover.hover.bottom="'Konflikt'" variant="danger"
              >2</b-badge
            >-->
          </b-list-group-item>
        </b-list-group>
      </p>
      <p>
        <b>Kommentarer</b>
      </p>
      <p>
        <b-list-group>
          <b-list-group-item
            v-for="comment in selectedNomination.comments"
            class="d-flex justify-content-between align-items-center"
            :key="comment.id"
          >
            {{ comment.message }}
            <b-badge variant="info">
              Botvid
            </b-badge>
          </b-list-group-item>
        </b-list-group>
        <b-form
          @submit="onSubmitCommend"
          @reset="onResetCommend"
        >
          <b-form-group label="Lägg till kommentar">
            <b-form-input />
          </b-form-group>
        </b-form>
      </p>
      <b-button
        variant="success"
        @click="confirmCommend()"
      >
        Slutför
      </b-button>
      <b-button variant="primary">
        Skicka
      </b-button>
      <b-button
        variant="danger"
        @click="confirmDeny()"
      >
        Neka
      </b-button>
    </b-modal>
  </div>
</template>
                  
<script>
import axios from "axios";
export default {
  data() {
    return {
      admin: true,
      connected: false,
      globalLists: [],
      globalAbbs: [],
      abbFields: [
        { abb: { label: "Förkortning", sortable: true } },
        { word: { label: "Ord", sortable: true } },
        { meta: { label: "Övrigt" } },
      ],
      listTabIndex: 0,
      stagedNominationFields: [
        { abb: { label: "Förkortning", sortable: true } },
        { word: { label: "Ord", sortable: true } },
      ],
      commendForm: {
        comment: "",
        selectedLevel: null,
        newList: false,
        newListName: "",
        newListType: false,
        selectedTargetList: null,
      },
      commendLevels: [
        { value: 0, text: "Global" },
        { value: 1, text: "Västerbottens läns landsting" },
      ],
      targetLists: [],
      commendRequests: [],
      commendRequestsFields: [
        { key: "comments", label: "Kommentarer" },
        { key: "abbs", label: "Förkortningar" },
      ],
      selectedNomination: {
        action: "",
        comments: [],
        list: {
          id: "",
          name: "",
          type: 0,
          status: 0,
        },
        abbs: {
          id: "",
          abb: "",
          word: "",
          conflicts: [],
        },
      },
      removeFromNomination: [],
    };
  },
  methods: {
    updateNominations() {
      this.stagedNominations = this.$store.state.stagedNominations;
      if (this.stagedNominations.length > 0) {
        this.tabIndex = 1;
      }

      if (this.commendRequests && this.commendRequests.length > 0) {
        this.tabIndex = 2;
      }
    },
    changeListTab(input) {
      if (input == -1) {
        return;
      }
      console.log(input)
      const listID = this.globalLists[input].id;
      this.getGlobalAbbs(listID);
    },
    fetchAndCommitCommends() {
      axios
        .get(this.$backend + "/api/commend/fetchandcommit")
        .then((resp) => {
          axios.get(this.$backend + "/api/commend").then((resp) => {
            if (resp.data != null) {
              this.commendRequests = resp.data;
            }
          });
        })
        .catch((err) => {
          console.log(err);
        });
    },
    getGlobalLists() {
      axios
        .get(this.$backend + "/api/abbs/global")
        .then((resp) => {
          if (resp.data != null) {
            this.globalLists = resp.data;

            this.targetLists = resp.data.map((l) => {
              return { value: l.id, text: l.name };
            });
          }
              this.updateNominations();

        })
        .catch((err) => {
          console.log(err);
        });
    },
    getGlobalAbbs(listID) {
      axios
        .get(this.$backend + "/api/abbs/global/" + listID)
        .then((response) => {
          if (response.data === null) {
            this.globalAbbs = [];
          } else {
            this.globalAbbs = response.data.sort((a, b) => {
              return new Date(b.updated) - new Date(a.updated);
            });
          }
          console.log(this.globalAbbs)
        })
        .catch((err) => {
          console.log(err);
        });
    },
    onSubmitCommend(evt) {
      evt.preventDefault();
      let list = {};
      if (this.commendForm.newList) {
        list.name = this.commendForm.newListName;
        list.type = this.commendForm.newListType ? 1 : 0;
      } else {
        list.id = this.commendForm.selectedTargetList;
      }
      axios
        .post(this.$backend + "/api/commend", {
          type: 1,
          abbs: this.stagedNominations,
          level: this.commendForm.selectedLevel,

          list: list,
          comments: [
            {
              message: this.commendForm.comment,
              user: this.$store.state.userData.id,
            },
          ],
        })
        .then(() => {
          this.$store.commit("clearStagedNominations");
          this.updateNominations();
        })
        .catch((err) => {
          console.log(err);
        });
    },
    onResetCommend() {
      this.stagedNominations = [];
      this.commendForm.name = "";
      this.commendForm.selectedLevel = null;
      this.$store.commit("clearStagedNominations");
      this.$emit("reset");
    },
    openCommendRequest(n) {
      console.log(n);
    },
    toggleRemoveAbb(abb) {
      console.log(abb);
    },
    toggleKeepAbb(abb) {
      console.log(abb);
    },
    showCommendModal(n, action) {
      console.log(this.commendRequests);
      console.log("show this:", n);
      this.selectedNomination = n;
      this.selectedNomination.action = action;
      this.$nextTick().then((_) => {
        this.$refs["commend-modal"].show();
      });
    },
    confirmCommend() {
      this.selectedNomination.type = 1;
      const data = JSON.stringify(this.selectedNomination);
      axios.post(
        this.$backend + "/api/commend/" + this.selectedNomination.id,
        data
      );
    },
    confirmDeny() {
      axios
        .delete(this.$backend + "/api/commend/" + this.selectedNomination.id)
        .then((response) => {
          console.log(response);
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
  computed: {
    stagedNominations: {
      get() {
        return this.$store.state.stagedNominations;
      },

      set(value) {
        this.$store.commit("setStagedNominations", value);
        return value;
      },
    },
    tabIndex: {
      get() {
        if (this.stagedNominations.length > 0) {
          return 1;
        }

        if (this.commendRequests && this.commendRequests.length > 0) {
          return 2;
        }
        return 0;
      },
      set(value) {
        return value;
      },
    },
  },
  watch: {
    stagedNominations(newValue) {
      if (newValue.length > 0) {
        this.tabIndex = 1;
      } else {
        this.tabIndex = 0;
      }
    },
  },
  mounted() {
    const connected = navigator.onLine ? true : false;
    if (connected == true) {
    }
    this.fetchAndCommitCommends();
    this.getGlobalLists();
  },
};
</script>
<style scoped>
.abbsInCommend {
  max-height: 450px;
  overflow-y: scroll;
}
.commendButtons {
  font-size: 1.7em;
  vertical-align: middle;
}
</style>
