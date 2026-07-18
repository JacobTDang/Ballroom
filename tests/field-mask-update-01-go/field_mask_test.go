package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSiblingUntouchedLeafApply(t *testing.T) {
	target := map[string]any{
		"name":    "Alice",
		"address": map[string]any{"city": "SF", "zip": "94107"},
		"tags":    map[string]any{"vip": "true"},
	}
	source := map[string]any{"address": map[string]any{"city": "NYC"}}

	if err := Update(target, source, []string{"address.city"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	wantAddress := map[string]any{"city": "NYC", "zip": "94107"}
	if !reflect.DeepEqual(target["address"], wantAddress) {
		t.Fatalf("address = %v, want %v (sibling zip was disturbed)", target["address"], wantAddress)
	}
	if target["name"] != "Alice" {
		t.Fatal("an untouched top-level sibling changed")
	}
}

func TestMultiPath(t *testing.T) {
	target := map[string]any{
		"name":    "Alice",
		"address": map[string]any{"city": "SF", "zip": "94107"},
		"tags":    map[string]any{"vip": "true"},
	}
	source := map[string]any{
		"name":    "Bob",
		"address": map[string]any{"zip": "10001"},
		"tags":    map[string]any{"vip": "false"},
	}

	if err := Update(target, source, []string{"name", "address.zip", "tags.vip"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if target["name"] != "Bob" {
		t.Fatal("name wasn't updated")
	}
	wantAddress := map[string]any{"city": "SF", "zip": "10001"}
	if !reflect.DeepEqual(target["address"], wantAddress) {
		t.Fatalf("address = %v, want %v (city sibling shouldn't move)", target["address"], wantAddress)
	}
	wantTags := map[string]any{"vip": "false"}
	if !reflect.DeepEqual(target["tags"], wantTags) {
		t.Fatalf("tags = %v, want %v", target["tags"], wantTags)
	}
}

func TestClearViaOmission(t *testing.T) {
	target := map[string]any{
		"name":    "Alice",
		"address": map[string]any{"city": "SF", "zip": "94107"},
	}
	source := map[string]any{"address": map[string]any{}}

	if err := Update(target, source, []string{"address.zip", "name"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	wantAddress := map[string]any{"city": "SF"}
	if !reflect.DeepEqual(target["address"], wantAddress) {
		t.Fatalf("address = %v, want %v (zip should have been cleared)", target["address"], wantAddress)
	}
	if _, exists := target["name"]; exists {
		t.Fatal("name should have been cleared")
	}
}

func TestMissingIntermediateError(t *testing.T) {
	target := map[string]any{"name": "Alice"} // no "address" key at all
	source := map[string]any{"address": map[string]any{"city": "NYC"}}

	err := Update(target, source, []string{"address.city"})
	if err == nil {
		t.Fatal("expected an error for a missing intermediate")
	}
	if !strings.Contains(err.Error(), "address") {
		t.Fatalf("error %q doesn't name the offending path", err.Error())
	}

	want := map[string]any{"name": "Alice"}
	if !reflect.DeepEqual(target, want) {
		t.Fatalf("target changed despite the error: got %v", target)
	}
}

func TestScalarIntermediateError(t *testing.T) {
	target := map[string]any{"name": "Alice", "address": "not-an-object"}
	source := map[string]any{"address": map[string]any{"city": "NYC"}}

	err := Update(target, source, []string{"address.city"})
	if err == nil {
		t.Fatal("expected an error for a scalar intermediate")
	}
	if !strings.Contains(err.Error(), "address") {
		t.Fatalf("error %q doesn't name the offending path", err.Error())
	}

	want := map[string]any{"name": "Alice", "address": "not-an-object"}
	if !reflect.DeepEqual(target, want) {
		t.Fatalf("target changed despite the error: got %v", target)
	}
}

func TestEmptyMaskError(t *testing.T) {
	target := map[string]any{"name": "Alice"}
	if err := Update(target, map[string]any{}, []string{}); err == nil {
		t.Fatal("expected an error for an empty mask")
	}
	want := map[string]any{"name": "Alice"}
	if !reflect.DeepEqual(target, want) {
		t.Fatalf("target changed despite the error: got %v", target)
	}
}
