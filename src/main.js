import { createApp } from 'vue'
import './tailwind.css'
import App from './App.vue'
import taskmanStore from './store/index.js'
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faUser } from '@fortawesome/free-solid-svg-icons'

library.add(faUser)

const app = createApp(App)
app.use(taskmanStore)
app.component('font-awesome-icon', FontAwesomeIcon)
app.mount('#app')
