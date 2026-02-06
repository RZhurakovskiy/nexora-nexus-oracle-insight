import type { InjectionKey } from 'vue'
import type { AppServices } from './services'

export const ServicesKey: InjectionKey<AppServices> = Symbol('AppServices')


