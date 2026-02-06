<script setup lang="ts">
import ConfirmModal from "@/components/confirm-modal/ConfirmModal.vue"
import BaseModal from "@/components/modal/BaseModal.vue"
import Button from "@/components/ui/Button.vue"
import "@/pages/port-manager/PortManager.css"
import { useNetworkStore } from "@/stores/network"
import { usePortsStore } from "@/stores/ports"
import { storeToRefs } from "pinia"
import { computed, onBeforeUnmount, onMounted, ref } from "vue"

type LaunchDraft = {
	command: string
	args: string
	cwd: string
}

const portsStore = usePortsStore()
const networkStore = useNetworkStore()

const { loading, error, starting, lastStartResult, killingPid } =
	storeToRefs(portsStore)
const { filteredConnections, topProcesses, interfaces } =
	storeToRefs(networkStore)

const showOnlyListening = computed({
	get: () => networkStore.showOnlyListen,
	set: v => (networkStore.showOnlyListen = v),
})
const isModalOpen = ref(false)
const newProcess = ref({
	command: "",
	args: "",
	cwd: "",
})

const STORAGE_KEY = "nexora.portManager.recentLaunches.v1"
const MAX_RECENTS = 8

type RecentLaunch = {
	id: string
	command: string
	args: string
	cwd: string
	createdAt: number
}

const presets = [
	{
		id: "vite",
		title: "Vite dev",
		hint: "vite --port=5173",
		value: { command: "vite", args: "--port=5173", cwd: "" },
	},
	{
		id: "npm-dev",
		title: "npm run dev",
		hint: "npm run dev -- --port=3000",
		value: { command: "npm", args: "run dev -- --port=3000", cwd: "" },
	},
	{
		id: "node",
		title: "node",
		hint: "node server.js --port=8080",
		value: { command: "node", args: "server.js --port=8080", cwd: "" },
	},
	{
		id: "python-http",
		title: "Python HTTP",
		hint: "python3 -m http.server 8000",
		value: { command: "python3", args: "-m http.server 8000", cwd: "" },
	},
] as const

const recentLaunches = ref<RecentLaunch[]>([])

/**
 * Прочитать историю запусков из localStorage, потом переделаю на хранение в БД.
 */
const loadRecents = (): void => {
	try {
		const raw = localStorage.getItem(STORAGE_KEY)
		if (!raw) {
			recentLaunches.value = []
			return
		}
		const parsed = JSON.parse(raw) as unknown
		if (!Array.isArray(parsed)) {
			recentLaunches.value = []
			return
		}
		recentLaunches.value = parsed
			.filter((x: any) => x && typeof x === "object")
			.slice(0, MAX_RECENTS)
			.map((x: any) => ({
				id: String(x.id ?? ""),
				command: String(x.command ?? ""),
				args: String(x.args ?? ""),
				cwd: String(x.cwd ?? ""),
				createdAt: Number(x.createdAt ?? Date.now()),
			}))
	} catch {
		recentLaunches.value = []
	}
}

/**
 * Записать историю запусков в localStorage.
 */
const saveRecents = (): void => {
	try {
		localStorage.setItem(
			STORAGE_KEY,
			JSON.stringify(recentLaunches.value.slice(0, MAX_RECENTS))
		)
	} catch {}
}

const addRecent = (draft: LaunchDraft): void => {
	const normalized = {
		command: draft.command.trim(),
		args: draft.args.trim(),
		cwd: draft.cwd.trim(),
	}
	if (!normalized.command) return

	const id = `${normalized.command}__${normalized.args}__${normalized.cwd}`
	const next: RecentLaunch = {
		id,
		command: normalized.command,
		args: normalized.args,
		cwd: normalized.cwd,
		createdAt: Date.now(),
	}

	recentLaunches.value = [
		next,
		...recentLaunches.value.filter(recent => recent.id !== id),
	].slice(0, MAX_RECENTS)
	saveRecents()
}

const applyDraft = (draft: LaunchDraft): void => {
	newProcess.value.command = draft.command
	newProcess.value.args = draft.args
	newProcess.value.cwd = draft.cwd
}

