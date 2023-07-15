import { createStore } from 'vuex'
import createPersistedState from 'vuex-persistedstate'
import taskman from './modules/taskman-store';
import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: 'https://localhost:8000',
  proxy: {
    host: 'localhost',
    port: 3000
  },
  // headers: {
  //   Authorization: "auth"
  // }
});

const taskmanStore = createStore({
  modules: {
    taskman
  },
  plugins: [createPersistedState()],
  state: {
    isAuthenticated: false,
    user: null,
    boards: [],
    selectedBoard: null,
    containers: [],
    selectedContainer: null,
    tasks: [],
    selectedTask: null,
    background: null,
    theme: null
  },
  mutations: {
    login(state, user) {
      state.isAuthenticated = true;
      state.user = user;
      console.log('User logged in')
    },
    signup(state, user) {
      state.isAuthenticated = true;
      state.user = user;
      console.log('User signed up')
    },
    logout(state) {
      state.isAuthenticated = false;
      state.user = null;
      state.boards = [];
      state.selectedBoard = null;
      state.containers = [];
      state.selectedContainer = null;
      state.tasks = [];
      state.selectedTask = null;
      state.background = null;
      state.theme = null;
      console.log('User logged out')
    },
    setBoards(state, boards) {
      state.boards = boards;
    },
    setSelectedBoard(state, boardId) {
      state.selectedBoard = state.boards.find(board => board.id === boardId);
    },
    setContainers(state, containers) {
      state.containers = containers;
    },
    setSelectedContainer(state, containerId) {
      state.selectedContainer = state.containers.find(container => container.id === containerId);
    },
    setTasks(state, tasks) {
      state.tasks = tasks;
    },
    setSelectedTask(state, taskId) {
      state.selectedTask = state.tasks.find(task => task.id === taskId);
    },
    setBackground(state, background) {
      state.background = background;
    },
    setTheme(state, theme) {
      state.theme = theme;
    }
  },
  actions: {
    async login({ commit, dispatch }, { username, password }) {
      try {
        const response = await axiosInstance.post('/login', { username, password });
        const user = response.data;
        commit('login', user);
        // Load the user's personalized data from the backend API
        await dispatch('loadUserData');
      } catch (error) {
        console.error(error);
        throw new Error('Invalid credentials');
      }
    },
    async signup({ commit, dispatch }, { username, email, password }) {
      try {
        const response = await axiosInstance.post('/signup', { username, email, password });
        const user = response.data;
        commit('signup', user);
        // Load the user's personalized data from the backend API
        await dispatch('loadUserData');
      } catch (error) {
        console.error(error);
        throw new Error('Signup failed');
      }
    },
    async logout({ commit }) {
      try {
        await axiosInstance.post('/logout');
        commit('logout');
      } catch (error) {
        console.error(error);
        throw new Error('Logout failed');
      }
    },
    async createBoard({ state, commit }, { name }) {
      try {
        const response = await axiosInstance.post('/boards', { name });
        const board = response.data;
        commit('setBoards', [...state.boards, board]);
      } catch (error) {
        console.error(error);
        throw new Error('Board creation failed');
      }
    },
    async updateBoard({ state, commit }, { id, name }) {
      try {
        await axiosInstance.put(`/boards/${id}`, { name });
        const updatedBoards = state.boards.map(board => board.id === id ? { ...board, name } : board);
        commit('setBoards', updatedBoards);
      } catch (error) {
        console.error(error);
        throw new Error('Board update failed');
      }
    },
    async deleteBoard({ state, commit }, { id }) {
      try {
        await axiosInstance.delete(`/boards/${id}`);
        const updatedBoards = state.boards.filter(board => board.id !== id);
        commit('setBoards', updatedBoards);
        commit('setSelectedBoard', null);
      } catch (error) {
        console.error(error);
        throw new Error('Board deletion failed');
      }
    },
    async loadUserData({ commit }) {
      try {
        // Load the user's personalized data from the backend API
        const response = await axiosInstance.get('/user-data');
        const { boards, containers, tasks, background } = response.data;
        commit('setBoards', boards);
        // commit('setSelectedBoard', selectedBoard);
        commit('setContainers', containers);
        // commit('setSelectedContainer', selectedContainer);
        commit('setTasks',tasks);
        // commit('setSelectedTask', selectedTask);
        commit('setBackground', background);
        // commit('setTheme', theme);
        // Load the user's taskman data from the backend API
        // const taskmanResponse = await axiosInstance.get('/taskman-data');
        // const taskmanData = taskmanResponse.data;
        // commit('taskman/SET_TASKMAN', taskmanData);
      } catch (error) {
        console.error(error);
        throw new Error('Failed to load user data');
      }
    }
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    user: state => state.user,
  }
});

export default taskmanStore;