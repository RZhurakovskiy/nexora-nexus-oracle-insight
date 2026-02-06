<script setup lang="ts">
import WebSocketStatusIndicator from "@/components/websocket-status/WebSocketStatusIndicator.vue"
import { ensureChartJsRegistered } from "@/lib/chartjs"
import { useMetricsStore } from "@/stores/metrics"
import { Chart } from "chart.js"
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue"
import "./memotyChart.css"

ensureChartJsRegistered()

const metrics = useMetricsStore()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let chart: Chart<"line"> | null = null

const isServerOffline = computed(() => {
	if (!metrics.monitoringEnabled) return true
	return (
		metrics.memoryWsStatus === "error" ||
		metrics.memoryWsStatus === "disconnected"
	)
})

const effectiveStatus = computed(() => {
	return !metrics.monitoringEnabled && metrics.memoryWsStatus === "connected"
		? "disconnected"
		: metrics.memoryWsStatus
})

const latestMemoryPercent = computed(() => {
	const last =
		metrics.memorySeries.length > 0
			? metrics.memorySeries[metrics.memorySeries.length - 1]
			: null
	return typeof last?.memoryUsage === "number" ? last.memoryUsage : 0
})

const getStatusColor = (value: number) => {
	if (value > 85) return "#ff453a"
	if (value > 65) return "#ff9500"
	return "#30d158"
}

const getStatusText = (value: number) => {
	if (value >= 85) return "Критическая"
	if (value >= 65) return "Высокая"
	if (value >= 40) return "Средняя"
	return "Низкая"
}

const formatMemoryValue = (valueInMB: number | null | undefined) => {
	if (valueInMB == null) return "—"
	if (valueInMB >= 1024) {
		const inGB = valueInMB / 1024
		return `${inGB.toFixed(inGB >= 10 ? 0 : 1)} ГБ`
	}
	return `${Math.round(valueInMB)} МБ`
}

const createChart = () => {
	if (!canvasRef.value) return
	if (chart) return

	const existing = Chart.getChart(canvasRef.value)
	if (existing) {
		existing.destroy()
	}

	chart = new Chart(canvasRef.value, {
		type: "line",
		data: {
			labels: [],
			datasets: [
				{
					label: "Использование ОЗУ (%)",
					data: [],
					borderColor: "rgb(0, 122, 255)",
					backgroundColor: "rgba(0, 122, 255, 0.15)",
					tension: 0.4,
					fill: true,
					pointRadius: 0,
					pointHoverRadius: 4,
					borderWidth: 2,
				},
			],
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			animation: false,
			interaction: {
				mode: "index",
				intersect: false,
			},
			plugins: {
				legend: {
					display: true,
					position: "top",
					labels: {
						color: "#e5e5e7",
						font: {
							size: 12,
							weight: 500,
						},
						padding: 12,
						usePointStyle: true,
					},
				},
				title: {
					display: true,
					text: "Использование ОЗУ",
					color: "#f5f5f7",
					font: {
						size: 16,
						weight: 600,
					},
					padding: {
						bottom: 20,
					},
				},
				tooltip: {
					backgroundColor: "rgba(28, 28, 30, 0.95)",
					titleColor: "#f5f5f7",
					bodyColor: "#e5e5e7",
					borderColor: "rgba(255, 255, 255, 0.1)",
					borderWidth: 1,
					padding: 12,
					cornerRadius: 8,
					displayColors: true,
					titleFont: {
						size: 13,
						weight: 600,
					},
					bodyFont: {
						size: 12,
					},
					callbacks: {
						label: context => {
							const point = metrics.memorySeries[context.dataIndex]
							const percent = context.parsed.y
							if (point) {
								const used = formatMemoryValue(point.usedMB)
								const total = formatMemoryValue(point.totalMemoryMB)
								return `Использовано: ${used} из ${total} (${Number(
									percent
								).toFixed(1)}%)`
							}
							return `Использование: ${Number(percent).toFixed(1)}%`
						},
					},
				},
			},
			scales: {
				x: {
					border: {
						display: false,
					},
					grid: {
						color: "rgba(255, 255, 255, 0.05)",
					},
					ticks: {
						color: "#8e8e93",
						font: {
							size: 11,
						},
						maxRotation: 45,
						minRotation: 45,
					},
					title: {
						display: true,
						text: "Время",
						color: "#8e8e93",
						font: {
							size: 12,
							weight: 500,
						},
					},
				},
				y: {
					border: {
						display: false,
					},
					grid: {
						color: "rgba(255, 255, 255, 0.05)",
					},
					ticks: {
						color: "#8e8e93",
						font: {
							size: 11,
						},
					},
					title: {
						display: true,
						text: "ОЗУ (%)",
						color: "#8e8e93",
						font: {
							size: 12,
							weight: 500,
						},
					},
					min: 0,
					max: 100,
				},
			},
		},
	})
}

