<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Шаблоны</h1>
      </template>

      <form class="grid gap-3 sm:grid-cols-[1fr_auto]" @submit.prevent="save">
        <UInput v-model="form.draft_name" placeholder="Название шаблона" />
        <UButton type="submit" color="primary">{{ form.draft_id ? 'Сохранить' : 'Добавить' }}</UButton>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2">Оборудование</th>
              <th class="py-2 w-56">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.drafts" :key="item.draft_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.draft_id }}</td>
              <td class="py-2">{{ item.draft_name }}</td>
              <td class="py-2">{{ item.equipment?.length || 0 }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="primary" variant="soft" :to="`/draft/${item.draft_id}`">Состав</UButton>
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.draft_id)">Удалить</UButton>
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
  draft_id: null,
  draft_name: ''
})

await crm.fetchDrafts()

function edit(item) {
  form.draft_id = item.draft_id
  form.draft_name = item.draft_name
}

async function save() {
  if (!form.draft_name.trim()) return

  if (form.draft_id) {
    await crm.updateDraft(form.draft_id, { draft_name: form.draft_name.trim() })
  } else {
    await crm.createDraft({ draft_name: form.draft_name.trim() })
  }

  form.draft_id = null
  form.draft_name = ''
}

async function remove(id) {
  await crm.deleteDraft(id)
  if (form.draft_id === id) {
    form.draft_id = null
    form.draft_name = ''
  }
}
</script>
