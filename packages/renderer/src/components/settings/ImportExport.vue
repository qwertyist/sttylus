<template>
  <b-row>
    <b-col cols="3">
      <h3>Importera förkortningslista</h3>
      <div v-if="!confirm">
        <b-form>
          <b-select v-model="from" @change="changeFrom" :disabled="to != 'dontExport'">
            <b-select-option selected value="dontImport">...</b-select-option>
            <b-select-option
              v-for="source in sources"
              :key="source.value"
              :value="source.value"
              :disabled="source.disabled"
            >{{ source.text }}</b-select-option>
          </b-select>
        </b-form>
        <br />
        <h3>Exportera förkortningslista</h3>
        <b-form>
          <b-select v-model="to" @change="changeTo" :disabled="from != 'dontImport'">
            <b-select-option selected value="dontExport">...</b-select-option>
            <b-select-option
              v-for="target in targets"
              v-bind:key="target.value"
              v-bind:value="target.value"
              :disabled="target.disabled"
            >{{ target.text }}</b-select-option>
          </b-select>
        </b-form>
      </div>
      <div v-if="confirm || safe">
        <hr v-if="safe" />
        <p>Antal förkortningar: {{ abbs.length }}</p>
        <b-form @submit.prevent="importAbbs" autocomplete="off">
          <b-form-checkbox v-model="listForm.addToExisting">Lägg till i befintlig lista</b-form-checkbox>
          <br />
          <div v-if="listForm.addToExisting">
            <b-form-group label="Välj förkortningslista">
              <b-select v-model="listForm.targetList">
                <b-select-option
                  v-bind:key="list.id"
                  v-bind:value="list"
                  v-for="list in lists"
                >{{ list.name }}</b-select-option>
              </b-select>
            </b-form-group>
          </div>
          <div v-else>
            <b-form-group label="Döp förkortningslistan">
              <b-form-input required placeholder="..." v-model="listForm.name" />
            </b-form-group>
          </div>Ämneslista
          <b-form-checkbox
            v-model="listForm.standard"
            inline
            name="standard-switch"
            switch
          >Standardlista</b-form-checkbox>
          <br />
          <br />
          <b-button @click="onCancelImport" variant="danger">Avbryt</b-button>
          <div v-if="!conflicts" class="float-right">
            <b-button variant="primary" type="submit">Lägg till</b-button>
          </div>
        </b-form>
      </div>
    </b-col>

    <b-col>
      <div v-if="exportList">
        <div v-if="to == 'textontop'">
          <h3>Exportera till TextOnTop</h3>
        </div>
        <div v-if="to == 'sttylus' || to == 'textfile'">
              <h3>Exportera till STTylus eller textfil</h3>
            <b-list-group>
              <b-list-group-item v-for="standard in standardLists" v-bind:key="standard.value">
                <b-row>
                  <b-col cols=4>
                    <b>{{ standard.text }}</b>
                  </b-col>
                  <b-col>
                    <small>{{ standard.counter }} förkortningar</small> 
                  </b-col>
                  <b-col>
                    <span style="font-size: 1.2em;"><b-badge size="md">{{ standard.value.split("-")[0] }}</b-badge></span>
                    <b-button class="float-right" size="sm" @click="exportSTTylus(standard.value)">Exportera</b-button>
                  </b-col>
                  <b-col>
                    <b-button 
                      size="sm" 
                      variant="primary" 
                      @click="getTextfile(standard)"
                    >
                      <b-icon icon="download" />
                      Hämta textfil
                    </b-button>
                  </b-col>
                </b-row>
              </b-list-group-item>
            <b-list-group-item v-for="addon in addonLists" v-bind:key="addon.value">
            <b-row>
              <b-col cols=4>
                {{ addon.text }} 
              </b-col>
              <b-col>
                <small>{{ addon.counter }} förkortningar</small>
              </b-col>
              <b-col>
                <span style="font-size: 1.2em;"><b-badge size="md">{{ addon.value.split("-")[0] }}</b-badge></span>
                <b-button class="float-right" size="sm" @click="exportSTTylus(addon.value)">Exportera</b-button>
              </b-col>
              <b-col>
                <b-button size="sm" variant="primary" @click="getTextfile(addon)">
                  <b-icon icon="download" />
                  Hämta textfil
                </b-button>
              </b-col>
            </b-row>
            </b-list-group-item>
          </b-list-group>
        </div>
        <div v-if="to != 'sttylus' && to != 'textfile'">
        <b-form @change="listenToEvent" @submit.prevent="submitExport">
          <b-form-radio-group v-model="exportForm.standard" :required="to == 'protype'">
            <b-list-group>
              <b-list-group-item v-for="standard in standardLists" v-bind:key="standard.value">
                <b-row>
                  <b-col cols="9">
                    <b-form-radio :value="standard.value"><b>{{ standard.text }}</b></b-form-radio>
                  </b-col>
                  <b-col>
                    <small>{{ standard.counter }} förkortningar</small>
                  </b-col>
                </b-row>
              </b-list-group-item>
            </b-list-group>
          </b-form-radio-group>
          <b-form-checkbox-group v-model="exportForm.addons">
            <b-list-group-item v-for="addon in addonLists" v-bind:key="addon.value">
              <b-row>
                <b-col cols="9">
                  <b-form-checkbox :value="addon.value">{{ addon.text }}</b-form-checkbox>
                </b-col>
                <b-col>
                  <small>{{ addon.counter }} förkortningar</small>
                </b-col>
              </b-row>
            </b-list-group-item>
          </b-form-checkbox-group>
          <hr />
          <b-button type="submit">Exportera</b-button>
        </b-form>
        </div>
      </div>
      <div v-if="!confirm">
        <div v-if="from == 'sttylus'">
          <h3>Importera från STTylus</h3>
          <b-form @submit.prevent="importSTTylus">
            Ange listans ID:
            <b-row>
              <b-col cols="4">
                <b-form-input :formatter="sttylusIdFormatter" v-model="form.id" />
              </b-col>
            </b-row>
            <br />
            <b-button type="submit">Hämta</b-button>
          </b-form>
        </div>

        <div v-if="from == 'protype'">
          <h3>Importera från ProType</h3>
          <b-jumbotron>
          <template #lead>
         <p>För import av listor och användare från ProType rekommenderas <a href="https://sttylus.se/utx/">UTX-programmets Förkortningsverktyg</a>.<br />
         Där kan du konvertera från flera filtyper till textfilsformat.</p>
          </template>
          Annars hittar du din användare under <code>Users/</code> där du har ProType installerat. För varje förkortningslista finns ytterligare en katalog och när du hittat din lista det filen som heter antingen <code>wordlist</code> eller <code>wordlist.dat</code> du ska välja för uppladdning i formuläret nedan.
          </b-jumbotron>
          <div v-if="online">
            <b-form @submit.prevent="uploadProType">
              <b-form-file v-model="form.file" accept=".dat" />
              <hr />
              <b-button type="submit">Ladda upp</b-button>
            </b-form>
          </div>
          <div v-else></div>
        </div>
        <div v-if="from == 'textontop' || from == 'illumitype'">
          <div v-if="!uploaded">
            <h3>
              Importera från
              <span v-if="from == 'textontop'">TextOnTop</span>
              <span v-else>IllumiType</span>
            </h3>
            <b-form @submit.prevent="uploadJSONList">
              <b-form-file v-model="form.file" accept=".json" />
              <hr />
              <b-button type="submit">Ladda upp</b-button>
            </b-form>
            {{ output }}
          </div>
          <div v-else>
            <h3>
              Välj lista från
              <span v-if="from == 'textontop'">TextOnTop</span>
              <span v-else>IllumiType</span>
            </h3>
            <b-form @submit.prevent="selectUploadedList">
              <b-select v-model="selectedJSONList" @change="changeSelectedJSONList">
                <b-select-option
                  v-bind:key="i + '_' + listName"
                  v-bind:value="listName"
                  v-for="(listName, i) in importableListNames"
                >{{ listName }}</b-select-option>
              </b-select>
              <br />
              <br />
              Det finns {{ countJSONAbbs }} förkortningar i denna lista.
              <br />
              <b-button type="submit">Välj lista</b-button>
            </b-form>
          </div>
        </div>
        <div v-if="from == 'textfile'">
          <h3>Importera från textfil</h3>

          <b-form @submit.prevent="uploadTxt">
            <b-form-textarea
              id="textarea"
              v-model="txt"
              placeholder="Klistra in förkortningar enligt formatet:
