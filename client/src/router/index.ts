import { useAuthStore } from "@/stores/auth"
import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: "/login",
			name: "login",
			component: () => import("@/pages/LoginPage.vue"),
		},
		{
			path: "/",
			name: "device-info",
			component: () => import("@/pages/DeviceInfoPage.vue"),
		},
		{
			path: "/metrics",
			name: "metrics",
			component: () => import("@/pages/MetricsPage.vue"),
		},
		{
			path: "/processes",
			name: "processes",
			component: () => import("@/pages/ProcessesPage.vue"),
		},
		{
			path: "/port-manager",
			name: "port-manager",
			component: () => import("@/pages/PortManagerPage.vue"),
		},
		{
			path: "/metrics-history",
			name: "metrics-history",
			component: () => import("@/pages/MetricsHistoryPage.vue"),
		},
	],
})

router.beforeEach((to, from, next) => {
	const authStore = useAuthStore()

	if (to.path === "/login") {
		if (authStore.isAuthenticated) {
			next("/")
		} else {
			next()
		}
		return
	}

	if (!authStore.isAuthenticated) {
		next("/login")
	} else {
		next()
	}
})

export default router
