<script setup lang="ts">
import BaseModal from "@/components/modal/BaseModal.vue"
import "./ConfirmModal.css"

type Props = {
	isOpen: boolean
	title: string
	message: string
	confirmText: string
	cancelText: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
	(e: "confirm"): void
	(e: "cancel"): void
}>()

const lines = () => props.message.split("\n")
</script>

<template>
	<BaseModal
		:is-open="isOpen"
		:title="title"
		width="520px"
		@close="emit('cancel')"
	>
		<div class="confirm-modal__header">
			<div class="confirm-modal__icon" aria-hidden="true">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none">
					<path
						d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2ZM13 17H11V15H13V17ZM13 13H11V7H13V13Z"
						fill="currentColor"
					/>
				</svg>
			</div>
		</div>

		<div class="confirm-modal__body">
			<p class="confirm-modal__message">
				<template v-for="(line, idx) in lines()" :key="idx">
					<span>{{ line }}</span>
					<br v-if="idx < lines().length - 1" />
				</template>
			</p>
		</div>

		<template #footer>
			<button
				type="button"
				class="confirm-modal__button confirm-modal__button--cancel"
				@click="emit('cancel')"
			>
				{{ cancelText }}
			</button>
			<button
				type="button"
				class="confirm-modal__button confirm-modal__button--confirm"
				@click="emit('confirm')"
			>
				{{ confirmText }}
			</button>
		</template>
	</BaseModal>
</template>
