/* Minimal typed event bus shared by framework parts. */

type Handler = (payload?: unknown) => void

const handlers = new Map<string, Set<Handler>>()

export function onAdminEvent(event: string, handler: Handler): () => void {
  if (!handlers.has(event)) handlers.set(event, new Set())
  handlers.get(event)!.add(handler)
  return () => handlers.get(event)?.delete(handler)
}

export function emitAdminEvent(event: string, payload?: unknown): void {
  handlers.get(event)?.forEach(h => h(payload))
}
