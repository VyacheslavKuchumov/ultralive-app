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
            <UInput v-model="form.date_of_purchase" size="lg" placeholder="Дата покупки YYYY-MM-DD" />
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
import { computed, reactive, ref, watch } from 'vue'
import { navigateTo, useRoute } from '#imports'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const route = useRoute()
const isFormOpen = ref(false)

const form = reactive({
  equipment_id: null,
  equipment_name: '',
  serial_number: '',
  description: '',
  warehouse_name: '',
  current_storage_name: '',
  needs_maintenance: false,
  date_of_purchase: '',
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

function resetForm() {
  form.equipment_id = null
  form.equipment_name = ''
  form.serial_number = ''
  form.description = ''
  form.warehouse_name = ''
  form.current_storage_name = ''
  form.needs_maintenance = false
  form.date_of_purchase = ''
  form.cost_of_purchase = ''
}

async function openCreate() {
  if (!currentSetName.value) {
    await navigateTo('/equipment_sets')
    return
  }

  if (!crm.warehouses.length) {
    await crm.fetchWarehouses()
  }

  resetForm()
  isFormOpen.value = true
}

async function load() {
  await Promise.all([crm.fetchEquipmentSets(), crm.fetchWarehouses()])

  if (!setId.value || !currentSet.value) {
    await navigateTo('/equipment_sets')
    return
  }

  await crm.fetchEquipment(setId.value)
}

await load()

watch(() => route.query.set, load)

async function edit(item) {
  if (!currentSetName.value) {
    await navigateTo('/equipment_sets')
    return
  }

  if (!crm.warehouses.length) {
    await crm.fetchWarehouses()
  }

  form.equipment_id = item.equipment_id
  form.equipment_name = item.equipment_name
  form.serial_number = item.serial_number
  form.description = item.description || ''
  form.warehouse_name = item.storage?.warehouse_name || ''
  form.current_storage_name = item.current_storage || ''
  form.needs_maintenance = item.needs_maintenance
  form.date_of_purchase = item.date_of_purchase || ''
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
    date_of_purchase: form.date_of_purchase.trim(),
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

  if (setId.value) {
    await crm.fetchEquipment(setId.value)
  }

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteEquipment(id)
  if (form.equipment_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
