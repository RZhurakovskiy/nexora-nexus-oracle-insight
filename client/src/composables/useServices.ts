import { ServicesKey } from "@/di/keys"
import { inject } from "vue"

/**
 * Доступ к DI-контейнеру сервисов из компонентов.
 *
 * Для store используется "store.$services".
 */
export function useServices() {
	const services = inject(ServicesKey, null)
	if (!services) {
		throw new Error(
			"AppServices не предоставлены. Проверь `main.ts` (app.provide)."
		)
	}
	return services
}
