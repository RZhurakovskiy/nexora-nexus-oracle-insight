/**
 * Конфигурация адресов бэкенда.
 *
 * Дефолт http
 */
export function getBackendBaseUrl(): string {
	const raw = (import.meta as any).env?.VITE_BACKEND_URL as string | undefined
	const fallback = "http://localhost:8080"
	const base = (raw ?? fallback).trim()
	return base.replace(/\/+$/, "")
}

/**
 * Строит WS base url из HTTP base url.
 *
 * - http
 * - https
 */
export function getBackendWsBaseUrl(): string {
	const httpBase = getBackendBaseUrl()
	if (httpBase.startsWith("https://"))
		return httpBase.replace(/^https:\/\//, "wss://")
	if (httpBase.startsWith("http://"))
		return httpBase.replace(/^http:\/\//, "ws://")
	return httpBase
}
