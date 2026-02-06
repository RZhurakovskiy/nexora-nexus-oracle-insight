<script setup lang="ts">
import ConfirmModal from "@/components/confirm-modal/ConfirmModal.vue"
import ProcessDetailsModal from "@/components/process-details/ProcessDetailsModal.vue"
import "@/components/processes/Processes.css"
import ProcessSearchInput from "@/components/processes/ProcessSearchInput.vue"
import "@/components/processes/Table.css"
import WebSocketStatusIndicator from "@/components/websocket-status/WebSocketStatusIndicator.vue"
import { useMetricsStore } from "@/stores/metrics"
import { useProcessesStore } from "@/stores/processes"
import type { ProcessInfo } from "@/types/processes"
import { storeToRefs } from "pinia"
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue"

const PAGE_SIZE = 100

const metrics = useMetricsStore()
const processesStore = useProcessesStore()
const { processes, wsStatus, isInitialLoading, error } =
	storeToRefs(processesStore)

const showUnknown = ref(false)
const showFilters = ref(false)
const searchQueries = ref({ pid: "", name: "", port: "" })
const renderCount = ref(PAGE_SIZE)
const filters = ref({
	username: "",
	cpuMin: "",
	cpuMax: "",
	memoryMin: "",
	memoryMax: "",
})
const sortBy = ref<null | "cpu" | "memory">(null)
const sortOrder = ref<"asc" | "desc">("desc")

const selectedProcess = ref<ProcessInfo | null>(null)
const isDetailsOpen = ref(false)

const isConfirmOpen = ref(false)
const confirmPid = ref<number | null>(null)
const confirmTitle = ref("Подтверждение завершения процесса")
const confirmMessage = ref("")
const confirmText = ref("Завершить процесс")
const cancelText = ref("Отмена")

const isServerOffline = computed(() => {
	if (!metrics.monitoringEnabled) return true
	return (
		wsStatus.value === "error" || wsStatus.value === "disconnected"
	)
})

const effectiveStatus = computed(() => {
	return !metrics.monitoringEnabled && wsStatus.value === "connected"
		? "disconnected"
		: wsStatus.value
})

const isProcessUnknown = (p: ProcessInfo) => {
	const unknownName = !p?.name || p.name?.toLowerCase() === "неизвестно"
	const emptyMeta =
		(!p?.exe || p.exe === "") &&
		(!p?.cmdline || p.cmdline === "") &&
		(!p?.username || p.username === "") &&
		(!p?.createTime || p.createTime === 0)
	return unknownName || emptyMeta
}

const getCpuPercentNumber = (p: ProcessInfo) => {
	const v = Number(p?.cpuPercent ?? 0)
	if (Number.isNaN(v)) return 0
	return Math.max(0, v)
}

const getMemoryMb = (p: ProcessInfo) => {
	const bytes = Number(p?.memoryRss ?? 0)
	if (Number.isNaN(bytes) || bytes <= 0) return 0
	return bytes / (1024 * 1024)
}

const formatPorts = (p: ProcessInfo) => {
	const ports = p?.ports || []
	if (ports.length === 0) return "—"
	return ports.join(", ")
}

const isActiveProcess = (p: ProcessInfo) => {
	return getCpuPercentNumber(p) > 0 || getMemoryMb(p) > 0
}

const getRowSeverity = (p: ProcessInfo) => {
	const cpu = getCpuPercentNumber(p)
	const memMb = getMemoryMb(p)
	if (cpu >= 60 || memMb >= 800) return "critical"
	if (cpu >= 30 || memMb >= 400) return "elevated"
	return "normal"
}

const visibleRowsMeta = computed(() => {
	const unknown = processes.value.filter(isProcessUnknown)
	const known = processes.value.filter(process => !isProcessUnknown(process))

	const active = known.filter(isActiveProcess)
	const inactive = known.filter(process => !isActiveProcess(process))

	const sortByLoad = (a: ProcessInfo, b: ProcessInfo) => {
		const cpuDiff = getCpuPercentNumber(b) - getCpuPercentNumber(a)
		if (cpuDiff !== 0) return cpuDiff
		const memDiff = getMemoryMb(b) - getMemoryMb(a)
		if (memDiff !== 0) return memDiff
		return String(a.name || "").localeCompare(String(b.name || ""))
	}

	active.sort(sortByLoad)
	inactive.sort(sortByLoad)

	const base = [...active, ...inactive]
	const withUnknown = showUnknown.value ? [...base, ...unknown] : base

	return {
		visibleRows: withUnknown,
		unknownCount: unknown.length,
		activeCount: active.length,
	}
})

