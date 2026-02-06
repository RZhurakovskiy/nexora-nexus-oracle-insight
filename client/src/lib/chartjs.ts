import {
	CategoryScale,
	Chart as ChartJS,
	Legend,
	LinearScale,
	LineController,
	LineElement,
	PointElement,
	Title,
	Tooltip,
} from "chart.js"

let isRegistered = false

/**
 * Регистрация модулей Chart.js.
 *
 * Я держу регистрацию в одном месте и с guard-флагом, чтобы:
 * - не дублировать код по компонентам
 * - не регистрировать плагины несколько раз
 */
export function ensureChartJsRegistered(): void {
	if (isRegistered) return
	ChartJS.register(
		LineController,
		CategoryScale,
		LinearScale,
		PointElement,
		LineElement,
		Title,
		Tooltip,
		Legend
	)
	isRegistered = true
}
