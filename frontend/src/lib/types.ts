export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
  api_key: string;
}
export interface Feed {
  id: string;
  name: string;
  url: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}
export interface FeedFollow {
  id: string;
  feed_id: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}
export interface Post {
  id: string;
  title: string;
  description: string | null;
  published_at: string;
  url: string;
  feed_id: string;
  created_at: string;
  updated_at: string;
}
export interface Page<T> {
  data: T[];
  next_cursor: string;
}
export interface AuthResponse {
  token: string;
  token_type: string;
}
export interface PostFilters {
  search: string;
  feed_id: string;
  unread: boolean;
}
