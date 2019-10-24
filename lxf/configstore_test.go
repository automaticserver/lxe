package lxf

import "testing"

func TestNewConfigStoreKeyReserved(t *testing.T) {
	cs := NewConfigStore().WithReserved("foo", "foo.bar")

	if !cs.IsReserved("foo") {
		t.Errorf("key foo should be reserved")
	}

	if !cs.IsReserved("foo.bar") {
		t.Errorf("key foo.bar should be reserved")
	}

	if cs.IsReserved("foo.baz") {
		t.Errorf("key foo.baz should not be reserved")
	}

	if cs.IsReserved("bar") {
		t.Errorf("key bar should not be reserved")
	}
}

func TestNewConfigStorePrefixReserved(t *testing.T) {
	cs := NewConfigStore().WithReservedPrefixes("foo", "zoo.bar")

	if !cs.IsReserved("foo") {
		t.Errorf("key foo should be reserved")
	}

	if !cs.IsReserved("foo.bar") {
		t.Errorf("key foo.bar should be reserved")
	}

	if cs.IsReserved("fool") {
		t.Errorf("key fool should not be reserved")
	}

	if cs.IsReserved("zoo") {
		t.Errorf("key zoo should not be reserved")
	}

	if !cs.IsReserved("zoo.bar") {
		t.Errorf("key zoo.bar.tree should be reserved")
	}

	if !cs.IsReserved("zoo.bar.tree") {
		t.Errorf("key zoo.bar.tree should be reserved")
	}
}

func TestConfigStoreUnreserved(t *testing.T) {
	cs := NewConfigStore().WithReserved("tree", "kanu.manu").WithReservedPrefixes("foo", "zoo.bar")

	unres := cs.UnreservedMap(map[string]string{
		"tree":             "no",
		"kanu.manu":        "no",
		"foo":              "no",
		"foo.lala":         "no",
		"foo.zoo.bar.yeah": "no",
		"foola":            "yes",
		"jouse":            "yes",
	})

	for k, v := range unres {
		if v != "yes" {
			t.Errorf("key %v should be in the unreserved result", k)
		}
	}
}
