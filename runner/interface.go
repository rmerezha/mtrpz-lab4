package runner

import "github.com/rmerezha/mtrpz-lab4/config"

type Runner interface {
	Run(container config.Container) error
	Stop(name string) error
	Kill(name string) error
	Restart(name string) error
	Remove(name string) error
}
