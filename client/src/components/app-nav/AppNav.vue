<script setup lang="ts">
import { useRoute } from "vue-router"
import "./AppNav.css"

const route = useRoute()

/**
 * Проверка активности маршрута для подсветки.
 */
const isActive = (path: string): boolean => {
	return route.path === path
}
</script>

<template>
	<nav class="app-nav">
		<div class="app-nav-container">
			<router-link
				to="/"
				class="nav-link"
				:class="{ 'nav-link--active': isActive('/') }"
			>
				<div class="nav-link-icon">
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<circle
							cx="10"
							cy="10"
							r="7.5"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>

						<line
							x1="10"
							y1="7"
							x2="10"
							y2="13"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>

						<circle cx="10" cy="5" r="1" fill="currentColor" />
					</svg>
				</div>
				<span class="nav-link-text">Общая информация</span>
				<div class="nav-link-indicator"></div>
			</router-link>
			<router-link
				to="/metrics"
				class="nav-link"
				:class="{ 'nav-link--active': isActive('/metrics') }"
			>
				<div class="nav-link-icon">
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M3.33333 16.6667H16.6667"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<path
							d="M5 13.3333L7.5 8.33333L10 11.6667L12.5 6.66667L15 13.3333"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							fill="none"
						/>
						<circle cx="5" cy="13.3333" r="1.25" fill="currentColor" />
						<circle cx="7.5" cy="8.33333" r="1.25" fill="currentColor" />
						<circle cx="10" cy="11.6667" r="1.25" fill="currentColor" />
						<circle cx="12.5" cy="6.66667" r="1.25" fill="currentColor" />
						<circle cx="15" cy="13.3333" r="1.25" fill="currentColor" />
					</svg>
				</div>
				<span class="nav-link-text">Метрики</span>
				<div class="nav-link-indicator"></div>
			</router-link>

			<router-link
				to="/processes"
				class="nav-link"
				:class="{ 'nav-link--active': isActive('/processes') }"
			>
				<div class="nav-link-icon">
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<rect
							x="3.33333"
							y="3.33333"
							width="13.3333"
							height="13.3333"
							rx="1.5"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
						<circle cx="6.66667" cy="7.5" r="1" fill="currentColor" />
						<path
							d="M9.16667 7.5H13.3333"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<circle cx="6.66667" cy="10.8333" r="1" fill="currentColor" />
						<path
							d="M9.16667 10.8333H13.3333"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
						<circle cx="6.66667" cy="14.1667" r="1" fill="currentColor" />
						<path
							d="M9.16667 14.1667H11.6667"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
					</svg>
				</div>
				<span class="nav-link-text">Процессы</span>
				<div class="nav-link-indicator"></div>
			</router-link>

			<router-link
				to="/port-manager"
				class="nav-link"
				:class="{ 'nav-link--active': isActive('/port-manager') }"
			>
				<div class="nav-link-icon">
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<rect
							x="4"
							y="5"
							width="12"
							height="10"
							rx="2"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
						<circle cx="6" cy="8" r="1" fill="currentColor" />
						<circle cx="6" cy="12" r="1" fill="currentColor" />
						<circle cx="14" cy="8" r="1" fill="currentColor" />
						<circle cx="14" cy="12" r="1" fill="currentColor" />
						<line
							x1="16"
							y1="10"
							x2="19"
							y2="10"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
						/>
					</svg>
				</div>
				<span class="nav-link-text">Сетевые порты</span>
				<div class="nav-link-indicator"></div>
			</router-link>

			<router-link
				to="/metrics-history"
				class="nav-link"
				:class="{ 'nav-link--active': isActive('/metrics-history') }"
			>
				<div class="nav-link-icon">
					<svg
						width="20"
						height="20"
						viewBox="0 0 20 20"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M3 3H17V17H3V3Z"
							stroke="currentColor"
							stroke-width="1.5"
							fill="none"
						/>
						<path d="M3 8H17M8 3V17" stroke="currentColor" stroke-width="1.5" />
					</svg>
				</div>
				<span class="nav-link-text">История метрик</span>
				<div class="nav-link-indicator"></div>
			</router-link>
		</div>
	</nav>
</template>
