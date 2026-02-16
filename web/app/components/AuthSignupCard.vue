<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">Регистрация в Ultralive CRM</h1>
        <p class="text-sm text-gray-600">Создайте аккаунт для доступа к CRM.</p>
      </div>
    </template>

    <UForm :state="state" class="space-y-4" @submit="onSubmit">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <UFormField label="Имя" required>
          <UInput v-model="state.firstName" class="w-full" placeholder="Иван" />
        </UFormField>

        <UFormField label="Фамилия" required>
          <UInput v-model="state.lastName" class="w-full" placeholder="Иванов" />
        </UFormField>
      </div>

      <UFormField label="Email" required>
        <UInput v-model="state.email" type="email" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Пароль" required>
        <UInput v-model="state.password" type="password" class="w-full" placeholder="Минимум 3 символа" />
      </UFormField>

      <p v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</p>

      <UButton type="submit" color="primary" block :loading="loading">
        Зарегистрироваться
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-gray-600">
        <span>Уже есть аккаунт?</span>
        <ULink to="/login">Войти</ULink>
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
  firstName: '',
  lastName: '',
  email: '',
  password: ''
})

async function onSubmit() {
  loading.value = true
  errorMessage.value = ''

  try {
    await auth.signup(state)
    await navigateTo('/')
  } catch (error) {
    errorMessage.value = error?.data?.statusMessage || error?.data?.message || error?.message || 'Ошибка регистрации'
    toast.add({
      title: 'Регистрация не выполнена',
      description: errorMessage.value,
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
