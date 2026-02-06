<script setup lang="ts">
import type { WsStatus } from "@/services/realtime/RealtimeSocket"
import { computed } from "vue"
import "./WebSocketStatusIndicator.css"

type Props = {
	status: WsStatus
	label?: string
}

const props = withDefaults(defineProps<Props>(), {
	label: "WebSocket",
})

type StatusConfig = {
	text: string
	color: string
	solid: boolean
}

const getStatusConfig = (status: WsStatus): StatusConfig => {
	switch (status) {
		case "connected":
			return { text: "Подключено", color: "#30d158", solid: true }
		case "connecting":
			return { text: "Подключение...", color: "#ff9500", solid: false }
		case "reconnecting":
			return { text: "Переподключение...", color: "#ff9500", solid: false }
		case "error":
			return { text: "Ошибка", color: "#ff453a", solid: true }
		case "disconnected":
		case "idle":
		default:
			return { text: "Отключено", color: "#8e8e93", solid: false }
	}
}

const config = computed(() => getStatusConfig(props.status))
</script>

<template>
	<div class="ws-status-indicator">
		<div
			class="ws-status-indicator__dot"
			:style="{ color: config.color }"
			:title="`${label}: ${config.text}`"
		>
			<svg
				v-if="config.solid"
				width="8"
				height="8"
				viewBox="0 0 8 8"
				fill="none"
			>
				<circle cx="4" cy="4" r="4" fill="currentColor" />
			</svg>
			<svg v-else width="8" height="8" viewBox="0 0 8 8" fill="none">
				<circle
					cx="4"
					cy="4"
					r="3"
					stroke="currentColor"
					stroke-width="2"
					fill="none"
				/>
			</svg>
		</div>
		<span class="ws-status-indicator__text">{{ config.text }}</span>
	</div>
</template>
