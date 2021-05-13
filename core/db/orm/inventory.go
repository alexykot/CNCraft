// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Inventory is an object representing the database table.
type Inventory struct {
	PlayerID   uuid.UUID `boil:"player_id" json:"player_id" toml:"player_id" yaml:"player_id"`
	SlotNumber int16     `boil:"slot_number" json:"slot_number" toml:"slot_number" yaml:"slot_number"`
	ItemID     int16     `boil:"item_id" json:"item_id" toml:"item_id" yaml:"item_id"`
	ItemCount  int16     `boil:"item_count" json:"item_count" toml:"item_count" yaml:"item_count"`

	R *inventoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L inventoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var InventoryColumns = struct {
	PlayerID   string
	SlotNumber string
	ItemID     string
	ItemCount  string
}{
	PlayerID:   "player_id",
	SlotNumber: "slot_number",
	ItemID:     "item_id",
	ItemCount:  "item_count",
}

// Generated where

type whereHelperuuid_UUID struct{ field string }

func (w whereHelperuuid_UUID) EQ(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelperuuid_UUID) NEQ(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperuuid_UUID) LT(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelperuuid_UUID) LTE(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperuuid_UUID) GT(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelperuuid_UUID) GTE(x uuid.UUID) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperint16 struct{ field string }

func (w whereHelperint16) EQ(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint16) NEQ(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint16) LT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint16) LTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint16) GT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint16) GTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint16) IN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint16) NIN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var InventoryWhere = struct {
	PlayerID   whereHelperuuid_UUID
	SlotNumber whereHelperint16
	ItemID     whereHelperint16
	ItemCount  whereHelperint16
}{
	PlayerID:   whereHelperuuid_UUID{field: "\"cncraft\".\"inventory\".\"player_id\""},
	SlotNumber: whereHelperint16{field: "\"cncraft\".\"inventory\".\"slot_number\""},
	ItemID:     whereHelperint16{field: "\"cncraft\".\"inventory\".\"item_id\""},
	ItemCount:  whereHelperint16{field: "\"cncraft\".\"inventory\".\"item_count\""},
}

// InventoryRels is where relationship names are stored.
var InventoryRels = struct {
	Player string
}{
	Player: "Player",
}

// inventoryR is where relationships are stored.
type inventoryR struct {
	Player *Player `boil:"Player" json:"Player" toml:"Player" yaml:"Player"`
}

// NewStruct creates a new relationship struct
func (*inventoryR) NewStruct() *inventoryR {
	return &inventoryR{}
}

// inventoryL is where Load methods for each relationship are stored.
type inventoryL struct{}

var (
	inventoryAllColumns            = []string{"player_id", "slot_number", "item_id", "item_count"}
	inventoryColumnsWithoutDefault = []string{"player_id"}
	inventoryColumnsWithDefault    = []string{"slot_number", "item_id", "item_count"}
	inventoryPrimaryKeyColumns     = []string{"player_id", "slot_number"}
)

type (
	// InventorySlice is an alias for a slice of pointers to Inventory.
	// This should generally be used opposed to []Inventory.
	InventorySlice []*Inventory

	inventoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	inventoryType                 = reflect.TypeOf(&Inventory{})
	inventoryMapping              = queries.MakeStructMapping(inventoryType)
	inventoryPrimaryKeyMapping, _ = queries.BindMapping(inventoryType, inventoryMapping, inventoryPrimaryKeyColumns)
	inventoryInsertCacheMut       sync.RWMutex
	inventoryInsertCache          = make(map[string]insertCache)
	inventoryUpdateCacheMut       sync.RWMutex
	inventoryUpdateCache          = make(map[string]updateCache)
	inventoryUpsertCacheMut       sync.RWMutex
	inventoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single inventory record from the query.
func (q inventoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Inventory, error) {
	o := &Inventory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "orm: failed to execute a one query for inventory")
	}

	return o, nil
}

// All returns all Inventory records from the query.
func (q inventoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (InventorySlice, error) {
	var o []*Inventory

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "orm: failed to assign all query results to Inventory slice")
	}

	return o, nil
}

// Count returns the count of all Inventory records in the query.
func (q inventoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to count inventory rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q inventoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "orm: failed to check if inventory exists")
	}

	return count > 0, nil
}

// Player pointed to by the foreign key.
func (o *Inventory) Player(mods ...qm.QueryMod) playerQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.PlayerID),
	}

	queryMods = append(queryMods, mods...)

	query := Players(queryMods...)
	queries.SetFrom(query.Query, "\"cncraft\".\"players\"")

	return query
}

