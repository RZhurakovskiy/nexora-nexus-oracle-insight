<script setup lang="ts">
import BaseModal from "@/components/modal/BaseModal.vue"
import type { ProcessInfo } from "@/types/processes"
import { computed } from "vue"
import "./ProcessDetailsModal.css"

type Props = {
	process: ProcessInfo | null
	isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{ (e: "close"): void }>()

const cpuPercent = computed(() => {
	const v = Number(props.process?.cpuPercent ?? 0)
	return Number.isFinite(v) ? Math.max(0, v) : 0
})

const memoryMb = computed(() => {
	const bytes = Number(props.process?.memoryRss ?? 0)
	if (!Number.isFinite(bytes) || bytes <= 0) return 0
	return bytes / (1024 * 1024)
})

const formatDate = (timestamp: number | null | undefined): string => {
	if (!timestamp) return "—"
	try {
		const ts = Number(timestamp)
		if (Number.isNaN(ts)) return "—"
		const date = new Date(ts < 1e12 ? ts * 1000 : ts)
		if (Number.isNaN(date.getTime())) return "—"
		return date.toLocaleString("ru-RU", {
			year: "numeric",
			month: "2-digit",
			day: "2-digit",
			hour: "2-digit",
			minute: "2-digit",
			second: "2-digit",
		})
	} catch {
		return "—"
	}
}

const formatBytes = (bytes: number | null | undefined): string => {
	if (!bytes || Number.isNaN(bytes)) return "—"
	const units = ["Б", "КБ", "МБ", "ГБ", "ТБ"]
	let value = Number(bytes)
	let unitIndex = 0
	while (value >= 1024 && unitIndex < units.length - 1) {
		value /= 1024
		unitIndex++
	}
	return `${value.toFixed(value >= 10 ? 0 : 1)} ${units[unitIndex]}`
}

const formatPorts = (ports: number[] | null | undefined): string => {
	if (!ports || ports.length === 0) return "—"
	return ports.join(", ")
}
</script>

<template>
	<BaseModal
		v-if="process"
		:is-open="isOpen"
		title="Информация о процессе"
		width="680px"
		@close="emit('close')"
	>
		<template #header>
			<div class="process-details-header-info">
				<p class="process-details-subtitle">
					Пользователь: <strong>{{ process.username || "—" }}</strong>
				</p>
			</div>
		</template>

		<div class="process-details-content">
			<section class="process-details-section">
				<h3 class="process-details-section-title">Основные сведения</h3>
				<div class="process-details-grid">
					<div class="process-details-item">
						<span class="process-details-label">PID</span>
						<span class="process-details-value">{{ process.pid }}</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Имя процесса</span>
						<span class="process-details-value">{{
							process.name || "—"
						}}</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">PPID</span>
						<span class="process-details-value">{{
							process.parentPid || "—"
						}}</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Время запуска</span>
						<span class="process-details-value">{{
							formatDate(process.createTime)
						}}</span>
					</div>
				</div>
			</section>

			<section class="process-details-section">
				<h3 class="process-details-section-title">Нагрузка</h3>
				<div class="process-details-grid">
					<div class="process-details-item">
						<span class="process-details-label">CPU</span>
						<span
							class="process-details-value process-details-value--highlight"
						>
							{{ cpuPercent.toFixed(2) }} %
						</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Память</span>
						<span
							class="process-details-value process-details-value--highlight"
						>
							{{ memoryMb.toFixed(1) }} МБ ({{
								formatBytes(process.memoryRss)
							}})
						</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Память (%)</span>
						<span class="process-details-value"
							>{{ (process.memoryPercent ?? 0).toFixed(2) }} %</span
						>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Порты</span>
						<span class="process-details-value">{{
							formatPorts(process.ports)
						}}</span>
					</div>
				</div>
			</section>

			<section class="process-details-section">
				<h3 class="process-details-section-title">Пути и команды</h3>
				<div class="process-details-grid process-details-grid--single">
					<div class="process-details-item">
						<span class="process-details-label"
							>Путь к исполняемому файлу</span
						>
						<span class="process-details-value process-details-value--code">{{
							process.exe || "—"
						}}</span>
					</div>
					<div class="process-details-item">
						<span class="process-details-label">Командная строка</span>
						<span class="process-details-value process-details-value--code">{{
							process.cmdline || "—"
						}}</span>
					</div>
				</div>
			</section>

			<section class="process-details-section">
				<h3 class="process-details-section-title">Дополнительно</h3>
				<div class="process-details-grid">
					<div class="process-details-item">
						<span class="process-details-label">Статус</span>
						<span class="process-details-value">{{
							process.status || "—"
						}}</span>
					</div>
				</div>
			</section>
		</div>
	</BaseModal>
</template>
