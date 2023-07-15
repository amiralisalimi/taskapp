<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  data() {
    return {
      showDropdown: false
    }
  },
  computed: {
    ...mapGetters(['user', 'isAuthenticated']),
  },
  methods: {
    ...mapActions(['logout']),
  },
}
</script>

<template>
  <div class="relative inline-block text-left">
    <div>
      <button @click="showDropdown = !showDropdown" type="button"
        class="inline-flex items-center justify-center w-full px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
        <svg class="h-5 w-5 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <span class="hidden md:inline-block">{{ user.username }}</span>
      </button>
    </div>

    <div v-show="showDropdown" @click.away="showDropdown = false" style="z-index: 999;"
      class="origin-top-right absolute right-0 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
      <div class="py-1" role="none">
        <div class="flex items-center px-4 py-2 text-sm text-gray-700" role="menuitem">
          <svg class="h-5 w-5 mr-2 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M19 14a2 2 0 11-4 0 2 2 0 014 0zM5 14a2 2 0 11-4 0 2 2 0 014 0zM12 14a2 2 0 11-4 0 2 2 0 014 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <div>
            <div class="font-medium text-gray-900">{{ user.username }}</div>
            <div class="text-sm text-gray-500">{{ user.email }}</div>
          </div>
        </div>
        <div class="border-t border-gray-100"></div>
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="logout">
          <svg class="h-5 w-5 mr-2 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M10 3L3 9M3 9l7 6M3 9l7-6m7 6h8a2 2 0 012 2v8m-2-10l-7 6m7-6l-7-6" />
          </svg>
          <span>Log Out</span>
        </div>
      </div>
    </div>
  </div>
</template>