class AuthApp {
    constructor() {
        this.initRouter()
        this.bindEvents()
        this.renderView()
    }

    initRouter() {
        this.routes = {
            '/': this.homeView,
            '/login': this.loginView,
            '/register': this.registerView,
            '/profile': this.profileView
        }
    }

    bindEvents() {
        window.addEventListener('popstate', () => this.renderView())
        document.addEventListener('DOMContentLoaded', () => {
            document.body.addEventListener('click', e => {
                if (e.target.matches('[data-link]')) {
                    e.preventDefault()
                    this.navigateTo(e.target.href)
                }
            })
        })
    }

    navigateTo(url) {
        history.pushState(null, null, url)
        this.renderView()
    }

    async renderView() {
        const path = window.location.pathname
        const view = this.routes[path] || this.notFoundView
        document.getElementById('app').innerHTML = await view.call(this)
        this.updateActiveLink()
    }

    updateActiveLink() {
        document.querySelectorAll('[data-link]').forEach(link => {
            link.classList.toggle('active', link.pathname === window.location.pathname)
        })
    }

    async homeView() {
        return `
            <div class="view home-view">
                <h1>Добро пожаловать</h1>
                <div class="actions">
                    <a href="/login" data-link class="btn">Войти</a>
                    <a href="/register" data-link class="btn secondary">Регистрация</a>
                </div>
            </div>
        `
    }

    async loginView() {
        return `
            <div class="view auth-view">
                <h1>Вход в систему</h1>
                <form id="loginForm">
                    <div class="form-group">
                        <input type="text" name="username" placeholder="Имя пользователя" required>
                    </div>
                    <div class="form-group">
                        <input type="password" name="password" placeholder="Пароль" required>
                    </div>
                    <button type="submit" class="btn">Войти</button>
                </form>
                <div class="auth-link">
                    Нет аккаунта? <a href="/register" data-link>Зарегистрируйтесь</a>
                </div>
            </div>
        `
    }

    // Остальные методы view (registerView, profileView, notFoundView)
    // Методы для работы с API (login, register, fetchProfile)
}

new AuthApp()