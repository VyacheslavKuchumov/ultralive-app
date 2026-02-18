<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h1 class="text-xl font-semibold">Шаблоны</h1>
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
              <th class="py-2">Оборудование</th>
              <th class="py-2 w-48">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.drafts" :key="item.draft_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.draft_id }}</td>
              <td class="py-2">{{ item.draft_name }}</td>
              <td class="py-2">{{ item.equipment?.length || 0 }}</td>
              <td class="py-2">
                <div class="flex gap-1 sm:gap-2">
                  <UButton size="xs" color="primary" variant="soft" icon="i-lucide-list" :to="`/draft/${item.draft_id}`" aria-label="Состав">
                    <span class="hidden sm:inline">Состав</span>
                  </UButton>
                  <UButton size="xs" color="neutral" variant="soft" icon="i-lucide-pencil" aria-label="Изменить" @click="edit(item)">
                    <span class="hidden sm:inline">Изменить</span>
                  </UButton>
                  <UButton size="xs" color="error" variant="soft" icon="i-lucide-trash-2" aria-label="Удалить" @click="remove(item.draft_id)">
                    <span class="hidden sm:inline">Удалить</span>
                  </UButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <UModal v-model:open="isFormOpen" :title="form.draft_id ? 'Редактировать шаблон' : 'Добавить шаблон'">
      <template #body>
        <form class="space-y-3" @submit.prevent="save">
          <UFormField label="Название шаблона" required>
            <UInput v-model="form.draft_name" placeholder="Название шаблона" required />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="isFormOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-save">{{ form.draft_id ? 'Сохранить' : 'Создать' }}</UButton>
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
  draft_id: null,
  draft_name: ''
})

await crm.fetchDrafts()

function resetForm() {
  form.draft_id = null
  form.draft_name = ''
}

function openCreate() {
  resetForm()
  isFormOpen.value = true
}

function edit(item) {
  form.draft_id = item.draft_id
  form.draft_name = item.draft_name
  isFormOpen.value = true
}

async function save() {
  if (!form.draft_name.trim()) return

  if (form.draft_id) {
    await crm.updateDraft(form.draft_id, { draft_name: form.draft_name.trim() })
  } else {
    await crm.createDraft({ draft_name: form.draft_name.trim() })
  }

  resetForm()
  isFormOpen.value = false
}

async function remove(id) {
  await crm.deleteDraft(id)
  if (form.draft_id === id) {
    resetForm()
    isFormOpen.value = false
  }
}
</script>
