<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Съёмки</h1>
      </template>

      <form class="grid gap-3 md:grid-cols-3" @submit.prevent="save">
        <UInput v-model="form.project_name" placeholder="Название съёмки" />
        <USelect v-model="form.project_type_name" :items="projectTypeOptions" placeholder="Площадка" />
        <USelect v-model="form.chief_engineer_name" :items="userOptions" placeholder="Главный инженер" />

        <UInput v-model="form.shooting_start_date" placeholder="Дата начала YYYY-MM-DD" />
        <UInput v-model="form.shooting_end_date" placeholder="Дата конца YYYY-MM-DD" />
        <UCheckbox v-model="form.archived" label="В архив" />

        <div class="md:col-span-3">
          <UButton type="submit" color="primary">{{ form.project_id ? 'Сохранить' : 'Добавить' }}</UButton>
        </div>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Съёмка</th>
              <th class="py-2">Площадка</th>
              <th class="py-2">Инженер</th>
              <th class="py-2">Период</th>
              <th class="py-2">Оборудование</th>
              <th class="py-2 w-56">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.projects" :key="item.project_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.project_id }}</td>
              <td class="py-2">{{ item.project_name }}</td>
              <td class="py-2">{{ item.type?.project_type_name || '-' }}</td>
              <td class="py-2">{{ item.chiefEngineer?.name || '-' }}</td>
              <td class="py-2">{{ item.shooting_start_date }} - {{ item.shooting_end_date }}</td>
              <td class="py-2">{{ item.equipment?.length || 0 }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="primary" variant="soft" :to="`/project/${item.project_id}`">Состав</UButton>
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
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
import { computed, reactive } from 'vue'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()

const form = reactive({
  project_id: null,
  project_name: '',
  project_type_name: '',
  chief_engineer_name: '',
  shooting_start_date: '',
  shooting_end_date: '',
  archived: false
})

await Promise.all([
  crm.fetchUsers(),
  crm.fetchProjectTypes(),
  crm.fetchProjects()
])

const projectTypeOptions = computed(() => crm.projectTypes.map((item) => item.project_type_name))
const userOptions = computed(() => crm.users.map((item) => item.name))

function edit(item) {
  form.project_id = item.project_id
  form.project_name = item.project_name
  form.project_type_name = item.type?.project_type_name || ''
  form.chief_engineer_name = item.chiefEngineer?.name || ''
  form.shooting_start_date = item.shooting_start_date
  form.shooting_end_date = item.shooting_end_date
  form.archived = !!item.archived
}

function payloadFromForm() {
  return {
    project_name: form.project_name.trim(),
    project_type_name: form.project_type_name,
    chief_engineer_name: form.chief_engineer_name,
    shooting_start_date: form.shooting_start_date.trim(),
    shooting_end_date: form.shooting_end_date.trim(),
    archived: !!form.archived
  }
}

async function save() {
  if (!form.project_name || !form.project_type_name || !form.chief_engineer_name || !form.shooting_start_date || !form.shooting_end_date) {
    return
  }

  const payload = payloadFromForm()

  if (form.project_id) {
    await crm.updateProject(form.project_id, payload)
  } else {
    await crm.createProject(payload)
  }

  await crm.fetchProjects()

  form.project_id = null
  form.project_name = ''
  form.project_type_name = ''
  form.chief_engineer_name = ''
  form.shooting_start_date = ''
  form.shooting_end_date = ''
  form.archived = false
}

async function remove(id) {
  await crm.deleteProject(id)
}
</script>
