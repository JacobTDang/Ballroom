def find_words(board: list[list[str]], words: list[str]) -> list[str]:
    """Return every word from words that can be traced out on board
    via sequentially adjacent cells, each cell used at most once per
    word."""
    root: dict = {}
    for w in words:
        node = root
        for c in w:
            node = node.setdefault(c, {})
        node["$"] = w

    rows, cols = len(board), len(board[0])
    res: list[str] = []

    def dfs(r: int, c: int, node: dict) -> None:
        if r < 0 or r >= rows or c < 0 or c >= cols:
            return
        ch = board[r][c]
        if ch == "#" or ch not in node:
            return
        nxt = node[ch]
        if "$" in nxt:
            res.append(nxt["$"])
            del nxt["$"]  # don't report the same word twice
        board[r][c] = "#"
        dfs(r + 1, c, nxt)
        dfs(r - 1, c, nxt)
        dfs(r, c + 1, nxt)
        dfs(r, c - 1, nxt)
        board[r][c] = ch

    for r in range(rows):
        for c in range(cols):
            dfs(r, c, root)
    return res
