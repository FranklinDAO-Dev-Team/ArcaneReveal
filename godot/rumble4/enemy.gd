class_name Enemy
extends RigidBody2D

# Enemy stats
var health = 100
var damage = 20

var move_force = 400 # Force to apply for movement
var move_distance = 32

var target_position = Vector2()
var is_moving = false
var wait_timer = 1.0 # Time to wait before the next move (in seconds)
var timer = 0.0 # Internal timer

func _ready():
	target_position = position

func _physics_process(delta):
	if is_moving:
		var direction = (target_position - position).normalized()
		if position.distance_to(target_position) > 10:
			apply_central_impulse(direction * move_force)
		else:
			# Arrived at target position, stoap moving
			position = target_position # Snap to grid
			is_moving = false
			timer = wait_timer # Reset the timer
	elif not is_moving and timer <= 0:
		choose_next_target()

	# Update the timer
	if timer > 0:
		timer -= delta

func choose_next_target():
	# Example: move to the right one tile
	target_position += Vector2(move_distance, 0)
	is_moving = true

func take_damage(amount):
	health -= amount
	if health <= 0:
		die()

func die():
	queue_free() # Remove the enemy from the scene

func _on_body_entered(body):
	if body.has_method("take_damage"):
		body.take_damage(damage)

func _on_ready():
	set_physics_process(true)

func _on_exit_tree():
	set_physics_process(false)


