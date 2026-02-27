<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Шаблоны</h1>
          <UButton color="primary" icon="i-lucide-plus" @click="openCreate">
            <span class="hidden sm:inline">Добавить</span>
          </UButton>
        </div>
      </template>
    </UCard>

    <UCard>
      <div class="space-y-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Поиск по шаблонам"
            class="md:max-w-sm"
          />

          <label class="flex items-center gap-2 text-sm text-gray-600">
            На странице
            <select v-model.number="perPage" class="rounded border border-gray-300 px-2 py-1 text-sm">
              <option v-for="option in perPageOptions" :key="option" :value="option">{{ option }}</option>
            </select>
          </label>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full min-w-max text-sm">
            <thead>
              <tr class="text-left border-b border-gray-200 whitespace-nowrap">
                <th class="py-2">ID</th>
                <th class="py-2">Название</th>
                <th class="py-2">Оборудование</th>
                <th class="py-2 w-48">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in crm.drafts" :key="item.draft_id" class="border-b border-gray-100">
                <td class="py-2">{{ item.draft_id }}</td>
                <td class="py-2">{{ item.draft_name }}</td>
                <td class="py-2">{{ item.equipment?.length || 0 }}</td>
                <td class="py-2">
                  <div class="flex gap-1 sm:gap-2">
                    <UButton size="xs" color="primary" variant="soft" icon="i-lucide-list" :to="`/draft/${item.draft_id}`" aria-label="Состав">
                      <span class="hidden sm:inline">Состав</span>
                    </UButton>
                    <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                      <span class="hidden sm:inline">Изменить</span>
                    </UButton>
                    <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.draft_id)">
                      <span class="hidden sm:inline">Удалить</span>
                    </UButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="!crm.drafts.length && !isLoading" class="text-sm text-gray-600">Ничего не найдено.</p>

        <div class="flex flex-col gap-3 border-t border-gray-100 pt-3 md:flex-row md:items-center md:justify-between">
          <p class="text-sm text-gray-600">Показано {{ from }}-{{ to }} из {{ pagination.total }}</p>

          <div class="flex items-center gap-2">
            <UButton size="xs" color="neutral" variant="soft" :disabled="page <= 1 || isLoading" @click="prevPage">Назад</UButton>
            <span class="text-sm text-gray-600">Стр. {{ page }} / {{ pagination.total_pages }}</span>
            <UButton
              size="xs"
              color="neutral"
              variant="soft"
              :disabled="page >= pagination.total_pages || isLoading"
              @click="nextPage"
            >
              Вперед
            </UButton>
          </div>
        </div>
      </div>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.draft_id ? 'Редактировать шаблон' : 'Добавить шаблон'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название шаблона" required>
            <UInput v-model="form.draft_name" placeholder="Название шаблона" required />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.draft_id ? 'Сохранить' : 'Создать' }}</UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { useServerList } from '~/composables/useServerList'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const isFormOpen = ref(false)
const perPageOptions = [10, 20, 50]

const form = reactive({
  draft_id: null,
  draft_name: ''
})

const {
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
} = useServerList(
  (params) => crm.fetchDrafts(params),
  computed(() => crm.pagination.drafts),
  { perPage: 10 }
)

function resetForm() {
  form.draft_id = null
  form.draft_name = ''
}

function openCreate() {
  resetForm()
  isFormOpen.value = true
}

function edit(item) {
  form.draft_id = item.draft_id
  form.draft_name = item.draft_name
  isFormOpen.value = true
}

async function save() {
  if (!form.draft_name.trim()) return

  if (form.draft_id) {
    await crm.updateDraft(form.draft_id, { draft_name: form.draft_name.trim() })
  } else {
    await crm.createDraft({ draft_name: form.draft_name.trim() })
  }

  await load()
  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteDraft(id)
  await load()

  if (form.draft_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
