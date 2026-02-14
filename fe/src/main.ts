import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import App from './App.vue'
import router from './router'
import { logger } from './api/logger'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

// Axios Interceptors for Global Logging
axios.interceptors.request.use(config => {
  logger.debug(`API Request: ${config.method?.toUpperCase()} ${config.url}`, config.data)
  return config
}, error => {
  logger.error('API Request Error', error)
  return Promise.reject(error)
})

axios.interceptors.response.use(response => {
  logger.debug(`API Response: ${response.status} ${response.config.url}`, response.data)
  return response
}, error => {
  const status = error.response?.status
  const url = error.config?.url
  logger.error(`API Error ${status}: ${url}`, error.response?.data || error.message)
  return Promise.reject(error)
})

app.use(pinia)
app.use(router)

app.mount('#app')

logger.info('Application started')
