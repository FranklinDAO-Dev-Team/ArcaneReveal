package query

import (
	comp "cinco-paus/component"
	"log"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type QueryGameIDByPersonaRequest struct {
	Persona string
}

type QueryGameIDByPersonaResponse struct {
	Success bool
	GameID  types.EntityID
}

func QueryGameIDByPersona(world cardinal.WorldContext, req *QueryGameIDByPersonaRequest) (*QueryGameIDByPersonaResponse, error) {
	log.Println("QueryGameByPersona() querying data for", req.Persona)
	var outsideErr error
	var resp *QueryGameIDByPersonaResponse = &QueryGameIDByPersonaResponse{
		Success: false,
	}

	searchErr := cardinal.NewSearch(
		world,
		filter.Contains(comp.Game{})).
		Each(func(id types.EntityID) bool {
			game, err := cardinal.GetComponent[comp.Game](world, id)
			if err != nil {
				log.Println("gameObj err: ", err)
				return false
			}

			if game.PersonaTag == req.Persona {
				resp = &QueryGameIDByPersonaResponse{
					Success: true,
					GameID:  id,
				}
				return false
			}

			return true
		})
	if searchErr != nil {
		return nil, searchErr
	}
	if outsideErr != nil {
		return nil, outsideErr
	}
	return resp, nil
}
