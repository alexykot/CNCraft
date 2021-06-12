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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Player is an object representing the database table.
type Player struct {
	ID            uuid.UUID   `boil:"id" json:"id" toml:"id" yaml:"id"`
	ConnID        null.String `boil:"conn_id" json:"conn_id,omitempty" toml:"conn_id" yaml:"conn_id,omitempty"`
	DimensionID   uuid.UUID   `boil:"dimension_id" json:"dimension_id" toml:"dimension_id" yaml:"dimension_id"`
	Username      string      `boil:"username" json:"username" toml:"username" yaml:"username"`
	PositionX     float64     `boil:"position_x" json:"position_x" toml:"position_x" yaml:"position_x"`
	PositionY     float64     `boil:"position_y" json:"position_y" toml:"position_y" yaml:"position_y"`
	PositionZ     float64     `boil:"position_z" json:"position_z" toml:"position_z" yaml:"position_z"`
	Yaw           float64     `boil:"yaw" json:"yaw" toml:"yaw" yaml:"yaw"`
	Pitch         float64     `boil:"pitch" json:"pitch" toml:"pitch" yaml:"pitch"`
	OnGround      bool        `boil:"on_ground" json:"on_ground" toml:"on_ground" yaml:"on_ground"`
	CurrentHotbar int16       `boil:"current_hotbar" json:"current_hotbar" toml:"current_hotbar" yaml:"current_hotbar"`
	CreatedAt     time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *playerR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L playerL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PlayerColumns = struct {
	ID            string
	ConnID        string
	DimensionID   string
	Username      string
	PositionX     string
	PositionY     string
	PositionZ     string
	Yaw           string
	Pitch         string
	OnGround      string
	CurrentHotbar string
	CreatedAt     string
}{
	ID:            "id",
	ConnID:        "conn_id",
	DimensionID:   "dimension_id",
	Username:      "username",
	PositionX:     "position_x",
	PositionY:     "position_y",
	PositionZ:     "position_z",
	Yaw:           "yaw",
	Pitch:         "pitch",
	OnGround:      "on_ground",
	CurrentHotbar: "current_hotbar",
	CreatedAt:     "created_at",
}

// Generated where

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperfloat64 struct{ field string }

func (w whereHelperfloat64) EQ(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat64) NEQ(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat64) LT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat64) LTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat64) GT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat64) GTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperfloat64) IN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperfloat64) NIN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var PlayerWhere = struct {
	ID            whereHelperuuid_UUID
	ConnID        whereHelpernull_String
	DimensionID   whereHelperuuid_UUID
	Username      whereHelperstring
	PositionX     whereHelperfloat64
	PositionY     whereHelperfloat64
	PositionZ     whereHelperfloat64
	Yaw           whereHelperfloat64
	Pitch         whereHelperfloat64
	OnGround      whereHelperbool
	CurrentHotbar whereHelperint16
	CreatedAt     whereHelpertime_Time
}{
	ID:            whereHelperuuid_UUID{field: "\"cncraft\".\"players\".\"id\""},
	ConnID:        whereHelpernull_String{field: "\"cncraft\".\"players\".\"conn_id\""},
	DimensionID:   whereHelperuuid_UUID{field: "\"cncraft\".\"players\".\"dimension_id\""},
	Username:      whereHelperstring{field: "\"cncraft\".\"players\".\"username\""},
	PositionX:     whereHelperfloat64{field: "\"cncraft\".\"players\".\"position_x\""},
	PositionY:     whereHelperfloat64{field: "\"cncraft\".\"players\".\"position_y\""},
	PositionZ:     whereHelperfloat64{field: "\"cncraft\".\"players\".\"position_z\""},
	Yaw:           whereHelperfloat64{field: "\"cncraft\".\"players\".\"yaw\""},
	Pitch:         whereHelperfloat64{field: "\"cncraft\".\"players\".\"pitch\""},
	OnGround:      whereHelperbool{field: "\"cncraft\".\"players\".\"on_ground\""},
	CurrentHotbar: whereHelperint16{field: "\"cncraft\".\"players\".\"current_hotbar\""},
	CreatedAt:     whereHelpertime_Time{field: "\"cncraft\".\"players\".\"created_at\""},
}

