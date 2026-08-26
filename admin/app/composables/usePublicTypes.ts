export interface PublicArticle {
  id: number
  title: string
  slug: string
  summary: string
  cover: string
  is_pinned: boolean
  published_at: string
  view_count: number
  comment_count: number
}

export interface PublicArticleDetail extends PublicArticle {
  content: string
  created_at: string
}

export interface PublicTopic {
  id: number
  title: string
  forum_category_id: number
  is_pinned: boolean
  is_featured: boolean
  reply_count: number
  like_count: number
  view_count: number
  last_reply_at: string
}

export interface PublicTopicDetail extends PublicTopic {
  content: string
  status: string
  best_reply_id: number
  replies: Array<{
    id: number
    user_id: number
    content: string
    floor: number
    like_count: number
    created_at: string
  }>
}

export interface PublicForum {
  id: number
  name: string
  slug: string
  description: string
  icon: string
  topic_count: number
}

export interface PaginatedResult<T> {
  items: T[]
  total: number
  page: number
  perPage: number
  totalPages: number
}
