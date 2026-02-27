<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Архив съёмок</h1>
      </template>

      <div class="space-y-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Поиск по архиву"
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
                <th class="py-2">Съёмка</th>
                <th class="py-2">Площадка</th>
                <th class="py-2">Инженер</th>
                <th class="py-2">Период</th>
                <th class="py-2 w-48">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in crm.archivedProjects" :key="item.project_id" class="border-b border-gray-100">
                <td class="py-2">{{ item.project_id }}</td>
                <td class="py-2">{{ item.project_name }}</td>
                <td class="py-2">{{ item.type?.project_type_name || '-' }}</td>
                <td class="py-2">{{ item.chiefEngineer?.name || '-' }}</td>
                <td class="py-2">{{ item.shooting_start_date }} - {{ item.shooting_end_date }}</td>
                <td class="py-2">
                  <div class="flex gap-2">
                    <UButton size="xs" color="neutral" variant="soft" @click="restore(item)">Восстановить</UButton>
                    <UButton size="xs" color="error" variant="soft" @click="remove(item.project_id)">Удалить</UButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="!crm.archivedProjects.length && !isLoading" class="text-sm text-gray-600">Ничего не найдено.</p>

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
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useServerList } from '~/composables/useServerList'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const perPageOptions = [10, 20, 50]

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
  (params) => crm.fetchArchivedProjects(params),
  computed(() => crm.pagination.archivedProjects),
  { perPage: 10 }
)

async function restore(item) {
  await crm.updateProject(item.project_id, {
    project_name: item.project_name,
    project_type_name: item.type?.project_type_name || '',
    chief_engineer_name: item.chiefEngineer?.name || '',
    shooting_start_date: item.shooting_start_date,
    shooting_end_date: item.shooting_end_date,
    archived: false
  })

  await load()
}

async function remove(id) {
  await crm.deleteProject(id)
  await load()
}
</script>
