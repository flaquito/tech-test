import { type Post } from "../types/Post";

const BASE_URL = "http://localhost:8080";

export async function fetchPosts(page: number): Promise<Post[]> {
  const res = await fetch(`${BASE_URL}/posts?page=${page}`);
  if (!res.ok) throw new Error("Failed to fetch posts");
  return res.json();
}

export async function createPost(
  file: File,
  text: string,
  tags: string,
): Promise<void> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("text", text);
  formData.append("tags", tags);

  const res = await fetch(`${BASE_URL}/uploads`, {
    method: "POST",
    body: formData,
  });

  if (!res.ok) {
    throw new Error("Failed to create post");
  }

  return;
}
