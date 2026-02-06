<script setup lang="ts">
import { onBeforeUnmount, watch } from "vue"
import "./BaseModal.css"

type Props = {
	isOpen: boolean
	title?: string
	width?: string
	/**
	 * Показывать подсказку внизу окна "Esc - закрыть".
	 */
	showEscHint?: boolean
}

const props = withDefaults(defineProps<Props>(), {
	title: "",
	width: "520px",
	showEscHint: true,
})

const emit = defineEmits<{
	(e: "close"): void
}>()

/**
 * - блокирует scroll страницы при открытом окне
 * - закрытие по клавише Escape
 * - overlay не закрывает модалку
 */
const onKeyDown = (e: KeyboardEvent) => {
	if (e.key === "Escape") emit("close")
}

watch(
	() => props.isOpen,
	isOpen => {
		document.body.style.overflow = isOpen ? "hidden" : ""
		if (isOpen) document.addEventListener("keydown", onKeyDown)
		else document.removeEventListener("keydown", onKeyDown)
	},
	{ immediate: true }
)

onBeforeUnmount(() => {
	document.body.style.overflow = ""
	document.removeEventListener("keydown", onKeyDown)
})
</script>

<template>
	<teleport to="body">
		<div v-if="isOpen" class="base-modal-backdrop">
			<div
				class="base-modal"
				role="dialog"
				aria-modal="true"
				:style="{ maxWidth: width }"
				@click.stop
			>
				<div class="base-modal__header">
					<div class="base-modal__title-wrap">
						<h3 v-if="title" class="base-modal__title">{{ title }}</h3>
						<slot name="header" />
					</div>
					<button
						class="base-modal__close"
						type="button"
						aria-label="Закрыть"
						@click="emit('close')"
					>
						<svg
							width="18"
							height="18"
							viewBox="0 0 20 20"
							fill="none"
							aria-hidden="true"
						>
							<path
								d="M15 5L5 15M5 5L15 15"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
							/>
						</svg>
					</button>
				</div>

				<div class="base-modal__body">
					<slot />
				</div>

				<div v-if="$slots.footer" class="base-modal__footer">
					<slot name="footer" />
				</div>

				<div v-if="showEscHint" class="base-modal__hint">
					<span class="base-modal__hint-key">Esc</span>
					<span class="base-modal__hint-text">— закрыть</span>
				</div>
			</div>
		</div>
	</teleport>
</template>
