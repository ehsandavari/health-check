package genericRepository

import (
	"fmt"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"gorm.io/gorm"
	"health-check/application/common"
	"health-check/infrastructure/postgres"
	"health-check/pkg/tracer"
)

type IGenericRepository[TE any] interface {
	Paginate(ctx *contextplus.Context, listQuery common.PaginateQuery) (int64, []TE, error)
	Count(ctx *contextplus.Context, specifications ...Specification) (int64, error)
	Exists(ctx *contextplus.Context, specifications ...Specification) (bool, error)
	All(ctx *contextplus.Context, specifications ...Specification) ([]TE, error)
	Find(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	First(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	FirstOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	Single(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	SingleOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	Last(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	LastOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error)
	CreateOrUpdate(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error)
	Create(ctx *contextplus.Context, entity *TE) error
	Creates(ctx *contextplus.Context, entity ...TE) ([]TE, error)
	Update(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error)
	UpdateColumn(ctx *contextplus.Context, column string, value any, specifications ...Specification) (*TE, error)
	Delete(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error)
}

type sGenericRepository[TE any] struct {
	logger   logger.ILogger
	tracer   tracer.ITracer
	postgres postgres.SPostgres
}

func NewGenericRepository[TE any](logger logger.ILogger, tracer tracer.ITracer, postgres postgres.SPostgres) IGenericRepository[TE] {
	return sGenericRepository[TE]{
		logger:   logger,
		tracer:   tracer,
		postgres: postgres,
	}
}

func (r sGenericRepository[TE]) Specification(ctx *contextplus.Context, specifications ...Specification) *gorm.DB {
	dbPreWarm := r.postgres.Database.WithContext(ctx)
	for _, s := range specifications {
		dbPreWarm = dbPreWarm.Where(s.GetQuery(), s.GetValues()...)
	}
	return dbPreWarm
}

func (r sGenericRepository[TE]) Paginate(ctx *contextplus.Context, paginateQuery common.PaginateQuery) (int64, []TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	query := r.postgres.Database.WithContext(ctx)
	if paginateQuery.Filters != nil {
		for _, filter := range paginateQuery.Filters {
			query = query.Where(fmt.Sprintf("%s %s", filter.Key, filter.Comparison), filter.Value)
		}
	}

	var entityObjects []TE

	var totalRows int64
	query.Model(entityObjects).Count(&totalRows)

	result := query.Offset(int(paginateQuery.GetOffset())).Limit(int(paginateQuery.GetLimit())).Order(paginateQuery.GetOrderBy()).Find(&entityObjects)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return 0, nil, result.Error
	}

	return totalRows, entityObjects, nil
}

func (r sGenericRepository[TE]) Count(ctx *contextplus.Context, specifications ...Specification) (int64, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var totalRows int64
	var entity TE
	result := r.Specification(ctx, specifications...).Model(&entity).Count(&totalRows)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return 0, result.Error
	}
	return totalRows, nil
}

func (r sGenericRepository[TE]) Exists(ctx *contextplus.Context, specifications ...Specification) (bool, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var totalRows int64
	var entity TE
	result := r.Specification(ctx, specifications...).Model(&entity).Count(&totalRows)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return false, result.Error
	}
	if totalRows == 0 {
		return false, nil
	}
	return true, nil
}

func (r sGenericRepository[TE]) All(ctx *contextplus.Context, specifications ...Specification) ([]TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entityObjects []TE
	result := r.Specification(ctx, specifications...).Find(&entityObjects)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entityObjects, nil
}

func (r sGenericRepository[TE]) Find(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Limit(1).Find(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (r sGenericRepository[TE]) Single(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Find(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var totalRows int64
	result.Count(&totalRows)
	if totalRows > 1 {
		return nil, ErrorMultipleRowsReturned
	}
	return entity, nil
}

func (r sGenericRepository[TE]) SingleOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Find(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	var totalRows int64
	result.Count(&totalRows)
	if totalRows > 1 {
		return nil, ErrorMultipleRowsReturned
	}
	return entity, nil
}

func (r sGenericRepository[TE]) First(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Order("created_at ASC").Take(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entity, nil
}

func (r sGenericRepository[TE]) FirstOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Order("created_at ASC").Limit(1).Find(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (r sGenericRepository[TE]) Last(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Order("created_at DESC").Take(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entity, nil
}

func (r sGenericRepository[TE]) LastOrDefault(ctx *contextplus.Context, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Order("created_at DESC").Limit(1).Find(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (r sGenericRepository[TE]) CreateOrUpdate(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	count, err := r.Count(ctx, specifications...)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("err", err)

		return nil, err
	}
	if count == 0 {
		result := r.postgres.Database.WithContext(ctx).Create(&entity)
		if result.Error != nil {
			span.SetTag("error", true)
			span.LogKV("err", result.Error)

			return nil, result.Error
		}
		return entity, nil
	}
	if count == 1 {
		result := r.Specification(ctx, specifications...).Updates(&entity)
		if result.Error != nil {
			span.SetTag("error", true)
			span.LogKV("err", result.Error)

			return nil, result.Error
		}
		return entity, nil
	}
	span.SetTag("error", true)
	span.LogKV("err", ErrorMultipleRowsReturned)

	return nil, ErrorMultipleRowsReturned
}

func (r sGenericRepository[TE]) Create(ctx *contextplus.Context, entity *TE) error {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	result := r.postgres.Database.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return result.Error
	}
	return nil
}

func (r sGenericRepository[TE]) Creates(ctx *contextplus.Context, entities ...TE) ([]TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	result := r.postgres.Database.WithContext(ctx).CreateInBatches(&entities, 5)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entities, nil
}

func (r sGenericRepository[TE]) Update(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	result := r.Specification(ctx, specifications...).Updates(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entity, nil
}

func (r sGenericRepository[TE]) UpdateColumn(ctx *contextplus.Context, column string, value any, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	var entity *TE
	result := r.Specification(ctx, specifications...).Model(&entity).Update(column, value)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entity, nil
}

func (r sGenericRepository[TE]) Delete(ctx *contextplus.Context, entity *TE, specifications ...Specification) (*TE, error) {
	span, ctx := r.tracer.SpanFromContext(ctx)
	defer span.Finish()

	result := r.Specification(ctx, specifications...).Delete(&entity)
	if result.Error != nil {
		span.SetTag("error", true)
		span.LogKV("err", result.Error)

		return nil, result.Error
	}
	return entity, nil
}
