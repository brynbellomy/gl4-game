package gameobjsys

import "math"

func DirectionFromRadians(radians float64) Direction {
	var direction Direction
	switch {
	case 0.25*math.Pi < radians && radians < 0.75*math.Pi:
		direction = Down

	case -0.25*math.Pi > radians && radians > -0.75*math.Pi:
		direction = Up

	case (0.75*math.Pi <= radians && radians <= math.Pi) ||
		(-0.75*math.Pi >= radians && radians >= -math.Pi):
		direction = Right

	case (0.25*math.Pi >= radians && radians >= 0) ||
		(-0.25*math.Pi <= radians && radians <= 0):
		direction = Left
	}

	return direction
}
