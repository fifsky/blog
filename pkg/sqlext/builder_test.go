package sqlext

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderDefaultsToSelectAll(t *testing.T) {
	b := NewBuilder()

	assert.Equal(t, "SELECT *", b.SQL())
	assert.Equal(t, "SELECT *", b.String())
	assert.Nil(t, b.Args())
}

func TestBuilderZeroValue(t *testing.T) {
	var b Builder
	b.Select("id").From("api").Where("id = ?", 1)

	assert.Equal(t, "SELECT id FROM api WHERE id = ?", b.SQL())
	assert.EqualValues(t, []any{1}, b.Args())
}

func TestBuilderSelectAndFrom(t *testing.T) {
	b := NewBuilder().
		Select("api.id, api.api_name").
		From("api")

	assert.Equal(t, "SELECT api.id, api.api_name FROM api", b.SQL())
	assert.Nil(t, b.Args())
}

func TestBuilderSelectAndFromAreJoinedWhenCalledMultipleTimes(t *testing.T) {
	b := NewBuilder().
		Select("api.id").
		Select("api.api_name").
		From("api").
		From("system_service s")

	assert.Equal(t, "SELECT api.id, api.api_name FROM api, system_service s", b.SQL())
}

func TestBuilderSelectAndFromIgnoreEmptyValues(t *testing.T) {
	b := NewBuilder().Select("").Select("  ").Select("id").From("")

	assert.Equal(t, "SELECT id", b.SQL())
}

func TestBuilderFromSupportsJoins(t *testing.T) {
	b := NewBuilder().
		Select("*").
		From("api LEFT JOIN system_service ON system_service.id = api.service_id AND system_service.type = ?", "api").
		Where("api.status = ?", "online")

	assert.Equal(t,
		"SELECT * FROM api LEFT JOIN system_service ON system_service.id = api.service_id AND system_service.type = ?"+
			" WHERE api.status = ?",
		b.SQL(),
	)
	assert.EqualValues(t, []any{"api", "online"}, b.Args())
}

func TestBuilderWhere(t *testing.T) {
	b := NewBuilder().
		From("api").
		Where("service_id = ? AND status = ? and api_type = ?", 1, 2, 3).
		Where("api_name = ?", "brett")

	assert.Equal(t, "SELECT * FROM api WHERE service_id = ? AND status = ? and api_type = ? AND api_name = ?", b.SQL())
	assert.EqualValues(t, []any{1, 2, 3, "brett"}, b.Args())
}

func TestBuilderWhereIf(t *testing.T) {
	b := NewBuilder().
		From("api").
		Where("a = ?", 1).
		WhereIf(false, "b = ?", 2).
		WhereIf(true, "(api_name LIKE ? or api_url LIKE ?)", "%a%", "%a%")

	assert.Equal(t, "SELECT * FROM api WHERE a = ? AND (api_name LIKE ? or api_url LIKE ?)", b.SQL())
	assert.EqualValues(t, []any{1, "%a%", "%a%"}, b.Args())
}

func TestBuilderOrWhere(t *testing.T) {
	b := NewBuilder().
		From("api").
		Where("a = ?", 1).
		OrWhere("b = ?", 2).
		OrWhereIf(false, "c = ?", 3).
		OrWhereIf(true, "d = ?", 4)

	assert.Equal(t, "SELECT * FROM api WHERE a = ? OR b = ? OR d = ?", b.SQL())
	assert.EqualValues(t, []any{1, 2, 4}, b.Args())
}

func TestBuilderOrWhereAsFirstCondition(t *testing.T) {
	b := NewBuilder().From("api").OrWhere("a = ?", 1)

	assert.Equal(t, "SELECT * FROM api WHERE a = ?", b.SQL())
}

func TestBuilderWhereIgnoresEmptyConditions(t *testing.T) {
	b := NewBuilder().From("api").Where("").Where("   ").Where("a = ?", 1)

	assert.Equal(t, "SELECT * FROM api WHERE a = ?", b.SQL())
	assert.EqualValues(t, []any{1}, b.Args())
}

func TestBuilderGroupAndHaving(t *testing.T) {
	b := NewBuilder().
		Select("service_id, COUNT(*) AS total").
		From("api").
		Where("status = ?", 1).
		GroupBy("service_id").
		Having("COUNT(*) > ?", 10).
		HavingIf(false, "COUNT(*) < ?", 100).
		OrHaving("SUM(price) > ?", 5)

	assert.Equal(t,
		"SELECT service_id, COUNT(*) AS total FROM api WHERE status = ?"+
			" GROUP BY service_id HAVING COUNT(*) > ? OR SUM(price) > ?",
		b.SQL(),
	)
	assert.EqualValues(t, []any{1, 10, 5}, b.Args())
}

