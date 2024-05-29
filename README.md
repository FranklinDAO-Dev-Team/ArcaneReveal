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

Use WASD to move the player. Moving into enemies triggers the default attack. Use number keys 1, 2, 3, 4 to select wands. Use arrow keys up, down, left, right to cast the wands in the corresponding direction. Each wand will have a randomized selection of spell abilities that are hidden at the beginning of each round. It will expire after one time of usage, and restore once player reaches the next level.

