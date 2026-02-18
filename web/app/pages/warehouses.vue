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
import { reactive, ref } from 'vue'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const isFormOpen = ref(false)

const form = reactive({
  warehouse_id: null,
  warehouse_name: '',
  warehouse_adress: ''
})

await crm.fetchWarehouses()

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

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteWarehouse(id)
  if (form.warehouse_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