// PlayerRels is where relationship names are stored.
var PlayerRels = struct {
	Inventories string
}{
	Inventories: "Inventories",
}

// playerR is where relationships are stored.
type playerR struct {
	Inventories InventorySlice `boil:"Inventories" json:"Inventories" toml:"Inventories" yaml:"Inventories"`
}

// NewStruct creates a new relationship struct
func (*playerR) NewStruct() *playerR {
	return &playerR{}
}

// playerL is where Load methods for each relationship are stored.
type playerL struct{}

var (
	playerAllColumns            = []string{"id", "conn_id", "dimension_id", "username", "position_x", "position_y", "position_z", "yaw", "pitch", "on_ground", "current_hotbar", "created_at"}
	playerColumnsWithoutDefault = []string{"id", "conn_id", "dimension_id", "username", "position_x", "position_y", "position_z", "yaw", "pitch", "created_at"}
	playerColumnsWithDefault    = []string{"on_ground", "current_hotbar"}
	playerPrimaryKeyColumns     = []string{"id"}
)

type (
	// PlayerSlice is an alias for a slice of pointers to Player.
	// This should generally be used opposed to []Player.
	PlayerSlice []*Player

	playerQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	playerType                 = reflect.TypeOf(&Player{})
	playerMapping              = queries.MakeStructMapping(playerType)
	playerPrimaryKeyMapping, _ = queries.BindMapping(playerType, playerMapping, playerPrimaryKeyColumns)
	playerInsertCacheMut       sync.RWMutex
	playerInsertCache          = make(map[string]insertCache)
	playerUpdateCacheMut       sync.RWMutex
	playerUpdateCache          = make(map[string]updateCache)
	playerUpsertCacheMut       sync.RWMutex
	playerUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single player record from the query.
func (q playerQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Player, error) {
	o := &Player{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "orm: failed to execute a one query for players")
	}

	return o, nil
}

// All returns all Player records from the query.
func (q playerQuery) All(ctx context.Context, exec boil.ContextExecutor) (PlayerSlice, error) {
	var o []*Player

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "orm: failed to assign all query results to Player slice")
	}

	return o, nil
}

// Count returns the count of all Player records in the query.
func (q playerQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to count players rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q playerQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "orm: failed to check if players exists")
	}

	return count > 0, nil
}

// Inventories retrieves all the inventory's Inventories with an executor.
func (o *Player) Inventories(mods ...qm.QueryMod) inventoryQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"cncraft\".\"inventory\".\"player_id\"=?", o.ID),
	)

	query := Inventories(queryMods...)
	queries.SetFrom(query.Query, "\"cncraft\".\"inventory\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"cncraft\".\"inventory\".*"})
	}

	return query
}

