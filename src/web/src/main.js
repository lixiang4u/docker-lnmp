import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'


axios.defaults.baseURL = 'http://127.0.0.1:8086';
axios.interceptors.request.use(config => {
    // console.log('[interceptors.config]', config);
    // 需要检查调用方有没有传params
    if (config.params === undefined) {
        config.params = [];
    }
    config.params['platform'] = 'web';
    return config;
}, error => {
    // console.log('[interceptors.error]', error);
    return Promise.reject(error);
});

const app = createApp(App)

app.use(ElementPlus)
app.use(router)

app.mount('#app')
