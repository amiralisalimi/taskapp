<template>
  <div class="rounded-lg bg-[#F2F3F9] bg-opacity-5 backdrop-blur-lg flex flex-col items-center justify-center min-h-screen py-6 bg-gray-50">
    <div class="w-full max-w-sm">
      <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <h2 class="text-center text-xl font-bold mb-6">{{ isLoginPage ? 'Log In' : 'Sign Up' }}</h2>
        <form>
          <div class="mb-4">
            <label class="block text-gray-700 font-bold mb-2" for="username">
              username
            </label>
            <input v-if="isLoginPage" v-model="login.username"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="username" type="username" autocomplete="username" required>
            <input v-else v-model="signup.username"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="username" type="username" autocomplete="username" required>
          </div>
          <div class="mb-4" v-if="!isLoginPage">
            <label class="block text-gray-700 font-bold mb-2" for="email">
              Email
            </label>
            <input v-model="signup.email"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="email" type="email" required>
          </div>
          <div class="mb-4">
            <label class="block text-gray-700 font-bold mb-2" for="password">
              Password
            </label>
            <input v-if="isLoginPage" v-model="login.password"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="password" type="password" required>
            <input v-else v-model="signup.password"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="password" type="password" required>
          </div>
          <div class="mb-4" v-if="!isLoginPage">
            <label class="block text-gray-700 font-bold mb-2" for="password_confirmation">
              Confirm Password
            </label>
            <input v-model="signup.passwordConfirmation"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="password_confirmation" type="password" required>
          </div>
          <div class="flex items-center justify-between">
            <button @click.prevent="submitForm"
              class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="submit">
              {{ isLoginPage ? 'Log In' : 'Sign Up' }}
            </button>
            <button @click.prevent="toggleMode"
              class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" type="button">
              {{ isLoginPage ? 'Create a new account' : 'Already have an account?' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      login: {
        username: '',
        password: ''
      },
      signup: {
        username: '',
        email: '',
        password: '',
        passwordConfirmation: '',
      },
      isLoginPage: false
    }
  },
  methods: {
    submitForm() {
      if (this.isLoginPage) {
        this.$store.dispatch('login', {
          username: this.login.username,
          password: this.login.password,
        })
      } else if (this.signup.password === this.signup.passwordConfirmation) {
        this.$store.dispatch('signup', {
          username: this.signup.username,
          email: this.signup.email,
          password: this.signup.password,
        })
      } else {
        alert('Password confirmation incorrect')
      }
    },
    toggleMode() {
      this.isLoginPage = !this.isLoginPage;
    }
  },
};
</script>