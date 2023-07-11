import { createApp } from 'vue'
import './tailwind.css'
import App from './App.vue'
import taskmanStore from './store/index.js'

const app = createApp(App)
app.use(taskmanStore)
app.mount('#app')
