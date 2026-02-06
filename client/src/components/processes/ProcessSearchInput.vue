<script setup lang="ts">
import type { ProcessInfo } from "@/types/processes"
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue"
import "./ProcessAutocomplete.css"

type SearchType = "auto" | "pid" | "name" | "port"

type Props = {
	processes: ProcessInfo[]
	modelValue: string
	placeholder: string
	searchType?: SearchType
	maxSuggestions?: number
}

const props = withDefaults(defineProps<Props>(), {
	searchType: "auto",
	maxSuggestions: 8,
})

const emit = defineEmits<{
	(e: "update:modelValue", value: string): void
}>()

type Suggestion = {
	pid: string
	name: string
	username: string
	ports: number[]
	displayValue: string
	type: Exclude<SearchType, "auto">
}

const rootRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const highlightedIndex = ref(-1)
const isFocused = ref(false)

const value = computed(() => props.modelValue ?? "")

const suggestions = computed<Suggestion[]>(() => {
	if (!value.value || value.value.trim().length === 0) return []

	const query = value.value.trim().toLowerCase()
	const isNumeric = /^\d+$/.test(query)
	const actualSearchType: Exclude<SearchType, "auto"> =
		props.searchType === "auto"
			? isNumeric
				? "pid"
				: "name"
			: props.searchType

	const processMap = new Map<string, Suggestion>()

	for (const proc of props.processes) {
		const pid = String((proc as any)?.pid ?? "")
		const name = String((proc as any)?.name ?? "").toLowerCase()
		const username = String((proc as any)?.username ?? "")
		const ports = Array.isArray((proc as any)?.ports)
			? ((proc as any).ports as number[])
			: []

		let shouldInclude = false
		let displayValue = ""

		if (actualSearchType === "pid") {
			if (pid.includes(query)) {
				if (pid === query) continue
				shouldInclude = true
				displayValue = pid
			}
		} else if (actualSearchType === "name") {
			if (name.includes(query)) {
				shouldInclude = true
				displayValue = (proc as any)?.name || "Неизвестно"
			}
		} else if (actualSearchType === "port") {
			const matchingPorts = ports.filter(port => String(port).includes(query))
			if (matchingPorts.length > 0) {
				shouldInclude = true
				displayValue = matchingPorts.join(", ")
			}
		}

		if (shouldInclude && !processMap.has(pid)) {
			processMap.set(pid, {
				pid,
				name: (proc as any)?.name || "Неизвестно",
				username,
				ports,
				displayValue,
				type: actualSearchType,
			})
		}
	}

	return Array.from(processMap.values())
		.slice(0, props.maxSuggestions)
		.sort((a, b) => {
			if (a.type !== b.type) {
				const order: Record<string, number> = { pid: 0, name: 1, port: 2 }
				return (order[a.type] ?? 0) - (order[b.type] ?? 0)
			}
			return a.name.localeCompare(b.name)
		})
})

watch(
	() => [suggestions.value.length, value.value] as const,
	() => {
		highlightedIndex.value = -1
		isOpen.value =
			isFocused.value &&
			suggestions.value.length > 0 &&
			value.value.trim().length > 0
	}
)

const close = () => {
	isOpen.value = false
	highlightedIndex.value = -1
}

const selectSuggestion = (s: Suggestion) => {
	if (s.type === "pid") {
		emit("update:modelValue", s.pid)
	} else if (s.type === "name") {
		emit("update:modelValue", s.name)
	} else if (s.type === "port") {
		const query = value.value.trim()
		const matchingPort = s.ports.find(port => String(port).includes(query))
		emit(
			"update:modelValue",
			matchingPort != null
				? String(matchingPort)
				: s.ports[0]
				? String(s.ports[0])
				: ""
		)
	}
	close()
}

const onKeyDown = (e: KeyboardEvent) => {
	if (!isOpen.value) return

	if (e.key === "ArrowDown") {
		e.preventDefault()
		highlightedIndex.value = Math.min(
			highlightedIndex.value + 1,
			suggestions.value.length - 1
		)
	} else if (e.key === "ArrowUp") {
		e.preventDefault()
		highlightedIndex.value = Math.max(highlightedIndex.value - 1, -1)
	} else if (e.key === "Enter" && highlightedIndex.value >= 0) {
		e.preventDefault()
		const s = suggestions.value[highlightedIndex.value]
		if (s) selectSuggestion(s)
	} else if (e.key === "Escape") {
		close()
	}
}

const onDocumentMouseDown = (e: MouseEvent) => {
	const root = rootRef.value
	if (!root) return
	if (e.target instanceof Node && !root.contains(e.target)) {
		close()
	}
}

onMounted(() => {
	document.addEventListener("mousedown", onDocumentMouseDown)
	document.addEventListener("keydown", onKeyDown)
})

onBeforeUnmount(() => {
	document.removeEventListener("mousedown", onDocumentMouseDown)
	document.removeEventListener("keydown", onKeyDown)
})

const onBlur = () => {
	isFocused.value = false

	setTimeout(() => close(), 0)
}
</script>

<template>
	<div class="processes-search-wrapper" ref="rootRef">
		<div class="processes-search">
			<input
				type="text"
				class="processes-search-input"
				:placeholder="placeholder"
				:value="modelValue"
				@input="
					emit('update:modelValue', ($event.target as HTMLInputElement).value)
				"
				@focus="
					() => {
						isFocused = true
						isOpen = suggestions.length > 0 && modelValue.trim().length > 0
					}
				"
				@blur="onBlur"
			/>

			<button
				v-if="modelValue"
				type="button"
				class="processes-search-clear"
				title="Очистить поиск"
				aria-label="Очистить поиск"
				@click="emit('update:modelValue', '')"
			>
				<svg width="16" height="16" viewBox="0 0 16 16" fill="none">
					<path
						d="M12 4L4 12M4 4L12 12"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
					/>
				</svg>
			</button>
		</div>

		<div v-if="isOpen" class="process-autocomplete" role="listbox">
			<div
				v-for="(item, idx) in suggestions"
				:key="`${item.pid}-${idx}-${item.type}`"
				class="process-autocomplete-item"
				:class="{
					'process-autocomplete-item--highlighted': idx === highlightedIndex,
				}"
				@mousedown.prevent="selectSuggestion(item)"
				@mouseenter="highlightedIndex = idx"
			>
				<div class="process-autocomplete-item__main">
					<span class="process-autocomplete-item__name">{{ item.name }}</span>
					<span
						v-if="item.username"
						class="process-autocomplete-item__username"
						>{{ item.username }}</span
					>
					<span
						v-if="item.type === 'port' && item.ports.length > 0"
						class="process-autocomplete-item__ports"
					>
						Порты: {{ item.ports.join(", ") }}
					</span>
				</div>
				<div class="process-autocomplete-item__meta">
					<div class="process-autocomplete-item__pid">PID: {{ item.pid }}</div>
				</div>
			</div>
		</div>
	</div>
</template>
