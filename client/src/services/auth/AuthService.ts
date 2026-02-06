import CryptoJS from "crypto-js"

const AUTH_TOKEN_KEY = "auth_token"

const ADMIN_LOGIN = "admin"
const ADMIN_PASSWORD_HASH = "b9264699701f57a49a1f845811c1f91c296f452398a80410eba3f255273eb00a"

/**
 * Хеширование пароля с использованием SHA-256.
 */
function hashPassword(password: string): string {
	return CryptoJS.SHA256(password).toString()
}

/**
 * Хеширование логина с использованием SHA-256.
 */
function hashLogin(login: string): string {
	return CryptoJS.SHA256(login.toLowerCase().trim()).toString()
}

export class AuthService {
	/**
	 * Выполнить вход в систему.
	 */
	async login(login: string, password: string): Promise<string> {
		await new Promise(resolve => setTimeout(resolve, 500))

		const loginHash = hashLogin(login)
		const passwordHash = hashPassword(password)
		const expectedLoginHash = hashLogin(ADMIN_LOGIN)

		if (loginHash === expectedLoginHash && passwordHash === ADMIN_PASSWORD_HASH) {
			const token = `nexora_token_${Date.now()}_${Math.random().toString(36).substring(2, 15)}`
			localStorage.setItem(AUTH_TOKEN_KEY, token)
			return token
		}

		throw new Error("Неверный логин или пароль")
	}

	/**
	 * Получить токен авторизации из localStorage.
	 */
	getToken(): string | null {
		return localStorage.getItem(AUTH_TOKEN_KEY)
	}

	/**
	 * Проверить наличие токена авторизации.
	 */
	isAuthenticated(): boolean {
		return this.getToken() !== null
	}

	/**
	 * Выйти из системы.
	 */
	logout(): void {
		localStorage.removeItem(AUTH_TOKEN_KEY)
	}
}

