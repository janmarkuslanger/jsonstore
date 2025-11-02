package jsonstore_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/janmarkuslanger/jsonstore"
)

func TestStore_Create(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	type Test int

	_, err := jsonstore.NewStore[Test](file)
	if err != nil {
		t.Fatalf("expected creation without errors, got: %v", err)
	}
}

func TestStore_Store_SimpleType(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	s, _ := jsonstore.NewStore[int](file)

	key := "eldenlord"
	val := 1

	s.Set(key, val)
	final, err := s.Get(key)

	if err != nil {
		t.Fatalf("expected get that val without errors but got: %v", err)
	}

	if final != val {
		t.Fatalf("expected get value to be the same but got %v", reflect.TypeOf(final))
	}
}

type MyType struct {
	Name string
	Age  int
}

func TestStore_Store_StructType(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	s, _ := jsonstore.NewStore[MyType](file)

	key := "eldenlord"
	val := MyType{
		Name: "Rainer",
		Age:  45,
	}

	s.Set(key, val)
	final, err := s.Get(key)

	if err != nil {
		t.Fatalf("expected get that val without errors but got: %v", err)
	}

	if final != val {
		t.Fatalf("expected get value to be the same but got %v", final)
	}
}

func TestStore_ConcurrentSetUniqueKeys(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "test.json")

	s, err := jsonstore.NewStore[int](file)
	if err != nil {
		t.Fatalf("new store failed: %v", err)
	}

	const n = 200
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("k%d", i)
			if err := s.Set(key, i); err != nil {
				t.Errorf("set %s faild: %v", key, err)
			}
		}()
	}

	wg.Wait()

	for i := 0; i < n; i++ {
		key := fmt.Sprintf("k%d", i)
		val, err := s.Get(key)
		if err != nil {
			t.Fatalf("expected key %s to exist, got err: %v", key, err)
		}

		if val != i {
			t.Fatalf("expected %d for key %s, got %d", i, key, val)
		}
	}
}

func TestStore_Delete(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "test.json")

	s, err := jsonstore.NewStore[int](file)
	if err != nil {
		t.Fatalf("new store failed: %v", err)
	}

	s.Set("one", 1)
	s.Set("two", 2)

	if err := s.Delete("one"); err != nil {
		t.Fatalf("expected deletion with no err gut got %v", err)
	}

	_, err = s.Get("one")
	if err == nil {
		t.Fatalf("expected one to be deleted but got %v", err)
	}
}
