<script setup lang="ts">
import { useVersionStore } from "@/stores/version"
import { storeToRefs } from "pinia"
import { computed, onMounted } from "vue"
import "./AppVersionBadge.css"

const version = useVersionStore()
const { clientVersion, serverVersion, loading, error } = storeToRefs(version)

const serverText = computed(() => {
	if (loading.value) return "..."
	if (serverVersion.value) return serverVersion.value
	if (error.value) return "недоступно"
	return "—"
})

onMounted(async () => {
	await version.init()
})
</script>

<template>
	<div
		class="app-version-badge"
		:title="error ? `Ошибка: ${error}` : 'Версии приложения'"
	>
		<span class="app-version-badge__label">Клиент</span>
		<span class="app-version-badge__value">{{ clientVersion }}</span>
		<span class="app-version-badge__sep">•</span>
		<span class="app-version-badge__label">Сервер</span>
		<span class="app-version-badge__value">{{ serverText }}</span>
	</div>
</template>
