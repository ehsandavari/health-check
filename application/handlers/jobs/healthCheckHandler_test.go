package jobs

import (
	"errors"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"health-check/application/interfaces"
	"health-check/domain/entities"
	"health-check/domain/enums"
	"health-check/pkg/genericRepository"
	"health-check/pkg/tracer"
	"testing"
	"time"
)

type sMockHealthCheckJobHandler struct {
	iLogger                       *logger.MockILogger
	iTracer                       *tracer.MockITracer
	iSpan                         *tracer.MockISpan
	iRedis                        *interfaces.MockIRedis
	iCron                         *interfaces.MockICron
	iRest                         *interfaces.MockIRest
	iNotification                 *interfaces.MockINotification
	iHealthCheckRepository        *interfaces.MockIHealthCheckRepository
	iHealthCheckRequestRepository *interfaces.MockIHealthCheckRequestRepository
	iUnitOfWork                   *interfaces.MockIUnitOfWork

	callAddJob              func(ctx *contextplus.Context, healthCheck entities.HealthCheck)
	callAddJobTimes         int
	callAddJobTimesExpected int

	callSubRedis              func(ctx *contextplus.Context)
	callSubRedisTimes         int
	callSubRedisTimesExpected int

	callSendRequest              func(ctx *contextplus.Context, healthCheck entities.HealthCheck)
	callSendRequestTimes         int
	callSendRequestTimesExpected int

	callSendNotification              func(ctx *contextplus.Context, subject string, msg string)
	callSendNotificationTimes         int
	callSendNotificationTimesExpected int
}

func setup(t *testing.T) (mock *sMockHealthCheckJobHandler) {
	mockController := gomock.NewController(t)
	mock = &sMockHealthCheckJobHandler{
		iLogger:                       logger.NewMockILogger(mockController),
		iTracer:                       tracer.NewMockITracer(mockController),
		iSpan:                         tracer.NewMockISpan(mockController),
		iRedis:                        interfaces.NewMockIRedis(mockController),
		iCron:                         interfaces.NewMockICron(mockController),
		iRest:                         interfaces.NewMockIRest(mockController),
		iNotification:                 interfaces.NewMockINotification(mockController),
		iHealthCheckRepository:        interfaces.NewMockIHealthCheckRepository(mockController),
		iHealthCheckRequestRepository: interfaces.NewMockIHealthCheckRequestRepository(mockController),
		iUnitOfWork:                   interfaces.NewMockIUnitOfWork(mockController),
	}
	t.Cleanup(func() {
		mock.callAddJob = nil
		mock.callAddJobTimes = 0
		mock.callAddJobTimesExpected = 0
		mock.callSubRedis = nil
		mock.callSubRedisTimes = 0
		mock.callSubRedisTimesExpected = 0
		mock.callSendRequest = nil
		mock.callSendRequestTimes = 0
		mock.callSendRequestTimesExpected = 0
		mock.callSendNotification = nil
		mock.callSendNotificationTimes = 0
		mock.callSendNotificationTimesExpected = 0
		mockController.Finish()
	})
	return mock
}

