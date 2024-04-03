extends RigidBody2D

# Enemy stats
var health = 100
var move_speed = 10  # Speed at which the enemy moves
var move_time = 1.0  # Time in seconds the enemy moves in one direction
var timer = 0.0  # Tracks time for movement

# Movement direction
var move_direction = Vector2.ZERO

func _ready():
	randomize()  # Initialize the random number generator
	timer = move_time
	choose_random_direction()

func _process(delta):
	timer -= delta
	if timer <= 0:
		timer = move_time
		choose_random_direction()

	move_local_x(move_direction.x * move_speed * delta)
	move_local_y(move_direction.y * move_speed * delta)

func choose_random_direction():
	var directions = [Vector2.UP, Vector2.DOWN, Vector2.LEFT, Vector2.RIGHT]
	move_direction = directions[randi() % directions.size()]

func take_damage(amount):
	health -= amount
	if health <= 0:
		die()

func die():
	queue_free()

