import {createApp} from 'vue'
import App from './App.vue'
import router from './router/index.js'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import VueKonva from 'vue-konva';
import 'element-plus/dist/index.css'
import 'bootstrap/dist/css/bootstrap.min.css'
import './assets/css/index.scss'

const app = createApp(App)

app.use(createPinia()).use(VueKonva).use(router).use(ElementPlus).mount('#app')

/*window.addEventListener('contextmenu', (e) => {
    e.preventDefault();
});*/