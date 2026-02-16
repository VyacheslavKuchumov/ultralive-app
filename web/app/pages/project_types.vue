<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Площадки</h1>
      </template>

      <form class="grid gap-3 sm:grid-cols-3" @submit.prevent="save">
        <UInput v-model="form.project_type_name" placeholder="Название площадки" />
        <UInput v-model="form.neaktor_id" placeholder="Neaktor ID (опционально)" />
        <UButton type="submit" color="primary">{{ form.project_type_id ? 'Сохранить' : 'Добавить' }}</UButton>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2">Neaktor</th>
              <th class="py-2 w-40">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.projectTypes" :key="item.project_type_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.project_type_id }}</td>
              <td class="py-2">{{ item.project_type_name }}</td>
              <td class="py-2">{{ item.neaktor_id || '-' }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.project_type_id)">Удалить</UButton>
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
  project_type_id: null,
  project_type_name: '',
  neaktor_id: ''
})

await crm.fetchProjectTypes()

function edit(item) {
  form.project_type_id = item.project_type_id
  form.project_type_name = item.project_type_name
  form.neaktor_id = item.neaktor_id || ''
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

  form.project_type_id = null
  form.project_type_name = ''
  form.neaktor_id = ''
}

async function remove(id) {
  await crm.deleteProjectType(id)
  if (form.project_type_id === id) {
    form.project_type_id = null
    form.project_type_name = ''
    form.neaktor_id = ''
  }
}
</script>
