extends Node

@onready var enemy1: Area2D = $Enemy1
@onready var enemy2: Area2D = $Enemy2
	
@onready var client : NakamaClient = Nakama.create_client("defaultkey", "127.0.0.1", 7350, "http")
@onready var socket = Nakama.create_socket_from(client)
# Get the System's unique device identifier
@onready var device_id = OS.get_unique_id()
@onready var session : NakamaSession = await client.authenticate_device_async(device_id)

var enemies_defeated = 0

func _on_enemy_exited(enemy: Area2D):
	if not enemy.is_inside_tree():
		enemies_defeated += 1
		if enemies_defeated == 2:
			# Display WIN text or perform any other desired action
			$GameWinLabel.visible = true


func _ready():
	enemy1.tree_exited.connect(_on_enemy_exited.bind(enemy1))
	enemy2.tree_exited.connect(_on_enemy_exited.bind(enemy2))

	# Authenticate with the Nakama server using Device Authentication
	var session_result: NakamaAsyncResult = await client.authenticate_device_async(device_id)
	if session_result.is_exception():
		print("An error occurred: %s" % session_result)
		return

	session = session_result
	print("Successfully authenticated: %s" % session)


	var connected_result: NakamaAsyncResult = await socket.connect_async(session)
	if connected_result.is_exception():
		print("An error occurred: %s" % connected_result)
		return
	print("Socket connected.")
	
	socket.received_notification.connect(self._on_notification)
	
	# Check whether account already has persona
	var resp = await client.rpc_async(session, "nakama/show-persona")
	if resp.is_exception():
		print("An error occured: %s", % resp)
	else:
		if JSON.parse_string(resp.payload)["status"] == "accepted":
			print("Device already has a persona, skipping creation")
		else:
			# Create persona	
			resp = await client.rpc_async(session, "nakama/claim-persona", JSON.stringify({"personaTag": "CoolMage"}))
			if resp.is_exception():
				print("An error occured: %s", % resp)
				return
			print("Created PersonaTag: CoolMage")
	
	# Make CreateGame TXN call
	var random = RandomNumberGenerator.new()
	resp = await client.rpc_async(session, "tx/game/request-game", JSON.stringify({"playerSource": str(random.randi_range(100000, 999999))}))
	if resp.is_exception():
		print("An error occured: %s", % resp)
		return
	print("Successfully created game: %s", % resp)
	

func _on_rpc_response(result: NakamaAsyncResult):
	if result.is_successful():
		print("RPC succeeded: ", result.payload)
	else:
		print("RPC failed: ", result.error)
		
func _on_notification(p_notification : NakamaAPI.ApiNotification):
	var notification = JSON.new()
	notification.parse(p_notification.content)
	print("[Notification]: ", notification.data)
