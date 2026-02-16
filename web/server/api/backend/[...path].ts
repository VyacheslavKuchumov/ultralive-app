import { getQuery, readBody } from 'h3'
import { callBackend } from '~~/server/utils/backend'

const BODY_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])

function toBackendPath(input: string | string[] | undefined) {
  if (!input) return '/'
  if (Array.isArray(input)) return `/${input.join('/')}`
  return `/${input}`
}

export default defineEventHandler(async (event) => {
  const method = (event.method || 'GET').toUpperCase() as 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  const path = toBackendPath(event.context.params?.path)
  const body = BODY_METHODS.has(method) ? await readBody(event) : undefined

  const isPublicAuthRoute = method === 'POST' && (path === '/login' || path === '/register')

  return callBackend(event, method, path, {
    body,
    query: getQuery(event),
    requireAuth: !isPublicAuthRoute
  })
})
