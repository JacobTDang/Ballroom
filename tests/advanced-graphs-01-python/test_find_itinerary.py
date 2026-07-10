from solution import find_itinerary


def test_dead_end():
    tickets = [["JFK", "SFO"], ["JFK", "ATL"], ["SFO", "ATL"], ["ATL", "JFK"], ["ATL", "SFO"]]
    assert find_itinerary(tickets) == ["JFK", "ATL", "JFK", "SFO", "ATL", "SFO"]


def test_simple():
    tickets = [["MUC", "LHR"], ["JFK", "MUC"], ["SFO", "SJC"], ["LHR", "SFO"]]
    assert find_itinerary(tickets) == ["JFK", "MUC", "LHR", "SFO", "SJC"]


def test_lexical_tie_break():
    tickets = [["JFK", "KUL"], ["JFK", "NRT"], ["NRT", "JFK"]]
    assert find_itinerary(tickets) == ["JFK", "NRT", "JFK", "KUL"]