const listeningPortsSet = computed(() => {
	const set = new Set<number>()
	for (const proc of networkStore.connections) {
		if (proc.status === "LISTEN" && typeof proc.port === "number") {
			set.add(proc.port)
		}
	}
	return set
})

/**
 * Извлекаю порт из args.
 *
 * Поддерживаю наиболее частые паттерны:
 * - --port=5173
 * - --port 5173
 * - -p 5173
 */
const extractPortFromArgs = (args: string): number | null => {
	const raw = args.trim()
	if (!raw) return null

	const m1 = raw.match(/--port\s*=\s*(\d{2,5})/i)
	if (m1?.[1]) return Number(m1[1])

	const m2 = raw.match(/--port\s+(\d{2,5})/i)
	if (m2?.[1]) return Number(m2[1])

	const m3 = raw.match(/(?:^|\s)-p\s+(\d{2,5})(?:\s|$)/i)
	if (m3?.[1]) return Number(m3[1])

	const parts = raw.split(/\s+/)
	const last = parts.length > 0 ? parts[parts.length - 1] : undefined
	if (last && /^\d{2,5}$/.test(last)) return Number(last)

	return null
}

const desiredPort = computed(() => extractPortFromArgs(newProcess.value.args))
const isDesiredPortBusy = computed(() => {
	const p = desiredPort.value
	if (!p) return false
	return listeningPortsSet.value.has(p)
})

const findNextFreePort = (start: number, attempts = 200): number | null => {
	let p = Math.max(1024, Math.min(start, 65535))
	for (let i = 0; i < attempts && p <= 65535; i += 1, p += 1) {
		if (!listeningPortsSet.value.has(p)) return p
	}
	return null
}

const suggestedPort = computed(() => {
	const p = desiredPort.value
	if (!p) return null
	if (!isDesiredPortBusy.value) return null
	return findNextFreePort(p + 1)
})

/**
 * Подставить порт в args, стараясь сохранить исходный формат.
 */
const replacePortInArgs = (args: string, nextPort: number): string => {
	if (!args.trim()) return `--port=${nextPort}`

	if (/--port\s*=\s*\d{2,5}/i.test(args)) {
		return args.replace(/(--port\s*=\s*)\d{2,5}/i, `$1${nextPort}`)
	}
	if (/--port\s+\d{2,5}/i.test(args)) {
		return args.replace(/(--port\s+)\d{2,5}/i, `$1${nextPort}`)
	}
	if (/(?:^|\s)-p\s+\d{2,5}(?:\s|$)/i.test(args)) {
		return args.replace(/((?:^|\s)-p\s+)\d{2,5}((?:\s|$))/i, `$1${nextPort}$2`)
	}

	const tokens = args.split(/\s+/)
	const last = tokens.length > 0 ? tokens[tokens.length - 1] : undefined
	if (last && /^\d{2,5}$/.test(last)) {
		tokens[tokens.length - 1] = String(nextPort)
		return tokens.join(" ")
	}

	return `${args.trim()} --port=${nextPort}`
}

const applySuggestedPort = (): void => {
	const next = suggestedPort.value
	if (!next) return
	newProcess.value.args = replacePortInArgs(newProcess.value.args, next)
}

const openModal = () => {
	isModalOpen.value = true
}

const closeModal = () => {
	isModalOpen.value = false
}

const resetForm = () => {
	newProcess.value = { command: "", args: "", cwd: "" }
}

const fetchPorts = async () => {
	await networkStore.refreshAll()
}

const isConfirmOpen = ref(false)
const confirmPid = ref<number | null>(null)
const confirmTitle = ref("Подтверждение завершения процесса")
const confirmMessage = ref("")
const confirmText = ref("Завершить процесс")
const cancelText = ref("Отмена")

const requestKillProcess = (pid: number, processName: string) => {
	confirmPid.value = pid
	confirmMessage.value = `Вы уверены, что хотите завершить процесс?\n\nPID: ${pid}\nПроцесс: ${
		processName || "—"
	}\n\nЭто действие нельзя отменить.`
	isConfirmOpen.value = true
}

