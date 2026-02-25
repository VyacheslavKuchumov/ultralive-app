<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">{{ currentSetName ? `Оборудование комплекта: ${currentSetName}` : 'Оборудование комплекта' }}</h1>
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
            placeholder="Поиск по оборудованию"
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
                <th class="py-2">Серия</th>
                <th class="py-2">Комплект</th>
                <th class="py-2">Склад</th>
                <th class="py-2">ТО</th>
                <th class="py-2 w-32">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in crm.equipment" :key="item.equipment_id" class="border-b border-gray-100">
                <td class="py-2">{{ item.equipment_id }}</td>
                <td class="py-2">{{ item.equipment_name }}</td>
                <td class="py-2">{{ item.serial_number }}</td>
                <td class="py-2">{{ item.equipment_set?.equipment_set_name || '-' }}</td>
                <td class="py-2">{{ item.storage?.warehouse_name || '-' }}</td>
                <td class="py-2">{{ item.needs_maintenance ? 'Да' : 'Нет' }}</td>
                <td class="py-2">
                  <div class="flex gap-1 sm:gap-2">
                    <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                      <span class="hidden sm:inline">Изменить</span>
                    </UButton>
                    <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.equipment_id)">
                      <span class="hidden sm:inline">Удалить</span>
                    </UButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="!crm.equipment.length && !isLoading" class="text-sm text-gray-600">Ничего не найдено.</p>

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

    <UModal v-model:open="isFormOpen" :title="form.equipment_id ? 'Редактировать оборудование' : 'Добавить оборудование'">
      <template #body>
        <form class="grid gap-3 md:grid-cols-2" @submit.prevent="save">
          <UFormField label="Название" required>
            <UInput v-model="form.equipment_name" size="lg" placeholder="Название" required />
          </UFormField>
          <UFormField label="Серийный номер" required>
            <UInput v-model="form.serial_number" size="lg" placeholder="Серийный номер" required />
          </UFormField>

          <UFormField label="Комплект" required>
            <UInput :model-value="currentSetName" size="lg" placeholder="Комплект" disabled />
          </UFormField>
          <UFormField label="Склад" required>
            <USelect v-model="form.warehouse_name" :items="warehouseOptions" :portal="false" size="lg" placeholder="Склад" />
          </UFormField>

          <UFormField label="Описание" class="md:col-span-2">
            <UInput v-model="form.description" size="lg" placeholder="Описание" />
          </UFormField>
          <UFormField label="Текущее хранение" class="md:col-span-2">
            <UInput v-model="form.current_storage_name" size="lg" placeholder="Текущее хранение" />
          </UFormField>

          <UFormField label="Дата покупки">
            <UInputDate v-model="form.date_of_purchase_value" class="w-full" />
          </UFormField>
          <UFormField label="Стоимость">
            <UInput v-model="form.cost_of_purchase" size="lg" placeholder="Стоимость" />
          </UFormField>

          <UCheckbox v-model="form.needs_maintenance" label="Требует обслуживания" class="md:col-span-2" />

          <div class="md:col-span-2 flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.equipment_id ? 'Сохранить' : 'Создать' }}</UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup>
import { CalendarDate } from '@internationalized/date'
import { computed, reactive, ref, watch } from 'vue'
import { navigateTo, useRoute } from '#imports'
import { useServerList } from '~/composables/useServerList'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const route = useRoute()
const isFormOpen = ref(false)
const perPageOptions = [10, 20, 50]

const form = reactive({
  equipment_id: null,
  equipment_name: '',
  serial_number: '',
  description: '',
  warehouse_name: '',
  current_storage_name: '',
  needs_maintenance: false,
  date_of_purchase_value: null,
  cost_of_purchase: ''
})

const setId = computed(() => {
  const value = Number(route.query.set || 0)
  return Number.isFinite(value) && value > 0 ? value : null
})