func TestStart(t *testing.T) {
	type (
		sIn struct {
			ctx *contextplus.Context
		}
		sOut struct {
			err error
		}
		sArg struct {
			in  sIn
			out sOut
		}
		sTableTest struct {
			name   string
			arg    sArg
			mock   func(mock *sMockHealthCheckJobHandler, arg sIn)
			assert func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut)
		}
	)

	tableTests := []sTableTest{
		{
			name: "error in r.iUnitOfWork.HealthCheckRepository().All",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				mock.iUnitOfWork.EXPECT().HealthCheckRepository().Return(mock.iHealthCheckRepository).Times(1)
				err := errors.New("error in all")
				mock.iHealthCheckRepository.EXPECT().All(arg.ctx, genericRepository.Equal("status", enums.StatusStart)).Return(nil, err).Times(1)

				mock.iSpan.EXPECT().SetTag("error", true).Times(1)
				mock.iSpan.EXPECT().LogKV("err", err).Times(1)

				mock.iLogger.EXPECT().WithError(err).Return(mock.iLogger).Times(1)
				mock.iLogger.EXPECT().Error(arg.ctx, "error in get all health checks").Times(1)
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
				assert.Error(t, arg.err)
			},
		},
		{
			name: "check all func is run",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				healthChecks := []entities.HealthCheck{
					{
						Id: 1,
					},
					{
						Id: 2,
					},
				}
				mock.iUnitOfWork.EXPECT().HealthCheckRepository().Return(mock.iHealthCheckRepository).Times(1)
				mock.iHealthCheckRepository.EXPECT().All(arg.ctx, genericRepository.Equal("status", enums.StatusStart)).Return(healthChecks, nil).Times(1)

				mock.callAddJobTimesExpected = len(healthChecks)
				mock.callAddJob = func(ctx *contextplus.Context, healthCheck entities.HealthCheck) {
					mock.callAddJobTimes++
				}

				mock.callSubRedisTimesExpected = 1
				mock.callSubRedis = func(ctx *contextplus.Context) {
					mock.callSubRedisTimes++
				}
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
				assert.NoError(t, arg.err)
				assert.Equal(t, mock.callAddJobTimesExpected, mock.callAddJobTimes)
				time.Sleep(10 * time.Millisecond)
				assert.Equal(t, mock.callSubRedisTimesExpected, mock.callSubRedisTimes)
			},
		},
	}

	for _, tableTest := range tableTests {
		t.Run(tableTest.name, func(t *testing.T) {
			mock := setup(t)
			healthCheckJobHandler := newHealthCheckJobHandler(
				mock.iLogger,
				mock.iTracer,
				mock.iRedis,
				mock.iCron,
				mock.iRest,
				mock.iNotification,
				mock.iUnitOfWork,
			)
			tableTest.mock(mock, tableTest.arg.in)
			healthCheckJobHandler.callAddJob = mock.callAddJob
			healthCheckJobHandler.callSubRedis = mock.callSubRedis
			healthCheckJobHandler.callSendRequest = mock.callSendRequest
			healthCheckJobHandler.callSendNotification = mock.callSendNotification

			tableTest.arg.out.err = healthCheckJobHandler.Start(tableTest.arg.in.ctx)
			tableTest.assert(mock, t, tableTest.arg.out)
		})
	}
}

