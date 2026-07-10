def daily_temperatures(temperatures: list[int]) -> list[int]:
    """Return, for each day, how many days until a warmer
    temperature, or 0 if there isn't one."""
    res = [0] * len(temperatures)
    stack: list[int] = []  # indices, decreasing temperature
    for i, temp in enumerate(temperatures):
        while stack and temperatures[stack[-1]] < temp:
            top = stack.pop()
            res[top] = i - top
        stack.append(i)
    return res
