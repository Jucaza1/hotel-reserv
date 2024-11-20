package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePatchUser(c *fiber.Ctx) error {
	var (
		userID      = c.Params("id")
		update      types.UpdateUser
		updateValid map[string]string
	)
	if len(userID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	if err := c.BodyParser(&update); err != nil {
		return types.ErrInvalidParams(err)
	}
	updateValid, errors := types.ValidateUserUpdate(update)
	if len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	h.userStore.UpdateUser(c.Context(), userID, updateValid)
	return c.JSON(types.MsgUpdated{Updated: userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if len(userID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(types.MsgDeleted{Deleted: userID})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return types.ErrInvalidParams(err)
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
func (h *UserHandler) HandlePostAdminUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return types.ErrInvalidParams(err)
	}
	user.IsAdmin = true
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetMyUser(c *fiber.Ctx) error {
	userID := c.Context().UserValue("user").(types.User).ID
	user, err := h.userStore.GetUserByID(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
func (h *UserHandler) HandlePatchMyUser(c *fiber.Ctx) error {
	var (
		userID      = c.Context().UserValue("user").(types.User).ID
		update      types.UpdateUser
		updateValid map[string]string
	)
	if err := c.BodyParser(&update); err != nil {
		return types.ErrInvalidParams(err)
	}
	updateValid, errors := types.ValidateUserUpdate(update)
	if len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	h.userStore.UpdateUser(c.Context(), userID, updateValid)
	return c.JSON(types.MsgUpdated{Updated: userID})
}
