<script setup lang="ts">
definePageMeta({ layout: 'public' })

interface ArticleDetail {
  id: number
  title: string
  slug: string
  content: string
  summary: string
  published_at: string
  view_count: number
  comment_count: number
}

interface Comment {
  id: number
  user_id: number
  content: string
  created_at: string
}

const route = useRoute()
const article = ref<ArticleDetail | null>(null)
const comments = ref<Comment[]>([])
const loading = ref(true)
const commentText = ref('')
const submitting = ref(false)
const submitMsg = ref('')

/* comments are paginated; identity comes from the shared auth state */
const commentPage = ref(1)
const commentTotalPages = ref(1)

const { user: me, refresh: refreshAuth } = useAuthUser()

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const rawFetch = $fetch as unknown as (url: string, opts?: Record<string, any>) => Promise<unknown>

async function fetchComments(page = 1) {
  try {
    const res = await rawFetch(`/api/v1/articles/${route.params.slug}/comments?page=${page}&perPage=10`) as { items: Comment[], totalPages: number }
    comments.value = res.items ?? []
    commentTotalPages.value = res.totalPages ?? 1
  } catch { /* silent */ }
}

async function goToCommentPage(p: number) {
  commentPage.value = p
  await fetchComments(p)
}

async function fetchData() {
  loading.value = true
  try {
    article.value = await rawFetch(`/api/v1/articles/${route.params.slug}`) as ArticleDetail
    await fetchComments(1)
    await loadFavoriteState()
  } catch { /* silent */ }
  loading.value = false
}

/* favorite toggle (article type) */
const favorited = ref(false)

async function loadFavoriteState() {
  if (!me.value || !article.value) {
    favorited.value = false
    return
  }
  try {
    const res = await rawFetch(`/api/v1/favorites/state?type=article&id=${article.value.id}`) as { favorited: boolean }
    favorited.value = res.favorited
  } catch { /* silent */ }
}

async function toggleFavorite() {
  if (!me.value || !article.value) return
  try {
    const res = await rawFetch('/api/v1/favorites/toggle', {
      method: 'POST',
      body: { type: 'article', id: article.value.id }
    }) as { favorited: boolean }
    favorited.value = res.favorited
  } catch { /* silent */ }
}

async function submitComment() {
  if (!me.value || !commentText.value.trim()) return
  submitting.value = true
  submitMsg.value = ''
  try {
    const res = await rawFetch(`/api/v1/articles/${route.params.slug}/comments`, {
      method: 'POST',
      body: { content: commentText.value }
    }) as { message?: string }
    submitMsg.value = res.message ?? '评论已提交，等待审核'
    commentText.value = ''
  } catch (e: unknown) {
    const err = e as { statusCode?: number, data?: { statusMessage?: string, message?: string } }
    if (err.statusCode === 401) {
      void refreshAuth()
      submitMsg.value = ''
    } else {
      submitMsg.value = err.data?.statusMessage ?? err.data?.message ?? '提交失败'
    }
  } finally {
    submitting.value = false
  }
}

onMounted(fetchData)

useHead(() => ({
  title: article.value?.title ?? 'Article',
  meta: article.value?.summary
    ? [{ name: 'description', content: article.value.summary }]
    : []
}))
</script>

<template>
  <div>
    <div
      v-if="loading"
      class="py-20 text-center text-gray-400"
    >
      Loading…
    </div>
    <article
      v-else-if="article"
      class="mx-auto max-w-3xl"
    >
      <h1 class="text-3xl font-bold text-gray-900">
        {{ article.title }}
      </h1>
      <div class="mt-2 flex items-center gap-4 text-sm text-gray-400">
        <time>{{ article.published_at }}</time>
        <span>{{ article.view_count }} 阅读</span>
        <span>{{ comments.length }} 评论</span>
        <button
          v-if="me"
          type="button"
          class="ml-auto rounded-md px-3 py-1 text-xs font-medium"
          :class="favorited
            ? 'bg-amber-100 text-amber-700 hover:bg-amber-200'
            : 'border border-gray-300 text-gray-500 hover:bg-gray-50'"
          @click="toggleFavorite"
        >
          {{ favorited ? '★ 已收藏' : '☆ 收藏' }}
        </button>
      </div>
      <div class="prose prose-gray mt-8 max-w-none">
        <pre class="whitespace-pre-wrap font-sans text-base leading-relaxed">{{ article.content }}</pre>
      </div>

      <!-- Comments -->
      <section class="mt-12">
        <h2 class="mb-4 text-lg font-semibold">
          {{ comments.length }} 条评论
        </h2>

        <div class="space-y-4">
          <div
            v-for="c in comments"
            :key="c.id"
            class="rounded-lg border bg-white p-4"
          >
            <div class="flex items-center justify-between text-xs text-gray-400">
              <NuxtLink
                :to="`/users/${c.user_id}`"
                class="font-medium text-gray-600 hover:text-blue-600 hover:underline"
              >用户 #{{ c.user_id }}</NuxtLink>
              <time>{{ c.created_at?.slice(0, 10) }}</time>
            </div>
            <p class="mt-2 text-sm text-gray-700">
              {{ c.content }}
            </p>
          </div>

          <p
            v-if="!comments.length"
            class="py-8 text-center text-sm text-gray-400"
          >
            暂无评论
          </p>
        </div>

        <PublicPagination
          :page="commentPage"
          :total-pages="commentTotalPages"
          @change="goToCommentPage"
        />

        <!-- Comment form (below the list, login-aware) -->
        <div class="mt-6 rounded-lg border bg-white p-4">
          <h3 class="mb-3 text-sm font-semibold">
            发表评论
          </h3>
          <template v-if="me">
            <p class="mb-2 text-xs text-gray-400">
              以 <span class="font-medium text-gray-600">{{ me.name }}</span> 身份评论 · 提交后等待审核
            </p>
            <textarea
              v-model="commentText"
              :rows="3"
              placeholder="写下你的评论…"
              class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
            <div class="mt-2 flex items-center justify-between">
              <span
                v-if="submitMsg"
                class="text-xs text-green-600"
              >{{ submitMsg }}</span>
              <button
                :disabled="submitting || !commentText.trim()"
                class="rounded-md bg-blue-600 px-4 py-1.5 text-sm text-white hover:bg-blue-700 disabled:opacity-50"
                @click="submitComment"
              >
                {{ submitting ? '提交中…' : '提交评论' }}
              </button>
            </div>
          </template>
          <div
            v-else
            class="flex items-center justify-between rounded-md bg-gray-50 px-3 py-2.5"
          >
            <p class="text-sm text-gray-500">
              登录后即可参与评论
            </p>
            <NuxtLink
              :to="`/login?redirect=/articles/${article.slug}`"
              class="rounded-md bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            >登录</NuxtLink>
          </div>
        </div>
      </section>
    </article>
    <p
      v-else
      class="py-20 text-center text-gray-400"
    >
      文章不存在
    </p>
  </div>
</template>
