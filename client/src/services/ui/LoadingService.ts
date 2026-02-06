import { computed, readonly, ref } from "vue"

/**
 * Глобальный трекер загрузки для HTTP запросов.
 *
 * Я держу это сервисом (а не store), потому что:
 * - состояние техническое и не относится к доменной модели
 * - его удобно обновлять из axios interceptors
 * - компонент-лоадер может подписаться на одно место
 */
export class LoadingService {
	private readonly _pending = ref(0)

	/**
	 * Кол-во активных запросов для отладки.
	 */
	readonly pending = readonly(this._pending)

	/**
	 * Флаг для UI.
	 */
	readonly isLoading = computed(() => this._pending.value > 0)

	begin(): void {
		this._pending.value += 1
	}

	end(): void {
		this._pending.value = Math.max(0, this._pending.value - 1)
	}
}
