import VueRouter from 'vue-router';
import Vue from 'vue';
Vue.use(VueRouter);
import { store } from "./store/index.js"
import Consumer from "./Consumer.vue";
import Tabula from './Tabula.vue';
import Settings from './Settings.vue';
import Login from './Login.vue';
import Presentation from "./components/Presentation.vue"

const ifNotAuthenticated = (to, from, next) => {
  if (!store.state.token) {
    next('/login');
    return;
  }
  next();
};

const ifAuthenticated = (to, from, next) => {
  if (store.state.token) {
    next();
    return;
  }
  next('/login');
};

function createRouter() {
  if (import.meta.env.VITE_STTYLUS_MODE == "desktop") {
    console.log("Create desktop router")
    return new VueRouter({
      mode: 'hash',
      base: '/',
      routes: [
        { path: '/', component: Tabula, beforeEnter: ifNotAuthenticated },
        {
          path: '/settings',
          component: Settings,
          beforeEnter: ifNotAuthenticated,
        },
        { path: '/login', component: Login },
        { path: "/presentation", component: Presentation },
      ],
    });
  }

  if (import.meta.env.VITE_STTYLUS_MODE == "api") {
    return new VueRouter({
      mode: 'hash',
      base: '/app2/',
      routes: [
        { path: '/', component: Tabula, beforeEnter: ifNotAuthenticated },
        {
          path: '/settings',
          component: Settings,
          beforeEnter: ifNotAuthenticated,
        },
        { path: '/login', component: Login },
        { path: "/view", component: Consumer }
      ],
    });
  }
  if (import.meta.env.PROD == 'production') {
    return new VueRouter({
      mode: 'history',
      base: '/app2/',
      routes: [
        { path: '/', component: Tabula, beforeEnter: ifNotAuthenticated },
        {
          path: '/settings',
          component: Settings,
          beforeEnter: ifNotAuthenticated,
        },
        { path: '/login', component: Login },
        { path: "/view", component: Consumer }

      ],
    });
  } else {
    return new VueRouter({
      mode: 'history',
      base: "/",
      routes: [
        { path: '/', component: Tabula, beforeEnter: ifNotAuthenticated },
        {
          path: '/settings',
          component: Settings,
          beforeEnter: ifNotAuthenticated,
        },
        { path: '/login', component: Login },
        { path: "/view", component: Consumer }
      ],
    });
  }
}

export const router = createRouter();
