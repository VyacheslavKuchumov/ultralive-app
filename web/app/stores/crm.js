import { useAuthStore } from '~/stores/auth'

let lastErrorToast = {
  key: '',
  at: 0
}

function extractBackendErrorMessage(error) {
  return (
    error?.data?.statusMessage ||
    error?.data?.error ||
    error?.data?.message ||
    error?.statusMessage ||
    error?.message ||
    'Не удалось выполнить запрос к API'
  )
}

function notifyBackendError(message) {
  if (!import.meta.client) return

  const now = Date.now()
  const key = message
  if (lastErrorToast.key === key && now - lastErrorToast.at < 4000) return

  lastErrorToast = { key, at: now }

  const toast = useToast()
  toast.add({
    title: 'Ошибка загрузки данных',
    description: message,
    color: 'error'
  })
}

function defaultPagination() {
  return {
    page: 1,
    per_page: 10,
    total: 0,
    total_pages: 1
  }
}

function fallbackListResponse(params = {}) {
  const page = Number(params?.page || 1)
  const perPage = Number(params?.per_page || 10)

  return {
    items: [],
    pagination: {
      page: Number.isFinite(page) && page > 0 ? page : 1,
      per_page: Number.isFinite(perPage) && perPage > 0 ? perPage : 10,
      total: 0,
      total_pages: 1
    }
  }
}

function normalizeListResponse(response) {
  if (response && Array.isArray(response.items)) {
    return {
      items: response.items,
      pagination: {
        ...defaultPagination(),
        ...(response.pagination || {})
      }
    }
  }

  if (Array.isArray(response)) {
    return {
      items: response,
      pagination: {
        page: 1,
        per_page: response.length || 10,
        total: response.length,
        total_pages: 1
      }
    }
  }

  return {
    items: [],
    pagination: defaultPagination()
  }
}

function applyListState(store, stateKey, paginationKey, response) {
  const normalized = normalizeListResponse(response)
  store[stateKey] = normalized.items
  if (paginationKey && store.pagination?.[paginationKey]) {
    store.pagination[paginationKey] = normalized.pagination
  }
  return normalized.items
}

async function backendRequest(path, options = {}) {
  const auth = useAuthStore()
  const method = options.method || 'GET'
  const throwOnError = options.throwOnError !== false

  try {
    return await $fetch(`/api/backend${path}`, {
      method,
      body: options.body,
      query: options.query,
      headers: auth.authHeader()
    })
  } catch (error) {
    const message = extractBackendErrorMessage(error)
    notifyBackendError(message)

    if (!throwOnError) {
      return options.fallback
    }

    throw error
  }
}

