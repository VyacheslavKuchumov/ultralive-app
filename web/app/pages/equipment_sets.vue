<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Комплекты оборудования</h1>
          <UButton color="primary" icon="i-lucide-plus" @click="openCreate">
            <span class="hidden sm:inline">Добавить</span>
          </UButton>
        </div>
      </template>
    </UCard>

    <template v-if="crm.equipmentSets.length">
      <div class="grid gap-4">
        <UCard
          v-for="item in crm.equipmentSets"
          :key="item.equipment_set_id"
          :data-testid="`equipment-set-card-${item.equipment_set_id}`"
        >
          <template #header>
            <div class="flex items-start justify-between gap-3">
              <div class="space-y-1">
                <h2 class="text-lg font-semibold">{{ item.equipment_set_name }}</h2>
                <p class="text-sm text-gray-600">Тип: {{ item.type?.set_type_name || '-' }}</p>
                <p class="text-sm text-gray-600">{{ item.description || 'Без описания' }}</p>
              </div>

              <div class="flex gap-1 sm:gap-2 shrink-0">
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-lucide-package"
                  :to="`/equipment?set=${item.equipment_set_id}`"
                  aria-label="Открыть комплект"
                >
                  <span class="hidden sm:inline">Открыть комплект</span>
                </UButton>
                <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                  <span class="hidden sm:inline">Изменить</span>
                </UButton>
                <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.equipment_set_id)">
                  <span class="hidden sm:inline">Удалить</span>
                </UButton>
              </div>
            </div>
          </template>

          <div v-if="equipmentInSet(item.equipment_set_id).length" class="overflow-x-auto">
            <table class="w-full min-w-max text-sm">
              <thead>
                <tr class="text-left border-b border-gray-200 whitespace-nowrap">
                  <th class="py-2">ID</th>
                  <th class="py-2">Название</th>
                  <th class="py-2">Серия</th>
                  <th class="py-2">Склад</th>
                  <th class="py-2">ТО</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="equipmentItem in equipmentInSet(item.equipment_set_id)" :key="equipmentItem.equipment_id" class="border-b border-gray-100">
                  <td class="py-2">{{ equipmentItem.equipment_id }}</td>
                  <td class="py-2">{{ equipmentItem.equipment_name }}</td>
                  <td class="py-2">{{ equipmentItem.serial_number }}</td>
                  <td class="py-2">{{ equipmentItem.storage?.warehouse_name || '-' }}</td>
                  <td class="py-2">{{ equipmentItem.needs_maintenance ? 'Да' : 'Нет' }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p v-else class="text-sm text-gray-600">В комплекте пока нет оборудования.</p>
        </UCard>
      </div>
    </template>

    <UCard v-else>
      <p class="text-sm text-gray-600">Комплекты пока не добавлены.</p>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.equipment_set_id ? 'Редактировать комплект' : 'Добавить комплект'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название комплекта" required>
            <UInput v-model="form.equipment_set_name" placeholder="Название комплекта" required />
          </UFormField>
          <UFormField label="Описание">
            <UInput v-model="form.description" placeholder="Описание" />
          </UFormField>
          <UFormField label="Вид комплекта" required>
            <USelect v-model="form.set_type_name" :items="setTypeOptions" :portal="false" placeholder="Вид комплекта" />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.equipment_set_id ? 'Сохранить' : 'Создать' }}</UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const isFormOpen = ref(false)
await Promise.all([crm.fetchSetTypes(), crm.fetchEquipmentSets(), crm.fetchEquipment()])

const form = reactive({
  equipment_set_id: null,
  equipment_set_name: '',
  description: '',
  set_type_name: ''
})

const setTypeOptions = computed(() => crm.setTypes.map((item) => item.set_type_name))

const equipmentBySetId = computed(() => {
  const grouped = {}

  for (const item of crm.equipment) {
    const id = item.equipment_set?.equipment_set_id || item.equipment_set_id
    if (!id) continue
    if (!grouped[id]) grouped[id] = []
    grouped[id].push(item)
  }

  return grouped
})

function equipmentInSet(setId) {
  return equipmentBySetId.value[setId] || []
}

function resetForm() {
  form.equipment_set_id = null
  form.equipment_set_name = ''
  form.description = ''
  form.set_type_name = ''
}

async function openCreate() {
  if (!crm.setTypes.length) {
    await crm.fetchSetTypes()
  }
  resetForm()
  isFormOpen.value = true
}

async function edit(item) {
  if (!crm.setTypes.length) {
    await crm.fetchSetTypes()
  }
  form.equipment_set_id = item.equipment_set_id
  form.equipment_set_name = item.equipment_set_name
  form.description = item.description || ''
  form.set_type_name = item.type?.set_type_name || ''
  isFormOpen.value = true
}

async function save() {
  if (!form.equipment_set_name.trim() || !form.set_type_name) return

  const payload = {
    equipment_set_name: form.equipment_set_name.trim(),
    description: form.description.trim(),
    set_type_name: form.set_type_name
  }

  if (form.equipment_set_id) {
    await crm.updateEquipmentSet(form.equipment_set_id, payload)
  } else {
    await crm.createEquipmentSet(payload)
  }

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteEquipmentSet(id)
  if (form.equipment_set_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