func TestAddJob(t *testing.T) {
	type (
		sIn struct {
			ctx         *contextplus.Context
			healthCheck entities.HealthCheck
		}
		sOut struct {
		}
		sArg struct {
			in  sIn
			out sOut
		}
		sTableTest struct {
			name   string
			arg    sArg
			mock   func(mock *sMockHealthCheckJobHandler, arg sIn)
			assert func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut)
		}
	)

	tableTests := []sTableTest{
		{
			name: "remove job when health check status is stop",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
					healthCheck: entities.HealthCheck{
						Id:     1,
						Status: enums.StatusStop,
					},
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				mock.iCron.EXPECT().RemoveJob(arg.healthCheck.Id).Times(1)
				mock.iCron.EXPECT().AddJob(arg.healthCheck.Id, arg.healthCheck.UpdatedAt, arg.healthCheck.Interval, func() {
					mock.callSendRequest = func(ctx *contextplus.Context, healthCheck entities.HealthCheck) {

					}
				}).Times(0)
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
			},
		},
		{
			name: "remove job when health check deleted at is valid",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
					healthCheck: entities.HealthCheck{
						Id: 1,
						Base3: entities.Base3{
							DeletedAt: gorm.DeletedAt{
								Time:  time.Now(),
								Valid: true,
							},
						},
					},
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				mock.iCron.EXPECT().RemoveJob(arg.healthCheck.Id).Times(1)
				mock.iCron.EXPECT().AddJob(arg.healthCheck.Id, arg.healthCheck.UpdatedAt, arg.healthCheck.Interval, func() {
					mock.callSendRequest = func(ctx *contextplus.Context, healthCheck entities.HealthCheck) {

					}
				}).Times(0)
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
			},
		},
		{
			name: "add job when health check status is start and deleted at is not valid",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
					healthCheck: entities.HealthCheck{
						Id:       1,
						Interval: "1s",
						Status:   enums.StatusStart,
						Base3: entities.Base3{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				mock.callSendRequestTimesExpected = 1
				mock.callSendRequest = func(ctx *contextplus.Context, healthCheck entities.HealthCheck) {
					mock.callSendRequestTimes++
				}

				mock.iCron.EXPECT().AddJob(arg.healthCheck.Id, arg.healthCheck.UpdatedAt, arg.healthCheck.Interval, gomock.AssignableToTypeOf(func() {})).
					DoAndReturn(func(key uint, createAt time.Time, interval string, job func()) error {
						job()
						return nil
					}).Times(1)
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
				assert.Equal(t, mock.callSendRequestTimesExpected, mock.callSendRequestTimes)
			},
		},
		{
			name: "add job return error when health check status is start and deleted at is not valid",
			arg: sArg{
				in: sIn{
					ctx: contextplus.Background(),
					healthCheck: entities.HealthCheck{
						Id:       1,
						Interval: "1s",
						Status:   enums.StatusStart,
						Base3: entities.Base3{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
				},
			},
			mock: func(mock *sMockHealthCheckJobHandler, arg sIn) {
				mock.iTracer.EXPECT().SpanFromContext(arg.ctx).Return(mock.iSpan, arg.ctx).Times(1)
				mock.iSpan.EXPECT().Finish().Times(1)

				mock.callSendRequestTimesExpected = 1
				mock.callSendRequest = func(ctx *contextplus.Context, healthCheck entities.HealthCheck) {
					mock.callSendRequestTimes++
				}

				err := errors.New("error in add job")
				mock.iCron.EXPECT().AddJob(arg.healthCheck.Id, arg.healthCheck.UpdatedAt, arg.healthCheck.Interval, gomock.AssignableToTypeOf(func() {})).
					DoAndReturn(func(key uint, createAt time.Time, interval string, job func()) error {
						job()
						return err
					}).Times(1)

				mock.iSpan.EXPECT().SetTag("error", true).Times(1)
				mock.iSpan.EXPECT().LogKV("err", err).Times(1)

				mock.iLogger.EXPECT().WithError(err).Return(mock.iLogger).Times(1)
				mock.iLogger.EXPECT().Error(arg.ctx, "error in add job").Times(1)
			},
			assert: func(mock *sMockHealthCheckJobHandler, t *testing.T, arg sOut) {
				assert.Equal(t, mock.callSendRequestTimesExpected, mock.callSendRequestTimes)
			},
		},
	}

	for _, tableTest := range tableTests {
		t.Run(tableTest.name, func(t *testing.T) {
			mock := setup(t)
			healthCheckJobHandler := newHealthCheckJobHandler(
				mock.iLogger,
				mock.iTracer,
				mock.iRedis,
				mock.iCron,
				mock.iRest,
				mock.iNotification,
				mock.iUnitOfWork,
			)
			tableTest.mock(mock, tableTest.arg.in)
			healthCheckJobHandler.callAddJob = mock.callAddJob
			healthCheckJobHandler.callSubRedis = mock.callSubRedis
			healthCheckJobHandler.callSendRequest = mock.callSendRequest
			healthCheckJobHandler.callSendNotification = mock.callSendNotification

			healthCheckJobHandler.addJob(tableTest.arg.in.ctx, tableTest.arg.in.healthCheck)
			tableTest.assert(mock, t, tableTest.arg.out)
		})
	}
}
