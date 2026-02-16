<template>
  <div class="space-y-6">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Виды комплектов</h1>
      </template>

      <form class="grid gap-3 sm:grid-cols-[1fr_auto]" @submit.prevent="save">
        <UInput v-model="form.set_type_name" placeholder="Название вида" />
        <UButton type="submit" color="primary">{{ form.set_type_id ? 'Сохранить' : 'Добавить' }}</UButton>
      </form>
    </UCard>

    <UCard>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2 w-40">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in crm.setTypes" :key="item.set_type_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.set_type_id }}</td>
              <td class="py-2">{{ item.set_type_name }}</td>
              <td class="py-2 flex gap-2">
                <UButton size="xs" color="neutral" variant="soft" @click="edit(item)">Изменить</UButton>
                <UButton size="xs" color="error" variant="soft" @click="remove(item.set_type_id)">Удалить</UButton>
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
  set_type_id: null,
  set_type_name: ''
})

await crm.fetchSetTypes()

function edit(item) {
  form.set_type_id = item.set_type_id
  form.set_type_name = item.set_type_name
}

async function save() {
  if (!form.set_type_name.trim()) return

  if (form.set_type_id) {
    await crm.updateSetType(form.set_type_id, { set_type_name: form.set_type_name.trim() })
  } else {
    await crm.createSetType({ set_type_name: form.set_type_name.trim() })
  }

  form.set_type_id = null
  form.set_type_name = ''
}

async function remove(id) {
  await crm.deleteSetType(id)
  if (form.set_type_id === id) {
    form.set_type_id = null
    form.set_type_name = ''
  }
}
</script>
