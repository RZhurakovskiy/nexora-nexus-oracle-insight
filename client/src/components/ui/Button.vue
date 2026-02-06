<script setup lang="ts">
import { computed } from "vue"

export interface ButtonProps {
	/**
	 * Цвета для кнопки
	 */
	variant?: "primary" | "danger" | "success" | "secondary" | "ghost"
	/**
	 * Размер кнопки
	 */
	size?: "sm" | "md" | "lg"
	/**
	 * Состояние загрузки. При "true" показывает индикатор загрузки и блокирует кнопку.
	 */
	loading?: boolean
	/**
	 * Отключена ли кнопка
	 */
	disabled?: boolean
	/**
	 * Тип кнопки
	 */
	type?: "button" | "submit" | "reset"
	/**
	 * Ширина кнопки. Если "full", занимает 100% ширины родителя.
	 */
	width?: string | "full"
}

const props = withDefaults(defineProps<ButtonProps>(), {
	variant: "primary",
	size: "md",
	loading: false,
	disabled: false,
	type: "button",
	width: undefined,
})

const emit = defineEmits<{
	click: [event: MouseEvent]
}>()

const buttonClasses = computed(() => {
	const classes = [
		"ui-button",
		`ui-button--${props.variant}`,
		`ui-button--${props.size}`,
	]

	if (props.loading) classes.push("ui-button--loading")
	if (props.disabled || props.loading) classes.push("ui-button--disabled")

	return classes.join(" ")
})

const buttonStyle = computed(() => {
	if (props.width === "full") {
		return { width: "100%" }
	}
	if (props.width && props.width !== "full") {
		return { width: props.width }
	}
	return {}
})

const handleClick = (event: MouseEvent) => {
	if (props.disabled || props.loading) return
	emit("click", event)
}
</script>

<template>
	<button
		:type="type"
		:class="buttonClasses"
		:style="buttonStyle"
		:disabled="disabled || loading"
		@click="handleClick"
	>
		<span v-if="loading" class="ui-button__spinner">
			<svg
				class="ui-button__spinner-icon"
				width="16"
				height="16"
				viewBox="0 0 16 16"
				fill="none"
				xmlns="http://www.w3.org/2000/svg"
			>
				<circle
					cx="8"
					cy="8"
					r="6"
					stroke="currentColor"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-dasharray="18.85"
					stroke-dashoffset="18.85"
					fill="none"
				>
					<animate
						attributeName="stroke-dashoffset"
						values="18.85;0;18.85"
						dur="1s"
						repeatCount="indefinite"
					/>
				</circle>
			</svg>
		</span>
		<span v-if="$slots.prefix && !loading" class="ui-button__prefix">
			<slot name="prefix" />
		</span>
		<span class="ui-button__content">
			<slot />
		</span>
		<span v-if="$slots.suffix && !loading" class="ui-button__suffix">
			<slot name="suffix" />
		</span>
	</button>
</template>

<style scoped>
.ui-button {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	font-family: inherit;
	font-weight: 500;
	border: none;
	border-radius: 6px;
	cursor: pointer;
	transition: all 0.2s ease;
	user-select: none;
	position: relative;
	white-space: nowrap;
}

.ui-button--sm {
	height: 32px;
	padding: 0 12px;
	font-size: 12px;
}

.ui-button--md {
	height: 40px;
	padding: 0 16px;
	font-size: 13px;
}

.ui-button--lg {
	height: 50px;
	padding: 0 20px;
	font-size: 14px;
}

.ui-button--primary {
	background: rgba(99, 102, 241, 0.2);
	border: 1px solid rgba(99, 102, 241, 0.4);
	color: rgba(255, 255, 255, 0.95);
}

.ui-button--primary:hover:not(.ui-button--disabled) {
	background: rgba(99, 102, 241, 0.3);
	border-color: rgba(99, 102, 241, 0.6);
	transform: translateY(-1px);
	box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.ui-button--danger {
	background: rgba(255, 59, 48, 0.15);
	border: 1px solid rgba(255, 59, 48, 0.3);
	color: rgba(255, 255, 255, 0.95);
}

.ui-button--danger:hover:not(.ui-button--disabled) {
	background: rgba(255, 59, 48, 0.25);
	border-color: rgba(255, 59, 48, 0.5);
	transform: translateY(-1px);
	box-shadow: 0 4px 12px rgba(255, 59, 48, 0.2);
}

.ui-button--success {
	background: rgba(82, 196, 26, 0.15);
	border: 1px solid rgba(82, 196, 26, 0.3);
	color: rgba(255, 255, 255, 0.95);
}

.ui-button--success:hover:not(.ui-button--disabled) {
	background: rgba(82, 196, 26, 0.25);
	border-color: rgba(82, 196, 26, 0.5);
	transform: translateY(-1px);
	box-shadow: 0 4px 12px rgba(82, 196, 26, 0.2);
}

.ui-button--secondary {
	background: rgba(142, 142, 147, 0.15);
	border: 1px solid rgba(142, 142, 147, 0.3);
	color: rgba(255, 255, 255, 0.85);
}

.ui-button--secondary:hover:not(.ui-button--disabled) {
	background: rgba(142, 142, 147, 0.25);
	border-color: rgba(142, 142, 147, 0.5);
	transform: translateY(-1px);
	box-shadow: 0 4px 12px rgba(142, 142, 147, 0.15);
}

.ui-button--ghost {
	background: transparent;
	border: 1px solid rgba(255, 255, 255, 0.1);
	color: rgba(255, 255, 255, 0.85);
}

.ui-button--ghost:hover:not(.ui-button--disabled) {
	background: rgba(255, 255, 255, 0.05);
	border-color: rgba(255, 255, 255, 0.2);
}

.ui-button:active:not(.ui-button--disabled) {
	transform: translateY(0);
}

.ui-button--disabled {
	opacity: 0.6;
	cursor: not-allowed;
	pointer-events: none;
}

.ui-button--loading {
	cursor: wait;
}

.ui-button__spinner,
.ui-button__prefix,
.ui-button__suffix {
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}

.ui-button__spinner-icon {
	width: 16px;
	height: 16px;
	animation: spin 1s linear infinite;
}

@keyframes spin {
	from {
		transform: rotate(0deg);
	}
	to {
		transform: rotate(360deg);
	}
}

.ui-button__content {
	display: flex;
	align-items: center;
	justify-content: center;
}

.ui-button svg {
	width: 16px;
	height: 16px;
	flex-shrink: 0;
}

.ui-button--sm svg {
	width: 14px;
	height: 14px;
}

.ui-button--lg svg {
	width: 18px;
	height: 18px;
}
</style>
