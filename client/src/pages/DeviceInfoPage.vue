<script setup lang="ts">
import { useDeviceInfoStore } from "@/stores/deviceInfo"
import { useDiskHealthStore } from "@/stores/diskHealth"
import { useRootStatusStore } from "@/stores/rootStatus"
import { useSystemInfoStore } from "@/stores/systemInfo"
import { storeToRefs } from "pinia"
import { computed, onMounted } from "vue"
const rootStatusStore = useRootStatusStore()
const deviceInfoStore = useDeviceInfoStore()
const { data: deviceInfo, loading, error } = storeToRefs(deviceInfoStore)

const systemInfoStore = useSystemInfoStore()
const {
	data: systemInfo,
	loading: systemLoading,
	error: systemError,
} = storeToRefs(systemInfoStore)

const diskHealthStore = useDiskHealthStore()
const {
	data: diskHealth,
	loading: diskHealthLoading,
	error: diskHealthError,
} = storeToRefs(diskHealthStore)

/** Словарь описаний флагов CPU */
const flagDescriptions: Record<string, string> = {
	sse: "Streaming SIMD Extensions - базовые инструкции для параллельной обработки данных (128-битные регистры)",
	sse2: "SSE2 - расширение SSE с поддержкой 64-битных операций с плавающей точкой",
	ssse3:
		"SSSE3 - Supplemental SSE3, дополнительные инструкции для обработки текста и криптографии",
	sse3: "SSE3 - расширение SSE с новыми инструкциями для обработки чисел",
	sse4_1: "SSE4.1 - инструкции для ускорения обработки графики, видео и текста",
	sse4_2: "SSE4.2 - инструкции для строковых операций и обработки текста",
	avx: "Advanced Vector Extensions - 256-битные регистры для ускорения вычислений с плавающей точкой",
	avx2: "AVX2 - расширение AVX с поддержкой 256-битных целочисленных операций",
	aes: "AES-NI - аппаратное ускорение алгоритма шифрования AES",
	fma: "Fused Multiply-Add - объединённая операция умножения и сложения для повышения точности",
	mmx: "MultiMedia eXtensions - базовые инструкции для обработки мультимедиа",
	vmx: "Intel Virtualization Technology - виртуализация Intel",
	svm: "AMD-V - виртуализация AMD",
}

const formatFrequency = (mhz: number): string => {
	if (mhz >= 1000) {
		return `${(mhz / 1000).toFixed(2)} GHz`
	}
	return `${mhz} MHz`
}

const formatCacheSize = (kb: number): string => {
	if (kb >= 1024) {
		return `${(kb / 1024).toFixed(2)} MB`
	}
	return `${kb} KB`
}

const getFlagDescription = (flag: string): string => {
	return flagDescriptions[flag.toLowerCase()] || "Инструкция процессора"
}

/**
 * Преобразовать байты.
 */
const formatBytes = (bytes: number): string => {
	const b = Number(bytes ?? 0)
	if (!Number.isFinite(b) || b <= 0) return "—"
	const units = ["Б", "КБ", "МБ", "ГБ", "ТБ"]
	let v = b
	let i = 0
	while (v >= 1024 && i < units.length - 1) {
		v /= 1024
		i += 1
	}
	return `${v.toFixed(v >= 10 ? 0 : 1)} ${units[i]}`
}

/**
 * Преобразовать uptime секунд в строку вида "HH:MM".
 */
const formatUptime = (uptimeSec: number): string => {
	const s = Math.max(0, Math.floor(Number(uptimeSec ?? 0)))
	if (!Number.isFinite(s) || s <= 0) return "—"
	const days = Math.floor(s / 86400)
	const hh = Math.floor((s % 86400) / 3600)
	const mm = Math.floor((s % 3600) / 60)
	return `${days > 0 ? `${days}д ` : ""}${String(hh).padStart(2, "0")}:${String(
		mm
	).padStart(2, "0")}`
}

const memorySummary = computed(() => {
	if (!systemInfo.value) return null
	const m = systemInfo.value.memory
	return {
		total: formatBytes(m.totalBytes),
		used: formatBytes(m.usedBytes),
		available: formatBytes(m.availableBytes),
		usedPercent: Number(m.usedPercent ?? 0),
		swap:
			m.swapTotalBytes > 0
				? `${formatBytes(m.swapUsedBytes)} / ${formatBytes(m.swapTotalBytes)}`
				: "—",
	}
})

/**
 * Форматировать часы работы.
 */
const formatHours = (hours: number | null | undefined): string => {
	const h = Number(hours ?? 0)
	if (!Number.isFinite(h) || h <= 0) return "—"
	const days = Math.floor(h / 24)
	const rest = h % 24
	return `${days > 0 ? `${days}д ` : ""}${rest}ч`
}

