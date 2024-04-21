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

