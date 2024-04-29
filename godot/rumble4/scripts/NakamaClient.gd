extends Node

@onready var enemy1: Area2D = $Enemy1
@onready var enemy2: Area2D = $Enemy2
	
@onready var client : NakamaClient
@onready var socket
@onready var session : NakamaSession
@onready var ray = $RayCast3D

var enemies_defeated = 0
var tile_size = 32

func _on_enemy_exited(enemy: Area2D):
	if not enemy.is_inside_tree():
		enemies_defeated += 1
		if enemies_defeated == 2:
			# Display WIN text or perform any other desired action
			$GameWinLabel.visible = true


func _ready():
	enemy1.tree_exited.connect(_on_enemy_exited.bind(enemy1))
	enemy2.tree_exited.connect(_on_enemy_exited.bind(enemy2))
	
	client = Nakama.create_client("defaultkey", "127.0.0.1", 7350, "http")
	socket = Nakama.create_socket_from(client)

	# Authenticate with the Nakama server using Device Authentication
	var device_id = OS.get_unique_id()
	session = await client.authenticate_device_async(device_id)
	if session.is_exception():
		print("An error occurred: %s" % session)
		return

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
		# Create persona	
		resp = await client.rpc_async(session, "nakama/claim-persona", JSON.stringify({"personaTag": "CoolMage"}))
		if resp.is_exception():
			print("An error occured while claiming persona: %s", % resp)
			return
		print("Created PersonaTag: CoolMage")
	else:
		if JSON.parse_string(resp.payload)["status"] == "accepted":
			print("Device already has a persona, skipping creation")
	
	while true:
		# wait until show persona succeed
		var personaResp = await client.rpc_async(session, "nakama/show-persona")
		if personaResp.is_exception():
			print("persona is not registered yet, waiting")
			await get_tree().create_timer(0.25)
			continue
		else:
			if personaResp.payload != null and JSON.parse_string(personaResp.payload)["status"] == "accepted":
				print("Device has a persona, continuing")
				break
	
	# Make CreateGame TXN call
	var random = RandomNumberGenerator.new()
	resp = await client.rpc_async(session, "tx/game/request-game", JSON.stringify({"playerSource": str(random.randi_range(100000, 999999))}))
	#print(resp)
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
	if notification.data.has("event") and notification.data["event"] == "player-turn":
		var turnLogs = notification.data["log"]
		print("caught player turn event")
