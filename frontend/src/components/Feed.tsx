import "./Feed.css";
import { type Post } from "../types/Post";

type Props = {
  posts: Post[];
  onTagClick: (tag: string) => void;
  selectedTags: string[];
};

export default function Feed({ posts, onTagClick, selectedTags }: Props) {
  return (
    <div className="feed">
      {posts.map((post) => (
        <div className="post">
          <img src={post.imageUrl}></img>
          <div>
            <p>{post.text}</p>
          </div>
          <div>
            {post.tags.map((tag) => (
              <button
                className={`tag ${selectedTags.includes(tag) ? "active" : ""}`}
                onClick={() => onTagClick(tag)}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
