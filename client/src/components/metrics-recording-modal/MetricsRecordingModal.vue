<script setup lang="ts">
import BaseModal from "@/components/modal/BaseModal.vue"
import { computed } from "vue"
import "./MetricsRecordingModal.css"

type StartPayload = {
	cpuThreshold: number
	ramThreshold: number
	duration: number
}

type Props = {
	isOpen: boolean
}

defineProps<Props>()
const emit = defineEmits<{
	(e: "close"): void
	(e: "start", payload: StartPayload): void
}>()

const cpuThreshold = defineModel<string>("cpuThreshold", { default: "70" })
const ramThreshold = defineModel<string>("ramThreshold", { default: "80" })
/**
 * Продолжительность в минутах для ui, в backend отправляю секунды.
 */
const durationMinutes = defineModel<string>("durationMinutes", { default: "5" })

const cpuThresholdNum = computed(
	() => Number.parseFloat(cpuThreshold.value) || 0
)
const ramThresholdNum = computed(
	() => Number.parseFloat(ramThreshold.value) || 0
)
const durationMinutesNum = computed(
	() => Number.parseInt(durationMinutes.value) || 0
)
const durationSecondsNum = computed(() => durationMinutesNum.value * 60)
const hasDurationInput = computed(
	() => String(durationMinutes.value ?? "").trim().length > 0
)

const isValid = computed(() => {
	return (
		cpuThresholdNum.value > 0 &&
		cpuThresholdNum.value <= 100 &&
		ramThresholdNum.value > 0 &&
		ramThresholdNum.value <= 100 &&
		durationMinutesNum.value >= 1
	)
})

const handleStart = () => {
	if (!isValid.value) return
	emit("start", {
		cpuThreshold: cpuThresholdNum.value,
		ramThreshold: ramThresholdNum.value,
		duration: durationSecondsNum.value,
	})
	emit("close")
}
</script>

<template>
	<BaseModal
		:is-open="isOpen"
		title="Настройка записи метрик"
		width="520px"
		@close="emit('close')"
	>
		<div class="metrics-recording-modal-content">
			<div class="metrics-recording-form-group">
				<label class="metrics-recording-label">Порог загрузки CPU (%)</label>
				<input
					v-model="cpuThreshold"
					type="text"
					class="metrics-recording-input"
					placeholder="0-100"
				/>
				<span class="metrics-recording-hint"
					>Процессы с загрузкой CPU выше этого значения будут записаны</span
				>
			</div>

			<div class="metrics-recording-form-group">
				<label class="metrics-recording-label"
					>Порог использования RAM (%)</label
				>
				<input
					v-model="ramThreshold"
					type="text"
					class="metrics-recording-input"
					placeholder="0-100"
				/>
				<span class="metrics-recording-hint"
					>Процессы с использованием RAM выше этого значения будут
					записаны</span
				>
			</div>

			<div class="metrics-recording-form-group">
				<label class="metrics-recording-label"
					>Продолжительность сессии (минуты)</label
				>
				<input
					v-model="durationMinutes"
					type="text"
					class="metrics-recording-input"
					placeholder="Минимум 1 минута"
				/>
				<span class="metrics-recording-hint">
					Время записи метрик (минимум 1 минута)
					<span
						v-if="hasDurationInput && durationMinutesNum < 1"
						class="metrics-recording-error"
					>
						• Минимум 1 минута</span
					>
				</span>
			</div>
		</div>

		<template #footer>
			<button
				class="metrics-recording-btn metrics-recording-btn--cancel"
				type="button"
				@click="emit('close')"
			>
				Отмена
			</button>
			<button
				class="metrics-recording-btn metrics-recording-btn--start"
				type="button"
				:disabled="!isValid"
				@click="handleStart"
			>
				Начать запись
			</button>
		</template>
	</BaseModal>
</template>
