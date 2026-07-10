def can_finish(num_courses: int, prerequisites: list[list[int]]) -> bool:
    adj: list[list[int]] = [[] for _ in range(num_courses)]
    for course, pre in prerequisites:
        adj[course].append(pre)

    UNVISITED, VISITING, VISITED = 0, 1, 2
    state = [UNVISITED] * num_courses

    def dfs(course: int) -> bool:
        if state[course] == VISITING:
            return False
        if state[course] == VISITED:
            return True
        state[course] = VISITING
        for pre in adj[course]:
            if not dfs(pre):
                return False
        state[course] = VISITED
        return True

    for c in range(num_courses):
        if not dfs(c):
            return False
    return True
