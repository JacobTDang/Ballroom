from collections import defaultdict


def find_itinerary(tickets: list[list[str]]) -> list[str]:
    graph = defaultdict(list)
    for src, dst in tickets:
        graph[src].append(dst)
    for src in graph:
        graph[src].sort()

    route: list[str] = []

    def visit(airport: str) -> None:
        dests = graph[airport]
        while dests:
            next_airport = dests.pop(0)
            visit(next_airport)
        route.append(airport)

    visit("JFK")
    route.reverse()
    return route
