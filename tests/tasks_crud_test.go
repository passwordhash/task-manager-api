package tests

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/ilyakaznacheev/cleanenv"
)

type testConfig struct {
	Port       string        `env:"PORT" env-default:"8080"`
	IoDuration time.Duration `env:"TASK_DURATION" env-required:"true"`
	Workers    int           `env:"WORKERS" env-required:"true"`
}

var (
	cfg testConfig
	u   url.URL
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Failed to read environment variables: " + err.Error())
	}
	u = url.URL{
		Scheme: "http",
		Host:   "localhost:" + cfg.Port,
	}
}

type createTaskResp struct {
	TaskUUID string `json:"task_uuid"`
}

func TestCreateTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	taskUUID := createTask(e)

	statusTask(e, taskUUID).
		Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "running")
	t.Cleanup(func() {
		cancelTask(e, taskUUID)
	})
}

func TestCreateTaskOverLimit(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	tasks := make([]string, cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		taskUUID := createTask(e)
		tasks[i] = taskUUID
	}
	t.Cleanup(func() {
		for _, taskUUID := range tasks {
			cancelTask(e, taskUUID).
				Expect().Status(http.StatusOK)
			statusTask(e, taskUUID).
				Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "canceled")
		}
	})

	// Attempt to create one more task than the limit
	taskUUID := createTask(e)
	e.GET("/api/v1/tasks/"+taskUUID+"/status").
		Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "pending")
}

func TestCancelTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	taskUUID := createTask(e)

	cancelTask(e, taskUUID).
		Expect().Status(http.StatusOK)
	statusTask(e, taskUUID).
		Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "canceled")
}

func TestCancelCanceledTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	taskUUID := createTask(e)

	cancelTask(e, taskUUID).
		Expect().Status(http.StatusOK)

	cancelTask(e, taskUUID).
		Expect().Status(http.StatusConflict)
}

func TestCancelNonExistentTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	cancelTask(e, "non-existing-uuid").
		Expect().Status(http.StatusNotFound)
}

func TestStatusOfCanceledTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	taskUUID := createTask(e)

	cancelTask(e, taskUUID).
		Expect().Status(http.StatusOK)

	cancelTask(e, taskUUID).
		Expect().Status(http.StatusConflict)
}

func TestStatusOfCompletedTask(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	taskUUID := createTask(e)

	e.GET("/api/v1/tasks/"+taskUUID+"/status").
		Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "running")

	time.Sleep(cfg.IoDuration)

	e.GET("/api/v1/tasks/"+taskUUID+"/status").
		Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "completed")
}

func createTask(e *httpexpect.Expect) string {
	var createResp createTaskResp

	e.POST("/api/v1/tasks/").
		Expect().Status(http.StatusOK).JSON().Object().
		ContainsKey("task_uuid").Decode(&createResp)

	return createResp.TaskUUID
}

func statusTask(e *httpexpect.Expect, taskUUID string) *httpexpect.Request {
	return e.GET("/api/v1/tasks/" + taskUUID + "/status")
}

func cancelTask(e *httpexpect.Expect, taskUUID string) *httpexpect.Request {
	return e.POST("/api/v1/tasks/" + taskUUID + "/cancel")
}
