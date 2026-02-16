<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">Вход в Ultralive CRM</h1>
        <p class="text-sm text-gray-600">Авторизуйтесь для работы с оборудованием и съёмками.</p>
      </div>
    </template>

    <UForm :state="state" class="space-y-4" @submit="onSubmit">
      <UFormField label="Email" required>
        <UInput v-model="state.email" type="email" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Пароль" required>
        <UInput v-model="state.password" type="password" class="w-full" placeholder="Ваш пароль" />
      </UFormField>

      <p v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</p>

      <UButton type="submit" color="primary" block :loading="loading">
        Войти
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-gray-600">
        <span>Нет аккаунта?</span>
        <ULink to="/signup">Зарегистрироваться</ULink>
      </div>
    </template>
  </UCard>
</template>

<script setup>
const auth = useAuthStore()
const toast = useToast()
const loading = ref(false)
const errorMessage = ref('')
const state = reactive({
  email: '',
  password: ''
})

async function onSubmit() {
  loading.value = true
  errorMessage.value = ''

  try {
    await auth.login(state)
    await navigateTo('/')
  } catch (error) {
    errorMessage.value = error?.data?.statusMessage || error?.data?.message || error?.message || 'Ошибка авторизации'
    toast.add({
      title: 'Вход не выполнен',
      description: errorMessage.value,
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
