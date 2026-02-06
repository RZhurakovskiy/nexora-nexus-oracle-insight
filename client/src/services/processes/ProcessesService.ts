import type { ApiClient } from "@/services/http/ApiClient"
import type { KillProcessResponse } from "@/types/processes"
import type { AxiosInstance } from "axios"

export type ExportFormat = "json" | "csv"

export class ProcessesService {
	private readonly api: ApiClient
	private readonly http: AxiosInstance

	constructor(api: ApiClient, http: AxiosInstance) {
		this.api = api
		this.http = http
	}

	/**
	 * Завершить процесс по PID.
	 */
	killProcessById(pid: number): Promise<KillProcessResponse> {
		if (!Number.isFinite(pid) || pid <= 0) {
			throw new Error("Некорректный PID")
		}
		return this.api.postJson<KillProcessResponse>("/api/kill-process-by-id", {
			pid,
		})
	}

	/**
	 * Экспортировать процессы в JSON/CSV.
	 *
	 * Сервер отдаёт "Content-Disposition: attachment; и т.д.",
	 * но на клиенте я всё равно задаю fallback-имя для предсказуемости.
	 */
	async exportProcesses(format: ExportFormat = "json"): Promise<void> {
		const res = await this.http.get(`/api/export/processes`, {
			params: { format },
			responseType: "blob",
		})

		const contentType =
			(res.headers?.["content-type"] as string | undefined) ??
			"application/octet-stream"
		const blob = new Blob([res.data], { type: contentType })
		const url = window.URL.createObjectURL(blob)

		const a = document.createElement("a")
		a.href = url
		a.download = `processes_${new Date()
			.toISOString()
			.slice(0, 19)
			.replace(/:/g, "-")}.${format}`
		document.body.appendChild(a)
		a.click()
		document.body.removeChild(a)

		window.URL.revokeObjectURL(url)
	}
}
