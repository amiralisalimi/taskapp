import axios from 'axios'

const state = {
  taskman: null
}

const getters = {
  gettaskmanDatas(state) {
    return state.taskman
  }
}

const actions = {
  settaskman({ commit }, data) {
    commit('SET_TASKMAN', data)
  }
}

const mutations = {
  SET_TASKMAN(state, data) {
    state.taskman = data
  }
}

const taskmanModule = {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

export default taskmanModule