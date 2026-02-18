<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Площадки</h1>
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
              <th class="py-2">Neaktor</th>
              <th class="py-2 w-32">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.projectTypes" :key="item.project_type_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.project_type_id }}</td>
              <td class="py-2">{{ item.project_type_name }}</td>
              <td class="py-2">{{ item.neaktor_id || '-' }}</td>
              <td class="py-2">
                <div class="flex gap-1 sm:gap-2">
                  <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                    <span class="hidden sm:inline">Изменить</span>
                  </UButton>
                  <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.project_type_id)">
                    <span class="hidden sm:inline">Удалить</span>
                  </UButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.project_type_id ? 'Редактировать площадку' : 'Добавить площадку'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название площадки" required>
            <UInput v-model="form.project_type_name" placeholder="Название площадки" required />
          </UFormField>
          <UFormField label="Neaktor ID">
            <UInput v-model="form.neaktor_id" placeholder="Neaktor ID (опционально)" />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.project_type_id ? 'Сохранить' : 'Создать' }}</UButton>
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
  project_type_id: null,
  project_type_name: '',
  neaktor_id: ''
})

await crm.fetchProjectTypes()

function resetForm() {
  form.project_type_id = null
  form.project_type_name = ''
  form.neaktor_id = ''
}

function openCreate() {
  resetForm()
  isFormOpen.value = true
}

function edit(item) {
  form.project_type_id = item.project_type_id
  form.project_type_name = item.project_type_name
  form.neaktor_id = item.neaktor_id || ''
  isFormOpen.value = true
}

async function save() {
  if (!form.project_type_name.trim()) return

  const payload = {
    project_type_name: form.project_type_name.trim(),
    neaktor_id: form.neaktor_id.trim()
  }

  if (form.project_type_id) {
    await crm.updateProjectType(form.project_type_id, payload)
  } else {
    await crm.createProjectType(payload)
  }

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteProjectType(id)
  if (form.project_type_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
