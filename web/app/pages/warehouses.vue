<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Склады</h1>
      </template>

      <form class="grid gap-3 sm:grid-cols-3" @submit.prevent="save">
        <UInput v-model="form.warehouse_name" placeholder="Название склада" />
        <UInput v-model="form.warehouse_adress" placeholder="Адрес" />
        <UButton type="submit" color="primary">{{ form.warehouse_id ? 'Сохранить' : 'Добавить' }}</UButton>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Склад</th>
              <th class="py-2">Адрес</th>
              <th class="py-2 w-40">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.warehouses" :key="item.warehouse_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.warehouse_id }}</td>
              <td class="py-2">{{ item.warehouse_name }}</td>
              <td class="py-2">{{ item.warehouse_adress || '-' }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.warehouse_id)">Удалить</UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()

const form = reactive({
  warehouse_id: null,
  warehouse_name: '',
  warehouse_adress: ''
})

await crm.fetchWarehouses()

function edit(item) {
  form.warehouse_id = item.warehouse_id
  form.warehouse_name = item.warehouse_name
  form.warehouse_adress = item.warehouse_adress || ''
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

  form.warehouse_id = null
  form.warehouse_name = ''
  form.warehouse_adress = ''
}

async function remove(id) {
  await crm.deleteWarehouse(id)
  if (form.warehouse_id === id) {
    form.warehouse_id = null
    form.warehouse_name = ''
    form.warehouse_adress = ''
  }
}
</script>
