package dic

import (
	"testing"
)

type (
	mockHandler struct {
		message string
	}
)

func TestContainerUsage(t *testing.T) {
	c := NewContainer()
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		c.Set("mock", func(c *Container) interface{} {
			return &mockHandler{
				message: "Hello Carthago",
			}
		})

		value, _ := c.Get("mock")
		handler, _ := value.(*mockHandler)
		if "Hello Carthago" != handler.message {
			t.Error("Wrong message")
		}
	})

	t.Run("already-exists", func(t *testing.T) {
		errAlready := c.Set("mock", func(c *Container) interface{} {
			return nil
		})
		if errAlready != ErrDependencyAlreadyExists {
			t.Error("Duplicated assignment error not detected!")
		}
	})

	t.Run("delete", func(t *testing.T) {
		c.Set("will be removed", func(c *Container) interface{} {
			return 10
		})
		c.Delete("will be removed")

		_, errNotFound := c.Get("will be removed")
		if errNotFound != ErrDependencyNotFound {
			t.Error("Dependency still exists after deletion!")
		}
	})

	t.Run("not-found", func(t *testing.T) {
		_, errNotFound := c.Get("xpto")
		if errNotFound != ErrDependencyNotFound {
			t.Error("Invalid dependency key error not detected!")
		}
	})
}
