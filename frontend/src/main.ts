import {createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import VueTransitions from '@morev/vue-transitions'
import '@morev/vue-transitions/styles'
import App from './App.vue'
// import {WindowMaximise} from '../wailsjs/runtime'
// try {
//   WindowMaximise()
// } catch (err) {
//   console.error(err)
// }
const app = createApp(App)
app.use(ElementPlus)
app.use(VueTransitions)
app.mount('#app')
