import { useEffect, useState, useMemo, useRef } from "react";
import { fetchPosts } from "./api/data";
import { type Post } from "./types/Post";
import "./App.css";
import Feed from "./components/Feed";
import HeaderBar from "./components/HeaderBar";
import PostCreationModal from "./components/PostCreationModal";

function App() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const loadingRef = useRef(false);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    loadMore();
  }, []);

  async function loadMore() {
    if (loadingRef.current || !hasMore) return;

    loadingRef.current = true;

    const nextPage = page;
    const newPosts = await fetchPosts(nextPage);

    setPosts((prev) => [...prev, ...newPosts]);
    setPage((prev) => prev + 1);

    if (newPosts.length < 10) {
      setHasMore(false);
    }

    loadingRef.current = false;
  }

  useEffect(() => {
    function handleScroll() {
      const scrollTop = window.scrollY;
      const windowHeight = window.innerHeight;
      const fullHeight = document.documentElement.scrollHeight;

      if (fullHeight <= windowHeight) return;
      const threshold = 100;
      const nearBottom = scrollTop + windowHeight >= fullHeight - threshold;

      if (nearBottom) {
        loadMore();
      }
    }

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [loadMore]);

  const filteredPosts = useMemo(() => {
    if (selectedTags.length === 0) return posts;

    return posts.filter((post) =>
      post.tags.some((tag) => selectedTags.includes(tag)),
    );
  }, [posts, selectedTags]);

  function toggleTag(tag: string) {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag],
    );
  }

  function openModal() {
    setIsModalOpen(true);
  }

  function closeModal() {
    setIsModalOpen(false);
  }

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onmessage = (event) => {
      const newPost: Post = JSON.parse(event.data);

      setPosts((prev) => {
        if (prev.some((p) => p.id === newPost.id)) return prev;
        return [newPost, ...prev];
      });
    };

    return () => socket.close();
  }, []);

  return (
    <>
      <HeaderBar onOpenPostCreationModal={openModal} />
      <Feed
        posts={filteredPosts}
        onTagClick={toggleTag}
        selectedTags={selectedTags}
      />
      {isModalOpen && <PostCreationModal onClose={closeModal} />}
    </>
  );
}

export default App;
