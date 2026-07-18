import pytest

from solution import update


def test_sibling_untouched_leaf_apply():
    target = {
        "name": "Alice",
        "address": {"city": "SF", "zip": "94107"},
        "tags": {"vip": "true"},
    }
    source = {"address": {"city": "NYC"}}

    update(target, source, ["address.city"])

    assert target["address"] == {"city": "NYC", "zip": "94107"}, "sibling zip was disturbed"
    assert target["name"] == "Alice", "an untouched top-level sibling changed"
    assert target["tags"] == {"vip": "true"}, "an untouched top-level sibling changed"


def test_multi_path():
    target = {
        "name": "Alice",
        "address": {"city": "SF", "zip": "94107"},
        "tags": {"vip": "true"},
    }
    source = {
        "name": "Bob",
        "address": {"zip": "10001"},
        "tags": {"vip": "false"},
    }

    update(target, source, ["name", "address.zip", "tags.vip"])

    assert target["name"] == "Bob"
    assert target["address"] == {"city": "SF", "zip": "10001"}, "city sibling shouldn't move"
    assert target["tags"] == {"vip": "false"}


def test_clear_via_omission():
    target = {
        "name": "Alice",
        "address": {"city": "SF", "zip": "94107"},
    }
    source = {"address": {}}  # address.zip and name both absent from source

    update(target, source, ["address.zip", "name"])

    assert target["address"] == {"city": "SF"}, "address.zip should have been cleared"
    assert "name" not in target, "name should have been cleared"


def test_missing_intermediate_error():
    target = {"name": "Alice"}  # no "address" key at all
    source = {"address": {"city": "NYC"}}

    with pytest.raises(ValueError, match="address"):
        update(target, source, ["address.city"])

    assert target == {"name": "Alice"}, "target changed despite the error"


def test_scalar_intermediate_error():
    target = {"name": "Alice", "address": "not-an-object"}
    source = {"address": {"city": "NYC"}}

    with pytest.raises(ValueError, match="address"):
        update(target, source, ["address.city"])

    assert target == {"name": "Alice", "address": "not-an-object"}, "target changed despite the error"


def test_empty_mask_error():
    target = {"name": "Alice"}
    with pytest.raises(ValueError):
        update(target, {}, [])
    assert target == {"name": "Alice"}, "target changed despite the error"
