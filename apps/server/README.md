# game-plan-api

Project is ruled by service manager which implement DI container for all other components.

# Getting started
Start web API by command
```bash
go run cmd/api/main.go
```
# Structure
Project is divided into 4 parts
- **internal** contains most of business logic
- **package** contains configuration for 3rd party libraries
- **utilities** common logic trick in golang
- **command/scripts/deployments** for devops to maintain the product

# Code convention
All components have their own dependencies and life-cycle
- Initialization
- Setup/Configuration
- Start/Run
- Stop

If you want to add this life-cycle setting for a instances, please implement `gobs.IService`
```go
type Product struct {
	log *logger.Logger
	orm *orm.Orm
	s3  *s3.S3
}

var _ gobs.IService = (*Product)(nil)

func (p *Product) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&orm.Orm{},
		&s3.S3{},
	}
	onSetup := func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		p.log = deps[0].(*logger.Logger)
		p.orm = deps[1].(*orm.Orm)
		p.s3 = deps[2].(*s3.S3)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}
```
then put this instance to the main thread at init step. All other components required this instance will find this instance with the order manner.
```go
sm := *gobs.Bootstrap
sm.AddDefault(&service.Product{}, "service.Product")
```
