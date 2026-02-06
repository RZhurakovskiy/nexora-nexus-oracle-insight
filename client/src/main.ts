import App from "@/App.vue"
import { ServicesKey } from "@/di/keys"
import { createAppServices } from "@/di/services"
import { createPiniaServicesPlugin } from "@/plugins/piniaServices"
import router from "@/router"
import { useAuthStore } from "@/stores/auth"
import { createPinia } from "pinia"
import { createApp } from "vue"
import "./style.css"

const pinia = createPinia()
const services = createAppServices()
pinia.use(createPiniaServicesPlugin(services))

const app = createApp(App)
app.provide(ServicesKey, services)
app.use(pinia)
app.use(router)

const authStore = useAuthStore()
authStore.checkAuth()

app.mount("#app")