// LoadPlayer allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (inventoryL) LoadPlayer(ctx context.Context, e boil.ContextExecutor, singular bool, maybeInventory interface{}, mods queries.Applicator) error {
	var slice []*Inventory
	var object *Inventory

	if singular {
		object = maybeInventory.(*Inventory)
	} else {
		slice = *maybeInventory.(*[]*Inventory)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &inventoryR{}
		}
		if !queries.IsNil(object.PlayerID) {
			args = append(args, object.PlayerID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &inventoryR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.PlayerID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.PlayerID) {
				args = append(args, obj.PlayerID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`cncraft.players`),
		qm.WhereIn(`cncraft.players.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Player")
	}

	var resultSlice []*Player
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Player")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for players")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for players")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Player = foreign
		if foreign.R == nil {
			foreign.R = &playerR{}
		}
		foreign.R.Inventories = append(foreign.R.Inventories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.PlayerID, foreign.ID) {
				local.R.Player = foreign
				if foreign.R == nil {
					foreign.R = &playerR{}
				}
				foreign.R.Inventories = append(foreign.R.Inventories, local)
				break
			}
		}
	}

	return nil
}

// SetPlayer of the inventory to the related item.
// Sets o.R.Player to related.
// Adds o to related.R.Inventories.
func (o *Inventory) SetPlayer(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Player) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"cncraft\".\"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"player_id"}),
		strmangle.WhereClause("\"", "\"", 2, inventoryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.PlayerID, o.SlotNumber}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.PlayerID, related.ID)
	if o.R == nil {
		o.R = &inventoryR{
			Player: related,
		}
	} else {
		o.R.Player = related
	}

	if related.R == nil {
		related.R = &playerR{
			Inventories: InventorySlice{o},
		}
	} else {
		related.R.Inventories = append(related.R.Inventories, o)
	}

	return nil
}

// Inventories retrieves all the records using an executor.
func Inventories(mods ...qm.QueryMod) inventoryQuery {
	mods = append(mods, qm.From("\"cncraft\".\"inventory\""))
	return inventoryQuery{NewQuery(mods...)}
}

// FindInventory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindInventory(ctx context.Context, exec boil.ContextExecutor, playerID uuid.UUID, slotNumber int16, selectCols ...string) (*Inventory, error) {
	inventoryObj := &Inventory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"cncraft\".\"inventory\" where \"player_id\"=$1 AND \"slot_number\"=$2", sel,
	)

	q := queries.Raw(query, playerID, slotNumber)

	err := q.Bind(ctx, exec, inventoryObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "orm: unable to select from inventory")
	}

	return inventoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Inventory) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("orm: no inventory provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(inventoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	inventoryInsertCacheMut.RLock()
	cache, cached := inventoryInsertCache[key]
	inventoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			inventoryAllColumns,
			inventoryColumnsWithDefault,
			inventoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(inventoryType, inventoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"cncraft\".\"inventory\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"cncraft\".\"inventory\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "orm: unable to insert into inventory")
	}

	if !cached {
		inventoryInsertCacheMut.Lock()
		inventoryInsertCache[key] = cache
		inventoryInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Inventory.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Inventory) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	inventoryUpdateCacheMut.RLock()
	cache, cached := inventoryUpdateCache[key]
	inventoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			inventoryAllColumns,
			inventoryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("orm: unable to update inventory, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"cncraft\".\"inventory\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, inventoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, append(wl, inventoryPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to update inventory row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by update for inventory")
	}

	if !cached {
		inventoryUpdateCacheMut.Lock()
		inventoryUpdateCache[key] = cache
		inventoryUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q inventoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to update all for inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to retrieve rows affected for inventory")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o InventorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("orm: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"cncraft\".\"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, inventoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to update all in inventory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to retrieve rows affected all in update all inventory")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Inventory) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("orm: no inventory provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(inventoryColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	inventoryUpsertCacheMut.RLock()
	cache, cached := inventoryUpsertCache[key]
	inventoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			inventoryAllColumns,
			inventoryColumnsWithDefault,
			inventoryColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			inventoryAllColumns,
			inventoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("orm: unable to upsert inventory, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(inventoryPrimaryKeyColumns))
			copy(conflict, inventoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"cncraft\".\"inventory\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(inventoryType, inventoryMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "orm: unable to upsert inventory")
	}

	if !cached {
		inventoryUpsertCacheMut.Lock()
		inventoryUpsertCache[key] = cache
		inventoryUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Inventory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Inventory) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("orm: no Inventory provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), inventoryPrimaryKeyMapping)
	sql := "DELETE FROM \"cncraft\".\"inventory\" WHERE \"player_id\"=$1 AND \"slot_number\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete from inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by delete for inventory")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q inventoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("orm: no inventoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete all from inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by deleteall for inventory")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o InventorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"cncraft\".\"inventory\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inventoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete all from inventory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by deleteall for inventory")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Inventory) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindInventory(ctx, exec, o.PlayerID, o.SlotNumber)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InventorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := InventorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"cncraft\".\"inventory\".* FROM \"cncraft\".\"inventory\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inventoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "orm: unable to reload all in InventorySlice")
	}

	*o = slice

	return nil
}

// InventoryExists checks if the Inventory row exists.
func InventoryExists(ctx context.Context, exec boil.ContextExecutor, playerID uuid.UUID, slotNumber int16) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"cncraft\".\"inventory\" where \"player_id\"=$1 AND \"slot_number\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, playerID, slotNumber)
	}
	row := exec.QueryRowContext(ctx, sql, playerID, slotNumber)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "orm: unable to check if inventory exists")
	}

	return exists, nil
}