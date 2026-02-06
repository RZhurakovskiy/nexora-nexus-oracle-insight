import type { AppServices } from "@/di/services"
import type { PiniaPluginContext } from "pinia"

/**
 * Pinia-плагин для DI.
 *
 * Я прокидываю сервисы в каждый store через "store.$services",
 * actions в store не создают зависимости сами.
 */
export function createPiniaServicesPlugin(services: AppServices) {
	return ({ store }: PiniaPluginContext) => {
		store.$services = services
	}
}
