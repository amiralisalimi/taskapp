import { createStore } from 'vuex'
import createPersistedState from 'vuex-persistedstate'
import taskman from './modules/taskman-store';
import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: 'https://localhost:8000',
  proxy: {
    host: 'localhost',
    port: 3000
  }
});

axiosInstance.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

const taskmanStore = createStore({
  modules: {
    taskman
  },
  plugins: [createPersistedState()],
  state: {
    token: localStorage.getItem('token'),
    isAuthenticated: false,
    user: null,
    boards: [],
    selectedBoard: null,
    boardId: null,
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
      state.token = user.token;
      localStorage.setItem('token', user.token);
      console.log('User logged in', state.user, state.token)
    },
    signup(state, user) {
      state.isAuthenticated = true;
      state.user = user;
      state.token = user.token;
      localStorage.setItem('token', user.token);
      console.log('User signed up')
      console.log(user)
      console.log(user.token)
    },
    logout(state) {
      localStorage.removeItem('token');
      state.token = null;
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
      state.boards = boards ?? [];
    },
    setSelectedBoard(state, boardId) {
      state.selectedBoard = state.boards.find(board => board.id === boardId);
      state.boardId = boardId
      console.log(boardId)
    },
    setContainers(state, containers) {
      state.containers = containers ?? [];
    },
    // setSelectedContainer(state, containerId) {
    //   state.selectedContainer = state.containers.find(container => container.id === containerId);
    // },
    setTasks(state, tasks) {
      state.tasks = tasks ?? [];
    },
    // setSelectedTask(state, taskId) {
    //   state.selectedTask = state.tasks.find(task => task.id === taskId);
    // },
    setBackground(state, background) {
      state.background = background;
    },
    // setTheme(state, theme) {
    //   state.theme = theme;
    // },
    setData(state, data) {
      state = data;
      axiosInstance.post('/update-user-data', this.getters.taskappData);
    },
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
    async createBoard({ state, commit }, { title }) {
      try {
        const response = await axiosInstance.post('/boards', { title });
        const board = response.data;
        commit('setBoards', [...state.boards, board]);
      } catch (error) {
        console.error(error);
        throw new Error('Board creation failed');
      }
    },
    async updateBoard({ state, commit }, { id, title }) {
      try {
        await axiosInstance.put(`/boards/${id}`, { title });
        const updatedBoards = state.boards.map(board => board.id === id ? { ...board, title } : board);
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
        console.log(containers)
        commit('setBoards', boards);
        commit('setSelectedBoard', boards[0].id);
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
    },
    async setData({ commit }, data) {
      commit('setData', data)
    }
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    user: state => state.user,
    taskappData: state => state,
  },
});

export default taskmanStore;