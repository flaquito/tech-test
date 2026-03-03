import "./PostCreationModal.css";
import { createPost } from "../api/data";
import { useState } from "react";

type Props = {
  onClose: () => void;
};

export default function PostCreationModal({ onClose }: Props) {
  const [file, setFile] = useState<File | null>(null);
  const [text, setText] = useState("");
  const [tags, setTags] = useState("");
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  const isValid = file != null && text.trim().length > 0;

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!file) return;

    await createPost(file, text, tags);
    onClose();
  };

  return (
    <div className="overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>Create Post</h2>
        <form onSubmit={handleSubmit} className="form">
          <input
            type="file"
            accept="image/*"
            onChange={(e) => {
              const selected = e.target.files?.[0] ?? null;
              setFile(selected);
              if (selected) {
                setPreviewUrl(URL.createObjectURL(selected));
              } else {
                setPreviewUrl(null);
              }
            }}
            required
          />
          {previewUrl && (
            <img src={previewUrl} alt="Preview" className="preview" />
          )}
          <textarea
            placeholder="Your post text"
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
          <input
            placeholder="Your post tags (comma separated)"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
          />
          <div className="actions">
            <button type="button" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" disabled={!isValid}>
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
