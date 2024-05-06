# Seismic_Cinco_Paus

# Dev Setup

## Setting up the world engine backend 
On the first set up, you may need to add the following line your .zshrc:
```
Add export PATH=$PATH:~/go/bin
```
From the project folder, run 
```
cd world_engine/cardinal
make
```
Then, after making sure docker is running,
```
cd .. 
world cardinal purge && world cardinal start --editor
```

## Starting the Game Client
Once the World Engine is running, you can start the godot/rumble4 folder in the Godot app. Press the Play button in the top right of the Godot game engine window (or Command B on Mac). It should create a separate pop-up window for the game, at which point you can play.

## Gameplay
Use WASD to move the player. Moving into enemies triggers the default attack. Use the number keys 1, 2, 3, 4 to select the Wand you would like to cast. After selecting, cast the wand using the arrow keys in the desired direction: up, down, right, left. Casting a wand does an unknown ability, and counts as a player turn. Wands are one time use per level, and reset after you reach the next level.
