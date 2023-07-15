<script setup>
import { ref, reactive, watch, computed, onBeforeMount } from 'vue'
import { useStore } from 'vuex'
import axios from 'axios'
import KanbanBoard from './components/kanban/KanbanBoard.vue'
import ContainerModal from '@/components/modals/ContainerModal.vue'
import CloseIcon from '@/components/icons/CloseIcon.vue'
import SaveIcon from '@/components/icons/SaveIcon.vue'
import GithubIcon from '@/components/icons/GithubIcon.vue'
import Login from '@/components/user/Login.vue'
import Button from '@/components/base/Button.vue'
import Dropdown from '@/components/user/Dropdown.vue'
import DropdownBackground from '@/components/base/DropdownBackground.vue'

const store = useStore()
const displayContainerModal = ref(false)
const displayCardModal = ref(false)
const state = reactive({
  is_editing_title: false,
  temp_title: null,
})
const payload = computed(() => {
  return store.getters['taskman/gettaskmanDatas']
})
const user = store.state.user

onBeforeMount(async () => {
  const data = store.getters['taskman/gettaskmanDatas']
  if (!data) {
    await axios.get('/sample-data.json').then(({ data }) => {
      store.dispatch('taskman/settaskman', data)
    })
  }
})

const handleEditTitle = (type) => {
  if (type === 'edit') {
    state.is_editing_title = true
    state.temp_title = payload.value.title
  } else if (type === 'save') {
    state.is_editing_title = false
    payload.value.last_modified = new Date().toLocaleString('fa-IR')
    store.dispatch('taskman/settaskman', payload.value)
  } else {
    state.is_editing_title = false
    payload.value.title = state.temp_title
  }
}
// const openRepo = () => {
//   window.open('https://github.com/kurnyaannn/taskman', '_blank')
// }
const isLoggedIn = () => {
  return store.getters.isAuthenticated
}
const logout = () => {
  store.dispatch('logout')
}
</script>

<template>
  <div v-if="isLoggedIn()" :class="backgroundImage">
    <div v-if="payload" class="flex h-screen flex-col p-4">
      <div class="flex justify-between rounded-lg text-white">
        <div>
          <Transition name="fade" mode="out-in">
            <div v-if="!state.is_editing_title">
              <span
                class="rounded-md px-2 text-3xl font-bold transition-all duration-300 ease-in-out hover:cursor-pointer hover:bg-slate-200"
                @click="handleEditTitle('edit')">
                {{ payload.title }}
              </span>
            </div>
            <div v-else class="flex place-items-center">
              <input v-model="payload.title" type="text"
                class="block w-[230px] rounded-lg border border-gray-300 bg-gray-50 p-2 text-xl text-gray-900 transition duration-300 ease-in-out focus:border-blue-500 focus:ring-blue-500"
                placeholder="Add Board Title" @keypress.enter="handleEditTitle('save')" />
              <div class="ml-2 flex place-items-center justify-center">
                <SaveIcon height="30px"
                  class="mr-2 cursor-pointer rounded-full bg-blue-500 p-1 text-white hover:bg-blue-700"
                  @click="handleEditTitle('save')" />
                <CloseIcon height="30px"
                  class="cursor-pointer rounded-full p-1 text-red-500 hover:bg-red-600 hover:text-white"
                  @click="handleEditTitle('cancel')" />
              </div>
            </div>
          </Transition>
          <h3 class="my-2 px-2 text-sm">
            Last Modified : {{ payload.last_modified }}
          </h3>
        </div>
        <div class="mt-px flex place-items-start justify-center">
          <Dropdown />
        </div>
        <DropdownBackground @back1="back1" @back2="back2" @back3="back3" @back4="back4" @back5="back5" />
      </div>
      <KanbanBoard :payload="payload" @addContainer="displayContainerModal = true" />
    </div>
    <!-- Container Modal -->

    <ContainerModal :value="displayContainerModal" @close="displayContainerModal = false" />
  </div>
  <div v-else class="flex h-screen flex-col p-4" :class="backgroundImage">
    <Login />
  </div>
</template>

<script>
  export default {
    data() {
      return {
        backgroundImage: "first",
        showPopup: false
      };
    },
    methods: {
      back1() {
        this.backgroundImage = "first"
      },
      back2() {
        this.backgroundImage = "second"
      },
      back3() {
        this.backgroundImage = "third"
      },
      back4() {
        this.backgroundImage = "forth"
      },
      back5() {
        this.backgroundImage = "fifth"
      }
    }
  };
</script>
