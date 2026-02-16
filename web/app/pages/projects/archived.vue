<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Архив съёмок</h1>
      </template>

      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Съёмка</th>
              <th class="py-2">Площадка</th>
              <th class="py-2">Инженер</th>
              <th class="py-2">Период</th>
              <th class="py-2 w-48">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.archivedProjects" :key="item.project_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.project_id }}</td>
              <td class="py-2">{{ item.project_name }}</td>
              <td class="py-2">{{ item.type?.project_type_name || '-' }}</td>
              <td class="py-2">{{ item.chiefEngineer?.name || '-' }}</td>
              <td class="py-2">{{ item.shooting_start_date }} - {{ item.shooting_end_date }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" @click="restore(item)">Восстановить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.project_id)">Удалить</UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()

await crm.fetchArchivedProjects()

async function restore(item) {
  await crm.updateProject(item.project_id, {
    project_name: item.project_name,
    project_type_name: item.type?.project_type_name || '',
    chief_engineer_name: item.chiefEngineer?.name || '',
    shooting_start_date: item.shooting_start_date,
    shooting_end_date: item.shooting_end_date,
    archived: false
  })

  await crm.fetchArchivedProjects()
}

async function remove(id) {
  await crm.deleteProject(id)
  await crm.fetchArchivedProjects()
}
</script>
