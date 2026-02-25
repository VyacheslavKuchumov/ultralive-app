import { computed, ref, watch } from 'vue'

function defaultPagination() {
  return {
    page: 1,
    per_page: 10,
    total: 0,
    total_pages: 1
  }
}

export function useServerList(fetchPage, paginationRef, options = {}) {
  const search = ref(options.search || '')
  const page = ref(options.page || 1)
  const perPage = ref(options.perPage || 10)
  const isLoading = ref(false)

  const pagination = computed(() => ({
    ...defaultPagination(),
    ...(paginationRef.value || {})
  }))

  const from = computed(() => {
    if (!pagination.value.total) return 0
    return (page.value - 1) * perPage.value + 1
  })

  const to = computed(() => {
    if (!pagination.value.total) return 0
    return Math.min(page.value * perPage.value, pagination.value.total)
  })

  async function load() {
    isLoading.value = true
    try {
      await fetchPage({
        search: search.value.trim(),
        page: page.value,
        per_page: perPage.value
      })

      const totalPages = Math.max(1, Number(pagination.value.total_pages) || 1)
      if (page.value > totalPages) {
        page.value = totalPages
      }
    } finally {
      isLoading.value = false
    }
  }

  function prevPage() {
    if (page.value > 1) {
      page.value -= 1
    }
  }

  function nextPage() {
    const totalPages = Math.max(1, Number(pagination.value.total_pages) || 1)
    if (page.value < totalPages) {
      page.value += 1
    }
  }

  watch(search, () => {
    page.value = 1
  })

  watch(perPage, () => {
    page.value = 1
  })

  watch([search, page, perPage], load, { immediate: true })

  return {
    search,
    page,
    perPage,
    isLoading,
    pagination,
    from,
    to,
    load,
    prevPage,
    nextPage
  }
}
