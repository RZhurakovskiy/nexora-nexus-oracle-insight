<script setup lang="ts">
import ConfirmModal from "@/components/confirm-modal/ConfirmModal.vue"
import "@/pages/metrics-history/MetricsHistory.css"
import { useMetricsHistoryStore } from "@/stores/metricsHistory"
import { useRecordingStore } from "@/stores/recording"
import { storeToRefs } from "pinia"
import { onMounted, ref, watch } from "vue"

const historyStore = useMetricsHistoryStore()
const recording = useRecordingStore()

const { items, loading, error, fromDate, toDate, limit } =
	storeToRefs(historyStore)
const { notifiedSessionEndId } = storeToRefs(recording)

const isConfirmOpen = ref(false)
const confirmTitle = ref("Подтверждение очистки метрик")
const confirmMessage = ref(
	"Вы уверены, что хотите удалить все метрики из базы данных?\n\nЭто действие нельзя отменить."
)
const confirmText = ref("Очистить")
const cancelText = ref("Отмена")

const openClearConfirm = () => {
	isConfirmOpen.value = true
}

const confirmClear = async () => {
	isConfirmOpen.value = false
	await historyStore.clearAll()
}

const cancelClear = () => {
	isConfirmOpen.value = false
}

const formatDate = (dateStr: string) => {
	if (!dateStr) return "—"
	return new Date(dateStr.replace(" ", "T")).toLocaleString("ru-RU")
}

onMounted(async () => {
	await historyStore.load()
})

watch(
	() => notifiedSessionEndId.value,
	async id => {
		if (!id) return
		await historyStore.fetchFromDb()
	}
)
</script>

<template>
	<div class="metrics-history-container">
		<div class="metrics-history-header">
			<div>
				<h2 class="metrics-history-title">История метрик</h2>
			</div>

			<div class="metrics-history-filters">
				<div class="filter-group">
					<label>От:</label>
					<div class="datetime-input-wrapper">
						<input
							v-model="fromDate"
							type="datetime-local"
							class="datetime-input"
						/>
					</div>
				</div>

				<div class="filter-group">
					<label>До:</label>
					<div class="datetime-input-wrapper">
						<input
							v-model="toDate"
							type="datetime-local"
							class="datetime-input"
						/>
					</div>
				</div>

				<div class="filter-group">
					<label>Лимит:</label>
					<input v-model.number="limit" type="number" min="1" max="10000" />
				</div>

				<button
					class="filter-apply-btn"
					type="button"
					:disabled="loading"
					@click="historyStore.load()"
				>
					Применить
				</button>

				<button
					class="filter-apply-btn"
					type="button"
					:disabled="loading"
					@click="historyStore.fetchFromDb()"
				>
					Выгрузить данные из БД
				</button>

				<button
					class="filter-apply-btn"
					type="button"
					:disabled="loading"
					@click="historyStore.exportLoaded('json')"
				>
					Экспорт JSON
				</button>

				<button
					class="filter-clear-btn"
					type="button"
					:disabled="loading"
					title="Очистить все метрики из базы данных"
					@click="openClearConfirm"
				>
					<svg
						width="16"
						height="16"
						viewBox="0 0 16 16"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M2 4H14M5 4V3C5 2.44772 5.44772 2 6 2H10C10.5523 2 11 2.44772 11 3V4M13 4V13C13 13.5523 12.5523 14 12 14H4C3.44772 14 3 13.5523 3 13V4H13Z"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
							stroke-linejoin="round"
						/>
					</svg>
					Очистить метрики
				</button>
			</div>
		</div>

		<div v-if="loading" class="loading-message">Загрузка истории метрик...</div>
		<div v-else-if="error" class="error-message">{{ error }}</div>
		<div v-else-if="items.length === 0" class="empty-message">
			{{
				fromDate || toDate
					? "Записи не найдены по заданным фильтрам"
					: "История метрик пуста"
			}}
		</div>
		<div v-else class="metrics-history-table-wrapper">
			<table class="metrics-history-table">
				<thead>
					<tr>
						<th>Время</th>
						<th>CPU %</th>
						<th>Память %</th>
						<th>Память использовано (МБ)</th>
						<th>Память всего (МБ)</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="(item, idx) in items" :key="`${item.timestamp}_${idx}`">
						<td>{{ formatDate(item.timestamp) }}</td>
						<td>{{ (item.cpuPercent ?? 0).toFixed(2) }}%</td>
						<td>{{ (item.memoryPercent ?? 0).toFixed(2) }}%</td>
						<td>{{ item.memoryUsedMB ?? "—" }}</td>
						<td>{{ item.memoryTotalMB ?? "—" }}</td>
					</tr>
				</tbody>
			</table>
			<div class="metrics-history-info">Всего записей: {{ items.length }}</div>
		</div>

		<ConfirmModal
			:is-open="isConfirmOpen"
			:title="confirmTitle"
			:message="confirmMessage"
			:confirm-text="confirmText"
			:cancel-text="cancelText"
			@confirm="confirmClear"
			@cancel="cancelClear"
		/>
	</div>
</template>
