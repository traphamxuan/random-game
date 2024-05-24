package servicemanager

import (
	"context"
	"fmt"
	"reflect"

	"github.com/traphamxuan/random-game/app/utils"
)

type IService interface {
	Setup(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type ServiceStatus int

const (
	StatusInit  ServiceStatus = 0
	StatusSetup ServiceStatus = 1
	StatusStart ServiceStatus = 2
	StatusStop  ServiceStatus = 3
)

type serviceBlock struct {
	service IService
	name    string
	status  ServiceStatus
}

type ServiceManager struct {
	services []serviceBlock
	keys     map[string]*serviceBlock
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		keys: make(map[string]*serviceBlock),
	}
}

func (sm *ServiceManager) Add(s IService, args ...string) {
	name := reflect.TypeOf(s).Name()
	key := name
	if len(args) > 0 {
		key = args[0]
	}
	sBlock := serviceBlock{
		service: s,
		name:    name,
		status:  StatusInit,
	}
	sm.keys[key] = &sBlock
	sm.services = append(sm.services, sBlock)
}

func GetService[T IService](sm *ServiceManager) (r T) {
	for _, sb := range sm.services {
		if v, ok := sb.service.(T); ok {
			return v
		}
	}
	return r
}

func GetServiceByKey[T IService](sm *ServiceManager, key string) (r T) {
	if sb, ok := sm.keys[key]; ok {
		if v, ok := sb.service.(T); ok {
			return v
		}
	}
	return r
}

func GetServiceOrPanic[T IService](sm *ServiceManager, args ...string) (r T) {
	arg_len := len(args)
	var (
		srcName string
		dstName string
	)
	if arg_len == 0 || args[0] == "" {
		srcName = ""
		dstName = ""
		r = GetService[T](sm)
	} else {
		srcName = args[0]
		if arg_len > 1 {
			dstName = args[1]
		}
		r = GetServiceByKey[T](sm, args[0])
	}

	if utils.IsNil(r) {
		panic(fmt.Sprintf("service %s not found. Required by %s", srcName, dstName))
	}
	return r
}

func (sm *ServiceManager) Setup(
	ctx context.Context,
	onSetup *func(ctx context.Context, s IService, key string, err error),
	priorityOrders ...string,
) error {
	for _, key := range priorityOrders {
		if sb, ok := sm.keys[key]; ok {
			if err := setup(ctx, sb, onSetup); err != nil {
				return err
			}
		}
	}

	for id := range sm.services {
		sb := &sm.services[id]
		if err := setup(ctx, sb, onSetup); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) Start(
	ctx context.Context,
	onStart *func(ctx context.Context, s IService, key string, err error),
	priorityOrders ...string,
) error {

	for _, key := range priorityOrders {
		if sb, ok := sm.keys[key]; ok {
			if err := start(ctx, sb, onStart); err != nil {
				return err
			}
		}
	}

	for id := range sm.services {
		sb := &sm.services[id]
		if err := start(ctx, sb, onStart); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) Stop(
	ctx context.Context,
	onStop *func(ctx context.Context, s IService, key string, err error),
	priorityOrders ...string,
) error {
	for _, key := range priorityOrders {
		if sb, ok := sm.keys[key]; ok {
			if err := stop(ctx, sb, onStop); err != nil {
				return err
			}
		}
	}

	for id := len(sm.services) - 1; id >= 0; id-- {
		sb := &sm.services[id]
		if err := stop(ctx, sb, onStop); err != nil {
			return err
		}
	}
	return nil
}

func setup(ctx context.Context, sb *serviceBlock,
	onSetup *func(ctx context.Context, s IService, key string, err error),
) error {
	if sb.status >= StatusSetup {
		return nil
	}
	err := sb.service.Setup(ctx)
	(*onSetup)(ctx, sb.service, sb.name, err)
	if err != nil {
		return err
	}
	return nil
}

func start(ctx context.Context, sb *serviceBlock,
	onStart *func(ctx context.Context, s IService, key string, err error),
) error {
	if sb.status >= StatusStart {
		return nil
	}
	err := sb.service.Start(ctx)
	(*onStart)(ctx, sb.service, sb.name, err)
	if err != nil {
		return err
	}
	return nil
}

func stop(ctx context.Context, sb *serviceBlock,
	onStop *func(ctx context.Context, s IService, key string, err error),
) error {
	if sb.status >= StatusStop {
		return nil
	}
	err := sb.service.Stop(ctx)
	(*onStop)(ctx, sb.service, sb.name, err)
	if err != nil {
		return err
	}
	return nil
}
