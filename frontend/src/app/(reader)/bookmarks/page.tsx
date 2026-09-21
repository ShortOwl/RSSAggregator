import { PostList } from "@/components/posts/post-list";
export const metadata = { title: "Bookmarks" };
export default function Bookmarks() {
  return <PostList mode="bookmarks" />;
}
