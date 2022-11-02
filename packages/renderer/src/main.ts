/* eslint-disable space-before-function-paren */
//import 'bootstrap/dist/css/bootstrap.css';
//import 'bootstrap-vue/dist/bootstrap-vue.css';
import './assets/custom.scss'
import('typeface-roboto-mono')
import 'typeface-roboto-mono/index.css';
import Vue from 'vue';
import Vuex from 'vuex';
import { router } from "./router.js"
import { store } from "./store/index.js"
import moment from 'moment';
import Vue2TouchEvents from 'vue2-touch-events';
declare module 'vue/types/vue' {
  interface Vue {
    $backend: string;
    userID: null;
    machineId: null;
    currentLists: any[];
    checkConnection: () => void;
    mode: string;
  }
}

import axios from 'axios';
axios.defaults.withCredentials = false;


import { ValidationProvider, extend, ValidationObserver } from 'vee-validate';
import { required, confirmed } from 'vee-validate/dist/rules';
import bootstrap from "./bootstrap.js"
Vue.use(bootstrap)
extend("confirmed", {
  ...confirmed,
  message: "Lösenordsbekräftelsen överensstämmer inte"
})
extend('required', {
  ...required,
  message: 'Detta fält är obligatoriskt',
});

Vue.component('ValidationObserver', ValidationObserver);
Vue.component('ValidationProvider', ValidationProvider);
import "./assets/toast.css"
import VueToast from 'vue-toast-notification';

import VueCookie from 'vue-cookie';

Vue.use(VueToast, { duration: 5000 });

//const isOnline = import('is-online');

Vue.use(Vue2TouchEvents);
Vue.use(Vuex);
Vue.use(VueCookie);

Vue.filter('formatDate', function (value: string) {
  if (value) {
    return moment(String(value)).format('YYYY/MM/DD HH:mm');
  }
});
Vue.filter('formatHour', function (value: string) {
  if (value) {
    return moment(String(value)).format('HH:mm');
  }
});
Vue.filter('formatChangeLogDate', function (value: string) {
  if (value) {
    return moment(value).format("DD MMMM YYYY")
  }
})

Vue.prototype.$lastUpdate = import.meta.env.VITE_STTYLUS_BUILD_DATE
Vue.config.productionTip = false;
Vue.prototype.$mode = import.meta.env.VITE_STTYLUS_MODE
if (import.meta.env.VITE_STTYLUS_MODE == "desktop") {
  Vue.prototype.$desktop = true
    console.log("Running as desktop app")
  Vue.prototype.$backend =
    import.meta.env.VITE_STTYLUS_LOCAL_BACKEND;
  Vue.prototype.$collabServer = "wss://sttylus.se/ws/"
  Vue.prototype.$collabAPI= "https://sttylus.se/ws/"
  Vue.prototype.$localCollab = import.meta.env.VITE_STTYLUS_LOCAL_COLLAB;
} else {
  Vue.prototype.$desktop = false

  if (import.meta.env.PROD) {
    console.log("Running in production")
    Vue.prototype.$backend = import.meta.env.VITE_STTYLUS_BACKEND
    Vue.prototype.$collabServer = "wss://sttylus.se/ws/"
    Vue.prototype.$collabAPI= "https://sttylus.se/ws/"
    Vue.prototype.$collabAPI = import.meta.env.VITE_STTYLUS_COLLAB_API;
  } else {
    console.log("Running in development")
    Vue.prototype.$backend =
      import.meta.env.VITE_STTYLUS_LOCAL_BACKEND;
  Vue.prototype.$collabServer = "wss://sttylus.se/ws/"
  Vue.prototype.$collabAPI= "https://sttylus.se/ws/"
    Vue.prototype.$localCollab = import.meta.env.VITE_STTYLUS_LOCAL_COLLAB;

  }
}
console.log("pointing to backend ", Vue.prototype.$backend)

import EventBus from './eventbus.js';


new Vue({
  // eslint-disable-line no-new
  router,
  store,
  template: `
  <div>
    <router-view class="view"></router-view>
  </div>
  `,

  data() {
    return {
      userId: null,
      machineId: null,
      currentLists: [],
    };
  },
  methods: {
    checkConnection() {
      let status = 'connecting';
      this.$store.commit('setConnectionStatus', 'connecting');
      this.$store.commit('setConnectionStatus', 'online');
      EventBus.$emit('networkStatusUpdate', status);
    },
    next() {
      if (this.mode == 'webapp') {
        if (!store.getters.isAuthenticated) {
          router.push('/login');
        }
      }
    }
  },
  computed: {
    mode(): string {
      let mode: string = ""
      if (typeof import.meta.env.VITE_STTYLUS_MODE === "string") {
        mode = import.meta.env.VITE_STTYLUS_MODE;
      }
      this.$store.commit('setMode', mode);
      return mode;
    },
  },

  mounted() {
    EventBus.$on('checkConnection', this.checkConnection);
    this.$store.commit('initState');
    EventBus.$on("stateReady", this.next)

  },
  beforeDestroy() {
    EventBus.$off("stateReady")
  }
}).$mount('#app').$nextTick(window.removeLoading);
