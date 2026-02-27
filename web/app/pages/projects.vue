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
      <div class="space-y-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Поиск по съемкам"
            class="md:max-w-sm"
          />

          <label class="flex items-center gap-2 text-sm text-gray-600">
            На странице
            <select v-model.number="perPage" class="rounded border border-gray-300 px-2 py-1 text-sm">
              <option v-for="option in perPageOptions" :key="option" :value="option">{{ option }}</option>
            </select>
          </label>
        </div>

        <div v-if="crm.projects.length" class="grid gap-4 md:grid-cols-2">
          <UCard
            v-for="item in crm.projects"
            :key="item.project_id"
            :data-testid="`project-card-${item.project_id}`"
          >
            <template #header>
              <div class="flex items-start justify-between gap-3">
                <div class="space-y-1">
                  <p class="text-xs text-gray-500">#{{ item.project_id }}</p>
                  <h2 class="text-lg font-semibold">{{ item.project_name }}</h2>
                </div>

                <div class="flex gap-1 sm:gap-2 shrink-0">
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
              </div>
            </template>

            <div class="space-y-1 text-sm text-gray-700">
              <p><span class="text-gray-500">Площадка:</span> {{ item.type?.project_type_name || '-' }}</p>
              <p><span class="text-gray-500">Инженер:</span> {{ item.chiefEngineer?.name || '-' }}</p>
              <p><span class="text-gray-500">Период:</span> {{ item.shooting_start_date }} - {{ item.shooting_end_date }}</p>
              <p><span class="text-gray-500">Оборудование:</span> {{ item.equipment?.length || 0 }}</p>
            </div>
          </UCard>
        </div>

        <p v-else-if="!isLoading" class="text-sm text-gray-600">Ничего не найдено.</p>

        <div class="flex flex-col gap-3 border-t border-gray-100 pt-3 md:flex-row md:items-center md:justify-between">
          <p class="text-sm text-gray-600">Показано {{ from }}-{{ to }} из {{ pagination.total }}</p>

          <div class="flex items-center gap-2">
            <UButton size="xs" color="neutral" variant="soft" :disabled="page <= 1 || isLoading" @click="prevPage">Назад</UButton>
            <span class="text-sm text-gray-600">Стр. {{ page }} / {{ pagination.total_pages }}</span>
            <UButton
              size="xs"
              color="neutral"
              variant="soft"
              :disabled="page >= pagination.total_pages || isLoading"
              @click="nextPage"
            >
              Вперед
            </UButton>
          </div>
        </div>
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

          <UFormField label="Период съёмки (DateRangePicker)" required class="md:col-span-2">
            <UPopover>
              <UButton color="neutral" variant="soft" class="w-full justify-start">
                {{ shootingDateRangeLabel }}
              </UButton>

              <template #content>
                <UCalendar v-model="shootingDateRangeModel" range />
              </template>
            </UPopover>
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
import { CalendarDate } from '@internationalized/date'
import { computed, reactive, ref } from 'vue'
import { useServerList } from '~/composables/useServerList'
import { useCRMStore } from '~/stores/crm'

const crm = useCRMStore()
const isFormOpen = ref(false)
const perPageOptions = [10, 20, 50]

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
  crm.fetchProjectTypes({ page: 1, per_page: 1000 })
])

const {
  search,
  page,
  perPage,
  isLoading,
  pagination,
  from,
  to,
  load,
  prevPage,
  nextPage
} = useServerList(
  (params) => crm.fetchProjects(params),
  computed(() => crm.pagination.projects),
  { perPage: 10 }
)

const projectTypeOptions = computed(() => crm.projectTypes.map((item) => item.project_type_name))
const userOptions = computed(() => crm.users.map((item) => item.name))

function parseDateValue(raw) {
  if (!raw || typeof raw !== 'string') return null
  const [year, month, day] = raw.split('-').map((item) => Number(item))
  if (!year || !month || !day) return null
  return new CalendarDate(year, month, day)
}

function formatDateValue(value) {
  if (!value || typeof value.year !== 'number') return ''
  const year = String(value.year).padStart(4, '0')
  const month = String(value.month).padStart(2, '0')
  const day = String(value.day).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const shootingDateRangeModel = computed({
  get() {
    const start = parseDateValue(form.shooting_start_date)
    const end = parseDateValue(form.shooting_end_date)

    if (!start || !end) return null

    return { start, end }
  },
  set(value) {
    form.shooting_start_date = formatDateValue(value?.start)
    form.shooting_end_date = formatDateValue(value?.end)
  }
})

const shootingDateRangeLabel = computed(() => {
  if (!form.shooting_start_date || !form.shooting_end_date) {
    return 'Выбрать период'
  }

  return `${form.shooting_start_date} - ${form.shooting_end_date}`
})

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
    await Promise.all([
      crm.fetchProjectTypes({ page: 1, per_page: 1000 }),
      crm.fetchUsers()
    ])
  }
  resetForm()
  isFormOpen.value = true
}

async function edit(item) {
  if (!crm.projectTypes.length || !crm.users.length) {
    await Promise.all([
      crm.fetchProjectTypes({ page: 1, per_page: 1000 }),
      crm.fetchUsers()
    ])
  }
  form.project_id = item.project_id
  form.project_name = item.project_name
  form.project_type_name = item.type?.project_type_name || ''
  form.chief_engineer_name = item.chiefEngineer?.name || ''
  form.shooting_start_date = item.shooting_start_date || ''
  form.shooting_end_date = item.shooting_end_date || ''
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
  const startDate = parseDateValue(form.shooting_start_date.trim())
  const endDate = parseDateValue(form.shooting_end_date.trim())

  if (
    !form.project_name.trim() ||
    !form.project_type_name ||
    !form.chief_engineer_name ||
    !startDate ||
    !endDate
  ) {
    return
  }

  const payload = payloadFromForm()

  if (form.project_id) {
    await crm.updateProject(form.project_id, payload)
  } else {
    await crm.createProject(payload)
  }

  await load()
  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteProject(id)
  await load()

  if (form.project_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