const confirmKill = async () => {
	const pid = confirmPid.value
	isConfirmOpen.value = false
	confirmPid.value = null
	if (!pid) return
	await portsStore.killProcess(pid)
}

const cancelKill = () => {
	isConfirmOpen.value = false
	confirmPid.value = null
}

const handleSubmit = async () => {
	/**
	 * Если порт занят LISTEN - показывается предупреждение.
	 * Пользователь может запускать процессы без порта или с динамическим портом.
	 */
	const timestamp = new Date().toISOString()
	const payload = {
		command: newProcess.value.command,
		args: newProcess.value.args,
		timestamp,
		...(newProcess.value.cwd && { cwd: newProcess.value.cwd }),
	}

	const res = await portsStore.startProcess(payload)
	if (res.ok) {
		addRecent({
			command: payload.command,
			args: payload.args,
			cwd: String(payload.cwd ?? ""),
		})
		closeModal()
		resetForm()
	}
}

onMounted(async () => {
	loadRecents()
	await portsStore.init()
	await networkStore.init()
})

onBeforeUnmount(() => {
	networkStore.setAutoRefresh(false)
})

const renderStatusTitle = (status: string): string => {
	return status === "LISTEN"
		? "Сервер: порт занят"
		: "Клиент: порт можно использовать"
}
</script>

