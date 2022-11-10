<template>
  <b-modal
    id="addAbb"
    size="lg"
    :modal-class="addModalClass"
    hide-footer
    hide-header
    hide-backdrop
    no-fade
    return-focus=".ql-editor"
    @show="showModal"
    @hide="closeModal"
    ref="addabbmodal"
  >
    <b-row>
      <b-form inline @submit.prevent="onSubmit" @keydown.113="nextList" autocomplete="off">
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
          >{{ list.name }}</option>
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
          style="position:relative;right: 10px"
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
import api from "../../api/api.js";
import EventBus from "../../eventbus.js";

export default {
  name: "AddAbbreviation",
  props: ["query"],
  data() {
    return {
      addModalClass: ["addModalClass"],
      currentLists: {},
      form: {
        abb: "",
        word: "",
        lists: [],
        selected: "",
        index: 0,
      },
      exists: null,
      existing: "",
      secondary: false,
    };
  },
  computed: {
    listColorID() {
      return { color: "red" }
    },
    sharedList() {
      return this.$store.state.sharedList
    }
  },
  methods: {
    catchEscape() {
      this.$toast.info("haha")
    },
    handleKeydown(e) {
      console.log(e)
    },
    checkExists() {
      if (this.form.abb !== "") {
        if (this.form.selected) {
          api.getAbb(this.form.selected, this.form.abb.trim()).then((resp) => {
            if (resp.data) {
              this.exists = false;
              this.existing = resp.data.word;
            } else {
              this.exists = null;
              this.existing = "";
            }
          });
        }
      }
    },
    showModal() {
      document.addEventListener("keyup", e => {
        if (e.key == "Escape") {
          this.$bvModal.hide("support")
        }
      })
      this.$store.commit("setModalOpen", true)
      EventBus.$emit("modalOpened");
      this.updateLists();
      console.log(this.sharedList)
      this.form.word = this.$store.state.selectedWord;
      this.form.selected = this.$store.state.targetList.id;
      this.index = this.$store.state.targetList.index;
    },
    closeModal() {
      document.removeEventListener("keyup", e => {
        if (e.key == "Escape") {
          this.$bvModal.hide("support")
        }
      })
      this.$store.commit("setModalOpen", false)
      this.$bvModal.hide("addAbb");
      EventBus.$emit("closeNav");
      EventBus.$emit("refocus", "");

      this.form.abb = "";
      this.form.word = "";
      this.exists = null;
      this.existing = "";
    },
    updateLists() {
      this.currentLists = this.$store.state.settings.selectedLists;
      api
        .getLists([this.currentLists.standard].concat(this.currentLists.addon))
        .then((resp) => {
          this.form.lists = resp.data;
        })
        .catch((err) => { });
    },
    onSubmit() {
      let targetListId = null;
      if (this.form.abb === "" || this.form.word === "") {
        this.$store.commit("setSelectedWord", "");
        this.$bvModal.hide("addAbb")
        return;
      }

      if (this.form.selected.length === 0) {
        targetListId = this.$store.state.selectedLists.standard.id;
      } else {
        targetListId = this.form.selected;
      }
      let data = {
        abb: this.form.abb.trim(),
        word: this.form.word.trim(),
        creator: this.$store.state.userData.id,
        targetListId,
      };
      api
        .createAbb(targetListId, {
          abb: data.abb,
          word: data.word,
          creator: data.creator,
        })
        .then((resp) => {
          EventBus.$emit("createdAbb", data);
          this.$store.commit("setSelectedWord", "");
          this.$bvModal.hide("addAbb")
        })
        .catch((err) => {
          console.log("Failed creating abbreviation:", err);
          this.$bvModal.hide("addAbb")
        });
    },
    nextList() {
      if (this.index < this.form.lists.length - 1) {
        this.index++;
        this.form.selected = this.form.lists[this.index].id;
      } else {
        this.index = 0;
        this.form.selected = this.form.lists[this.index].id;
      }
      let targetList = {
        index: this.index,
        id: this.form.selected,
      };
      this.$store.commit("setTargetList", targetList);
      this.checkExists();
    },
  },

  mounted() {
    EventBus.$on("showTextView", () => {
      this.$bvModal.hide("addAbb");
      this.$bvModal.hide("support");
    })
    this.currentLists = this.$store.state.settings.selectedLists;
    api
      .getLists([this.currentLists.standard].concat(this.currentLists.addon))
      .then((resp) => {
        this.form.lists = resp.data;
      });
  },
  beforeDestroy() {
    EventBus.$off("showTextView")
  },
};
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
