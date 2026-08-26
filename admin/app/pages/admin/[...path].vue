<script setup lang="ts">
import type { Component } from 'vue'
import ResourceListPage from '~/admin/framework/ResourceListPage.vue'
import ResourceFormPage from '~/admin/framework/ResourceFormPage.vue'
import ResourceViewPage from '~/admin/framework/ResourceViewPage.vue'

/* =============================================================
 * Resource router:
 *   /admin/{resource}              -> List
 *   /admin/{resource}/create       -> Create
 *   /admin/{resource}/{id}         -> View
 *   /admin/{resource}/{id}/edit    -> Edit
 * ============================================================= */

const route = useRoute()
const allow = useCan()

const raw = route.params.path
const segments = (Array.isArray(raw) ? raw : [raw]).map(String)

const resource = getResource(segments[0] ?? '')
if (!resource) {
  throw createError({ statusCode: 404, statusMessage: `Unknown resource "${segments[0]}"`, fatal: true })
}

if (!allow(`${resource.permissionPrefix}.view`)) {
  throw createError({ statusCode: 403, statusMessage: 'You do not have permission to view this resource.', fatal: true })
}

const id = segments[1]
const action = segments[2]

/* schema presence checks: pages cannot render without their schemas */
if ((id === 'create' || action === 'edit') && !resource.form) {
  throw createError({ statusCode: 400, statusMessage: `Resource "${resource.name}" does not define a form schema.`, fatal: true })
}
if (!id && !resource.table) {
  throw createError({ statusCode: 400, statusMessage: `Resource "${resource.name}" does not define a table schema.`, fatal: true })
}

let component: Component | null = null
const bind = reactive<Record<string, unknown>>({ resource })

if (!id) {
  component = ResourceListPage
} else if (id === 'create') {
  if (!allow(`${resource.permissionPrefix}.create`)) {
    throw createError({ statusCode: 403, statusMessage: 'Create permission required.', fatal: true })
  }
  component = ResourceFormPage
  bind.mode = 'create'
} else if (action === 'edit') {
  if (!allow(`${resource.permissionPrefix}.edit`)) {
    throw createError({ statusCode: 403, statusMessage: 'Edit permission required.', fatal: true })
  }
  component = ResourceFormPage
  bind.mode = 'edit'
  bind.id = id
} else {
  if (!allow(`${resource.permissionPrefix}.view`)) {
    throw createError({ statusCode: 403, statusMessage: 'View permission required.', fatal: true })
  }
  component = ResourceViewPage
  bind.id = id
}

useHead({ title: resource.labelPlural })
</script>

<template>
  <component
    :is="component"
    v-bind="bind"
    :key="route.fullPath"
  />
</template>
