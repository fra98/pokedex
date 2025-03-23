package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		err := httperror.NewHTTPError("unable to retrieve pokemon info", http.StatusNotFound)
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
		err := httperror.NewHTTPError("unable to retrieve translated pokemon info", http.StatusNotFound)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}
