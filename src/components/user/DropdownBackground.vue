<script>

  export default {
    data() {
      return {
        showDropdown: false,
        backgroundImage: "first",
        showPopup: false
      };
    },
    methods: {
      back1() {
        this.$emit('back1')
      },
      back2() {
        this.$emit('back2')
      },
      back3() {
        this.$emit('back3')
      },
      back4() {
        this.$emit('back4')
      },
      back5() {
        this.$emit('back5')
      }
    }
  };
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
        <span class="hidden md:inline-block">Background Image</span>
      </button>
    </div>

    <div v-show="showDropdown" @click.away="showDropdown = false" style="z-index: 999;"
      class="origin-top-right absolute right-0 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
      <div class="py-1" role="none">
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="back1">
          <svg class="h-7 w-10 mr-2 text-gray-400 first" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
          <span>Image 1</span>
        </div>
        <div class="border-t border-gray-100"></div>
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="back2">
          <svg class="h-7 w-10 mr-2 text-gray-400 second" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
          <span>Image 2</span>
        </div>
        <div class="border-t border-gray-100"></div>
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="back3">
          <svg class="h-7 w-10 mr-2 text-gray-400 third" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
          <span>Image 3</span>
        </div>
        <div class="border-t border-gray-100"></div>
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="back4">
          <svg class="h-7 w-10 mr-2 text-gray-400 forth" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
          <span>Image 4</span>
        </div>
        <div class="border-t border-gray-100"></div>
        <div
          class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem" @click="back5">
          <svg class="h-7 w-10 mr-2 text-gray-400 fifth" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
          </svg>
          <span>Image 5</span>
        </div>
      </div>
    </div>
  </div>
</template>