const hasUnknownProcess = computed(() => visibleRowsMeta.value.unknownCount > 0)

const filteredRows = computed(() => {
	let result = visibleRowsMeta.value.visibleRows

	if (searchQueries.value.pid.trim()) {
		const pidQuery = searchQueries.value.pid.trim().toLowerCase()
		result = result.filter(process => String(process?.pid || "").includes(pidQuery))
	}

	if (searchQueries.value.name.trim()) {
		const nameQuery = searchQueries.value.name.trim().toLowerCase()
		result = result.filter(process =>
			String(process?.name || "")
				.toLowerCase()
				.includes(nameQuery)
		)
	}

	if (searchQueries.value.port.trim()) {
		const portQuery = searchQueries.value.port.trim().toLowerCase()
		result = result.filter(process =>
			(process?.ports || []).some(port => String(port).includes(portQuery))
		)
	}

	if (filters.value.username) {
		const usernameFilter = filters.value.username.trim().toLowerCase()
		result = result.filter(process =>
			String(process?.username || "")
				.toLowerCase()
				.includes(usernameFilter)
		)
	}

	const cpuMin = filters.value.cpuMin ? parseFloat(filters.value.cpuMin) : null
	const cpuMax = filters.value.cpuMax ? parseFloat(filters.value.cpuMax) : null
	if (cpuMin !== null || cpuMax !== null) {
		result = result.filter(process => {
			const cpu = getCpuPercentNumber(process)
			if (cpuMin !== null && cpu < cpuMin) return false
			if (cpuMax !== null && cpu > cpuMax) return false
			return true
		})
	}

	const memMin = filters.value.memoryMin
		? parseFloat(filters.value.memoryMin)
		: null
	const memMax = filters.value.memoryMax
		? parseFloat(filters.value.memoryMax)
		: null
	if (memMin !== null || memMax !== null) {
		result = result.filter(process => {
			const mem = getMemoryMb(process)
			if (memMin !== null && mem < memMin) return false
			if (memMax !== null && mem > memMax) return false
			return true
		})
	}

	if (sortBy.value) {
		const column = sortBy.value
		result = [...result].sort((a, b) => {
			const aValue = column === "cpu" ? getCpuPercentNumber(a) : getMemoryMb(a)
			const bValue = column === "cpu" ? getCpuPercentNumber(b) : getMemoryMb(b)
			return sortOrder.value === "asc" ? aValue - bValue : bValue - aValue
		})
	}

	return result
})

watch(
	() => filteredRows.value.length,
	() => {
		renderCount.value = PAGE_SIZE
	}
)

const handleSort = (column: "cpu" | "memory") => {
	if (sortBy.value === column) {
		sortOrder.value = sortOrder.value === "asc" ? "desc" : "asc"
	} else {
		sortBy.value = column
		sortOrder.value = "desc"
	}
}

const tableScrollRef = ref<HTMLElement | null>(null)
const onScrollLoadMore = () => {
	const el = tableScrollRef.value
	if (!el) return
	const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 48
	if (nearBottom) {
		renderCount.value = Math.min(
			renderCount.value + PAGE_SIZE,
			filteredRows.value.length
		)
	}
}

const getBarTone = (value: number, elevated: number, critical: number) => {
	if (value >= critical) return "critical"
	if (value >= elevated) return "elevated"
	return "ok"
}

const openProcessDetails = (p: ProcessInfo) => {
	selectedProcess.value = p
	isDetailsOpen.value = true
}

const closeProcessDetails = () => {
	isDetailsOpen.value = false
	selectedProcess.value = null
}

const requestKillProcess = (p: ProcessInfo) => {
	const name = p?.name || "Неизвестный процесс"
	const pid = Number(p?.pid || 0)
	confirmPid.value = pid
	confirmMessage.value = `Вы уверены, что хотите завершить процесс?\n\nPID: ${pid}\nНазвание: ${name}\n\nЭто действие нельзя отменить.`
	isConfirmOpen.value = true
}

const confirmKill = async () => {
	const pid = confirmPid.value
	isConfirmOpen.value = false
	confirmPid.value = null
	if (!pid) return
	await processesStore.killProcess(pid)
}

const cancelKill = () => {
	isConfirmOpen.value = false
	confirmPid.value = null
}

const clearAllFilters = () => {
	filters.value = {
		username: "",
		cpuMin: "",
		cpuMax: "",
		memoryMin: "",
		memoryMax: "",
	}
}

const exportJson = async () => {
	await processesStore.exportProcesses("json")
}