export const useCRMStore = defineStore('crm', {
  state: () => ({
    users: [],
    setTypes: [],
    projectTypes: [],
    warehouses: [],
    equipmentSets: [],
    equipment: [],
    projects: [],
    archivedProjects: [],
    drafts: [],
    currentProject: null,
    currentDraft: null,
    projectBoard: null,
    draftBoard: null,
    pagination: {
      setTypes: defaultPagination(),
      projectTypes: defaultPagination(),
      warehouses: defaultPagination(),
      equipmentSets: defaultPagination(),
      equipment: defaultPagination(),
      projects: defaultPagination(),
      archivedProjects: defaultPagination(),
      drafts: defaultPagination()
    }
  }),
  actions: {
    async fetchUsers() {
      this.users = await backendRequest('/users', { throwOnError: false, fallback: [] })
      return this.users
    },

    async fetchSetTypes(params = {}) {
      const response = await backendRequest('/set_types', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'setTypes', 'setTypes', response)
    },

    async createSetType(payload) {
      this.setTypes = await backendRequest('/set_types', { method: 'POST', body: payload })
      return this.setTypes
    },

    async updateSetType(id, payload) {
      this.setTypes = await backendRequest(`/set_types/${id}`, { method: 'PUT', body: payload })
      return this.setTypes
    },

    async deleteSetType(id) {
      this.setTypes = await backendRequest(`/set_types/${id}`, { method: 'DELETE' })
      return this.setTypes
    },

    async fetchProjectTypes(params = {}) {
      const response = await backendRequest('/project_types', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'projectTypes', 'projectTypes', response)
    },

    async createProjectType(payload) {
      this.projectTypes = await backendRequest('/project_types', { method: 'POST', body: payload })
      return this.projectTypes
    },

    async updateProjectType(id, payload) {
      this.projectTypes = await backendRequest(`/project_types/${id}`, { method: 'PUT', body: payload })
      return this.projectTypes
    },

    async deleteProjectType(id) {
      this.projectTypes = await backendRequest(`/project_types/${id}`, { method: 'DELETE' })
      return this.projectTypes
    },

    async fetchWarehouses(params = {}) {
      const response = await backendRequest('/warehouse', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'warehouses', 'warehouses', response)
    },

    async createWarehouse(payload) {
      this.warehouses = await backendRequest('/warehouse', { method: 'POST', body: payload })
      return this.warehouses
    },

    async updateWarehouse(id, payload) {
      this.warehouses = await backendRequest(`/warehouse/${id}`, { method: 'PUT', body: payload })
      return this.warehouses
    },

    async deleteWarehouse(id) {
      this.warehouses = await backendRequest(`/warehouse/${id}`, { method: 'DELETE' })
      return this.warehouses
    },

    async fetchEquipmentSets(params = {}) {
      const response = await backendRequest('/equipment_set', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'equipmentSets', 'equipmentSets', response)
    },

    async createEquipmentSet(payload) {
      this.equipmentSets = await backendRequest('/equipment_set', { method: 'POST', body: payload })
      return this.equipmentSets
    },

    async updateEquipmentSet(id, payload) {
      this.equipmentSets = await backendRequest(`/equipment_set/${id}`, { method: 'PUT', body: payload })
      return this.equipmentSets
    },

    async deleteEquipmentSet(id) {
      this.equipmentSets = await backendRequest(`/equipment_set/${id}`, { method: 'DELETE' })
      return this.equipmentSets
    },

    async fetchEquipment(setId = null, params = {}) {
      const path = setId ? `/equipment/set/${setId}` : '/equipment'
      const response = await backendRequest(path, {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'equipment', 'equipment', response)
    },

    async createEquipment(payload) {
      this.equipment = await backendRequest('/equipment', { method: 'POST', body: payload })
      return this.equipment
    },

    async updateEquipment(id, payload) {
      this.equipment = await backendRequest(`/equipment/${id}`, { method: 'PUT', body: payload })
      return this.equipment
    },

    async deleteEquipment(id) {
      await backendRequest(`/equipment/${id}`, { method: 'DELETE' })
      this.equipment = this.equipment.filter((item) => item.equipment_id !== id)
      return this.equipment
    },

    async fetchProjects(params = {}) {
      const response = await backendRequest('/projects', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'projects', 'projects', response)
    },

    async fetchArchivedProjects(params = {}) {
      const response = await backendRequest('/projects/archived', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'archivedProjects', 'archivedProjects', response)
    },

    async fetchProjectById(id) {
      this.currentProject = await backendRequest(`/projects/search/${id}`, { throwOnError: false, fallback: null })
      return this.currentProject
    },

    async createProject(payload) {
      this.projects = await backendRequest('/projects', { method: 'POST', body: payload })
      return this.projects
    },

    async updateProject(id, payload) {
      const updated = await backendRequest(`/projects/${id}`, { method: 'PUT', body: payload })
      if (payload.archived) {
        this.archivedProjects = updated
      } else {
        this.projects = updated
      }
      return updated
    },

    async deleteProject(id) {
      this.projects = await backendRequest(`/projects/${id}`, { method: 'DELETE' })
      return this.projects
    },

    async fetchDrafts(params = {}) {
      const response = await backendRequest('/drafts', {
        throwOnError: false,
        query: params,
        fallback: fallbackListResponse(params)
      })
      return applyListState(this, 'drafts', 'drafts', response)
    },

    async fetchDraftById(id) {
      this.currentDraft = await backendRequest(`/drafts/search/${id}`, { throwOnError: false, fallback: null })
      return this.currentDraft
    },

    async createDraft(payload) {
      this.drafts = await backendRequest('/drafts', { method: 'POST', body: payload })
      return this.drafts
    },

    async updateDraft(id, payload) {
      this.drafts = await backendRequest(`/drafts/${id}`, { method: 'PUT', body: payload })
      return this.drafts
    },

    async deleteDraft(id) {
      this.drafts = await backendRequest(`/drafts/${id}`, { method: 'DELETE' })
      return this.drafts
    },

    async fetchProjectBoard(projectId) {
      this.projectBoard = await backendRequest(`/equipment_in_project/${projectId}`, { throwOnError: false, fallback: null })
      return this.projectBoard
    },

    async addEquipmentToProject(payload) {
      this.projectBoard = await backendRequest('/equipment_in_project/add', { method: 'POST', body: payload })
      return this.projectBoard
    },

    async removeEquipmentFromProject(payload) {
      this.projectBoard = await backendRequest('/equipment_in_project/del', { method: 'PUT', body: payload })
      return this.projectBoard
    },

    async addSetToProject(payload) {
      this.projectBoard = await backendRequest('/equipment_in_project/add_set', { method: 'POST', body: payload })
      return this.projectBoard
    },

    async removeSetFromProject(payload) {
      this.projectBoard = await backendRequest('/equipment_in_project/del_set', { method: 'PUT', body: payload })
      return this.projectBoard
    },

    async addDraftToProject(payload) {
      this.projectBoard = await backendRequest('/equipment_in_project/add_draft', { method: 'POST', body: payload })
      return this.projectBoard
    },

    async resetEquipmentInProject(projectId) {
      await backendRequest(`/equipment_in_project/reset/${projectId}`, { method: 'DELETE' })
      return this.fetchProjectBoard(projectId)
    },

    async fetchConflictingEquipment(projectId) {
      return backendRequest('/equipment_in_project/conflicting', {
        method: 'POST',
        body: { project_id: projectId },
        throwOnError: false,
        fallback: []
      })
    },

    async fetchConflictingProjects() {
      return backendRequest('/equipment_in_project/conflicting_projects', {
        method: 'POST',
        throwOnError: false,
        fallback: []
      })
    },

    async fetchDraftBoard(draftId) {
      this.draftBoard = await backendRequest(`/equipment_in_draft/${draftId}`, { throwOnError: false, fallback: null })
      return this.draftBoard
    },

    async addEquipmentToDraft(payload) {
      this.draftBoard = await backendRequest('/equipment_in_draft/add', { method: 'POST', body: payload })
      return this.draftBoard
    },

    async removeEquipmentFromDraft(payload) {
      this.draftBoard = await backendRequest('/equipment_in_draft/del', { method: 'PUT', body: payload })
      return this.draftBoard
    },

    async addSetToDraft(payload) {
      this.draftBoard = await backendRequest('/equipment_in_draft/add_set', { method: 'POST', body: payload })
      return this.draftBoard
    },

    async removeSetFromDraft(payload) {
      this.draftBoard = await backendRequest('/equipment_in_draft/del_set', { method: 'PUT', body: payload })
      return this.draftBoard
    },

    async refreshLookups() {
      await Promise.all([
        this.fetchUsers(),
        this.fetchSetTypes(),
        this.fetchProjectTypes(),
        this.fetchWarehouses(),
        this.fetchEquipmentSets()
      ])
    }
  }
})
