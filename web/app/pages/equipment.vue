<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Оборудование</h1>
      </template>

      <form class="grid gap-3 md:grid-cols-3" @submit.prevent="save">
        <UInput v-model="form.equipment_name" placeholder="Название" />
        <UInput v-model="form.serial_number" placeholder="Серийный номер" />
        <USelect v-model="form.equipment_set_name" :items="equipmentSetOptions" placeholder="Комплект" />

        <UInput v-model="form.description" placeholder="Описание" />
        <USelect v-model="form.warehouse_name" :items="warehouseOptions" placeholder="Склад" />
        <UInput v-model="form.current_storage_name" placeholder="Текущее хранение" />

        <UInput v-model="form.date_of_purchase" placeholder="Дата покупки YYYY-MM-DD" />
        <UInput v-model="form.cost_of_purchase" placeholder="Стоимость" />
        <UCheckbox v-model="form.needs_maintenance" label="Требует обслуживания" />

        <div class="md:col-span-3">
          <UButton type="submit" color="primary">{{ form.equipment_id ? 'Сохранить' : 'Добавить' }}</UButton>
        </div>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2">Серия</th>
              <th class="py-2">Комплект</th>
              <th class="py-2">Склад</th>
              <th class="py-2">ТО</th>
              <th class="py-2 w-40">Действия</th>
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
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.equipment_id)">Удалить</UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
import { computed, reactive, watch } from 'vue'
import { useRoute } from '#imports'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const route = useRoute()

const form = reactive({
  equipment_id: null,
  equipment_name: '',
  serial_number: '',
  equipment_set_name: '',
  description: '',
  warehouse_name: '',
  current_storage_name: '',
  needs_maintenance: false,
  date_of_purchase: '',
  cost_of_purchase: ''
})

const equipmentSetOptions = computed(() => crm.equipmentSets.map((item) => item.equipment_set_name))
const warehouseOptions = computed(() => crm.warehouses.map((item) => item.warehouse_name))

async function load() {
  await Promise.all([crm.fetchEquipmentSets(), crm.fetchWarehouses()])
  const setId = Number(route.query.set || 0)
  await crm.fetchEquipment(Number.isFinite(setId) && setId > 0 ? setId : null)

  if (setId > 0) {
    const set = crm.equipmentSets.find((item) => item.equipment_set_id === setId)
    if (set) {
      form.equipment_set_name = set.equipment_set_name
    }
  }
}

await load()

watch(() => route.query.set, load)

function edit(item) {
  form.equipment_id = item.equipment_id
  form.equipment_name = item.equipment_name
  form.serial_number = item.serial_number
  form.equipment_set_name = item.equipment_set?.equipment_set_name || ''
  form.description = item.description || ''
  form.warehouse_name = item.storage?.warehouse_name || ''
  form.current_storage_name = item.current_storage || ''
  form.needs_maintenance = item.needs_maintenance
  form.date_of_purchase = item.date_of_purchase || ''
  form.cost_of_purchase = item.cost_of_purchase || ''
}

function payloadFromForm() {
  return {
    equipment_name: form.equipment_name.trim(),
    serial_number: form.serial_number.trim(),
    equipment_set_name: form.equipment_set_name,
    description: form.description.trim(),
    warehouse_name: form.warehouse_name,
    current_storage_name: form.current_storage_name.trim(),
    needs_maintenance: !!form.needs_maintenance,
    date_of_purchase: form.date_of_purchase.trim(),
    cost_of_purchase: form.cost_of_purchase ? Number(form.cost_of_purchase) : null
  }
}

async function save() {
  if (!form.equipment_name || !form.serial_number || !form.equipment_set_name || !form.warehouse_name) {
    return
  }

  const payload = payloadFromForm()

  if (form.equipment_id) {
    await crm.updateEquipment(form.equipment_id, payload)
  } else {
    await crm.createEquipment(payload)
  }

  const currentSet = Number(route.query.set || 0)
  await crm.fetchEquipment(currentSet > 0 ? currentSet : null)

  form.equipment_id = null
  form.equipment_name = ''
  form.serial_number = ''
  form.description = ''
  form.current_storage_name = ''
  form.needs_maintenance = false
  form.date_of_purchase = ''
  form.cost_of_purchase = ''
}

async function remove(id) {
  await crm.deleteEquipment(id)
}
</script>
