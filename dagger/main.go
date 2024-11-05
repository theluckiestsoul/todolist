package main

import (
	"context"
	"dagger/todolist/internal/dagger"
)

type Todolist struct{}

// Build compiles the Todolist server binary in a golang:1 container.
func (m *Todolist) Build(src *dagger.Directory) *dagger.File {
	return dag.Container().
		From("golang:1").
		WithDirectory("/app", src).
		WithWorkdir("/app").
		WithEnvVariable("CGO_ENABLED", "0").
		//WithEnvVariable("GOARCH", arch). // Set the architecture
		WithExec([]string{"go", "build", "-o", "todolist", "./cmd/server"}).
		File("/app/todolist")
}

// Test runs the tests in a golang:1 container.
func (m *Todolist) Test(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1").
		WithDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"go", "test", "-v", "./..."})
}

// Lint runs golangci-lint in a golang:1 container.
func (m *Todolist) Lint(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1").
		WithDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})
}

func (m *Todolist) OsInfo(ctx context.Context, ctr *dagger.Container) (string, error) {
	return ctr.
		WithExec([]string{"uname", "-a"}).
		Stdout(ctx)
}

func (m *Todolist) GoBuilder(ctx context.Context, src *dagger.Directory, arch string, os string) *dagger.Directory {
	return dag.Container().
		From("golang:1.21").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "build/"}).
		Directory("/src/build")
}

func (m *Todolist) HttpService() *dagger.Service {
	return dag.Container().
		From("python").
		WithWorkdir("/srv").
		WithNewFile("index.html", "Hello, world!").
		WithExec([]string{"python", "-m", "http.server", "8080"}).
		WithExposedPort(8080).
		AsService()
}

func (m *Todolist) Get(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine").
		WithServiceBinding("www", m.HttpService()).
		WithExec([]string{"wget", "-O-", "http://www:8080"}).
		Stdout(ctx)
}