func TestBuilderGroupAlias(t *testing.T) {
	b := NewBuilder().Select("*").From("api").Group("service_id, status")

	assert.Equal(t, "SELECT * FROM api GROUP BY service_id, status", b.SQL())
}

func TestBuilderOrderAndLimit(t *testing.T) {
	b := NewBuilder().
		From("api").
		OrderBy("created_at desc").
		Order("api_name asc, id desc").
		Limit(10).
		Offset(20)

	assert.Equal(t, "SELECT * FROM api ORDER BY created_at desc, api_name asc, id desc LIMIT 10 OFFSET 20", b.SQL())
}

func TestBuilderLimitAndOffsetAreIgnoredWhenNotPositive(t *testing.T) {
	b := NewBuilder().From("api").Limit(0).Offset(-1)

	assert.Equal(t, "SELECT * FROM api", b.SQL())

	b.Limit(5).Offset(0)
	assert.Equal(t, "SELECT * FROM api LIMIT 5", b.SQL())
}

func TestBuilderArgsFollowSQLPlaceholdersNotCallOrder(t *testing.T) {
	// from is called after where but the placeholders of the from clause come
	// first in the generated SQL
	b := NewBuilder().
		Where("api.status = ?", "online").
		From("api LEFT JOIN system_service ON system_service.type = ?", "api").
		GroupBy("api.id").
		Having("COUNT(*) > ?", 3)

	assert.Equal(t,
		"SELECT * FROM api LEFT JOIN system_service ON system_service.type = ?"+
			" WHERE api.status = ? GROUP BY api.id HAVING COUNT(*) > ?",
		b.SQL(),
	)
	assert.EqualValues(t, []any{"api", "online", 3}, b.Args())
}

func TestBuilderBuild(t *testing.T) {
	b := NewBuilder().Select("id").From("api").Where("id = ?", 1).Limit(1)

	query, args := b.Build()
	assert.Equal(t, "SELECT id FROM api WHERE id = ? LIMIT 1", query)
	assert.Equal(t, b.SQL(), query)
	assert.EqualValues(t, []any{1}, args)
}

func TestBuilderReuse(t *testing.T) {
	b := NewBuilder().From("api")
	b.Where("a = ?", 1)

	assert.Equal(t, "SELECT * FROM api WHERE a = ?", b.SQL())
	assert.EqualValues(t, []any{1}, b.Args())

	// additional conditions are appended to the same builder
	b.Where("b = ?", 2)
	assert.Equal(t, "SELECT * FROM api WHERE a = ? AND b = ?", b.SQL())
	assert.EqualValues(t, []any{1, 2}, b.Args())
}

func ExampleBuilder() {
	keyword := "log"
	menuID := 10

	q := NewBuilder().
		Select("api.id, api.api_name").
		From("api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id").
		Where("api.service_id = ? AND api.status = ? and api.api_type = ?", 1, 2, 3).
		WhereIf(keyword != "", "(api.api_name LIKE ? or api.api_url LIKE ?)", "%"+keyword+"%", "%"+keyword+"%").
		WhereIf(menuID > 0, "api.id NOT IN (SELECT api_id FROM system_menu_api WHERE menu_id = ?)", menuID).
		OrderBy("api.created_at desc").
		Limit(20)

	fmt.Println(q.SQL())
	fmt.Println(q.Args())
	// Output:
	// SELECT api.id, api.api_name FROM api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id WHERE api.service_id = ? AND api.status = ? and api.api_type = ? AND (api.api_name LIKE ? or api.api_url LIKE ?) AND api.id NOT IN (SELECT api_id FROM system_menu_api WHERE menu_id = ?) ORDER BY api.created_at desc LIMIT 20
	// [1 2 3 %log% %log% 10]
}

func TestBuilderSelectAndOrderBySupportArgs(t *testing.T) {
	match := "MATCH(a.title, a.content) AGAINST(? IN NATURAL LANGUAGE MODE)"
	b := NewBuilder().
		Select("a.id, SUBSTRING(a.content, GREATEST(1, LOCATE(?, a.content) - 50), 150) AS snippet", "log").
		From("articles a").
		Where("a.status = ?", 1).
		OrderBy(match+" DESC", "log").
		OrderBy("a.updated_at DESC")

	assert.Equal(t,
		"SELECT a.id, SUBSTRING(a.content, GREATEST(1, LOCATE(?, a.content) - 50), 150) AS snippet"+
			" FROM articles a WHERE a.status = ?"+
			" ORDER BY "+match+" DESC, a.updated_at DESC",
		b.SQL(),
	)
	// select args come first, order args last, following the SQL placeholder order
	assert.EqualValues(t, []any{"log", 1, "log"}, b.Args())
}

