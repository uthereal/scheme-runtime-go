package column

import (
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Test_Column_ColumnName tests the ColumnName method of Column.
func Test_Column_ColumnName(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	assert.Equal(t, "id", c.ColumnName())
}

// Test_Column_Asc tests the Asc method of Column.
func Test_Column_Asc(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	o := c.Asc()
	assert.NotNil(t, o)
}

// Test_Column_Desc tests the Desc method of Column.
func Test_Column_Desc(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	o := c.Desc()
	assert.NotNil(t, o)
}

// Test_Column_Eq tests the Eq method of Column.
func Test_Column_Eq(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Eq("foo")
	assert.NotNil(t, w)
}

// Test_Column_Neq tests the Neq method of Column.
func Test_Column_Neq(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Neq("foo")
	assert.NotNil(t, w)
}

// Test_Column_In tests the In method of Column.
func Test_Column_In(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.In("foo", "bar")
	assert.NotNil(t, w)
}

// Test_Column_NotIn tests the NotIn method of Column.
func Test_Column_NotIn(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.NotIn("foo", "bar")
	assert.NotNil(t, w)
}

// Test_Column_InQuery tests the InQuery method of Column.
func Test_Column_InQuery(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	type mockQueryStateProvider struct {
		contract.QueryStateProvider
	}
	q := mockQueryStateProvider{}
	w := c.InQuery(q)
	assert.NotNil(t, w)
}

// Test_Column_NotInQuery tests the NotInQuery method of Column.
func Test_Column_NotInQuery(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	type mockQueryStateProvider struct {
		contract.QueryStateProvider
	}
	q := mockQueryStateProvider{}
	w := c.NotInQuery(q)
	assert.NotNil(t, w)
}

// Test_Column_Between tests the Between method of Column.
func Test_Column_Between(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.Between("a", "z")
	assert.NotNil(t, w)
}

// Test_Column_NotBetween tests the NotBetween method of Column.
func Test_Column_NotBetween(t *testing.T) {
	c := Column[any, string]{
		Name: "id",
	}
	w := c.NotBetween("a", "z")
	assert.NotNil(t, w)
}

func Test_Column_PostgresCast(t *testing.T) {
	assert.Equal(t, "::text[]", Column[any, string]{}.PostgresCast())
	assert.Equal(t, "::integer[]", Column[any, int32]{}.PostgresCast())
	assert.Equal(t, "::bigint[]", Column[any, int64]{}.PostgresCast())
	assert.Equal(t, "::bigint[]", Column[any, uint64]{}.PostgresCast())
	assert.Equal(t, "::integer[]", Column[any, uint32]{}.PostgresCast())
	assert.Equal(t, "::smallint[]", Column[any, uint16]{}.PostgresCast())
	assert.Equal(t, "::boolean[]", Column[any, bool]{}.PostgresCast())
	assert.Equal(t, "::double precision[]", Column[any, float64]{}.PostgresCast())
	assert.Equal(t, "::timestamp with time zone[]", Column[any, time.Time]{}.PostgresCast())
	assert.Equal(t, "::interval[]", Column[any, pgtype.Interval]{}.PostgresCast())
	assert.Equal(t, "::interval[]", Column[any, *pgtype.Interval]{}.PostgresCast())
	assert.Equal(t, "::numeric[]", Column[any, pgtype.Numeric]{}.PostgresCast())
	assert.Equal(t, "::point[]", Column[any, pgtype.Point]{}.PostgresCast())
	assert.Equal(t, "::varbit[]", Column[any, pgtype.Bits]{}.PostgresCast())
	assert.Equal(t, "::bytea[]", Column[any, []byte]{}.PostgresCast())
	assert.Equal(t, "::inet[]", Column[any, netip.Addr]{}.PostgresCast())
	assert.Equal(t, "::inet[]", Column[any, *netip.Addr]{}.PostgresCast())
	assert.Equal(t, "::cidr[]", Column[any, netip.Prefix]{}.PostgresCast())
	assert.Equal(t, "::cidr[]", Column[any, *netip.Prefix]{}.PostgresCast())
	assert.Equal(t, "::macaddr[]", Column[any, net.HardwareAddr]{}.PostgresCast())
	assert.Equal(t, "::macaddr[]", Column[any, *net.HardwareAddr]{}.PostgresCast())
}

func Test_Column_ToTypedSlice_NetIP(t *testing.T) {
	c := Column[any, netip.Addr]{Name: "ip"}
	ip1 := netip.MustParseAddr("192.168.1.1")
	ip2 := netip.MustParseAddr("10.0.0.1")

	slice := []any{ip1, ip2}
	typed := c.ToTypedSlice(slice)
	assert.Equal(t, []netip.Addr{ip1, ip2}, typed)

	ptrSlice := []any{&ip1, &ip2}
	typedPtr := c.ToTypedSlice(ptrSlice)
	assert.Equal(t, []*netip.Addr{&ip1, &ip2}, typedPtr)
}

