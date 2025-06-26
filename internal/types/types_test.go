package types_test

import (
	"database/sql/driver"
	"encoding"
	"testing"

	entfield "entgo.io/ent/schema/field"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	entfield.ValueScanner
	entfield.Validator
	gomock.Matcher
} = (*types.UserID)(nil)

func TestParse(t *testing.T) {
	_, err := types.Parse[types.UserID]("abra-cadabra")
	require.Error(t, err)

	userID, err := types.Parse[types.UserID]("f0317e88-bbfe-11ed-8728-461e464ebed8")
	require.NoError(t, err)
	assert.Equal(t, "f0317e88-bbfe-11ed-8728-461e464ebed8", userID.String())
}

func TestMustParse(t *testing.T) {
	assert.Panics(t, func() {
		types.MustParse[types.UserID]("abra-cadabra")
	})

	assert.NotPanics(t, func() {
		userID := types.MustParse[types.UserID]("f0317e88-bbfe-11ed-8728-461e464ebed8")
		assert.Equal(t, "f0317e88-bbfe-11ed-8728-461e464ebed8", userID.String())
	})
}

func TestChatIDNil(t *testing.T) {
	t.Log(types.UserIDNil)
	assert.Equal(t, types.UserIDNil.String(), uuid.Nil.String())
}

func TestChatID_String(t *testing.T) {
	id := types.NewUserID()
	require.NotEmpty(t, id.String())
	assert.Equal(t, uuid.MustParse(id.String()).String(), id.String())
}

func TestChatID_Scan(t *testing.T) {
	const src = "5c9de646-529c-11ed-81ba-461e464ebed9"

	t.Run("from string and bytes", func(t *testing.T) {
		var id1, id2 types.UserID
		{
			err := id1.Scan(src)
			require.NoError(t, err)
		}
		{
			err := id2.Scan([]byte(src))
			require.NoError(t, err)
		}
		assert.Equal(t, id1.String(), id2.String())
		assert.Equal(t, getValueAsString(t, id1), getValueAsString(t, id2))
	})

	t.Run("from NULL", func(t *testing.T) {
		for _, src := range []any{nil, []byte(nil), []byte{}, ""} {
			t.Run("", func(t *testing.T) {
				var id types.UserID
				err := id.Scan(src)
				require.NoError(t, err)
				assert.Equal(t, types.UserIDNil.String(), id.String())
				assert.Equal(t, types.UserIDNil.String(), getValueAsString(t, id))
			})
		}
	})
}

func TestChatID_MarshalText(t *testing.T) {
	userID := types.MustParse[types.UserID]("f0317e88-bbfe-11ed-8728-461e464ebed8")
	v, err := userID.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "f0317e88-bbfe-11ed-8728-461e464ebed8", string(v))

	var userID2 types.UserID
	err = userID2.UnmarshalText(v)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), userID2.String())
}

func TestChatID_IsZero(t *testing.T) {
	id := types.NewUserID()
	assert.False(t, id.IsZero())
	assert.True(t, types.UserIDNil.IsZero())
	assert.Equal(t, uuid.Nil.String(), types.UserIDNil.String())
}

func TestChatID_Matches(t *testing.T) {
	id := types.NewUserID()
	id2 := types.MustParse[types.UserID](id.String())
	// Matched.
	assert.Equal(t, id, id2)
	assert.True(t, id.Matches(id2))
	// Not matched.
	assert.NotEqual(t, id, id2.String())
	assert.False(t, id.Matches(id2.String()))
	assert.NotEqual(t, id, types.NewMessageID())
	assert.False(t, id.Matches(types.NewMessageID()))
}

func TestChatID_Validate(t *testing.T) {
	assert.NoError(t, types.NewUserID().Validate())
}

func getValueAsString(t *testing.T, valuer driver.Valuer) string {
	t.Helper()

	v, err := valuer.Value()
	require.NoError(t, err)
	vv, ok := v.(string)
	require.True(t, ok)
	return vv
}
