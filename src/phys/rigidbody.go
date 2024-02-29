package phys

type RigidBodyType int

const (
	RIGIDBODY_NONE RigidBodyType = iota
	RIGIDBODY_STATIC
	RIGIDBODY_DYNAMIC
	RIGIDBODY_KINEMATIC
)

type RigidBodyState int

const (
	RBSTATE_STILL RigidBodyState = iota
	RBSTATE_MOVING
	RBSTATE_GROUNDED
	RBSTATE_JUMPING
	RBSTATE_FALLING
)

type RigidBody struct {
	RBtype RigidBodyType
	State  RigidBodyState
}