förkortning=fras
förkortning=fras"
              rows="20"
              max-rows="6"
            ></b-form-textarea>
            <hr />
            <b-button type="submit">Importera</b-button>
          </b-form>
        </div>
      </div>
      <div v-if="confirm || safe">
        <b-row>
          <b-col cols="8">
            <h3 v-if="safe">Redigera importerade förkortningar</h3>
            <h3 v-else>Kontrollera importerade förkortningar</h3>
          </b-col>
          <b-col>
          <span v-if="nConflicts > 0">
              <b-badge class="float-right" variant="warning">
                Det finns
                <span v-if="nConflicts == 1">ett</span>
                <span v-else>{{ nConflicts }}</span> fel eller
                <span v-if="nConflicts == 1">konflikt</span>
                <span v-else>konflikter</span> att korrigera
              </b-badge>
            </span>
            <span class="float-right">
                Ersätt alla 
                <b-form-checkbox class="pt-1" inline model="replaceAll" @change="changeOverwriteAll"></b-form-checkbox>
            </span>
          </b-col>
          <b-col>
            <b-pagination
              class="float-right"
              v-model="currentPage"
              :totalRows="abbs.length"
              :perPage="10"
              aria-controls="abbs"
            />
          </b-col>
        </b-row>
        <b-form autocomplete="off">
          <div
            style="height: 68vh !important; overflow-y: auto;"
          >
            <b-table
              small
              responsive
              hover
              striped
              :perPage="10"
              :currentPage="currentPage"
              :items="abbs"
              :fields="abbFields"
              :sort-by="sortBy"
              :sort-desc="sortDesc"
              @sort-changed="sortChanged"
            >
              <template v-slot:cell(abb)="row">
                <div style="min-width: 110px; max-width: 100px; white-space: normal">
                  <b-form-input
                    v-model="row.item.abb"
                    size="md"
                    class="mb-0"
                    @change="onUpdateAbb(row.item)"
                    :state="row.item.validabb"
                  />
                </div>
              </template>
              <template v-slot:cell(word)="row">
                <b-form-input
                  v-model="row.item.word"
                  size="md"
                  class="mb-0"
                  @change="onUpdateWord(row.item)"
                  :state="row.item.validword"
                />
              </template>
              <template v-slot:cell(issue)="row">
                <span v-if="debug">
                  <b-button size="sm" :id="'popover' + row.item.id">?</b-button>
                  <b-popover :target="'popover' + row.item.id" triggers="hover" placement="left">
                    validabb: {{ row.item.validabb }}
                    <br />
                    validword: {{ row.item.validword }}
                    <br />
                    dup: {{ row.item.notduplicate }}
                    <br />
                    issue: {{ row.item.issue }}
                    <br />
                  </b-popover>
                </span>
                <span v-if="row.item.notduplicate != null">
                  <b-badge variant="warning" v-if="!row.item.notduplicate">Konflikt</b-badge>
                  <br />
                  <span v-if="row.item.issue == '__textduplicate__'">Dubblett i textfil</span>
                  <span v-else-if="row.item.issue == '__importduplicate__'">Dubblett vid importering</span>
                  <span v-else>
                    <b-form-checkbox
                      v-model="row.item.overwrite"
                      @change="changeOverwriteAbb(row.item)"
                    >
                      Ersätt
                      <i>{{ row.item.issue }}</i>
                    </b-form-checkbox>
                  </span>
                </span>
                <span
                  v-if="
                    row.item.notduplicate == null &&
                    (row.item.validword != null || row.item.validabb != null)
                  "
                >
                  <b-badge variant="danger">Fel format</b-badge>
                  <br />
                  <span v-if="row.item.issue == '__formaterror__'">Fel format vid importering</span>
                  <span v-else>{{ row.item.issue }}</span>
                </span>
              </template>
              <template v-slot:cell(delete)="row">
                <b-button class="float-right" variant="danger" @click="onDeleteAbb(row.item)">
                  <b-icon icon="trash" />
                </b-button>
              </template>
            </b-table>
          </div>
          <b-button @click="resolveConflicts(null)">Kontrollera fel/konflikter</b-button>
        </b-form>
      </div>
    </b-col>
    <ResolveConflicts />
  </b-row>