// LoadInventories allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (playerL) LoadInventories(ctx context.Context, e boil.ContextExecutor, singular bool, maybePlayer interface{}, mods queries.Applicator) error {
	var slice []*Player
	var object *Player

	if singular {
		object = maybePlayer.(*Player)
	} else {
		slice = *maybePlayer.(*[]*Player)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &playerR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &playerR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`cncraft.inventory`),
		qm.WhereIn(`cncraft.inventory.player_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load inventory")
	}

	var resultSlice []*Inventory
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice inventory")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on inventory")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for inventory")
	}

	if singular {
		object.R.Inventories = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &inventoryR{}
			}
			foreign.R.Player = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.PlayerID) {
				local.R.Inventories = append(local.R.Inventories, foreign)
				if foreign.R == nil {
					foreign.R = &inventoryR{}
				}
				foreign.R.Player = local
				break
			}
		}
	}

	return nil
}

// AddInventories adds the given related objects to the existing relationships
// of the player, optionally inserting them as new records.
// Appends related to o.R.Inventories.
// Sets related.R.Player appropriately.
func (o *Player) AddInventories(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Inventory) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.PlayerID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"cncraft\".\"inventory\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"player_id"}),
				strmangle.WhereClause("\"", "\"", 2, inventoryPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.PlayerID, rel.SlotNumber}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.PlayerID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &playerR{
			Inventories: related,
		}
	} else {
		o.R.Inventories = append(o.R.Inventories, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &inventoryR{
				Player: o,
			}
		} else {
			rel.R.Player = o
		}
	}
	return nil
}

// Players retrieves all the records using an executor.
func Players(mods ...qm.QueryMod) playerQuery {
	mods = append(mods, qm.From("\"cncraft\".\"players\""))
	return playerQuery{NewQuery(mods...)}
}

// FindPlayer retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPlayer(ctx context.Context, exec boil.ContextExecutor, iD uuid.UUID, selectCols ...string) (*Player, error) {
	playerObj := &Player{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"cncraft\".\"players\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, playerObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "orm: unable to select from players")
	}

	return playerObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Player) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("orm: no players provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(playerColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	playerInsertCacheMut.RLock()
	cache, cached := playerInsertCache[key]
	playerInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			playerAllColumns,
			playerColumnsWithDefault,
			playerColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(playerType, playerMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(playerType, playerMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"cncraft\".\"players\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"cncraft\".\"players\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "orm: unable to insert into players")
	}

	if !cached {
		playerInsertCacheMut.Lock()
		playerInsertCache[key] = cache
		playerInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Player.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Player) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	playerUpdateCacheMut.RLock()
	cache, cached := playerUpdateCache[key]
	playerUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			playerAllColumns,
			playerPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("orm: unable to update players, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"cncraft\".\"players\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, playerPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(playerType, playerMapping, append(wl, playerPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "orm: unable to update players row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by update for players")
	}

	if !cached {
		playerUpdateCacheMut.Lock()
		playerUpdateCache[key] = cache
		playerUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q playerQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to update all for players")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to retrieve rows affected for players")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PlayerSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), playerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"cncraft\".\"players\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, playerPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to update all in player slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to retrieve rows affected all in update all player")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Player) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("orm: no players provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(playerColumnsWithDefault, o)

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

	playerUpsertCacheMut.RLock()
	cache, cached := playerUpsertCache[key]
	playerUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			playerAllColumns,
			playerColumnsWithDefault,
			playerColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			playerAllColumns,
			playerPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("orm: unable to upsert players, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(playerPrimaryKeyColumns))
			copy(conflict, playerPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"cncraft\".\"players\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(playerType, playerMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(playerType, playerMapping, ret)
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
		return errors.Wrap(err, "orm: unable to upsert players")
	}

	if !cached {
		playerUpsertCacheMut.Lock()
		playerUpsertCache[key] = cache
		playerUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Player record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Player) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("orm: no Player provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), playerPrimaryKeyMapping)
	sql := "DELETE FROM \"cncraft\".\"players\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete from players")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by delete for players")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q playerQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("orm: no playerQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete all from players")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by deleteall for players")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PlayerSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), playerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"cncraft\".\"players\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, playerPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "orm: unable to delete all from player slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "orm: failed to get rows affected by deleteall for players")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Player) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPlayer(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PlayerSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PlayerSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), playerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"cncraft\".\"players\".* FROM \"cncraft\".\"players\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, playerPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "orm: unable to reload all in PlayerSlice")
	}

	*o = slice

	return nil
}

// PlayerExists checks if the Player row exists.
func PlayerExists(ctx context.Context, exec boil.ContextExecutor, iD uuid.UUID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"cncraft\".\"players\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "orm: unable to check if players exists")
	}

	return exists, nil
}
