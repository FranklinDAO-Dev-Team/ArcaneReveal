extends Node

@onready var enemy1: Area2D = $Enemy1
@onready var enemy2: Area2D = $Enemy2

var enemies_defeated = 0

func _ready():
	enemy1.tree_exited.connect(_on_enemy_exited.bind(enemy1))
	enemy2.tree_exited.connect(_on_enemy_exited.bind(enemy2))

func _on_enemy_exited(enemy: Area2D):
	if not enemy.is_inside_tree():
		enemies_defeated += 1
		if enemies_defeated == 2:
			# Display WIN text or perform any other desired action
			$GameWinLabel.visible = true 

#
#@onready var client : NakamaClient = Nakama.create_client("defaultkey", "shiny-pens-wave.loca.lt", 443, "https")
#@onready var socket = Nakama.create_socket_from(client)
#
#func _ready():
#
	## Get the System's unique device identifier
	#var device_id = OS.get_unique_id()
	#print('test')
#
	## Authenticate with the Nakama server using Device Authentication
	#var session : NakamaSession = await client.authenticate_device_async(device_id)
	#if session.is_exception():
		#print("An error occurred: %s" % session)
		#return
	#print("Successfully authenticated: %s" % session)
#
	#var connected : NakamaAsyncResult = await socket.connect_async(session)
	#if connected.is_exception():
		#print("An error occurred: %s" % connected)
		#return
	#print("Socket connected.")
	#
	#socket.received_notification.connect(self._on_notification)
#
#func _on_notification(p_notification : NakamaAPI.ApiNotification):
	#var notification = JSON.new()
	#notification.parse(p_notification.content)
	#print(notification.data)