</template>

<script>
import axios from "axios";
import api from "../../api/api.js";
import ResolveConflicts from "../modals/ResolveConflicts.vue";
import EventBus from "../../eventbus";
import { parseTxt, resolveIssues } from "../functions/resolveConflicts.js";
import { get } from "vue-cookie";
export default {
  components: {
    ResolveConflicts,
  },
  data() {
    return {
      output: "",
      debug: false,
      sources: [
        { value: "sttylus", text: "från STTylus", disabled: false },
        { value: "protype", text: "från ProType", disabled: false },
        { value: "textontop", text: "från TextOnTop", disabled: false },
        { value: "illumitype", text: "från IllumiType", disabled: false },
        { value: "textfile", text: "från textfil", disabled: false },
        { value: "shared", text: "från delad lista", disabled: true },
      ],
      targets: [
        { value: "sttylus", text: "till STTylus", disable: false },
        { value: "textontop", text: "till TextOnTop", disabled: false },
        { value: "protype", text: "till Protype", disabled: false },
        { value: "textfile", text: "till textfil", disabled: false },
      ],
      from: "dontImport",
      to: "dontExport",
      importList: true,
      exportList: false,
      online: true,
      uploaded: false,
      importableListNames: [],
      importableLists: {},
      selectedTextOnTopList: "",
      countJSONAbbs: 0,
      safe: false,
      confirm: false,
      conflicts: true,
      nConflicts: -1,
      replaceAll: false,
      form: {
        file: null,
        id: "",
      },
      sortBy: "issue",
      sortDesc: true,
      txt: "",
      abbs: [],
      abbFields: [
        {
          key: "abb",
          label: "Förkortning",
          sortable: true,
          thStyle: { width: "10%" },
          tdClass: [ "conflictTableCell" ]
        },
        {
          key: "word",
          label: "Text/Fras",
          sortable: true,
          thStyle: { textAlign: "center", width: "50%" },
          tdClass: [ "conflictTableCell" ]
        },
        {
          key: "issue",
          label: "Fel/Konflikt",
          sortable: true,
          tdClass: [ "conflictInfo", "conflictTableCell" ]
          /*          thStyle: { textAlign: "center", width: "30%" },

          tdClass: "text-center",
          */
        },
        {
          key: "delete",
          label: "Ta bort",
          sortable: false,
          thStyle: { textAlign: "right" },
          tdClass: [ "conflictTableCell" ]
        },
      ],
      lists: [],
      exportForm: {
        standard: "",
        addons: [],
        counter: 0,
      },
      listForm: {
        addToExisting: false,
        listName: "",
        targetList: {},
        standard: false,
      },
      currentPage: 1,
      regex: /([åäöÅÄÖA-Za-z0-9]*)=(.*)/iu,
    };
  },
  watch: {
    listForm: {
      deep: true,
      handler() {
        this.resolveConflicts();
      },
    },
  },
  computed: {
    standardLists() {
      return this.lists.filter((list) => {
        return list.type == 0
      }).map(standard => {
        return { value: standard.id, text: standard.name, counter: standard.counter }
      })
    },
    addonLists() {
      return this.lists.filter((list) => {
        return list.type == 1
      }).map(addon => {
        return { value: addon.id, text: addon.name, counter: addon.counter }
      })
    },
    abbsToConfirm() {
    }
  },
  methods: {
    listenToEvent(value) {
    },
    sttylusIdFormatter(value) {
      return value.trim().toLowerCase();
    },
    importSTTylus(short_id) {
      this.$toast.info("Försöker hämta lista med kort-ID:" + this.form.id)
      if(this.$mode =="desktop") {
        axios.get("https://sttylus.se/api2/abbs/public/" + this.form.id)
        .then(resp => {
          this.$toast.info("Hämtar listan...")
          this.confirm = true;
          this.populateAbbs(resp.data.abbs);
        })
        return
      }
      api.importSTTylus(this.form.id)
      .then(resp => {
        this.$toast.info("Hämtar listan...")
        this.confirm = true;
        this.populateAbbs(resp.data.abbs);
      })
    },
    exportSTTylus(list) {
      api.getList(list)
      .then(resp => {
        api.exportSTTylus(resp.data)
        .then(resp => {
          this.$toast.success("Listan går att importera online/offline med ID-numret till vänster")
        })
        .catch(err => {
          this.$toast.warning("Listan kunde inte exporteras:" + err, { duration: 5000})
        })
      })
      .catch(err => {
        this.$toast.error("Något gick fel:" + err, { duration: 5000})
      })
    },
    getTextfile(list) {
      this.$toast.default("get textfiles..."+ list.text) 
      api.getAbbs(list.value)
      .then(resp => {
        this.$toast.success("textfile abbs" + resp.data.length) 
        const text = resp.data.map(abb => [abb.abb+"="+abb.word]).join("\n")

        const a = document.createElement("a");
        a.style.display = "none";
        document.body.appendChild(a);
        a.href = window.URL.createObjectURL(
          new Blob([text], { type: "text/plain" })
        );
        a.setAttribute("download", list.text+".txt")
        a.click();
        window.URL.revokeObjectURL(a.href);
        document.body.removeChild(a);
      })
      .catch(err => {
        console.error("failed converting to textfile", err)
        this.$toast.warning("textfile failed", err) 
      })
    },
    submitExport() {
      if (this.exportForm.standard == "" && this.exportForm.addons.length == 0) {
        this.$toast.warning("Du måste välja minst en lista att exportera")
        return
      }
      if(this.to == "textfile") {
        this.$toast.warning("handskas med textfil...")
        return
      }
      api.exportLists(this.to, this.exportForm).then(resp => {
        if (this.to == "protype") {
          const url = window.URL.createObjectURL(new Blob([resp.data], { type: "application/zip" }));
          const link = document.createElement("a");
          link.href = url;
          link.setAttribute("download", this.$store.state.userData.name.split(" ")[0] + ".zip")
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
        } else if (this.to == "textontop") {
          const url = window.URL.createObjectURL(new Blob([resp.data], { type: "application/json" }));
          const link = document.createElement("a");
          link.href = url;
          link.setAttribute("download", this.$store.state.userData.name.split(" ")[0] + ".json")
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
        }
      })
    },
    addEventListeners() {
      EventBus.$on("importSharedAbbs", () => {
        this.getSharedAbbs(true);
      });
    },
    removeEventListeners() {
      EventBus.$off("importSharedAbbs");
    },
    onUpdateAbb(val) {
      this.resolveConflicts(val);
    },
    onUpdateWord(val) {
      this.resolveConflicts(val);
    },
    changeOverwriteAbb(abb) {
      console.log(
        abb.overwrite ? "overwrite" : "edit",
        " ",
        abb.abb,
        "(" + abb.issue + ")"
      );
    },
    onDeleteAbb(abb) {
      let i = this.abbs.findIndex((a) => a.id === abb.id);
      this.abbs.splice(i, 1);
      this.resolveConflicts();
    },
    changeOverwriteAll(val) {
      console.log(this.abbs)
      this.replaceAll = val
      this.abbs = this.abbs.map(a => { 
        a.overwrite = val
        return a
      })
    },
    changeFrom(value) {
      switch (this.from) {
        case "dontImport":
          return
        case "textfile":
          this.safe = false;
          break;
        case "protype":
          this.safe = false;
          break;
        case "shared":
          this.getSharedAbbs(true);
          break;
        case "sttylus":
          this.safe = false;
          break;
      }
    },
    changeTo(value) {
      switch (this.to) {
        case "dontExport":
          this.exportList = false;
          return
        case "textontop":
          this.exportForm.standard = ""
          break;
        case "protype":
          if (this.exportForm.standard == "") {
            this.exportForm.standard = this.standardLists[0].value;
          }
          break;
        case "shortform":
          this.exportForm.standard = ""
          break;
        default:
          break;
      }
      this.exportList = true;
    },
    sorter() { },
    reset() {
      this.$store.commit("setSharedAbbs", []);
      this.from = "dontImport";
      this.to = "dontExport";
      this.txt = "";
      this.abbs = [];
      this.confirm = false;
      this.safe = false;
      this.form.file = null;
      this.conflicts = false;
      this.nConflicts = -1;
    },
    populateAbbs(abbs) {
      this.abbs = abbs;
    },
    getSharedAbbs(focus) {
      if (this.$store.state.sharedAbbs.length > 0) {
        let i = this.sources.findIndex((m) => m.value == "shared");
        this.sources[i].disabled = false;
        this.from = "shared";

        if (focus) {
          let abbs = this.$store.state.sharedAbbs.filter((a) => {
            return a.abb != a.word;
          });
          api.importAbbs(abbs).then((resp) => {
            this.from = "shared";
            this.safe = true;
            this.abbs = resp.data;
            this.resolveConflicts();
          });
        }
      }
    },
    sortChanged(val) { 

    },
    uploadProType() {
      let data = new FormData();
      data.append("file", this.form.file);
      if (this.form.file == null) {
        return;
      }
      api.uploadProtype(data).then((resp) => {
        this.confirm = true;
        this.populateAbbs(resp.data);
      });
    },
    uploadIllumiType() {
      let data = new FormData();
      data.append("file", this.form.file);
      if (this.form.file == null) {
        return;
      }
    },
    uploadJSONList() {
      console.log("upload json list")
      let data = new FormData();
      data.append("file", this.form.file);
      if (this.form.file == null) {
        console.log("No file selected");
        return;
      }
      if (this.from == "textontop") {
        api.uploadTextOnTop(data).then((resp) => {
          this.importableListNames = Object.keys(resp.data)
          console.log("got tot lists: ", resp.data)
          this.importableLists = resp.data
          this.selectedJSONList = this.importableListNames[0]
          this.changeSelectedJSONList()
          this.uploaded = true
        })
      } else if (this.from == "illumitype") {
        api.uploadIllumiType(data).then((resp) => {
          this.importableListNames = Object.keys(resp.data)
          this.importableLists = resp.data
          this.selectedJSONList = this.importableListNames[0]
          this.changeSelectedJSONList()
          this.output = resp.data
          this.uploaded = true
        })
      }
    },
    selectUploadedList() {
      this.abbs = this.importableLists[this.selectedJSONList]
      this.confirm = true
    },
    changeSelectedJSONList() {
      console.log(this.importableLists)
      console.log(this.selectedJSONList)
      this.countJSONAbbs = this.importableLists[this.selectedJSONList].length
    },
    uploadTxt() {
      let parsed = parseTxt(this.txt);
      api.uploadTxt(parsed).then((resp) => {
        this.abbs = resp.data.map((a) => {
          let c = JSON.parse(a.comment);
          a.comment = "";
          return { ...a, ...c };
        });
        this.resolveConflicts();
        this.confirm = true;
      });
    },
    resolveConflicts(abb) {
      let abbs = this.abbs;
      let resolved = resolveIssues(abbs);
      abbs = resolved.abbs;
      let issues = resolved.issues;
      if (issues > 0) {
        this.$toast.warning(
          "Det finns fel eller konflikter i listan som du importerar"
        );
      }
      if (this.listForm.addToExisting) {
        api
          .resolveConflicts(this.listForm.targetList.id, abbs)
          .then((resp) => {
            let conflicts = resp.data;
            conflicts.map((c) => {
              let i = this.abbs.findIndex((a) => a.abb == c.abb);
              this.abbs[i].validabb = false;
              this.abbs[i].notduplicate = false;
              this.abbs[i].issue = c.old;
            });
            this.$nextTick(() => {
              console.log("====================");
              console.log("calculate conflicts:");
              console.log("====================");
              this.nConflicts = this.abbs.filter((a) => {
                if (a.issue != null) {
                  console.log("abb:", a.abb, "word:", a.word);
                  console.log("issue:", a.issue);
                  if (a.overwrite) {
                    console.log(":::::::::::::::::::::");
                    console.log("overwrite:", a.overwrite);
                    console.log(":::::::::::::::::::::");
                  } else {
                    console.log("====================");
                  }
                }
                return a.issue != null && a.overwrite != true;
              }).length;
              console.log("conflicts:", this.nConflicts);
              console.log("\n");
              if (this.nConflicts > 0) {
                this.sortBy = "issue";
                this.conflicts = true;
              } else {
                this.conflicts = false;
              }
            });
          })
          .catch((err) => {
            console.log(err);
          });
      } else {
        this.$nextTick(() => {
          this.nConflicts = this.abbs.filter(
            (a) => a.issue != null && a.overwrite == false
          ).length;
          if (this.nConflicts == 0) {
            this.sortBy = "issue";
            this.conflicts = false;
          } else {
            this.conflicts = true;
          }
        });
      }
    },
    onCancelImport() {
      this.$bvModal
        .msgBoxConfirm("Är du säker på att du vill avbryta importeringen?", {
          okTitle: "Ja",
          okVariant: "danger",
          cancelTitle: "Fortsätt importera",
          cancelVariant: "primary",
        })
        .then((value) => {
          if (value) {
            if (this.from == "sharedAbbs") {
              this.$store.commit("setSharedAbbs", []);
            }
            this.reset();
          } else {
            this.$toast.info("Fortsätt importera, gör ingenting");
          }
        });
    },
    importAbbs() {
      this.resolveConflicts();
      if (this.conflict) {
        return;
      }
      if (this.listForm.addToExisting) {
        if (this.listForm.targetList != "") {
          api
            .importList(this.listForm.targetList.id, this.abbs)
            .then((resp) => {
              EventBus.$emit("createdList", this.listForm.targetList);
              this.reset();
            });
          return;
        }
      } else {
        let list = {
          name: this.listForm.name,
          type: this.listForm.standard ? 0 : 1,
        };
        api
          .createList(list)
          .then((resp) => {
            api.importList(resp.data.id, this.abbs)
              .then((resp) => {
                EventBus.$emit("createdList");
                this.reset();
              });
          })
          .catch((err) => {
            console.log("import couldnt create list", err);
          });
      }
    },
  },
  mounted() {
    this.addEventListeners();
    this.getSharedAbbs(true);
    api.getUserLists().then((resp) => {
      this.lists = resp.data;
      this.listForm.targetList = this.lists[0];
    });
  },
  beforeDestroy() {
    this.removeEventListeners();
  },
};
</script>
<style scoped>
::v-deep .table > tbody > tr > td {
  padding: 0.1em;
  padding-top: 0.3em;
}

::v-deep .conflictInfo {
  white-space: nowrap;
  padding: 0em !important;
}
</style>
