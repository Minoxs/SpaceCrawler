import os
from pathlib import Path
from stat import S_ISDIR

from rich.style import Style
from rich.text import Text, TextType
from textual.app import App, ComposeResult
from textual.widgets import Input, Header, Footer, DirectoryTree, Tree
from textual.widgets._directory_tree import DirEntry
from textual.widgets._tree import TreeNode


def GetSizeDir(path: Path | str, recursive: bool) -> int:
    size = 0

    for file in os.listdir(path):
        size += GetSize(Path(path) / file, recursive)

    return size


def GetSize(path: Path | str, recursive: bool) -> int:
    info = os.stat(path)
    if S_ISDIR(info.st_mode):
        if recursive:
            return GetSizeDir(path, recursive)
        else:
            return 0
    else:
        return info.st_size


def GetSizeH(path: Path | str, recursive: bool = True) -> str:
    size = GetSize(path, recursive)
    if size <= 1024:
        return f"{size} bytes"

    mag = ["bytes", "Kb", "Mb", "Gb"]
    idx = 0
    while size > 1024 and idx < len(mag):
        size /= 1024
        idx += 1

    return f"{size:.02f} {mag[idx]}"


class DiskStats(DirectoryTree):
    def __init__(self, root: str):
        super().__init__(root)
        self.recursive = False

    def render_label(
            self, node: TreeNode[DirEntry], base_style: Style, style: Style
    ) -> Text:
        if not hasattr(node.data, "size"):
            node.data.size = GetSizeH(node.data.path, self.recursive)

        text = super().render_label(node, base_style, style)

        if node.is_root:
            text += os.path.abspath(node.data.path)

        text += " " + node.data.size

        return text

    # async def _on_tree_node_expanded(self, event: Tree.NodeExpanded[DirEntry]) -> None:
    #     await super()._on_tree_node_expanded(event)
    #     if hasattr(event.node.data, "size") and event.node.data.size == "0 bytes":
    #         event.node.data.size = GetSizeH(event.node.data.path, True)



class DiskUsage(App):
    def compose(self) -> ComposeResult:
        yield Header()
        yield DiskStats(".")
        yield Footer()


if __name__ == '__main__':
    DiskUsage().run()
