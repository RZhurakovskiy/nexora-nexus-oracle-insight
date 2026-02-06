<script setup lang="ts">
import { useServices } from "@/composables/useServices"
import { computed } from "vue"
import "./appLoader.css"

const { loading } = useServices()

const isLoading = computed(() => loading.isLoading.value)
</script>

<template>
	<transition name="loader-fade">
		<div
			v-if="isLoading"
			class="app-loader"
			aria-live="polite"
			aria-busy="true"
		>
			<div class="app-loader__backdrop" />
			<div class="app-loader__content">
				<div class="app-loader__spinner" />
				<div class="app-loader__text">Загрузка…</div>
			</div>
		</div>
	</transition>
</template>
