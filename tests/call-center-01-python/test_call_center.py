from solution import CallCenter


def test_first_call_goes_to_a_respondent():
    c = CallCenter(1, 1, 1)
    assert c.dispatch(1) == "respondent"


def test_escalates_to_manager_when_respondents_busy():
    c = CallCenter(1, 1, 1)
    c.dispatch(1)
    assert c.dispatch(2) == "manager"


def test_escalates_to_director_when_managers_busy():
    c = CallCenter(1, 1, 1)
    c.dispatch(1)
    c.dispatch(2)
    assert c.dispatch(3) == "director"


def test_queues_when_everyone_is_busy():
    c = CallCenter(1, 0, 0)
    c.dispatch(1)
    assert c.dispatch(2) == "queued"
    assert c.handler_of(2) == "queued"


def test_ending_a_call_assigns_the_queued_call():
    c = CallCenter(1, 1, 0)
    c.dispatch(1)
    c.dispatch(2)
    c.dispatch(3)
    assert c.end_call(1) is True
    assert c.handler_of(3) == "respondent"


def test_queued_calls_are_assigned_fifo():
    c = CallCenter(2, 0, 0)
    c.dispatch(1)
    c.dispatch(2)
    c.dispatch(3)
    c.dispatch(4)
    c.end_call(2)
    assert c.handler_of(3) == "respondent"
    assert c.handler_of(4) == "queued"


def test_end_call_frees_the_employee_when_queue_is_empty():
    c = CallCenter(1, 0, 0)
    c.dispatch(1)
    c.end_call(1)
    assert c.dispatch(2) == "respondent"


def test_end_unknown_call_returns_false():
    c = CallCenter(1, 1, 1)
    assert c.end_call(99) is False


def test_end_call_twice_returns_false():
    c = CallCenter(1, 0, 0)
    c.dispatch(1)
    assert c.end_call(1) is True
    assert c.end_call(1) is False


def test_abandoning_a_queued_call_removes_it():
    c = CallCenter(1, 0, 0)
    c.dispatch(1)
    c.dispatch(2)
    c.dispatch(3)
    assert c.end_call(2) is True
    c.end_call(1)
    assert c.handler_of(3) == "respondent"
    assert c.handler_of(2) == ""


def test_handler_of_unknown_call_is_empty():
    c = CallCenter(1, 1, 1)
    assert c.handler_of(42) == ""


def test_handler_of_reports_the_active_level():
    c = CallCenter(1, 1, 1)
    c.dispatch(1)
    c.dispatch(2)
    assert c.handler_of(1) == "respondent"
    assert c.handler_of(2) == "manager"
