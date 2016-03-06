package mainscene

type Collider uint64

const (
	HeroCollider Collider = 1 << iota
	EnemyCollider
	ProjectileCollider
)
