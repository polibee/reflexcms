export function usePublicApi() {
  async function get<T>(path: string, query?: Record<string, string | number>): Promise<T> {
    return await $fetch<T>(`/api/v1/${path}`, { query })
  }
  return { get }
}
