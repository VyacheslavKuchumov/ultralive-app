<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Виды комплектов</h1>
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
              <th class="py-2 w-32">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.setTypes" :key="item.set_type_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.set_type_id }}</td>
              <td class="py-2">{{ item.set_type_name }}</td>
              <td class="py-2">
                <div class="flex gap-1 sm:gap-2">
                  <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                    <span class="hidden sm:inline">Изменить</span>
                  </UButton>
                  <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.set_type_id)">
                    <span class="hidden sm:inline">Удалить</span>
                  </UButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.set_type_id ? 'Редактировать вид комплекта' : 'Добавить вид комплекта'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название вида" required>
            <UInput v-model="form.set_type_name" placeholder="Название вида" required />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.set_type_id ? 'Сохранить' : 'Создать' }}</UButton>
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
  set_type_id: null,
  set_type_name: ''
})

await crm.fetchSetTypes()

function resetForm() {
  form.set_type_id = null
  form.set_type_name = ''
}

function openCreate() {
  resetForm()
  isFormOpen.value = true
}

function edit(item) {
  form.set_type_id = item.set_type_id
  form.set_type_name = item.set_type_name
  isFormOpen.value = true
}

async function save() {
  if (!form.set_type_name.trim()) return

  if (form.set_type_id) {
    await crm.updateSetType(form.set_type_id, { set_type_name: form.set_type_name.trim() })
  } else {
    await crm.createSetType({ set_type_name: form.set_type_name.trim() })
  }

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteSetType(id)
  if (form.set_type_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
