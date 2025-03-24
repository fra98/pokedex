package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/fra98/pokedex/pkg/errors"
	"github.com/fra98/pokedex/pkg/server/httperror"
	"github.com/fra98/pokedex/pkg/service"
)

// PokemonHandler handles the Pokemon API endpoints.
type PokemonHandler struct {
	pokemonService service.Pokemon
}

// NewPokemonHandler creates a new PokemonHandler with the given PokemonService.
func NewPokemonHandler(pokemonService service.Pokemon) *PokemonHandler {
	return &PokemonHandler{pokemonService: pokemonService}
}

// GetPokemon returns the information of a Pokemon given its name.
func (h *PokemonHandler) GetPokemon(c *gin.Context) {
	name := c.Param("name")

	pokemon, err := h.pokemonService.GetPokemonInfo(c.Request.Context(), name)
	if err != nil {
		err := httperror.NewHTTPError("unable to retrieve pokemon info", getStatusCode(err))
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

// GetTranslatedPokemon returns the information of a Pokemon given its name with a translated description.
func (h *PokemonHandler) GetTranslatedPokemon(c *gin.Context) {
	name := c.Param("name")

	pokemon, err := h.pokemonService.GetTranslatedPokemonInfo(c.Request.Context(), name)
	if err != nil {
		err := httperror.NewHTTPError("unable to retrieve translated pokemon info", getStatusCode(err))
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrResourceNotFound):
		return http.StatusNotFound
	default:
		return http.StatusServiceUnavailable
	}
	// Internal server errors are handled by the ErrorHandler middleware
}
