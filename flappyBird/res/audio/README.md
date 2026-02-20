# Flappy Bird Audio Setup

## Required Audio Files

Place the following audio files in the `res/audio/` directory:

### Music Files (OGG format recommended):
- `menu_music.ogg` - Background music for the menu screen
- `game_music.ogg` - Background music during gameplay
- `dead_music.ogg` - Music for the game over screen

### Sound Effects (WAV format):
- `jump.wav` - Sound when the bird jumps/flaps
- `dead.wav` - Sound when the bird dies

## SDL2_mixer Requirement

The audio system uses SDL2_mixer. Make sure you have:

1. **SDL2_mixer.dll** in your system PATH or next to SDL2.dll

### Download SDL2_mixer:
- Go to: https://github.com/libsdl-org/SDL_mixer/releases
- Download the Windows runtime (e.g., `SDL2_mixer-2.x.x-win32-x64.zip`)
- Extract `SDL2_mixer.dll` to the same location as `SDL2.dll`

## Notes:
- If audio files are missing, the game will still run without sound
- If SDL2_mixer is not available, the game will run silently
- OGG format is recommended for music (smaller file size, good quality)
- WAV format is recommended for short sound effects (instant playback)

