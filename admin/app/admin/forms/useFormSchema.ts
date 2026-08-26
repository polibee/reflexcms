import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import type { SchemaNode } from '../core/types'
import { schemaToZod } from './schemaToZod'

interface UseFormSchemaOptions {
  schema: () => SchemaNode[]
  initialValues?: Record<string, unknown>
  onSubmit: (values: Record<string, unknown>) => Promise<void> | void
}

/**
 * Form engine: wires a declarative schema into vee-validate + zod.
 * Returns everything FormField components need for binding.
 */
export function useFormSchema(options: UseFormSchemaOptions) {
  const zodSchema = toTypedSchema(schemaToZod(options.schema()))

  const form = useForm<Record<string, unknown>>({
    validationSchema: zodSchema as never,
    initialValues: options.initialValues as never
  })

  const submitting = ref(false)

  const submit = form.handleSubmit(async (values) => {
    submitting.value = true
    try {
      await options.onSubmit(values)
    } finally {
      submitting.value = false
    }
  })

  return { form, submit, submitting }
}
