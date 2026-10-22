import { useState } from "react";
import "./styles.css";

type NodeKind = "folder" | "file";

type TreeNode = {
  id: string;
  name: string;
  kind: NodeKind;
  children: string[]; // ids; always empty for files
};

type Nodes = Record<string, TreeNode>;

const ROOT_ID = "root";

const SEED_NODES: Nodes = {
  root: { id: "root", name: "my-project", kind: "folder", children: ["src", "readme"] },
  src: { id: "src", name: "src", kind: "folder", children: ["app", "styles"] },
  app: { id: "app", name: "App.tsx", kind: "file", children: [] },
  styles: { id: "styles", name: "styles.css", kind: "file", children: [] },
  readme: { id: "readme", name: "README.md", kind: "file", children: [] },
};

export default function App() {
  const [nodes] = useState<Nodes>(SEED_NODES);

  return (
    <main className="app">
      <h1>File explorer</h1>
      <p>Root folder: {nodes[ROOT_ID].name}</p>
    </main>
  );
}
