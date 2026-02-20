package game.audio;

import java.lang.foreign.Arena;
import java.lang.foreign.FunctionDescriptor;
import java.lang.foreign.Linker;
import java.lang.foreign.MemoryLayout;
import java.lang.foreign.MemorySegment;
import java.lang.foreign.SymbolLookup;
import java.lang.invoke.MethodHandle;

import static java.lang.foreign.ValueLayout.ADDRESS;
import static java.lang.foreign.ValueLayout.JAVA_INT;

public final class Audio {
    private static final Linker LINKER = Linker.nativeLinker();
    private static SymbolLookup LOOKUP;
    private static SymbolLookup SDL_LOOKUP;
    private static boolean initialized = false;
    private static boolean available = false;

    private static final boolean IS_WINDOWS = System.getProperty("os.name").toLowerCase().contains("win");
    private static final boolean IS_LINUX = System.getProperty("os.name").toLowerCase().contains("linux");
    private static final boolean IS_MAC = System.getProperty("os.name").toLowerCase().contains("mac");

    private static MemorySegment menuMusic;
    private static MemorySegment gameMusic;
    private static MemorySegment deadMusic;

    private static MemorySegment jumpSound;
    private static MemorySegment deadSound;

    private static String currentMusic = "";

    private static MethodHandle Mix_OpenAudio;
    private static MethodHandle Mix_CloseAudio;
    private static MethodHandle Mix_LoadMUS;
    private static MethodHandle Mix_LoadWAV_RW;
    private static MethodHandle Mix_PlayMusic;
    private static MethodHandle Mix_HaltMusic;
    private static MethodHandle Mix_PlayChannel;
    private static MethodHandle Mix_FreeMusic;
    private static MethodHandle Mix_FreeChunk;
    private static MethodHandle Mix_VolumeMusic;
    private static MethodHandle Mix_Volume;

    private static MethodHandle fn(String name, MemoryLayout ret, MemoryLayout... args) {
        return LINKER.downcallHandle(
                LOOKUP.find(name).orElseThrow(() -> new UnsatisfiedLinkError(name)),
                FunctionDescriptor.of(ret, args)
        );
    }

    private static MethodHandle fnVoid(String name, MemoryLayout... args) {
        return LINKER.downcallHandle(
                LOOKUP.find(name).orElseThrow(() -> new UnsatisfiedLinkError(name)),
                FunctionDescriptor.ofVoid(args)
        );
    }

    /**
     * Initialize the audio system. Call this after SDL is initialized.
     */
    public static void init() {
        if (initialized) return;
        initialized = true;

        try {
            String mixerLibName;
            String sdlLibName;

            if (IS_WINDOWS) {
                mixerLibName = "SDL2_mixer";
                sdlLibName = "SDL2";
            } else if (IS_MAC) {
                mixerLibName = "SDL2_mixer";
                sdlLibName = "SDL2";
            } else {
                // Linux
                mixerLibName = "libSDL2_mixer.so";
                sdlLibName = "libSDL2.so";
            }

            System.out.println("Platform: " + System.getProperty("os.name"));
            System.out.println("Attempting to load " + mixerLibName + "...");

            LOOKUP = SymbolLookup.libraryLookup(mixerLibName, Arena.global());
            System.out.println(mixerLibName + " loaded successfully");

            SDL_LOOKUP = SymbolLookup.libraryLookup(sdlLibName, Arena.global());
            System.out.println(sdlLibName + " loaded successfully");

            Mix_OpenAudio = fn("Mix_OpenAudio", JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT);
            Mix_CloseAudio = fnVoid("Mix_CloseAudio");
            Mix_LoadMUS = fn("Mix_LoadMUS", ADDRESS, ADDRESS);
            Mix_LoadWAV_RW = fn("Mix_LoadWAV_RW", ADDRESS, ADDRESS, JAVA_INT);
            Mix_PlayMusic = fn("Mix_PlayMusic", JAVA_INT, ADDRESS, JAVA_INT);
            Mix_HaltMusic = fn("Mix_HaltMusic", JAVA_INT);
            Mix_PlayChannel = fn("Mix_PlayChannel", JAVA_INT, JAVA_INT, ADDRESS, JAVA_INT);
            Mix_FreeMusic = fnVoid("Mix_FreeMusic", ADDRESS);
            Mix_FreeChunk = fnVoid("Mix_FreeChunk", ADDRESS);
            Mix_VolumeMusic = fn("Mix_VolumeMusic", JAVA_INT, JAVA_INT);
            Mix_Volume = fn("Mix_Volume", JAVA_INT, JAVA_INT, JAVA_INT);

            System.out.println("Opening audio device...");
            int result = (int) Mix_OpenAudio.invoke(44100, 0x8010, 2, 2048);
            if (result < 0) {
                System.err.println("Failed to open audio device (error code: " + result + ")");
                return;
            }

            available = true;
            System.out.println("Audio system initialized successfully");

            loadAudioFiles();

        } catch (UnsatisfiedLinkError e) {
            System.err.println("===========================================");
            if (IS_WINDOWS) {
                System.err.println("SDL2_mixer.dll not found!");
            } else if (IS_LINUX) {
                System.err.println("libSDL2_mixer.so not found!");
                System.err.println("Install with: sudo apt install libsdl2-mixer-2.0-0");
            } else if (IS_MAC) {
                System.err.println("SDL2_mixer library not found!");
                System.err.println("Install with: brew install sdl2_mixer");
            }
            System.err.println("Audio will be disabled.");
            System.err.println("");
            System.err.println("To enable audio:");
            System.err.println("1. Download SDL2_mixer from:");
            System.err.println("   https://github.com/libsdl-org/SDL_mixer/releases");
            if (IS_WINDOWS) {
                System.err.println("2. Extract SDL2_mixer.dll to your SDL2 folder");
            }
            System.err.println("===========================================");
            available = false;
        } catch (Throwable t) {
            System.err.println("Audio system not available: " + t.getMessage());
            t.printStackTrace();
            available = false;
        }
    }