<template>
	<div class="portManager">
		<div class="portManagerHeader">
			<div class="portManagerTitleWrapper">
				<svg
					class="portManagerTitleIcon"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					xmlns="http://www.w3.org/2000/svg"
				>
					<rect
						x="4"
						y="5"
						width="12"
						height="10"
						rx="2"
						stroke="currentColor"
						stroke-width="1.5"
						fill="none"
					/>
					<circle cx="6" cy="8" r="1" fill="currentColor" />
					<circle cx="6" cy="12" r="1" fill="currentColor" />
					<circle cx="14" cy="8" r="1" fill="currentColor" />
					<circle cx="14" cy="12" r="1" fill="currentColor" />
					<line
						x1="16"
						y1="10"
						x2="19"
						y2="10"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
					/>
				</svg>
				<h2>Сетевые соединения</h2>
			</div>

			<div class="portManagerActions">
				<button
					class="portManagerButton"
					type="button"
					title="Обновить сетевую статистику"
					@click="fetchPorts"
				>
					<svg
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<polyline points="23 4 23 10 17 10" />
						<polyline points="1 20 1 14 7 14" />
						<path
							d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"
						/>
					</svg>
					<span>Обновить</span>
				</button>
				<button
					class="portManagerButton portManagerButtonPrimary"
					type="button"
					@click="openModal"
				>
					<svg
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<line x1="12" y1="5" x2="12" y2="19" />
						<line x1="5" y1="12" x2="19" y2="12" />
					</svg>
					<span>Запустить процесс</span>
				</button>
			</div>
		</div>

		<div class="processesCard" style="margin-bottom: 16px">
			<div
				style="
					padding: 14px 16px;
					display: flex;
					gap: 12px;
					align-items: center;
					justify-content: space-between;
					flex-wrap: wrap;
				"
			>
				<div style="font-weight: 800; color: rgba(245, 245, 247, 0.92)">
					Сетевые интерфейсы (I/O)
				</div>
				<label class="portFilterLabel" style="margin: 0">
					<input
						:checked="networkStore.autoRefresh"
						@change="
							networkStore.setAutoRefresh(
								($event.target as HTMLInputElement).checked
							)
						"
						type="checkbox"
						class="portFilterCheckbox"
					/>
					<span class="portFilterCheckmark" />
					<span class="portFilterText">Автообновление (5с)</span>
				</label>
			</div>
			<div class="processesTableWrapper">
				<table class="processesTable">
					<thead>
						<tr>
							<th>Интерфейс</th>
							<th>RX</th>
							<th>TX</th>
							<th>Ошибки</th>
							<th>Drop</th>
						</tr>
					</thead>
					<tbody class="processesTableBody">
						<tr v-for="int in interfaces" :key="int.name">
							<td>
								<span class="portAddress">{{ int.name }}</span>
							</td>
							<td>{{ (int.bytesRecv / 1024 / 1024).toFixed(1) }} MB</td>
							<td>{{ (int.bytesSent / 1024 / 1024).toFixed(1) }} MB</td>
							<td>{{ int.errIn }} / {{ int.errOut }}</td>
							<td>{{ int.dropIn }} / {{ int.dropOut }}</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<div class="portFilters">
			<label class="portFilterLabel">
				<input
					v-model="showOnlyListening"
					type="checkbox"
					class="portFilterCheckbox"
				/>
				<span class="portFilterCheckmark" />
				<span class="portFilterText">Только серверы (LISTEN)</span>
				<span v-if="filteredConnections.length > 0" class="portFilterCount"
					>({{ filteredConnections.length }})</span
				>
			</label>
		</div>

		<div v-if="loading" class="processesCard">
			<div class="spinnerOverlay">
				<div class="spinnerSystem" />
				<div class="spinnerText">Загрузка сетевой статистики…</div>
			</div>
		</div>

		<div v-else-if="filteredConnections.length === 0" class="portManagerEmpty">
			<svg
				class="portManagerEmptyIcon"
				width="64"
				height="64"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.5"
			>
				<rect x="4" y="5" width="12" height="10" rx="2" opacity="0.3" />
				<line x1="4" y1="10" x2="16" y2="10" opacity="0.3" />
				<circle cx="6" cy="8" r="1" />
				<circle cx="6" cy="12" r="1" />
				<circle cx="14" cy="8" r="1" />
				<circle cx="14" cy="12" r="1" />
			</svg>
			<p class="portManagerEmptyText">
				{{
					showOnlyListening ? "Нет активных серверов" : "Нет сетевых соединений"
				}}
			</p>
			<p class="portManagerEmptyHint">
				{{
					showOnlyListening
						? "Попробуйте снять фильтр, чтобы увидеть все соединения"
						: "Запустите процесс, чтобы увидеть сетевые соединения"
				}}
			</p>
			<p v-if="error" class="portManagerEmptyHint">
				<b>Ошибка:</b> {{ error }}
			</p>
		</div>

		<div v-else class="processesCard">
			<div
				style="
					padding: 12px 16px;
					display: flex;
					gap: 10px;
					align-items: center;
					flex-wrap: wrap;
					border-bottom: 1px solid rgba(255, 255, 255, 0.05);
				"
			>
				<input
					v-model="networkStore.filter.query"
					class="addProcessFieldInput"
					style="max-width: 340px"
					type="text"
					placeholder="Фильтр: PID / процесс / порт / адрес / статус"
				/>
				<input
					v-model="networkStore.filter.status"
					class="addProcessFieldInput"
					style="max-width: 180px"
					type="text"
					placeholder="Статус (LISTEN/ESTABLISHED)"
				/>
				<Button
					variant="ghost"
					size="sm"
					@click="networkStore.filter = { query: '', status: '' }"
					>Сброс</Button
				>
			</div>
			<div class="processesTableWrapper">
				<table class="processesTable">
					<thead>
						<tr>
							<th>Локальный адрес</th>
							<th>Удалённый адрес</th>
							<th>Протокол</th>
							<th>Статус</th>
							<th>PID</th>
							<th>Процесс</th>
							<th style="text-align: right">Действия</th>
						</tr>
					</thead>

					<tbody class="processesTableBody">
						<tr
							v-for="item in filteredConnections"
							:key="`${item.localAddr}-${item.remoteAddr}-${item.pid}-${item.status}`"
							:class="
								item.status === 'LISTEN' ? 'portRowListen' : 'portRowOther'
							"
						>
							<td>
								<span class="portAddress">{{ item.localAddr }}</span>
							</td>
							<td>
								<span class="portAddress">{{
									item.remoteAddr === "-" ? "—" : item.remoteAddr
								}}</span>
							</td>
							<td>
								<span class="portProtocol">{{ item.protocol }}</span>
							</td>
							<td>
								<div class="portStatusWrapper">
									<span
										class="portStatus"
										:class="
											item.status === 'LISTEN'
												? 'portStatusListen'
												: 'portStatusOther'
										"
										:title="renderStatusTitle(item.status)"
									/>
									<span class="portStatusText">{{ item.status }}</span>
								</div>
							</td>
							<td>
								<span class="portPid">{{ item.pid }}</span>
							</td>
							<td>
								<span class="portProcess">{{ item.process }}</span>
							</td>
							<td>
								<div class="portActions">
									<button
										type="button"
										class="portActionBtn portActionBtn--kill"
										:disabled="killingPid != null"
										title="Завершить процесс"
										@click="requestKillProcess(item.pid, item.process)"
									>
										{{ killingPid === item.pid ? "Завершение…" : "Завершить" }}
									</button>
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>

			<div class="portLegend">
				<div class="legendItem">
					<span class="portStatus portStatusListen" />
					<span>Сервер (порт занят — нельзя запустить другой)</span>
				</div>
				<div class="legendItem">
					<span class="portStatus portStatusOther" />
					<span>Клиент (временное соединение — порт можно использовать)</span>
				</div>
			</div>
		</div>

		<div class="processesCard" style="margin-top: 16px">
			<div
				style="
					padding: 14px 16px;
					font-weight: 800;
					color: rgba(245, 245, 247, 0.92);
					border-bottom: 1px solid rgba(255, 255, 255, 0.05);
				"
			>
				Топ процессов по сетевой активности (по числу соединений)
			</div>
			<div class="processesTableWrapper">
				<table class="processesTable">
					<thead>
						<tr>
							<th>PID</th>
							<th>Процесс</th>
							<th>Пользователь</th>
							<th>ESTABLISHED</th>
							<th>LISTEN</th>
							<th>Другие</th>
							<th>Всего</th>
						</tr>
					</thead>
					<tbody class="processesTableBody">
						<tr v-for="proc in topProcesses" :key="proc.pid">
							<td>
								<span class="portPid">{{ proc.pid }}</span>
							</td>
							<td>
								<span class="portProcess">{{ proc.process }}</span>
							</td>
							<td>
								<span class="portAddress">{{ proc.username || "—" }}</span>
							</td>
							<td>{{ proc.established }}</td>
							<td>{{ proc.listening }}</td>
							<td>{{ proc.otherStates }}</td>
							<td>
								<b>{{ proc.connections }}</b>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<BaseModal
			:is-open="isModalOpen"
			title="Запустить новый процесс"
			width="760px"
			@close="closeModal"
		>
			<form
				id="launchProcessForm"
				class="processModalBody"
				@submit.prevent="handleSubmit"
			>
				<div class="launchSection">
					<div class="launchSectionHeader">
						<div class="launchSectionTitle">Пресеты</div>
						<div class="launchSectionHint">Быстро заполняют поля команды</div>
					</div>
					<div class="launchPresets">
						<button
							v-for="preset in presets"
							:key="preset.id"
							type="button"
							class="launchPreset"
							:title="preset.hint"
							@click="applyDraft(preset.value)"
						>
							<div class="launchPresetTitle">{{ preset.title }}</div>
							<div class="launchPresetHint">{{ preset.hint }}</div>
						</button>
					</div>
				</div>

				<div v-if="recentLaunches.length > 0" class="launchSection">
					<div class="launchSectionHeader">
						<div class="launchSectionTitle">Недавние</div>
						<div class="launchSectionHint">
							Последние запуски на этой машине
						</div>
					</div>
					<div class="launchRecents">
						<button
							v-for="recent in recentLaunches"
							:key="recent.id"
							type="button"
							class="launchRecent"
							@click="
								applyDraft({
									command: recent.command,
									args: recent.args,
									cwd: recent.cwd,
								})
							"
						>
							<div class="launchRecentMain">
								<span class="launchRecentCmd">{{ recent.command }}</span>
								<span v-if="recent.args" class="launchRecentArgs">
									{{ recent.args }}</span
								>
							</div>
							<div
								v-if="recent.cwd"
								class="launchRecentCwd"
								:title="recent.cwd"
							>
								cwd: {{ recent.cwd }}
							</div>
						</button>
					</div>
				</div>

				<div class="addProcessField">
					<label class="addProcessFieldLabel">
						<svg
							class="addProcessFieldIcon"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
						>
							<path d="M4 16v4a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-4" />
							<polyline points="8 12 12 16 16 12" />
							<line x1="12" y1="8" x2="12" y2="16" />
						</svg>
						Команда *
					</label>
					<input
						v-model="newProcess.command"
						type="text"
						name="command"
						placeholder="vite, node, python3, npm и т.д."
						class="addProcessFieldInput"
						required
					/>
				</div>

				<div class="addProcessField">
					<label class="addProcessFieldLabel">
						<svg
							class="addProcessFieldIcon"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
						>
							<path
								d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
							/>
							<polyline points="14 2 14 8 20 8" />
							<line x1="16" y1="13" x2="8" y2="13" />
							<line x1="16" y1="17" x2="8" y2="17" />
							<polyline points="10 9 9 9 8 9" />
						</svg>
						Аргументы
					</label>
					<input
						v-model="newProcess.args"
						type="text"
						name="args"
						placeholder="--port=8081 или run dev -- --port=5173"
						class="addProcessFieldInput"
					/>
				</div>

				<div
					v-if="desiredPort"
					class="portCheck"
					:class="{ 'portCheck--busy': isDesiredPortBusy }"
				>
					<div class="portCheckMain">
						<div class="portCheckTitle">
							Порт из аргументов: <b>{{ desiredPort }}</b>
						</div>
						<div v-if="isDesiredPortBusy" class="portCheckText">
							Порт занят сервером (LISTEN). Лучше выбрать другой.
						</div>
						<div v-else class="portCheckText portCheckText--ok">
							Порт свободен (по текущему списку LISTEN).
						</div>
					</div>
					<button
						v-if="suggestedPort"
						type="button"
						class="portCheckAction"
						@click="applySuggestedPort"
						:title="`Подставить свободный порт ${suggestedPort}`"
					>
						Подставить {{ suggestedPort }}
					</button>
				</div>

				<div class="addProcessField">
					<label class="addProcessFieldLabel">
						<svg
							class="addProcessFieldIcon"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
						>
							<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
							<path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
						</svg>
						Рабочая директория (необязательно)
					</label>
					<input
						v-model="newProcess.cwd"
						type="text"
						name="cwd"
						placeholder="Пример: /home/user/my-app или C:\\my-project"
						class="addProcessFieldInput"
					/>
				</div>

				<div class="addProcessHint">
					Примеры:
					<br />
					<code>vite</code> + <code>--port=8081</code> → запустит dev-сервер на
					порту 8081
					<br />
					<code>npm</code> + <code>run dev -- --port=3000</code> → для
					npm-проектов
					<br />
					Если порт занят <strong>сервером (LISTEN)</strong> — выберите другой.
				</div>

				<div
					v-if="portsStore.error"
					class="addProcessHint"
					style="color: #ff453a"
				>
					<b>Ошибка</b>: {{ portsStore.error }}
				</div>
				<div
					v-else-if="lastStartResult"
					class="addProcessHint"
					style="color: rgba(96, 255, 165, 0.95)"
				>
					Процесс запущен. PID: <b>{{ lastStartResult.pid }}</b>
				</div>
			</form>
			<template #footer>
				<button
					type="button"
					class="processModalButton processModalButtonCancel"
					@click="closeModal"
				>
					Отмена
				</button>
				<button
					type="submit"
					form="launchProcessForm"
					class="processModalButton processModalButtonConfirm"
					:disabled="starting"
				>
					{{ starting ? "Запуск…" : "Запустить" }}
				</button>
			</template>
		</BaseModal>

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
