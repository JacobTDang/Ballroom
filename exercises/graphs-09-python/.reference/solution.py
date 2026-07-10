def find_order(num_courses: int, prerequisites: list[list[int]]) -> list[int]:
    adj: list[list[int]] = [[] for _ in range(num_courses)]
    for course, pre in prerequisites:
        adj[course].append(pre)

    UNVISITED, VISITING, VISITED = 0, 1, 2
    state = [UNVISITED] * num_courses
    order: list[int] = []

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
        order.append(course)
        return True

    for c in range(num_courses):
        if not dfs(c):
            return []
    return order
