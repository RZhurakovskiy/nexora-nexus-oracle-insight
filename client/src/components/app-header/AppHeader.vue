<script setup lang="ts">
import systemMonitorLogo from "@/assets/system-monitor-logo.svg"
import { useAuthStore } from "@/stores/auth"
import { useHostInfoStore } from "@/stores/hostInfo"
import { useMetricsStore } from "@/stores/metrics"
import { storeToRefs } from "pinia"
import { computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import Button from "../ui/Button.vue"
import "./appHeader.css"

const metrics = useMetricsStore()
const hostInfo = useHostInfoStore()
const authStore = useAuthStore()
const router = useRouter()
const { displayUsername, displayHostname } = storeToRefs(hostInfo)

const isMonitoringActive = computed(() => metrics.monitoringEnabled)
const isLoading = computed(() => metrics.starting || metrics.stopping)

const toggleMonitoring = async () => {
	if (isLoading.value) return
	if (isMonitoringActive.value) await metrics.stopMonitoring()
	else await metrics.startMonitoring()
}

const handleLogout = async () => {
	await authStore.logoutWithCleanup()
	router.push("/login")
}

onMounted(async () => {
	await metrics.init()
	await hostInfo.initHostInfo()
})
</script>

<template>
	<header class="app-header">
		<div class="app-header-content">
			<img :src="systemMonitorLogo" alt="Nexora" class="app-logo" />
			<div class="app-tite-item">
				<h1 class="app-title">Nexora</h1>
				<p class="app-subtitle">
					Мониторинг процессов, CPU, памяти и сети в реальном времени
				</p>
			</div>
		</div>
		<div class="account-item">
			<div class="account-actions">
				<Button
					:variant="isMonitoringActive ? 'danger' : 'success'"
					size="lg"
					:loading="isLoading"
					:disabled="isLoading"
					width="200px"
					@click="toggleMonitoring"
				>
					<template v-if="!isLoading && isMonitoringActive" #prefix>
						<!-- иконка stop - выключить мониторинг -->
						<svg
							width="16"
							height="16"
							viewBox="0 0 16 16"
							fill="none"
							xmlns="http://www.w3.org/2000/svg"
						>
							<rect
								x="3.5"
								y="3.5"
								width="9"
								height="9"
								rx="2"
								stroke="currentColor"
								stroke-width="1.6"
							/>
							<rect
								x="5.5"
								y="5.5"
								width="5"
								height="5"
								rx="1.2"
								fill="currentColor"
							/>
						</svg>
					</template>
					<template v-else-if="!isLoading" #prefix>
						<!-- иконка play - включить мониторинг -->
						<svg
							width="16"
							height="16"
							viewBox="0 0 16 16"
							fill="none"
							xmlns="http://www.w3.org/2000/svg"
						>
							<circle
								cx="8"
								cy="8"
								r="6.5"
								stroke="currentColor"
								stroke-width="1.6"
								opacity="0.9"
							/>
							<path d="M7 5.5L11 8L7 10.5V5.5Z" fill="currentColor" />
						</svg>
					</template>
					<template v-if="isLoading">Загрузка...</template>
					<template v-else-if="isMonitoringActive"
						>Выключить мониторинг</template
					>
					<template v-else>Включить мониторинг</template>
				</Button>
			</div>

			<div class="account-info">
				<div class="account-field">
					<svg
						class="account-icon"
						width="16"
						height="16"
						viewBox="0 0 16 16"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M8 8C10.2091 8 12 6.20914 12 4C12 1.79086 10.2091 0 8 0C5.79086 0 4 1.79086 4 4C4 6.20914 5.79086 8 8 8Z"
							fill="currentColor"
						/>
						<path
							d="M8 10C4.68629 10 2 12.6863 2 16H14C14 12.6863 11.3137 10 8 10Z"
							fill="currentColor"
						/>
					</svg>
					<div class="account-field-content">
						<span class="account-label">Пользователь</span>
						<span class="account-value">{{ displayUsername }}</span>
					</div>
				</div>
				<div class="account-divider"></div>
				<div class="account-field">
					<svg
						class="account-icon"
						width="16"
						height="16"
						viewBox="0 0 16 16"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<rect
							x="2"
							y="3"
							width="12"
							height="9"
							rx="1"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
						<path
							d="M2 6H14"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<path
							d="M5 9H11"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<path
							d="M6 12H10"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<circle cx="8" cy="7.5" r="0.5" fill="currentColor" />
					</svg>
					<div class="account-field-content">
						<span class="account-label">Хост</span>
						<span class="account-value">{{ displayHostname }}</span>
					</div>
				</div>
			</div>
			<div class="account-logout">
				<Button
					variant="ghost"
					size="md"
					@click="handleLogout"
				>
					<template #prefix>
						<svg
							width="16"
							height="16"
							viewBox="0 0 16 16"
							fill="none"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								d="M6 12L10 8L6 4"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
							<path
								d="M10 8H2"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
							/>
						</svg>
					</template>
					Выйти
				</Button>
			</div>
		</div>
	</header>
</template>
