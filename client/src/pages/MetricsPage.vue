<script setup lang="ts">
import CPUChart from "@/components/cpu-chart/CPUChart.vue"

import MemoryChart from "@/components/memory-chart/MemoryChart.vue"
import Button from "@/components/ui/Button.vue"
import { useMetricsStore } from "@/stores/metrics"
import { useMetricsHistoryStore } from "@/stores/metricsHistory"
import { onBeforeUnmount, onMounted } from "vue"

const metrics = useMetricsStore()
const history = useMetricsHistoryStore()

onMounted(async () => {
	await metrics.init()
	/**
	 * WS держим на время нахождения на странице:
	 * - если мониторинг выключен, WS остаётся подключённым и UI показывает "Мониторинг выключен"
	 * - если сервер упал, RealtimeSocket будет пытаться переподключиться, и UI покажет "Сервер выключен"
	 */
	metrics.connectSockets()
})

onBeforeUnmount(() => {
	/**
	 * Отключаются только ws, без влияния на серверный флаг мониторинга,
	 * иначе уход со страницы будет выключать мониторинг для всей системы.
	 */
	metrics.disconnectSockets()
})
</script>

<template>
	<main style="padding: 16px">
		<div
			style="
				display: flex;
				align-items: center;
				justify-content: space-between;
				gap: 12px;
				margin: 0 0 16px 0;
			"
		>
			<h2 style="margin: 0">Системные метрики</h2>
			<Button
				variant="secondary"
				size="md"
				@click="
					;async () => {
						await history.fetchFromDb()
						history.exportLoaded('json')
					}
				"
			>
				Выгрузить из БД
			</Button>
		</div>

		<section style="display: grid; gap: 16px">
			<CPUChart />
			<MemoryChart />

			<div v-if="metrics.error" style="color: #b00020">
				<b>Ошибка</b>: {{ metrics.error }}
			</div>
		</section>
	</main>
</template>
