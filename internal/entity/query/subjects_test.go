package query

import (
	"testing"

	"github.com/photoprism/photoprism/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeople(t *testing.T) {
	entity.ValidateFixtures(t)
	if results, err := People(); err != nil {
		t.Fatal(err)
	} else {
		assert.LessOrEqual(t, 3, len(results))
		t.Logf("people: %#v", results)
	}
}

func TestPeopleCount(t *testing.T) {
	entity.ValidateFixtures(t)
	if result, err := PeopleCount(); err != nil {
		t.Fatal(err)
	} else {
		assert.LessOrEqual(t, 3, result)
		t.Logf("there are %d people", result)
	}
}

func TestSubjects(t *testing.T) {
	entity.ValidateFixtures(t)
	results, err := Subjects(3, 0)

	if err != nil {
		t.Fatal(err)
	}

	assert.GreaterOrEqual(t, len(results), 1)

	for _, val := range results {
		assert.IsType(t, entity.Subject{}, val)
	}
}

func TestSubjectMap(t *testing.T) {
	entity.ValidateFixtures(t)
	results, err := SubjectMap()

	if err != nil {
		t.Fatal(err)
	}

	assert.GreaterOrEqual(t, len(results), 1)

	for _, val := range results {
		assert.IsType(t, entity.Subject{}, val)
	}
}

func TestRemoveOrphanSubjects(t *testing.T) {
	entity.ValidateFixtures(t)
	affected, err := RemoveOrphanSubjects()

	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		require.NoError(t, Db().Create(entity.SubjectFixtures.Pointer("dangling")).Error)
	})

	assert.Equal(t, int64(1), affected)
}

func TestCreateMarkerSubjects(t *testing.T) {
	entity.ValidateFixtures(t)
	affected, err := CreateMarkerSubjects()

	assert.NoError(t, err)
	assert.LessOrEqual(t, int64(0), affected)
}
