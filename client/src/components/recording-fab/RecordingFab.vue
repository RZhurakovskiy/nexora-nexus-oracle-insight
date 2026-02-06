<script setup lang="ts">
import MetricsRecordingModal from "@/components/metrics-recording-modal/MetricsRecordingModal.vue"
import { useMetricsStore } from "@/stores/metrics"
import { useNotificationsStore } from "@/stores/notifications"
import { useRecordingStore } from "@/stores/recording"
import { storeToRefs } from "pinia"
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
import "./RecordingFab.css"

const metrics = useMetricsStore()
const recording = useRecordingStore()
const notifications = useNotificationsStore()

const { monitoringEnabled } = storeToRefs(metrics)

const isModalOpen = ref(false)
const cpuThreshold = ref("70")
const ramThreshold = ref("80")
const durationMinutes = ref("5")

const isBusy = computed(() => recording.starting || recording.stopping)
const canStart = computed(
	() => monitoringEnabled.value && !recording.active && !isBusy.value
)

const open = () => {
	if (!monitoringEnabled.value) {
		notifications.info(
			"Перед запуском записи метрик необходимо включить мониторинг.",
			"Запись метрик"
		)
		return
	}
	isModalOpen.value = true
}

const start = async (payload: {
	cpuThreshold: number
	ramThreshold: number
	duration: number
}) => {
	await recording.start(payload)
}

onMounted(async () => {
	await metrics.init()
	await recording.init()
})

onBeforeUnmount(() => {
	/**
	 * Очистка таймера, если приложение размонтируется.
	 */
	recording.stopTick()
})
</script>

<template>
	<div class="recording-fab">
		<div v-if="recording.active" class="recording-fab__panel">
			<div class="recording-fab__row">
				<span class="recording-fab__dot" />
				<span class="recording-fab__text"
					>Запись • осталось {{ recording.remainingText }}</span
				>
			</div>
			<button
				class="recording-fab__btn recording-fab__btn--stop"
				type="button"
				:disabled="isBusy"
				@click="recording.stop()"
			>
				Остановить
			</button>
		</div>

		<button
			v-else
			class="recording-fab__btn recording-fab__btn--start"
			type="button"
			:disabled="!canStart"
			@click="open"
			:title="
				monitoringEnabled
					? 'Начать запись метрик'
					: 'Включите мониторинг, чтобы начать запись'
			"
		>
			Запись метрик
		</button>

		<MetricsRecordingModal
			:is-open="isModalOpen"
			v-model:cpuThreshold="cpuThreshold"
			v-model:ramThreshold="ramThreshold"
			v-model:durationMinutes="durationMinutes"
			@close="isModalOpen = false"
			@start="start"
		/>
	</div>
</template>
