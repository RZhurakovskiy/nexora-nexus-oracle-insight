import { defineStore } from "pinia"

export type NotificationType = "success" | "error" | "info"

export type AppNotification = {
	id: string
	type: NotificationType
	title: string
	message: string
	createdAt: number
	ttlMs: number
}

type NotificationsState = {
	items: AppNotification[]
}

const DEFAULT_TTL_MS = 4500

export const useNotificationsStore = defineStore("notifications", {
	state: (): NotificationsState => ({
		items: [],
	}),

	actions: {
		/**
		 * Показать уведомление.
		 *
		 * Я храню уведомления в Pinia, а авто-закрытие делаю таймером.
		 */
		push(input: {
			type: NotificationType
			title?: string
			message: string
			ttlMs?: number
		}): string {
			const id = `${Date.now()}_${Math.random().toString(16).slice(2)}`
			const ttlMs =
				typeof input.ttlMs === "number" ? input.ttlMs : DEFAULT_TTL_MS
			const title =
				typeof input.title === "string" && input.title.trim().length > 0
					? input.title
					: input.type === "success"
					? "Готово"
					: input.type === "error"
					? "Ошибка"
					: "Информация"

			const item: AppNotification = {
				id,
				type: input.type,
				title,
				message: input.message,
				createdAt: Date.now(),
				ttlMs,
			}

			this.items = [item, ...this.items].slice(0, 5)

			window.setTimeout(() => {
				try {
					this.remove(id)
				} catch {}
			}, ttlMs)

			return id
		},

		/**
		 * Удалить уведомление.
		 */
		remove(id: string): void {
			this.items = this.items.filter(notification => notification.id !== id)
		},

		success(message: string, title?: string): string {
			return this.push({ type: "success", title, message })
		},

		error(message: string, title?: string): string {
			return this.push({ type: "error", title, message, ttlMs: 6500 })
		},

		info(message: string, title?: string): string {
			return this.push({ type: "info", title, message })
		},

		clear(): void {
			this.items = []
		},
	},
})
