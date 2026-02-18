<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Съёмки</h1>
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
              <th class="py-2">Съёмка</th>
              <th class="py-2">Площадка</th>
              <th class="py-2">Инженер</th>
              <th class="py-2">Период</th>
              <th class="py-2">Оборудование</th>
              <th class="py-2 w-48">Действия</th>
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
              <td class="py-2">
                <div class="flex gap-1 sm:gap-2">
                  <UButton size="xs" color="primary" variant="soft" icon="i-lucide-list" :to="`/project/${item.project_id}`" aria-label="Состав">
                    <span class="hidden sm:inline">Состав</span>
                  </UButton>
                  <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                    <span class="hidden sm:inline">Изменить</span>
                  </UButton>
                  <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.project_id)">
                    <span class="hidden sm:inline">Удалить</span>
                  </UButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.project_id ? 'Редактировать съёмку' : 'Добавить съёмку'">
      <template #body>
        <form class="grid gap-3 md:grid-cols-2" @submit.prevent="save">
          <UFormField label="Название съёмки" required class="md:col-span-2">
            <UInput v-model="form.project_name" placeholder="Название съёмки" required />
          </UFormField>
          <UFormField label="Площадка" required>
            <USelect v-model="form.project_type_name" :items="projectTypeOptions" :portal="false" placeholder="Площадка" />
          </UFormField>
          <UFormField label="Главный инженер" required>
            <USelect v-model="form.chief_engineer_name" :items="userOptions" :portal="false" placeholder="Главный инженер" />
          </UFormField>
          <UFormField label="Дата начала" required>
            <UInput v-model="form.shooting_start_date" placeholder="Дата начала YYYY-MM-DD" required />
          </UFormField>
          <UFormField label="Дата окончания" required>
            <UInput v-model="form.shooting_end_date" placeholder="Дата конца YYYY-MM-DD" required />
          </UFormField>
          <UCheckbox v-model="form.archived" label="В архив" class="md:col-span-2" />

          <div class="md:col-span-2 flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.project_id ? 'Сохранить' : 'Создать' }}</UButton>
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

function resetForm() {
  form.project_id = null
  form.project_name = ''
  form.project_type_name = ''
  form.chief_engineer_name = ''
  form.shooting_start_date = ''
  form.shooting_end_date = ''
  form.archived = false
}

async function openCreate() {
  if (!crm.projectTypes.length || !crm.users.length) {
    await Promise.all([crm.fetchProjectTypes(), crm.fetchUsers()])
  }
  resetForm()
  isFormOpen.value = true
}

async function edit(item) {
  if (!crm.projectTypes.length || !crm.users.length) {
    await Promise.all([crm.fetchProjectTypes(), crm.fetchUsers()])
  }
  form.project_id = item.project_id
  form.project_name = item.project_name
  form.project_type_name = item.type?.project_type_name || ''
  form.chief_engineer_name = item.chiefEngineer?.name || ''
  form.shooting_start_date = item.shooting_start_date
  form.shooting_end_date = item.shooting_end_date
  form.archived = !!item.archived
  isFormOpen.value = true
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
  if (!form.project_name.trim() || !form.project_type_name || !form.chief_engineer_name || !form.shooting_start_date || !form.shooting_end_date) {
    return
  }

  const payload = payloadFromForm()

  if (form.project_id) {
    await crm.updateProject(form.project_id, payload)
  } else {
    await crm.createProject(payload)
  }

  await crm.fetchProjects()
  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteProject(id)
  if (form.project_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
