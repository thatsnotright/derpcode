export type NodeKind = "folder" | "file";

export type TreeNode = {
  id: string;
  name: string;
  kind: NodeKind;
  children: string[];
};

export type Nodes = Record<string, TreeNode>;

export const ROOT_ID = "root";

export const SEED_NODES: Nodes = {
  root: { id: "root", name: "my-project", kind: "folder", children: ["src", "readme"] },
  src: { id: "src", name: "src", kind: "folder", children: ["app", "styles"] },
  app: { id: "app", name: "App.tsx", kind: "file", children: [] },
  styles: { id: "styles", name: "styles.css", kind: "file", children: [] },
  readme: { id: "readme", name: "README.md", kind: "file", children: [] },
};
