import { getBackendBaseUrl } from "@/config/backend"
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios"
import axios from "axios"

export class ApiError extends Error {
	readonly status: number
	readonly payload: unknown

	constructor(message: string, status: number, payload: unknown) {
		super(message)
		this.name = "ApiError"
		this.status = status
		this.payload = payload
	}
}

/**
 * HTTP-клиент на axios.
 *
 * Я держу единый слой для:
 * - baseURL/таймаутов
 * - предсказуемых ошибок (ApiError)
 * - типизированных JSON ответов
 */
export class ApiClient {
	private readonly http: AxiosInstance

	constructor(http?: AxiosInstance) {
		this.http =
			http ??
			axios.create({
				baseURL: getBackendBaseUrl(),
				timeout: 15_000,
			})
	}

	async getJson<T>(path: string, config: AxiosRequestConfig = {}): Promise<T> {
		return this.requestJson<T>({ ...config, method: "GET", url: path })
	}

	async postJson<T>(
		path: string,
		body: unknown,
		config: AxiosRequestConfig = {}
	): Promise<T> {
		return this.requestJson<T>({
			...config,
			method: "POST",
			url: path,
			data: body,
		})
	}

	private async requestJson<T>(config: AxiosRequestConfig): Promise<T> {
		try {
			const res: AxiosResponse<T> = await this.http.request<T>(config)
			return res.data
		} catch (err: any) {
			const status = err?.response?.status ?? 0
			const payload = err?.response?.data ?? null
			const message =
				typeof payload?.message === "string"
					? payload.message
					: typeof err?.message === "string"
					? err.message
					: `Ошибка запроса ${String(
							config.method ?? "GET"
					  ).toUpperCase()} ${String(config.url ?? "")}`

			throw new ApiError(message, status, payload)
		}
	}
}
