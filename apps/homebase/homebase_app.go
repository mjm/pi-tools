package homebase

import (
	"context"
	"fmt"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

type App struct {
	name string
}

func New(name string) *App {
	return &App{
		name: name,
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Install(ctx context.Context, clients nomadic.Clients) error {
	if err := a.installBotVaultPolicy(ctx, clients); err != nil {
		return err
	}

	if err := a.installConfigEntries(ctx, clients); err != nil {
		return err
	}

	job := nomadic.NewJob(a.name, 50)
	a.addWebTaskGroup(job)
	a.addAPITaskGroup(job)
	a.addBotTaskGroup(job)

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := a.uninstallBotVaultPolicy(ctx, clients); err != nil {
		return err
	}

	if err := a.uninstallConfigEntries(ctx, clients); err != nil {
		return err
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

func (a *App) installConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	if err := a.installWebConfigEntries(ctx, clients); err != nil {
		return err
	}

	if err := a.installAPIConfigEntries(ctx, clients); err != nil {
		return err
	}

	if err := a.installBotConfigEntries(ctx, clients); err != nil {
		return err
	}

	return nil
}

func (a *App) uninstallConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	if err := a.uninstallWebConfigEntries(ctx, clients); err != nil {
		return err
	}

	if err := a.uninstallAPIConfigEntries(ctx, clients); err != nil {
		return err
	}

	if err := a.uninstallBotConfigEntries(ctx, clients); err != nil {
		return err
	}

	return nil
}
