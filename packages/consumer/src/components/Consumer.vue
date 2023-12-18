<template>
  <div>
    <Login />
    <TextView />
  </div>
</template>
<script>
import Login from "./modals/Login.vue"
import TextView from "./TextView.vue"
import Captions from "./Captions.vue"
export default {
    components: {
        TextView,
        Login,
        Captions
    }
}
</script>
