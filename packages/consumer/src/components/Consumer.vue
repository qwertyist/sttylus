<template>
    <div>
        <Login />
        <TextView />
    </div>
</template>
<script>
import Login from "./modals/Login.vue"
import TextView from "./TextView.vue"
export default {
    components: {
        TextView,
        Login,
    }
}
</script>