/**
 * Получение S.M.A.R.T. / health и получение статуса - rootStatus для условного отображения раздела.
 */
const getDiskHealth = () => {
	diskHealthStore.refresh()
	rootStatusStore.init()
}

onMounted(async () => {
	await deviceInfoStore.init()
	await systemInfoStore.init()
	await diskHealthStore.init()
})
</script>

<template>
	<main class="device-info-page">
		<h1 class="page-title">Информация о системе</h1>

		<div v-if="loading || systemLoading" class="loading">
			Загрузка информации...
		</div>

		<div v-else-if="error || systemError" class="error">
			<b>Ошибка</b>: {{ error || systemError }}
		</div>

		<div v-else class="device-info-container">
			<section class="info-section">
				<h2 class="section-title">Операционная система</h2>
				<div v-if="systemInfo?.host" class="info-grid">
					<div class="info-item">
						<span class="info-label">Хост:</span>
						<span class="info-value">{{ systemInfo.host.hostname }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Платформа:</span>
						<span class="info-value"
							>{{ systemInfo.host.platform }}
							{{ systemInfo.host.platformVersion }}</span
						>
					</div>
					<div class="info-item">
						<span class="info-label">Ядро:</span>
						<span class="info-value"
							>{{ systemInfo.host.kernelVersion }} ({{
								systemInfo.host.kernelArch
							}})</span
						>
					</div>
					<div class="info-item">
						<span class="info-label">Uptime:</span>
						<span class="info-value">{{
							formatUptime(systemInfo.host.uptimeSec)
						}}</span>
					</div>
					<div v-if="systemInfo.load" class="info-item">
						<span class="info-label">Load avg:</span>
						<span class="info-value"
							>{{ systemInfo.load.load1.toFixed(2) }} /
							{{ systemInfo.load.load5.toFixed(2) }} /
							{{ systemInfo.load.load15.toFixed(2) }}</span
						>
					</div>
				</div>
				<div v-else class="info-item">
					<span class="info-label">Информация об ОС</span>
					<span class="info-value">Данные об ОС недоступны</span>
				</div>
			</section>

			<section v-if="systemInfo?.memory && memorySummary" class="info-section">
				<h2 class="section-title">Память</h2>
				<div class="info-grid">
					<div class="info-item">
						<span class="info-label">Всего:</span>
						<span class="info-value">{{ memorySummary.total }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Использовано:</span>
						<span class="info-value"
							>{{ memorySummary.used }} ({{
								memorySummary.usedPercent.toFixed(1)
							}}%)</span
						>
					</div>
					<div class="info-item">
						<span class="info-label">Доступно:</span>
						<span class="info-value">{{ memorySummary.available }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Swap:</span>
						<span class="info-value">{{ memorySummary.swap }}</span>
					</div>
				</div>
			</section>

			<section class="info-section">
				<h2 class="section-title">Диски и разделы</h2>
				<template v-if="systemInfo?.disks && systemInfo.disks.length > 0">
					<div
						v-for="disk in systemInfo.disks"
						:key="`${disk.device}:${disk.mountpoint}`"
						style="margin-bottom: 24px"
					>
						<h3
							style="
								font-size: 16px;
								font-weight: 600;
								color: #f5f5f7;
								margin: 0 0 16px 0;
								padding-bottom: 8px;
								border-bottom: 1px solid rgba(255, 255, 255, 0.1);
							"
						>
							{{ disk.mountpoint }}
						</h3>
						<div class="info-grid">
							<div class="info-item">
								<span class="info-label">Устройство:</span>
								<span class="info-value">{{ disk.device }}</span>
							</div>
							<div class="info-item">
								<span class="info-label">Файловая система:</span>
								<span class="info-value">{{ disk.fstype }}</span>
							</div>
							<div class="info-item">
								<span class="info-label">Всего:</span>
								<span class="info-value">{{
									formatBytes(disk.totalBytes)
								}}</span>
							</div>
							<div class="info-item">
								<span class="info-label">Использовано:</span>
								<span class="info-value"
									>{{ formatBytes(disk.usedBytes) }} ({{
										disk.usedPercent.toFixed(1)
									}}%)</span
								>
							</div>
							<div class="info-item">
								<span class="info-label">Свободно:</span>
								<span class="info-value">{{
									formatBytes(disk.freeBytes)
								}}</span>
							</div>
						</div>
					</div>
				</template>
				<div v-else class="info-item">
					<span class="info-label">Информация о дисках</span>
					<span class="info-value">Данные о дисках недоступны</span>
				</div>
			</section>

			<section class="info-section">
				<div
					style="
						display: flex;
						align-items: center;
						justify-content: space-between;
						gap: 12px;
						flex-wrap: wrap;
					"
				>
					<h2 class="section-title" style="margin: 0; border: 0; padding: 0">
						S.M.A.R.T. / Health
					</h2>
					<button
						type="button"
						style="
							padding: 10px 14px;
							border-radius: 10px;
							border: 1px solid rgba(255, 255, 255, 0.12);
							background: rgba(92, 141, 255, 0.16);
							color: rgba(245, 245, 247, 0.95);
							font-weight: 800;
							cursor: pointer;
						"
						:disabled="diskHealthLoading"
						@click="getDiskHealth()"
					>
						{{ diskHealthLoading ? "Обновление…" : "Обновить" }}
					</button>
				</div>

				<div v-if="diskHealthError" class="error" style="margin-top: 12px">
					<b>Ошибка</b>: {{ diskHealthError }}
				</div>

				<div
					v-else-if="diskHealth && !diskHealth.supported"
					class="info-item"
					style="margin-top: 12px"
				>
					<span class="info-label">Недоступно</span>
					<span class="info-value">{{ diskHealth.message }}</span>
				</div>

				<div
					v-else-if="diskHealth && rootStatusStore.rootStatus"
					style="margin-top: 12px"
				>
					<div class="info-item" style="margin-bottom: 12px">
						<span class="info-label">Источник</span>
						<span class="info-value">{{ diskHealth.message }}</span>
					</div>

					<div v-if="diskHealth.devices.length === 0" class="info-item">
						<span class="info-label">Диски</span>
						<span class="info-value">Не найдено устройств для S.M.A.R.T.</span>
					</div>

					<div v-else class="info-grid">
						<div
							v-for="device in diskHealth.devices"
							:key="device.device"
							class="info-item"
						>
							<span class="info-label"
								>{{ device.device }}
								<span v-if="device.busType">({{ device.busType }})</span></span
							>
							<span class="info-value">
								<span
									v-if="device.smartPassed === true"
									style="color: rgba(96, 255, 165, 0.95); font-weight: 900"
									>OK</span
								>
								<span
									v-else-if="device.smartPassed === false"
									style="color: rgba(255, 69, 58, 0.95); font-weight: 900"
									>FAILED</span
								>
								<span
									v-else
									style="color: rgba(142, 142, 147, 0.9); font-weight: 800"
									>—</span
								>
								<span style="opacity: 0.9"> • {{ device.model || "—" }}</span>
							</span>
							<div class="info-hint">
								Serial: {{ device.serial || "—" }}<br />
								Temp:
								{{
									device.temperatureC != null ? `${device.temperatureC}°C` : "—"
								}}
								• Power-on: {{ formatHours(device.powerOnHours) }}<br />
								<span v-if="device.nvmePercentUsed != null"
									>NVMe wear: {{ device.nvmePercentUsed }}%</span
								>
								<span v-else>NVMe wear: —</span>
								<span v-if="device.unsafeShutdowns != null">
									• Unsafe shutdowns: {{ device.unsafeShutdowns }}</span
								>
							</div>
							<div
								v-if="device.warnings && device.warnings.length > 0"
								class="info-hint"
								style="color: rgba(255, 149, 0, 0.95); margin-top: 8px"
							>
								<b>Предупреждения:</b>
								<ul style="margin: 6px 0 0 16px; padding: 0">
									<li v-for="(warning, idx) in device.warnings" :key="idx">
										{{ warning }}
									</li>
								</ul>
							</div>
						</div>
					</div>
				</div>
				<div class="alert-message" v-if="!rootStatusStore.rootStatus">
					<span>
						Приложение запущено без root прав. Для просмотра раздела запустите
						приложение с root правами
					</span>
				</div>
			</section>

			<!-- Основная информация -->
			<section v-if="deviceInfo" class="info-section">
				<h2 class="section-title">Процессор</h2>
				<div class="info-grid">
					<div class="info-item">
						<span class="info-label">Процессор:</span>
						<span class="info-value">{{ deviceInfo.processor_name }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Производитель:</span>
						<span class="info-value">{{ deviceInfo.vendor }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Архитектура:</span>
						<span class="info-value">{{ deviceInfo.architecture }}</span>
					</div>
				</div>
			</section>

			<!-- Характеристики производительности -->
			<section v-if="deviceInfo" class="info-section">
				<h2 class="section-title">Характеристики производительности</h2>
				<div class="info-grid">
					<div class="info-item">
						<span class="info-label">Физические ядра:</span>
						<span class="info-value">{{ deviceInfo.physical_cores }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Логические процессоры:</span>
						<span class="info-value">{{ deviceInfo.logical_processors }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Частота:</span>
						<span class="info-value">{{
							formatFrequency(deviceInfo.frequency_mhz)
						}}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Кэш:</span>
						<span class="info-value">{{
							formatCacheSize(deviceInfo.cache_size_kb)
						}}</span>
					</div>
				</div>
			</section>

			<!-- Технические характеристики -->
			<section v-if="deviceInfo" class="info-section">
				<h2 class="section-title">Технические характеристики</h2>
				<div class="info-grid">
					<div class="info-item">
						<span class="info-label">Семейство:</span>
						<span class="info-value">{{ deviceInfo.family }}</span>
					</div>
					<div class="info-item">
						<span class="info-label">Модель:</span>
						<span class="info-value">{{ deviceInfo.model }}</span>
					</div>
				</div>
			</section>

			<!-- Поддерживаемые инструкции из словаря -->
			<section v-if="deviceInfo" class="info-section flags-section">
				<h2 class="section-title">Поддерживаемые инструкции процессора</h2>
				<div
					v-if="
						deviceInfo.supported_flags && deviceInfo.supported_flags.length > 0
					"
					class="flags-grid"
				>
					<div
						v-for="flag in deviceInfo.supported_flags"
						:key="flag"
						class="flag-item"
					>
						<div class="flag-header">
							<span class="flag-name">{{ flag.toUpperCase() }}</span>
						</div>
						<p class="flag-description">{{ getFlagDescription(flag) }}</p>
					</div>
				</div>
				<div v-else class="no-flags">
					Информация о поддерживаемых инструкциях недоступна
				</div>
			</section>
		</div>
	</main>
</template>

<style scoped>
.device-info-page {
	padding: 24px;
	max-width: 1200px;
	margin: 0 auto;
}

.page-title {
	font-size: 32px;
	font-weight: 700;
	color: #f5f5f7;
	margin: 0 0 32px 0;
}

.loading,
.error {
	padding: 24px;
	border-radius: 12px;
	background: rgba(28, 28, 30, 0.6);
	border: 1px solid rgba(255, 255, 255, 0.1);
	text-align: center;
	color: #f5f5f7;
	font-size: 16px;
}

.error {
	color: #ff453a;
	background: rgba(255, 69, 58, 0.1);
}

.device-info-container {
	display: flex;
	flex-direction: column;
	gap: 24px;
}

.info-section {
	padding: 24px;
	border-radius: 12px;
	background: rgba(28, 28, 30, 0.6);
	border: 1px solid rgba(255, 255, 255, 0.1);
	backdrop-filter: blur(20px);
}

.section-title {
	font-size: 20px;
	font-weight: 600;
	color: #f5f5f7;
	margin: 0 0 20px 0;
	padding-bottom: 12px;
	border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: 16px;
}

.info-item {
	display: flex;
	flex-direction: column;
	gap: 6px;
	padding: 12px;
	border-radius: 8px;
	background: rgba(0, 0, 0, 0.2);
}

.info-label {
	font-size: 13px;
	color: #8e8e93;
	font-weight: 500;
	text-transform: uppercase;
	letter-spacing: 0.5px;
}

.info-value {
	font-size: 16px;
	color: #f5f5f7;
	font-weight: 500;
}

.info-hint {
	font-size: 12px;
	color: rgba(142, 142, 147, 0.85);
	margin-top: 6px;
}

.flags-section {
	padding: 24px;
}

.flags-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
	gap: 16px;
}

.flag-item {
	padding: 16px;
	border-radius: 8px;
	background: rgba(0, 0, 0, 0.2);
	border: 1px solid rgba(255, 255, 255, 0.05);
	transition: all 0.2s ease;
}

.flag-item:hover {
	background: rgba(0, 0, 0, 0.3);
	border-color: rgba(0, 122, 255, 0.3);
	transform: translateY(-2px);
}

.flag-header {
	margin-bottom: 8px;
}

.flag-name {
	font-size: 14px;
	font-weight: 600;
	color: #007aff;
	letter-spacing: 1px;
}

.flag-description {
	font-size: 13px;
	color: #8e8e93;
	line-height: 1.5;
	margin: 0;
}

.no-flags {
	padding: 16px;
	text-align: center;
	color: #8e8e93;
	font-size: 14px;
}

@media (max-width: 768px) {
	.device-info-page {
		padding: 16px;
	}

	.page-title {
		font-size: 24px;
		margin-bottom: 24px;
	}

	.info-grid {
		grid-template-columns: 1fr;
	}

	.flags-grid {
		grid-template-columns: 1fr;
	}

	.info-section {
		padding: 16px;
	}
}
.alert-message {
	color: #fe5901;
	font-size: 14px;
}
</style>