    private static void loadAudioFiles() {
        if (!available) return;

        System.out.println("Loading audio files...");
        System.out.println("Working directory: " + System.getProperty("user.dir"));

        try {
            menuMusic = loadMusic("res/audio/menu_music.ogg");
            if (menuMusic == null) menuMusic = loadMusic("res/audio/menu_music.wav");

            gameMusic = loadMusic("res/audio/game_music.ogg");
            if (gameMusic == null) gameMusic = loadMusic("res/audio/game_music.wav");

            deadMusic = loadMusic("res/audio/dead_music.ogg");
            if (deadMusic == null) deadMusic = loadMusic("res/audio/dead_music.wav");

            jumpSound = loadSound("res/audio/jump.wav");
            deadSound = loadSound("res/audio/dead.wav");

            int loaded = 0;
            if (menuMusic != null) loaded++;
            if (gameMusic != null) loaded++;
            if (deadMusic != null) loaded++;
            if (jumpSound != null) loaded++;
            if (deadSound != null) loaded++;

            System.out.println("Audio files loaded: " + loaded + "/5");

            if (loaded == 0) {
                System.err.println("===========================================");
                System.err.println("No audio files found in res/audio/");
                System.err.println("Required files:");
                System.err.println("  - menu_music.ogg");
                System.err.println("  - game_music.ogg");
                System.err.println("  - dead_music.ogg");
                System.err.println("  - jump.wav");
                System.err.println("  - dead.wav");
                System.err.println("===========================================");
            }
        } catch (Throwable t) {
            System.err.println("Failed to load some audio files: " + t.getMessage());
        }
    }

    private static MemorySegment loadMusic(String path) {
        try (Arena arena = Arena.ofConfined()) {
            java.io.File file = new java.io.File(path);
            System.out.println("Loading music: " + path + " (exists: " + file.exists() + ", size: " + file.length() + ")");

            MemorySegment cstr = arena.allocateUtf8String(path);
            MemorySegment music = (MemorySegment) Mix_LoadMUS.invoke(cstr);
            if (music == null || music.address() == 0) {
                System.err.println("Failed to load music: " + path + " (Mix_LoadMUS returned null)");
                return null;
            }
            System.out.println("Successfully loaded music: " + path);
            return music;
        } catch (Throwable t) {
            System.err.println("Error loading music " + path + ": " + t.getMessage());
            return null;
        }
    }

