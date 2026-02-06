import type { AppServices } from "@/di/services"
import "pinia"

declare module "pinia" {
	export interface PiniaCustomProperties {
		$services: AppServices
	}
}
