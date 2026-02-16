function parseTokenPayload(token) {
  try {
    const payload = token.split('.')[1]
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
    const decoded =
      typeof atob === 'function'
        ? atob(normalized)
        : Buffer.from(normalized, 'base64').toString('utf8')
    return JSON.parse(decoded)
  } catch {
    return null
  }
}

function getTokenExpiryMs(token) {
  const payload = parseTokenPayload(token)
  const expSeconds = Number(payload?.expiredAt)
  if (!Number.isFinite(expSeconds)) return 0
  return expSeconds * 1000
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    userId: null,
    profile: null
  }),
  persist: true,
  getters: {
    isAuthenticated(state) {
      if (!state.token) return false
      const expiryMs = getTokenExpiryMs(state.token)
      if (!expiryMs) return false
      return Date.now() < expiryMs
    },
    displayName(state) {
      return state.profile?.name || state.profile?.email || 'Пользователь'
    }
  },
  actions: {
    authHeader() {
      if (!this.token) return {}
      return { Authorization: `Bearer ${this.token}` }
    },

    hydrateFromToken() {
      if (!this.token) {
        this.userId = null
        this.profile = null
        return
      }

      const payload = parseTokenPayload(this.token)
      const parsedId = Number(payload?.userID)
      this.userId = Number.isFinite(parsedId) && parsedId > 0 ? parsedId : null

      if (!this.isAuthenticated) {
        this.logout(false)
      }
    },

    async login({ email, password }) {
      const response = await $fetch('/api/backend/login', {
        method: 'POST',
        body: { email, password }
      })

      this.token = response.token
      this.hydrateFromToken()
      await this.fetchProfile()
    },

    async signup({ firstName, lastName, email, password }) {
      await $fetch('/api/backend/register', {
        method: 'POST',
        body: { firstName, lastName, email, password }
      })

      await this.login({ email, password })
    },

    async fetchProfile() {
      if (!this.token) return null

      const profile = await $fetch('/api/backend/profile', {
        headers: this.authHeader()
      })

      this.profile = profile
      this.userId = Number(profile?.id) || this.userId
      return profile
    },

    async updateProfile({ firstName, lastName, email, currentPassword }) {
      if (!this.token) return null

      const profile = await $fetch('/api/backend/profile', {
        method: 'PUT',
        body: { firstName, lastName, email, currentPassword },
        headers: this.authHeader()
      })

      this.profile = profile
      this.userId = Number(profile?.id) || this.userId
      return profile
    },

    async updatePassword({ currentPassword, newPassword }) {
      if (!this.token) return null

      await $fetch('/api/backend/profile/password', {
        method: 'PUT',
        body: { currentPassword, newPassword },
        headers: this.authHeader()
      })
    },

    logout(redirect = true) {
      this.token = null
      this.userId = null
      this.profile = null

      if (redirect) {
        navigateTo('/login')
      }
    }
  }
})