    private static MemorySegment loadSound(String path) {
        try (Arena arena = Arena.ofConfined()) {
            java.io.File file = new java.io.File(path);
            System.out.println("Loading sound: " + path + " (exists: " + file.exists() + ", size: " + file.length() + ")");

            MemorySegment pathStr = arena.allocateUtf8String(path);
            MemorySegment modeStr = arena.allocateUtf8String("rb");

            MethodHandle rwFromFile = LINKER.downcallHandle(
                SDL_LOOKUP.find("SDL_RWFromFile").orElseThrow(),
                FunctionDescriptor.of(ADDRESS, ADDRESS, ADDRESS)
            );

            MemorySegment rwops = (MemorySegment) rwFromFile.invoke(pathStr, modeStr);
            if (rwops == null || rwops.address() == 0) {
                System.err.println("Failed to open file: " + path + " (SDL_RWFromFile returned null)");
                return null;
            }

            MemorySegment chunk = (MemorySegment) Mix_LoadWAV_RW.invoke(rwops, 1);
            if (chunk == null || chunk.address() == 0) {
                System.err.println("Failed to load sound: " + path + " (Mix_LoadWAV_RW returned null)");
                return null;
            }
            System.out.println("Successfully loaded sound: " + path);
            return chunk;
        } catch (Throwable t) {
            System.err.println("Error loading sound " + path + ": " + t.getMessage());
            t.printStackTrace();
            return null;
        }
    }

    /**
     * Play menu background music (loops).
     */
    public static void playMenuMusic() {
        if (!available || menuMusic == null || "menu".equals(currentMusic)) return;
        stopMusic();
        playMusic(menuMusic);
        currentMusic = "menu";
    }

    /**
     * Play game background music (loops).
     */
    public static void playGameMusic() {
        if (!available || gameMusic == null || "game".equals(currentMusic)) return;
        stopMusic();
        playMusic(gameMusic);
        currentMusic = "game";
    }

    /**
     * Play dead/game over music (loops).
     */
    public static void playDeadMusic() {
        if (!available || deadMusic == null || "dead".equals(currentMusic)) return;
        stopMusic();
        playMusic(deadMusic);
        currentMusic = "dead";
    }

    private static void playMusic(MemorySegment music) {
        if (music == null || music.address() == 0) {
            System.err.println("Cannot play music: null or invalid segment");
            return;
        }
        try {
            System.out.println("Playing music, setting volume...");
            Mix_VolumeMusic.invoke(64);  // Set volume to 50%
            System.out.println("Calling Mix_PlayMusic...");
            int result = (int) Mix_PlayMusic.invoke(music, -1);  // -1 = loop forever
            System.out.println("Mix_PlayMusic result: " + result);
        } catch (Throwable t) {
            System.err.println("Error playing music: " + t.getMessage());
            t.printStackTrace();
        }
    }

    /**
     * Stop any playing music.
     */
    public static void stopMusic() {
        if (!available) return;
        try {
            Mix_HaltMusic.invoke();
            currentMusic = "";
        } catch (Throwable ignored) {}
    }

    /**
     * Play jump sound effect.
     */
    public static void playJumpSound() {
        if (!available || jumpSound == null) return;
        playSound(jumpSound);
    }

    /**
     * Play dead sound effect.
     */
    public static void playDeadSound() {
        if (!available || deadSound == null) return;
        playSound(deadSound);
    }

    private static void playSound(MemorySegment chunk) {
        if (chunk == null || chunk.address() == 0) {
            System.err.println("Cannot play sound: null or invalid segment");
            return;
        }
        try {
            System.out.println("Playing sound effect...");
            Mix_Volume.invoke(-1, 96);  // Set volume to 75%
            int result = (int) Mix_PlayChannel.invoke(-1, chunk, 0);  // -1 = first free channel, 0 = no loop
            System.out.println("Mix_PlayChannel result: " + result);
        } catch (Throwable t) {
            System.err.println("Error playing sound: " + t.getMessage());
            t.printStackTrace();
        }
    }

    public static void cleanup() {
        if (!available) return;
        try {
            stopMusic();

            if (menuMusic != null && menuMusic.address() != 0) Mix_FreeMusic.invoke(menuMusic);
            if (gameMusic != null && gameMusic.address() != 0) Mix_FreeMusic.invoke(gameMusic);
            if (deadMusic != null && deadMusic.address() != 0) Mix_FreeMusic.invoke(deadMusic);
            if (jumpSound != null && jumpSound.address() != 0) Mix_FreeChunk.invoke(jumpSound);
            if (deadSound != null && deadSound.address() != 0) Mix_FreeChunk.invoke(deadSound);

            Mix_CloseAudio.invoke();
        } catch (Throwable ignored) {}
    }

    public static boolean isAvailable() {
        return available;
    }

    private Audio() {}
}

