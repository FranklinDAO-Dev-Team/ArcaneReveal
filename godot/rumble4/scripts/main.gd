extends Node2D

var is_on_fire = false
var fire_duration = 3.0  # Duration of the fire effect in seconds
var fire_damage = 10  # Damage dealt by the fire effect

func _ready():
	# Initialize the tile
	is_on_fire = false

func apply_fire_effect():
	if not is_on_fire:
		is_on_fire = true
		# Play fire animation or visual effect
		# You can add an AnimatedSprite or Particles2D node as a child of the tile to represent the fire
		# For example, if you have an AnimatedSprite named "FireAnimation":
		$FireAnimation.visible = true
		$FireAnimation.play("fire")
		
		# Start a timer to control the duration of the fire effect
		var timer = Timer.new()
		timer.connect("timeout", Callable(self, "_on_fire_timer_timeout"))
		timer.set_wait_time(fire_duration)
		timer.set_one_shot(true)
		add_child(timer)
		timer.start()

func _on_fire_timer_timeout():
	is_on_fire = false
	# Stop fire animation or visual effect
	# For example, if you have an AnimatedSprite named "FireAnimation":
	$FireAnimation.visible = false
	$FireAnimation.stop()
	
	# Deal damage to any enemies standing on the tile (implement this based on your game logic)
	# You can add a function to deal damage to enemies here