const currentSet = computed(() => {
  if (!setId.value) return null
  return crm.equipmentSets.find((item) => item.equipment_set_id === setId.value) || null
})

const currentSetName = computed(() => currentSet.value?.equipment_set_name || '')
const warehouseOptions = computed(() => crm.warehouses.map((item) => item.warehouse_name))

await Promise.all([
  crm.fetchEquipmentSets({ page: 1, per_page: 1000 }),
  crm.fetchWarehouses({ page: 1, per_page: 1000 })
])

if (!setId.value || !currentSet.value) {
  await navigateTo('/equipment_sets')
}

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
  (params) => crm.fetchEquipment(setId.value, params),
  computed(() => crm.pagination.equipment),
  { perPage: 10 }
)

function parseDateValue(raw) {
  if (!raw || typeof raw !== 'string') return null
  const [year, month, day] = raw.split('-').map((item) => Number(item))
  if (!year || !month || !day) return null
  return new CalendarDate(year, month, day)
}

function formatDateValue(value) {
  if (!value || typeof value.year !== 'number') return ''
  const year = String(value.year).padStart(4, '0')
  const month = String(value.month).padStart(2, '0')
  const day = String(value.day).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function resetForm() {
  form.equipment_id = null
  form.equipment_name = ''
  form.serial_number = ''
  form.description = ''
  form.warehouse_name = ''
  form.current_storage_name = ''
  form.needs_maintenance = false
  form.date_of_purchase_value = null
  form.cost_of_purchase = ''
}

async function ensureSetContext() {
  await crm.fetchEquipmentSets({ page: 1, per_page: 1000 })

  if (!setId.value || !currentSet.value) {
    await navigateTo('/equipment_sets')
    return false
  }

  return true
}

watch(
  () => route.query.set,
  async () => {
    const valid = await ensureSetContext()
    if (!valid) return

    page.value = 1
    await load()
  }
)

async function openCreate() {
  if (!currentSetName.value) {
    await navigateTo('/equipment_sets')
    return
  }

  if (!crm.warehouses.length) {
    await crm.fetchWarehouses({ page: 1, per_page: 1000 })
  }

  resetForm()
  isFormOpen.value = true
}

async function edit(item) {
  if (!currentSetName.value) {
    await navigateTo('/equipment_sets')
    return
  }

  if (!crm.warehouses.length) {
    await crm.fetchWarehouses({ page: 1, per_page: 1000 })
  }

  form.equipment_id = item.equipment_id
  form.equipment_name = item.equipment_name
  form.serial_number = item.serial_number
  form.description = item.description || ''
  form.warehouse_name = item.storage?.warehouse_name || ''
  form.current_storage_name = item.current_storage || ''
  form.needs_maintenance = item.needs_maintenance
  form.date_of_purchase_value = parseDateValue(item.date_of_purchase || '')
  form.cost_of_purchase = item.cost_of_purchase || ''
  isFormOpen.value = true
}

function payloadFromForm() {
  return {
    equipment_name: form.equipment_name.trim(),
    serial_number: form.serial_number.trim(),
    equipment_set_name: currentSetName.value,
    description: form.description.trim(),
    warehouse_name: form.warehouse_name,
    current_storage_name: form.current_storage_name.trim(),
    needs_maintenance: !!form.needs_maintenance,
    date_of_purchase: formatDateValue(form.date_of_purchase_value),
    cost_of_purchase: form.cost_of_purchase ? Number(form.cost_of_purchase) : null
  }
}

async function save() {
  if (!currentSetName.value || !form.equipment_name.trim() || !form.serial_number.trim() || !form.warehouse_name) {
    return
  }

  const payload = payloadFromForm()

  if (form.equipment_id) {
    await crm.updateEquipment(form.equipment_id, payload)
  } else {
    await crm.createEquipment(payload)
  }

  await load()
  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteEquipment(id)
  await load()

  if (form.equipment_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
