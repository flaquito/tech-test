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
      {posts.map((p) => (
        <div className="post">
          <img src={p.imageUrl}></img>
          <div>
            <p>{p.text}</p>
          </div>
          <div>
            {p.tags.map((tag) => (
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
