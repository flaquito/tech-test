import "./HeaderBar.css";

type Props = {
  onOpenPostCreationModal: () => void;
};

export default function HeaderBar({ onOpenPostCreationModal }: Props) {
  return (
    <div className="header-bar">
      <button onClick={onOpenPostCreationModal}>Create post</button>
      <header className="header-bar-title">
        <h2>Feed of Pretty Posts</h2>
      </header>
    </div>
  );
}
