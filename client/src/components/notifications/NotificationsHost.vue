<script setup lang="ts">
import { useNotificationsStore } from "@/stores/notifications"
import { storeToRefs } from "pinia"
import "./NotificationsHost.css"

const notifications = useNotificationsStore()
const { items } = storeToRefs(notifications)

const iconPath = (type: "success" | "error" | "info") => {
	if (type === "success")
		return "M6.5 11.2 3.8 8.5 2.7 9.6 6.5 13.4 13.3 6.6 12.2 5.5z"
	if (type === "error") return "M6.1 6.1 9.9 9.9M9.9 6.1 6.1 9.9"
	return "M8 3.2a.9.9 0 1 1 0 1.8.9.9 0 0 1 0-1.8zM7.2 6.4h1.6v6.4H7.2z"
}
</script>

<template>
	<div class="toasts" aria-live="polite" aria-relevant="additions">
		<transition-group name="toast" tag="div" class="toasts__stack">
			<div
				v-for="value in items"
				:key="value.id"
				class="toast"
				:class="`toast--${value.type}`"
			>
				<div class="toast__icon" aria-hidden="true">
					<svg width="18" height="18" viewBox="0 0 16 16" fill="none">
						<path
							:d="iconPath(value.type)"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
						/>
					</svg>
				</div>
				<div class="toast__content">
					<div class="toast__title">{{ value.title }}</div>
					<div class="toast__message">{{ value.message }}</div>
				</div>
			</div>
		</transition-group>
	</div>
</template>
