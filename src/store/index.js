import { createStore } from 'vuex'
import createPersistedState from 'vuex-persistedstate'
import taskman from './modules/taskman-store'

const taskmanStore = createStore({
  modules: {
    taskman
  },
  plugins: [createPersistedState()],
  state: {
    isAuthenticated: false,
    user: null,
  },
  mutations: {
    login(state, user) {
      state.isAuthenticated = true;
      state.user = user;
    },
    signup(state, user) {
      state.isAuthenticated = true;
      state.user = user;
    },
    logout(state) {
      state.isAuthenticated = false;
      state.user = null;
    },
  },
  actions: {
    login({ commit }, { username, password }) {
      // Call authentication API here
      // If successful, commit the login mutation
      commit('login', { username });
    },
    signup({ commit }, { username, email, password }) {
      commit('signup', { username, email, password });
    },
    logout({ commit }) {
      // Call logout API here
      // If successful, commit the logout mutation
      commit('logout');
    },
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    user: state => state.user,
  },
})

export default taskmanStore