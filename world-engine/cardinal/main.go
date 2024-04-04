package main

import (
	"errors"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal"

	"cinco-paus/component"
	"cinco-paus/msg"
	"cinco-paus/system"
)

func main() {
	w, err := cardinal.NewWorld(cardinal.WithDisableSignatureVerification())
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	// Register components
	// NOTE: You must register your components here for it to be accessible.
	Must(
		cardinal.RegisterComponent[component.PendingGame](w),
		cardinal.RegisterComponent[component.Game](w),

		cardinal.RegisterComponent[component.Collidable](w),
		cardinal.RegisterComponent[component.Player](w),
		cardinal.RegisterComponent[component.Monster](w),
		cardinal.RegisterComponent[component.Wall](w),
		cardinal.RegisterComponent[component.WandCore](w),
		cardinal.RegisterComponent[component.Available](w),
		cardinal.RegisterComponent[component.Spell](w),
		cardinal.RegisterComponent[component.AwaitingReveal](w),
		cardinal.RegisterComponent[component.Health](w),
		cardinal.RegisterComponent[component.Position](w),
	)

	// Register messages (user action)
	// NOTE: You must register your transactions here for it to be executed.
	Must(
		// cardinal.RegisterMessage[msg.CreatePlayerMsg, msg.CreatePlayerResult](w, "create-player"),
		// cardinal.RegisterMessage[msg.AttackPlayerMsg, msg.AttackPlayerMsgReply](w, "attack-player"),
		cardinal.RegisterMessage[msg.RequestGameMsg, msg.RequestGameMsgResult](w, "request-game"),
		cardinal.RegisterMessage[msg.FulfillCreateGameMsg, msg.FulfillCreateGameMsgResult](w, "fulfill-create-game"),
		cardinal.RegisterMessage[msg.FulfillCastMsg, msg.FulfillCastMsgResult](w, "fulfill-cast"),
		cardinal.RegisterMessage[msg.PlayerTurnMsg, msg.PlayerTurnResult](w, "player-turn"),
	)

	// Register queries
	// NOTE: You must register your queries here for it to be accessible.
	Must(
	// cardinal.RegisterQuery[query.PlayerHealthRequest, query.PlayerHealthResponse](w, "player-health", query.PlayerHealth),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		system.RequestGameSystem,
		system.FulfillCreateGameSystem,
		system.FulfillCastSystem,
		system.PlayerTurnSystem,
		// system.MonsterTurnSystem,
	))

	Must(cardinal.RegisterInitSystems(w,
		system.SpawnPlayerSystem,
		system.SpawnWandsSystem,
		system.PopulateBoardSystem,
	))

	seismicClient := system.Initialize(w)
	seismicClient.Start()

	Must(w.StartGame())
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
