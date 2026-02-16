<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Комплекты оборудования</h1>
      </template>

      <form class="grid gap-3 md:grid-cols-4" @submit.prevent="save">
        <UInput v-model="form.equipment_set_name" placeholder="Название комплекта" />
        <UInput v-model="form.description" placeholder="Описание" />
        <USelect v-model="form.set_type_name" :items="setTypeOptions" placeholder="Вид комплекта" />
        <UButton type="submit" color="primary">{{ form.equipment_set_id ? 'Сохранить' : 'Добавить' }}</UButton>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Комплект</th>
              <th class="py-2">Тип</th>
              <th class="py-2">Описание</th>
              <th class="py-2 w-56">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.equipmentSets" :key="item.equipment_set_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.equipment_set_id }}</td>
              <td class="py-2">{{ item.equipment_set_name }}</td>
              <td class="py-2">{{ item.type?.set_type_name || '-' }}</td>
              <td class="py-2">{{ item.description || '-' }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" :to="`/equipment?set=${item.equipment_set_id}`">Оборудование</UButton>
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.equipment_set_id)">Удалить</UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
await Promise.all([crm.fetchSetTypes(), crm.fetchEquipmentSets()])

const form = reactive({
  equipment_set_id: null,
  equipment_set_name: '',
  description: '',
  set_type_name: ''
})

const setTypeOptions = computed(() => crm.setTypes.map((item) => item.set_type_name))

function edit(item) {
  form.equipment_set_id = item.equipment_set_id
  form.equipment_set_name = item.equipment_set_name
  form.description = item.description || ''
  form.set_type_name = item.type?.set_type_name || ''
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

  form.equipment_set_id = null
  form.equipment_set_name = ''
  form.description = ''
  form.set_type_name = ''
}

async function remove(id) {
  await crm.deleteEquipmentSet(id)
  if (form.equipment_set_id === id) {
    form.equipment_set_id = null
    form.equipment_set_name = ''
    form.description = ''
    form.set_type_name = ''
  }
}
</script>
