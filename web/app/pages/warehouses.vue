<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Склады</h1>
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
            placeholder="Поиск по складам"
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
                <th class="py-2">Склад</th>
                <th class="py-2">Адрес</th>
                <th class="py-2 w-32">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in crm.warehouses" :key="item.warehouse_id" class="border-b border-gray-100">
                <td class="py-2">{{ item.warehouse_id }}</td>
                <td class="py-2">{{ item.warehouse_name }}</td>
                <td class="py-2">{{ item.warehouse_adress || '-' }}</td>
                <td class="py-2">
                  <div class="flex gap-1 sm:gap-2">
                    <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                      <span class="hidden sm:inline">Изменить</span>
                    </UButton>
                    <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.warehouse_id)">
                      <span class="hidden sm:inline">Удалить</span>
                    </UButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="!crm.warehouses.length && !isLoading" class="text-sm text-gray-600">Ничего не найдено.</p>

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

    <UModal v-model:open="isFormOpen" :title="form.warehouse_id ? 'Редактировать склад' : 'Добавить склад'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название склада" required>
            <UInput v-model="form.warehouse_name" placeholder="Название склада" required />
          </UFormField>
          <UFormField label="Адрес">
            <UInput v-model="form.warehouse_adress" placeholder="Адрес" />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.warehouse_id ? 'Сохранить' : 'Создать' }}</UButton>
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
  warehouse_id: null,
  warehouse_name: '',
  warehouse_adress: ''
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
  (params) => crm.fetchWarehouses(params),
  computed(() => crm.pagination.warehouses),
  { perPage: 10 }
)

function resetForm() {
  form.warehouse_id = null
  form.warehouse_name = ''
  form.warehouse_adress = ''
}

function openCreate() {
  resetForm()
  isFormOpen.value = true
}

function edit(item) {
  form.warehouse_id = item.warehouse_id
  form.warehouse_name = item.warehouse_name
  form.warehouse_adress = item.warehouse_adress || ''
  isFormOpen.value = true
}

async function save() {
  if (!form.warehouse_name.trim()) return

  const payload = {
    warehouse_name: form.warehouse_name.trim(),
    warehouse_adress: form.warehouse_adress.trim()
  }

  if (form.warehouse_id) {
    await crm.updateWarehouse(form.warehouse_id, payload)
  } else {
    await crm.createWarehouse(payload)
  }

  await load()
  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteWarehouse(id)
  await load()

  if (form.warehouse_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
