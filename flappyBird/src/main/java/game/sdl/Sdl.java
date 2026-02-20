package game.sdl;

import java.lang.foreign.Arena;
import java.lang.foreign.FunctionDescriptor;
import java.lang.foreign.Linker;
import java.lang.foreign.MemoryLayout;
import java.lang.foreign.MemorySegment;
import java.lang.foreign.SymbolLookup;
import java.lang.invoke.MethodHandle;

import static java.lang.foreign.ValueLayout.ADDRESS;
import static java.lang.foreign.ValueLayout.JAVA_INT;

/** Minimal SDL2 bindings via Panama for windowing, GL context creation, and event polling. */
public final class Sdl {
    private static final Linker LINKER = Linker.nativeLinker();
    private static final SymbolLookup LOOKUP = SymbolLookup.libraryLookup(System.getProperty("os.name").toLowerCase().contains("linux") ? "libSDL2.so" : "SDL2", Arena.global());

    public static final int SDL_INIT_VIDEO = 0x00000020;
    public static final int SDL_INIT_AUDIO = 0x00000010;
    public static final int SDL_WINDOW_OPENGL = 0x00000002;
    public static final int SDL_WINDOW_SHOWN = 0x00000004;
    public static final int SDL_WINDOWPOS_CENTERED = 0x2FFF0000;

    public static final int SDL_GL_CONTEXT_MAJOR_VERSION = 17;
    public static final int SDL_GL_CONTEXT_MINOR_VERSION = 18;
    public static final int SDL_GL_CONTEXT_PROFILE_MASK = 21;
    public static final int SDL_GL_CONTEXT_PROFILE_CORE = 0x0001;
    public static final int SDL_GL_DOUBLEBUFFER = 5;

    public static final int SDL_QUIT = 0x100;
    public static final int SDL_KEYDOWN = 0x300;
    public static final int SDL_KEYUP = 0x301;

    private static final MethodHandle SDL_Init = fn("SDL_Init", JAVA_INT, JAVA_INT);
    private static final MethodHandle SDL_Quit = fnVoid("SDL_Quit");
    private static final MethodHandle SDL_CreateWindow = fn("SDL_CreateWindow", ADDRESS, ADDRESS, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT);
    private static final MethodHandle SDL_DestroyWindow = fnVoid("SDL_DestroyWindow", ADDRESS);
    private static final MethodHandle SDL_PollEvent = fn("SDL_PollEvent", JAVA_INT, ADDRESS);
    private static final MethodHandle SDL_GL_SetAttribute = fn("SDL_GL_SetAttribute", JAVA_INT, JAVA_INT, JAVA_INT);
    private static final MethodHandle SDL_GL_CreateContext = fn("SDL_GL_CreateContext", ADDRESS, ADDRESS);
    private static final MethodHandle SDL_GL_DeleteContext = fnVoid("SDL_GL_DeleteContext", ADDRESS);
    private static final MethodHandle SDL_GL_SwapWindow = fnVoid("SDL_GL_SwapWindow", ADDRESS);
    private static final MethodHandle SDL_GL_GetProcAddress = fn("SDL_GL_GetProcAddress", ADDRESS, ADDRESS);

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

    public static void initVideo() {
        try {
            int rc = (int) SDL_Init.invoke(SDL_INIT_VIDEO | SDL_INIT_AUDIO);
            if (rc != 0) throw new IllegalStateException("SDL_Init failed: " + rc);
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    public static void quit() {
        try {
            SDL_Quit.invoke();
        } catch (Throwable ignored) {
        }
    }

    public static MemorySegment createWindow(String title, int w, int h) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment cstr = arena.allocateUtf8String(title);
            return (MemorySegment) SDL_CreateWindow.invoke(cstr, SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, w, h, SDL_WINDOW_OPENGL | SDL_WINDOW_SHOWN);
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    public static void destroyWindow(MemorySegment win) {
        try {
            SDL_DestroyWindow.invoke(win);
        } catch (Throwable ignored) {
        }
    }

    public static MemorySegment createGlContext(MemorySegment win, int major, int minor) {
        try {
            SDL_GL_SetAttribute.invoke(SDL_GL_DOUBLEBUFFER, 1);
            SDL_GL_SetAttribute.invoke(SDL_GL_CONTEXT_MAJOR_VERSION, major);
            SDL_GL_SetAttribute.invoke(SDL_GL_CONTEXT_MINOR_VERSION, minor);
            SDL_GL_SetAttribute.invoke(SDL_GL_CONTEXT_PROFILE_MASK, SDL_GL_CONTEXT_PROFILE_CORE);
            return (MemorySegment) SDL_GL_CreateContext.invoke(win);
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    public static void deleteGlContext(MemorySegment ctx) {
        try {
            SDL_GL_DeleteContext.invoke(ctx);
        } catch (Throwable ignored) {
        }
    }

    public static MemorySegment getProcAddress(String name) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment cstr = arena.allocateUtf8String(name);
            return (MemorySegment) SDL_GL_GetProcAddress.invoke(cstr);
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    public static void swapWindow(MemorySegment win) {
        try {
            SDL_GL_SwapWindow.invoke(win);
        } catch (Throwable ignored) {
        }
    }

    /** Poll a single event; returns type, and if key event fills keyOut[0] with scancode. */
    public static int pollEvent(int[] keyOut) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment evt = arena.allocate(56); // SDL_Event size on 64-bit
            int rc = (int) SDL_PollEvent.invoke(evt);
            if (rc == 0) return 0;
            int type = evt.get(JAVA_INT, 0);
            if (type == SDL_KEYDOWN || type == SDL_KEYUP) {
                // SDL_KeyboardEvent structure:
                // Uint32 type (4 bytes) - offset 0
                // Uint32 timestamp (4 bytes) - offset 4
                // Uint32 windowID (4 bytes) - offset 8
                // Uint8 state (1 byte) - offset 12
                // Uint8 repeat (1 byte) - offset 13
                // Uint8 padding2 (1 byte) - offset 14
                // Uint8 padding3 (1 byte) - offset 15
                // SDL_Keysym keysym - offset 16
                //   SDL_Scancode scancode (4 bytes) - offset 16
                int scancode = evt.get(JAVA_INT, 16);
                keyOut[0] = scancode;
            }
            return type;
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    private Sdl() {
    }
}
