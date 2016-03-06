package physicssys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brynbellomy/gl4-game/entity"
)

// points must be specified in clockwise order
type BoundingBox [4]mgl32.Vec2

type Collision struct {
	EntityA, EntityB entity.ID
}