const destroyChart = () => {
	chart?.destroy()
	chart = null
}

const syncChartData = () => {
	if (isServerOffline.value) return
	if (!chart) createChart()
	if (!chart) return
	chart.data.labels = metrics.memorySeries.map(point => point.timestamp)
	chart.data.datasets[0]!.data = metrics.memorySeries.map(point => point.memoryUsage)
	chart.update("none")
}

onMounted(() => {
	watch(
		[() => canvasRef.value, () => isServerOffline.value],
		async ([canvas, offline]) => {
			if (offline || !canvas) {
				destroyChart()
				return
			}
			await nextTick()
			createChart()
			syncChartData()
		},
		{ immediate: true, flush: "post" }
	)

	watch(
		() => metrics.memorySeries,
		() => syncChartData(),
		{ deep: false }
	)
})

onBeforeUnmount(() => {
	destroyChart()
})
</script>

<template>
	<div class="cpu-chart-container">
		<div class="cpu-chart-header">
			<WebSocketStatusIndicator :status="effectiveStatus" label="Память" />
		</div>

		<div v-if="isServerOffline" class="server-offline-message">
			<div class="server-offline-icon">
				<svg width="64" height="64" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
					<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none" opacity="0.15"/>
					<path d="M12 7V11M12 15H12.01" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
					<circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.5" fill="none"/>
				</svg>
			</div>
			<h2 class="server-offline-title">
				{{
					!metrics.monitoringEnabled && metrics.memoryWsStatus === "connected"
						? "Мониторинг выключен"
						: "Сервер передачи данных выключен"
				}}
			</h2>
			<p class="server-offline-description">
				{{
					!metrics.monitoringEnabled && metrics.memoryWsStatus === "connected"
						? "Мониторинг выключен. Включите мониторинг для получения данных."
						: "Не удается подключиться к серверу мониторинга."
				}}
				<br />
				<template
					v-if="
						!metrics.monitoringEnabled && metrics.memoryWsStatus === 'connected'
					"
				>
					Используйте кнопку "Включить мониторинг" в верхней части страницы.
				</template>
				<template v-else>
					Убедитесь, что сервер запущен и доступен по адресу
					<code>http://localhost:8080</code>.
				</template>
			</p>
		</div>

		<div v-else class="cpu-chart-content">
			<div class="cpu-chart-wrapper">
				<canvas ref="canvasRef" />
			</div>

			<div class="cpu-indicator">
				<div
					class="cpu-indicator-value"
					:style="{ color: getStatusColor(latestMemoryPercent) }"
				>
					{{ latestMemoryPercent.toFixed(1) }}%
				</div>
				<div class="cpu-indicator-label">Текущая нагрузка памяти</div>
				<div
					class="cpu-indicator-status"
					:style="{ color: getStatusColor(latestMemoryPercent) }"
				>
					{{ getStatusText(latestMemoryPercent) }}
				</div>

				<div class="cpu-indicator-details">
					<template v-if="metrics.latestMemory">
						<div class="memory-detail memory-detail-used">
							<span class="memory-detail-label">Использовано</span>
							<span class="memory-detail-value">{{
								formatMemoryValue(metrics.latestMemory.usedMB)
							}}</span>
						</div>
						<div class="memory-detail memory-detail-total">
							<span class="memory-detail-label">Всего</span>
							<span class="memory-detail-value">{{
								formatMemoryValue(metrics.latestMemory.totalMemoryMB)
							}}</span>
						</div>
					</template>
					<template v-else>
						<div>Данные ещё загружаются…</div>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>