const exportCsv = async () => {
	await processesStore.exportProcesses("csv")
}

onMounted(() => {
	processesStore.connectSockets()
})

onBeforeUnmount(() => {
	processesStore.disconnectSockets()
})

watch(
	() => metrics.monitoringEnabled,
	enabled => {
		/**
		 * Сервер пушит список процессов раз в 5 секунд.
		 * Когда пользователь включает мониторинг, мы не должны показывать empty-state
		 * до первого списка — поэтому включаем режим "ожидание данных".
		 */
		if (enabled) {
			processesStore.setWaitingForFirstPayload()
		}
	}
)
</script>

<template>
	<div class="processes-container">
		<div class="processes-header">
			<div class="processes-header-left">
				<h2 class="processes-title">Запущенные процессы</h2>
				<WebSocketStatusIndicator :status="effectiveStatus" label="Процессы" />
			</div>

			<div class="processes-toolbar">
				<div class="processes-search-group">
					<ProcessSearchInput
						v-model="searchQueries.pid"
						:processes="visibleRowsMeta.visibleRows"
						placeholder="Поиск по PID..."
						search-type="pid"
					/>
					<ProcessSearchInput
						v-model="searchQueries.name"
						:processes="visibleRowsMeta.visibleRows"
						placeholder="Поиск по имени..."
						search-type="name"
					/>
					<ProcessSearchInput
						v-model="searchQueries.port"
						:processes="visibleRowsMeta.visibleRows"
						placeholder="Поиск по порту..."
						search-type="port"
					/>
				</div>

				<button
					type="button"
					class="processes-filters-toggle"
					@click="showFilters = !showFilters"
					title="Показать/скрыть фильтры"
				>
					<svg width="16" height="16" viewBox="0 0 16 16" fill="none">
						<path
							d="M2 4H14M4 8H12M6 12H10"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
						/>
					</svg>
					Фильтры
					<span
						v-if="
							filters.username ||
							filters.cpuMin ||
							filters.cpuMax ||
							filters.memoryMin ||
							filters.memoryMax
						"
						class="processes-filters-badge"
					/>
				</button>

				<div class="processes-count">
					Активные: {{ visibleRowsMeta.activeCount }}
				</div>

				<div class="processes-export-group">
					<button
						type="button"
						class="processes-export-btn"
						:disabled="processesStore.exporting"
						title="Экспортировать процессы в JSON"
						@click="exportJson"
					>
						<svg width="16" height="16" viewBox="0 0 16 16" fill="none">
							<path
								d="M8 2V10M8 10L5 7M8 10L11 7M2 12V14C2 14.5523 2.44772 15 3 15H13C13.5523 15 14 14.5523 14 14V12"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
						</svg>
						JSON
					</button>
					<button
						type="button"
						class="processes-export-btn"
						:disabled="processesStore.exporting"
						title="Экспортировать процессы в CSV"
						@click="exportCsv"
					>
						<svg width="16" height="16" viewBox="0 0 16 16" fill="none">
							<path
								d="M8 2V10M8 10L5 7M8 10L11 7M2 12V14C2 14.5523 2.44772 15 3 15H13C13.5523 15 14 14.5523 14 14V12"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
						</svg>
						CSV
					</button>
				</div>

				<button
					v-if="hasUnknownProcess"
					type="button"
					class="toggle-unknown-btn"
					:title="
						showUnknown
							? 'Скрыть неизвестные процессы'
							: 'Показать скрытые процессы'
					"
					@click="showUnknown = !showUnknown"
				>
					{{
						showUnknown
							? "Скрыть неизвестные"
							: `Показать скрытые процессы (${visibleRowsMeta.unknownCount})`
					}}
				</button>
			</div>
		</div>

		<div v-if="showFilters" class="processes-filters-panel">
			<div class="processes-filters-header">
				<h3 class="processes-filters-title">Фильтры</h3>
				<button
					type="button"
					class="processes-filters-clear"
					title="Очистить все фильтры"
					@click="clearAllFilters"
				>
					Очистить все
				</button>
			</div>
			<div class="processes-filters-grid">
				<div class="processes-filter-group">
					<label class="processes-filter-label">Пользователь</label>
					<input
						v-model="filters.username"
						type="text"
						class="processes-filter-input"
						placeholder="Имя пользователя..."
					/>
				</div>

				<div class="processes-filter-group">
					<label class="processes-filter-label">CPU (%)</label>
					<div class="processes-filter-range">
						<input
							v-model="filters.cpuMin"
							type="number"
							class="processes-filter-input"
							placeholder="Мин"
							min="0"
							max="100"
							step="0.1"
						/>
						<span class="processes-filter-separator">—</span>
						<input
							v-model="filters.cpuMax"
							type="number"
							class="processes-filter-input"
							placeholder="Макс"
							min="0"
							max="100"
							step="0.1"
						/>
					</div>
				</div>

				<div class="processes-filter-group">
					<label class="processes-filter-label">Память (МБ)</label>
					<div class="processes-filter-range">
						<input
							v-model="filters.memoryMin"
							type="number"
							class="processes-filter-input"
							placeholder="Мин"
							min="0"
							step="0.1"
						/>
						<span class="processes-filter-separator">—</span>
						<input
							v-model="filters.memoryMax"
							type="number"
							class="processes-filter-input"
							placeholder="Макс"
							min="0"
							step="0.1"
						/>
					</div>
				</div>
			</div>
		</div>

		<div v-if="isServerOffline" class="cpu-chart-container">
			<div class="server-offline-message">
				<div class="server-offline-icon">
					<svg width="64" height="64" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
						<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none" opacity="0.15"/>
						<path d="M12 7V11M12 15H12.01" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
						<circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.5" fill="none"/>
					</svg>
				</div>
				<h2 class="server-offline-title">
					{{
						!metrics.monitoringEnabled && wsStatus === "connected"
							? "Мониторинг выключен"
							: "Сервер передачи данных выключен"
					}}
				</h2>
				<p class="server-offline-description">
					{{
						!metrics.monitoringEnabled && wsStatus === "connected"
							? "Мониторинг выключен. Включите мониторинг для получения данных."
							: "Не удается подключиться к серверу мониторинга."
					}}
					<br />
					<template
						v-if="!metrics.monitoringEnabled && wsStatus === 'connected'"
					>
						Используйте кнопку "Включить мониторинг" в верхней части страницы.
					</template>
				</p>
			</div>
		</div>

		<div v-else class="processes-card">
			<div
				v-if="
					isInitialLoading ||
					(metrics.monitoringEnabled && !processesStore.hasReceivedFirstPayload)
				"
				class="spinner-overlay"
			>
				<div class="spinner-system" />
				<div class="spinner-text">Загрузка процессов…</div>
			</div>

			<div
				v-if="
					!isInitialLoading &&
					processesStore.hasReceivedFirstPayload &&
					filteredRows.length === 0
				"
				class="empty-state"
			>
				<div class="empty-title">Процесс не найден</div>
				<div class="empty-subtitle">Измените запрос поиска.</div>
			</div>

			<div v-else class="processes-table-wrapper">
				<table class="processes-table processes-table--header">
					<thead>
						<tr>
							<th>PID</th>
							<th>Имя</th>
							<th>Путь</th>
							<th>Пользователь</th>
							<th>Порт</th>
							<th>
								<div class="table-header-with-sort">
									<span>CPU</span>
									<button
										class="table-header-sort-button"
										title="Сортировать по CPU"
										@click="handleSort('cpu')"
									>
										<svg
											v-if="sortBy !== 'cpu'"
											class="sort-icon sort-icon--inactive"
											width="16"
											height="16"
											viewBox="0 0 16 16"
											fill="none"
										>
											<path
												d="M8 3L11 6H9V13H7V6H5L8 3Z"
												fill="currentColor"
												opacity="0.4"
											/>
											<path
												d="M8 13L5 10H7V3H9V10H11L8 13Z"
												fill="currentColor"
												opacity="0.4"
											/>
										</svg>
										<svg
											v-else
											:class="`sort-icon sort-icon--${sortOrder}`"
											width="16"
											height="16"
											viewBox="0 0 16 16"
											fill="none"
										>
											<path
												v-if="sortOrder === 'asc'"
												d="M8 3L11 6H9V13H7V6H5L8 3Z"
												fill="currentColor"
											/>
											<path
												v-else
												d="M8 13L5 10H7V3H9V10H11L8 13Z"
												fill="currentColor"
											/>
										</svg>
									</button>
								</div>
							</th>
							<th>
								<div class="table-header-with-sort">
									<span>Память</span>
									<button
										class="table-header-sort-button"
										title="Сортировать по памяти"
										@click="handleSort('memory')"
									>
										<svg
											v-if="sortBy !== 'memory'"
											class="sort-icon sort-icon--inactive"
											width="16"
											height="16"
											viewBox="0 0 16 16"
											fill="none"
										>
											<path
												d="M8 3L11 6H9V13H7V6H5L8 3Z"
												fill="currentColor"
												opacity="0.4"
											/>
											<path
												d="M8 13L5 10H7V3H9V10H11L8 13Z"
												fill="currentColor"
												opacity="0.4"
											/>
										</svg>
										<svg
											v-else
											:class="`sort-icon sort-icon--${sortOrder}`"
											width="16"
											height="16"
											viewBox="0 0 16 16"
											fill="none"
										>
											<path
												v-if="sortOrder === 'asc'"
												d="M8 3L11 6H9V13H7V6H5L8 3Z"
												fill="currentColor"
											/>
											<path
												v-else
												d="M8 13L5 10H7V3H9V10H11L8 13Z"
												fill="currentColor"
											/>
										</svg>
									</button>
								</div>
							</th>
							<th>Действия</th>
						</tr>
					</thead>
				</table>

				<div
					ref="tableScrollRef"
					class="processes-table-body-wrapper"
					@scroll="onScrollLoadMore"
				>
					<table class="processes-table processes-table--body">
						<tbody>
							<tr
								v-for="row in filteredRows.slice(0, renderCount)"
								:key="row.pid"
								:class="`row--${getRowSeverity(row)}`"
							>
								<td>{{ row.pid }}</td>
								<td :title="row.cmdline || ''">{{ row.name || "—" }}</td>
								<td class="path-cell" :title="row.exe || ''">
									{{ row.exe || "—" }}
								</td>
								<td>{{ row.username || "—" }}</td>
								<td>{{ formatPorts(row) }}</td>
								<td>
									<div class="loadbar-wrap">
										<div class="loadbar">
											<div
												class="loadbar__fill"
												:class="`loadbar__fill--${getBarTone(
													getCpuPercentNumber(row),
													30,
													60
												)}`"
												:style="{
													width: `${Math.max(
														0,
														Math.min(getCpuPercentNumber(row), 100)
													)}%`,
												}"
											/>
										</div>
										<div class="loadbar__label--outside">
											{{ getCpuPercentNumber(row).toFixed(1) }} %
										</div>
									</div>
								</td>
								<td>
									<div class="loadbar-wrap">
										<div class="loadbar">
											<div
												class="loadbar__fill"
												:class="`loadbar__fill--${getBarTone(
													getMemoryMb(row),
													400,
													800
												)}`"
												:style="{
													width: `${Math.max(
														0,
														Math.min((getMemoryMb(row) / 16000) * 100, 100)
													)}%`,
												}"
											/>
										</div>
										<div class="loadbar__label--outside">
											{{ getMemoryMb(row).toFixed(1) }} MB
										</div>
									</div>
								</td>
								<td>
									<div class="process-actions">
										<button
											class="process-action-btn process-action-btn--view"
											title="Просмотреть детали процесса"
											@click="openProcessDetails(row)"
										>
											<svg
												width="16"
												height="16"
												viewBox="0 0 16 16"
												fill="none"
											>
												<path
													d="M8 2.5C4.5 2.5 1.73 4.61 0 7.5C1.73 10.39 4.5 12.5 8 12.5C11.5 12.5 14.27 10.39 16 7.5C14.27 4.61 11.5 2.5 8 2.5ZM8 10.5C6.34 10.5 5 9.16 5 7.5C5 5.84 6.34 4.5 8 4.5C9.66 4.5 11 5.84 11 7.5C11 9.16 9.66 10.5 8 10.5ZM8 6C7.17 6 6.5 6.67 6.5 7.5C6.5 8.33 7.17 9 8 9C8.83 9 9.5 8.33 9.5 7.5C9.5 6.67 8.83 6 8 6Z"
													fill="currentColor"
												/>
											</svg>
											<span>Посмотреть</span>
										</button>
										<button
											class="process-action-btn process-action-btn--kill"
											title="Завершить процесс"
											@click="requestKillProcess(row)"
										>
											Завершить
										</button>
									</div>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<div v-if="error" style="padding: 12px 14px; color: #ff453a">
				<b>Ошибка</b>: {{ error }}
			</div>
		</div>

		<div v-if="hasUnknownProcess && !showUnknown" class="admin-hint">
			<b>Часть процессов отображается как «неизвестно».</b> Для корректного
			получения списка процессов запустите сервер от имени администратора.
		</div>

		<ProcessDetailsModal
			:process="selectedProcess"
			:is-open="isDetailsOpen"
			@close="closeProcessDetails"
		/>

		<ConfirmModal
			:is-open="isConfirmOpen"
			:title="confirmTitle"
			:message="confirmMessage"
			:confirm-text="confirmText"
			:cancel-text="cancelText"
			@confirm="confirmKill"
			@cancel="cancelKill"
		/>
	</div>
</template>
