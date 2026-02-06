<script setup lang="ts">
import Button from "@/components/ui/Button.vue"
import { useAuthStore } from "@/stores/auth"
import { ref } from "vue"
import { useRouter } from "vue-router"
import "./LoginPage.css"

const authStore = useAuthStore()
const router = useRouter()

const login = ref("")
const password = ref("")
const error = ref("")
const showPassword = ref(false)

const handleSubmit = async (e: Event): Promise<void> => {
	e.preventDefault()
	error.value = ""

	if (!login.value.trim() || !password.value.trim()) {
		error.value = "Заполните все поля"
		return
	}

	try {
		await authStore.login(login.value.trim(), password.value)
		router.push("/")
	} catch (err) {
		error.value =
			err instanceof Error ? err.message : "Произошла ошибка при входе"
	}
}
</script>

<template>
	<div class="login-page">
		<div class="login-container">
			<div class="login-card">
				<div class="login-header">
					<div class="login-logo">
						<svg
							width="64"
							height="64"
							viewBox="0 0 24 24"
							fill="none"
							xmlns="http://www.w3.org/2000/svg"
						>
							<rect
								x="4"
								y="11"
								width="16"
								height="9"
								rx="2"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
								stroke-linejoin="round"
								fill="none"
							/>
							<path
								d="M7 11V7C7 4.79086 8.79086 3 11 3H13C15.2091 3 17 4.79086 17 7V11"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
								stroke-linejoin="round"
								fill="none"
							/>
							<circle
								cx="12"
								cy="15.5"
								r="1.5"
								fill="currentColor"
							/>
						</svg>
					</div>
					<h1 class="login-title">Вход в систему</h1>
					<p class="login-subtitle">Введите учетные данные для доступа</p>
				</div>

				<form class="login-form" @submit="handleSubmit">
					<div class="login-field">
						<label class="login-label">
							<svg
								class="login-label-icon"
								width="16"
								height="16"
								viewBox="0 0 16 16"
								fill="none"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path
									d="M8 8C10.2091 8 12 6.20914 12 4C12 1.79086 10.2091 0 8 0C5.79086 0 4 1.79086 4 4C4 6.20914 5.79086 8 8 8Z"
									fill="currentColor"
								/>
								<path
									d="M8 10C4.68629 10 2 12.6863 2 16H14C14 12.6863 11.3137 10 8 10Z"
									fill="currentColor"
								/>
							</svg>
							Логин
						</label>
						<input
							v-model="login"
							type="text"
							class="login-input"
							placeholder="Введите логин"
							autocomplete="username"
							:disabled="authStore.loading"
						/>
					</div>

					<div class="login-field">
						<label class="login-label">
							<svg
								class="login-label-icon"
								width="16"
								height="16"
								viewBox="0 0 16 16"
								fill="none"
								xmlns="http://www.w3.org/2000/svg"
							>
								<rect
									x="3"
									y="7"
									width="10"
									height="7"
									rx="1.5"
									stroke="currentColor"
									stroke-width="1.5"
									fill="none"
								/>
								<path
									d="M5 7V5C5 3.34315 6.34315 2 8 2C9.65685 2 11 3.34315 11 5V7"
									stroke="currentColor"
									stroke-width="1.5"
									stroke-linecap="round"
								/>
							</svg>
							Пароль
						</label>
						<div class="login-input-wrapper">
							<input
								v-model="password"
								:type="showPassword ? 'text' : 'password'"
								class="login-input"
								placeholder="Введите пароль"
								autocomplete="current-password"
								:disabled="authStore.loading"
							/>
							<button
								type="button"
								class="login-password-toggle"
								@click="showPassword = !showPassword"
								:disabled="authStore.loading"
							>
								<svg
									v-if="showPassword"
									width="16"
									height="16"
									viewBox="0 0 16 16"
									fill="none"
									xmlns="http://www.w3.org/2000/svg"
								>
									<path
										d="M1 8C1 8 3.5 4 8 4C12.5 4 15 8 15 8C15 8 12.5 12 8 12C3.5 12 1 8 1 8Z"
										stroke="currentColor"
										stroke-width="1.5"
										stroke-linecap="round"
										stroke-linejoin="round"
										fill="none"
									/>
									<circle
										cx="8"
										cy="8"
										r="2"
										stroke="currentColor"
										stroke-width="1.5"
										fill="none"
									/>
									<path
										d="M2 2L14 14"
										stroke="currentColor"
										stroke-width="1.5"
										stroke-linecap="round"
									/>
								</svg>
								<svg
									v-else
									width="16"
									height="16"
									viewBox="0 0 16 16"
									fill="none"
									xmlns="http://www.w3.org/2000/svg"
								>
									<path
										d="M1 8C1 8 3.5 4 8 4C12.5 4 15 8 15 8C15 8 12.5 12 8 12C3.5 12 1 8 1 8Z"
										stroke="currentColor"
										stroke-width="1.5"
										stroke-linecap="round"
										stroke-linejoin="round"
										fill="none"
									/>
									<circle
										cx="8"
										cy="8"
										r="2"
										stroke="currentColor"
										stroke-width="1.5"
										fill="none"
									/>
								</svg>
							</button>
						</div>
					</div>

					<div v-if="error" class="login-error">
						<svg
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
								fill="none"
							/>
							<path
								d="M8 5V8M8 11H8.01"
								stroke="currentColor"
								stroke-width="1.5"
								stroke-linecap="round"
							/>
						</svg>
						<span>{{ error }}</span>
					</div>

					<Button
						type="submit"
						variant="primary"
						size="lg"
						width="full"
						:loading="authStore.loading"
						:disabled="authStore.loading"
					>
						Войти
					</Button>
				</form>
			</div>
		</div>
	</div>
</template>

