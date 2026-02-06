<script setup lang="ts">
import AppHeader from "@/components/app-header/AppHeader.vue"
import AppLoader from "@/components/app-loader/AppLoader.vue"
import AppNav from "@/components/app-nav/AppNav.vue"
import AppVersionBadge from "@/components/app-version/AppVersionBadge.vue"
import NotificationsHost from "@/components/notifications/NotificationsHost.vue"
import RecordingFab from "@/components/recording-fab/RecordingFab.vue"
import { useAuthStore } from "@/stores/auth"
import { useRootStatusStore } from "@/stores/rootStatus"
import { computed, onMounted } from "vue"
import { useRoute } from "vue-router"

const authStore = useAuthStore()
const rootStatus = useRootStatusStore()
const route = useRoute()

const isLoginPage = computed(() => route.path === "/login")

onMounted(async () => {
	if (authStore.isAuthenticated) {
		await rootStatus.init()
	}
})
</script>

<template>
	<div class="app">
		<AppLoader />
		<NotificationsHost />
		<template v-if="!isLoginPage">
			<RecordingFab />
			<AppVersionBadge />
			<AppHeader />
			<AppNav />
		</template>
		<main class="app-main" :class="{ 'app-main--login': isLoginPage }">
			<div class="monitor-container">
				<router-view />
			</div>
		</main>
	</div>
</template>