func TestBuilderSelectIgnoresEmptyValuesButKeepsArgsOfOthers(t *testing.T) {
	b := NewBuilder().Select("", "dropped").Select("id")

	assert.Equal(t, "SELECT id", b.SQL())
	assert.Empty(t, b.Args())
}

func TestBuilderWhereExpandsSliceArgs(t *testing.T) {
	b := NewBuilder().
		From("api").
		Where("id IN (?)", []int{1, 2, 3}).
		Where("api_name = ?", "brett")

	assert.Equal(t, "SELECT * FROM api WHERE id IN (?,?,?) AND api_name = ?", b.SQL())
	assert.EqualValues(t, []any{1, 2, 3, "brett"}, b.Args())
}

func TestBuilderWhereExpandsSliceArgsInsideCondition(t *testing.T) {
	b := NewBuilder().
		From("api").
		Where("service_id = ? AND api_type IN (?) AND status = ?", 1, []string{"a", "b"}, 2)

	assert.Equal(t, "SELECT * FROM api WHERE service_id = ? AND api_type IN (?,?) AND status = ?", b.SQL())
	assert.EqualValues(t, []any{1, "a", "b", 2}, b.Args())
}

func TestBuilderWhereExpandsArrayArgs(t *testing.T) {
	b := NewBuilder().From("api").Where("id IN (?)", [3]int{1, 2, 3})

	assert.Equal(t, "SELECT * FROM api WHERE id IN (?,?,?)", b.SQL())
	assert.EqualValues(t, []any{1, 2, 3}, b.Args())
}

func TestBuilderWhereKeepsSingleElementSlice(t *testing.T) {
	b := NewBuilder().From("api").Where("id IN (?)", []int{7})

	assert.Equal(t, "SELECT * FROM api WHERE id IN (?)", b.SQL())
	assert.EqualValues(t, []any{7}, b.Args())
}

func TestBuilderWhereExpandsEmptySliceToNull(t *testing.T) {
	// an empty list is not valid SQL and IN (NULL) matches nothing
	b := NewBuilder().From("api").Where("id IN (?)", []int{})

	assert.Equal(t, "SELECT * FROM api WHERE id IN (NULL)", b.SQL())
	assert.Empty(t, b.Args())
}

func TestBuilderWhereKeepsByteSliceAsOneValue(t *testing.T) {
	b := NewBuilder().From("api").Where("hash = ?", []byte("ab"))

	assert.Equal(t, "SELECT * FROM api WHERE hash = ?", b.SQL())
	assert.EqualValues(t, []any{[]byte("ab")}, b.Args())
}

// joinedIDs converts itself into a single comma joined value, like a column
// storing a list of ids.
type joinedIDs []int

func (joinedIDs) Value() (driver.Value, error) { return nil, nil }

func TestBuilderWhereKeepsDriverValuerSliceAsOneValue(t *testing.T) {
	ids := joinedIDs{1, 2}
	b := NewBuilder().From("api").Where("share_ids = ?", ids)

	assert.Equal(t, "SELECT * FROM api WHERE share_ids = ?", b.SQL())
	assert.EqualValues(t, []any{ids}, b.Args())
}

func TestBuilderWhereKeepsNilAsPlaceholder(t *testing.T) {
	b := NewBuilder().From("api").Where("deleted_at = ?", nil)

	assert.Equal(t, "SELECT * FROM api WHERE deleted_at = ?", b.SQL())
	assert.EqualValues(t, []any{nil}, b.Args())
}

func TestBuilderHavingExpandsSliceArgs(t *testing.T) {
	b := NewBuilder().
		Select("service_id, COUNT(*) AS total").
		From("api").
		GroupBy("service_id").
		Having("service_id IN (?)", []int{1, 2})

	assert.Equal(t,
		"SELECT service_id, COUNT(*) AS total FROM api GROUP BY service_id HAVING service_id IN (?,?)",
		b.SQL(),
	)
	assert.EqualValues(t, []any{1, 2}, b.Args())
}

func TestBuilderWhereIfAndOrWhereExpandSliceArgs(t *testing.T) {
	b := NewBuilder().
		From("api").
		WhereIf(false, "id IN (?)", []int{1, 2}).
		WhereIf(true, "service_id IN (?)", []int{3, 4}).
		OrWhere("status IN (?)", []string{"online"})

	assert.Equal(t, "SELECT * FROM api WHERE service_id IN (?,?) OR status IN (?)", b.SQL())
	assert.EqualValues(t, []any{3, 4, "online"}, b.Args())
}
