<template>
  <b-modal id="dashboard" hideHeader hideFooter>
    <b-card
      :title="'STTylus ' + $version"
      :subTitle="'Välkommen tillbaka ' + userName"
    >
      <hr />
      Laddar listor: {{ lists.length }}
      <p>Nyheter i programmet:</p>
    </b-card>
  </b-modal>
</template>
<script>
import db from '../../store/db.js'
export default {
  computed: {
    userName() {
      return this.$store.state.userData.name
    },
    lists() {
      return db.getLists()
    },
  },
}
</script>
