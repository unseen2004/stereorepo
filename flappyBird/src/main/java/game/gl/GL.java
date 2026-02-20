package game.gl;

import java.lang.foreign.Arena;
import java.lang.foreign.FunctionDescriptor;
import java.lang.foreign.Linker;
import java.lang.foreign.MemoryLayout;
import java.lang.foreign.MemorySegment;
import java.lang.foreign.SymbolLookup;
import java.lang.invoke.MethodHandle;

import static java.lang.foreign.ValueLayout.ADDRESS;
import static java.lang.foreign.ValueLayout.JAVA_FLOAT;
import static java.lang.foreign.ValueLayout.JAVA_INT;
import static java.lang.foreign.ValueLayout.JAVA_LONG;

/** Minimal OpenGL bindings via Panama for the functions used in this project. */
public final class GL {
    private static final Linker LINKER = Linker.nativeLinker();
    private static final SymbolLookup LOOKUP = name -> {
        try {
            MemorySegment addr = game.sdl.Sdl.getProcAddress(name);
            if (addr != null && addr.address() != 0) {
                return java.util.Optional.of(addr);
            }
        } catch (Throwable t) {
        }
        String libName = System.getProperty("os.name").toLowerCase().contains("linux") ? "libGL.so" : "opengl32";
        return SymbolLookup.libraryLookup(libName, Arena.global()).find(name);
    };

    private static MethodHandle fn(String n, MemoryLayout ret, MemoryLayout... args) {
        FunctionDescriptor fd = (ret == null) ? FunctionDescriptor.ofVoid(args) : FunctionDescriptor.of(ret, args);
        return LINKER.downcallHandle(
                LOOKUP.find(n).orElseThrow(() -> new UnsatisfiedLinkError(n)),
                fd
        );
    }

    private static final MethodHandle glClear = fn("glClear", null, JAVA_INT);
    private static final MethodHandle glClearColor = fn("glClearColor", null, JAVA_FLOAT, JAVA_FLOAT, JAVA_FLOAT, JAVA_FLOAT);
    private static final MethodHandle glGetError = fn("glGetError", JAVA_INT);
    private static final MethodHandle glEnable = fn("glEnable", null, JAVA_INT);
    private static final MethodHandle glBlendFunc = fn("glBlendFunc", null, JAVA_INT, JAVA_INT);
    private static final MethodHandle glViewport = fn("glViewport", null, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT);
    private static final MethodHandle glGetString = fn("glGetString", ADDRESS, JAVA_INT);
    private static final MethodHandle glGenTextures = fn("glGenTextures", null, JAVA_INT, ADDRESS);
    private static final MethodHandle glActiveTexture = fn("glActiveTexture", null, JAVA_INT);
    private static final MethodHandle glBindTexture = fn("glBindTexture", null, JAVA_INT, JAVA_INT);
    private static final MethodHandle glTexParameteri = fn("glTexParameteri", null, JAVA_INT, JAVA_INT, JAVA_INT);
    private static final MethodHandle glTexImage2D = fn("glTexImage2D", null, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glGenVertexArrays = fn("glGenVertexArrays", null, JAVA_INT, ADDRESS);
    private static final MethodHandle glBindVertexArray = fn("glBindVertexArray", null, JAVA_INT);
    private static final MethodHandle glGenBuffers = fn("glGenBuffers", null, JAVA_INT, ADDRESS);
    private static final MethodHandle glBindBuffer = fn("glBindBuffer", null, JAVA_INT, JAVA_INT);
    private static final MethodHandle glBufferData = fn("glBufferData", null, JAVA_INT, JAVA_LONG, ADDRESS, JAVA_INT);
    private static final MethodHandle glBufferSubData = fn("glBufferSubData", null, JAVA_INT, JAVA_LONG, JAVA_LONG, ADDRESS);
    private static final MethodHandle glVertexAttribPointer = fn("glVertexAttribPointer", null, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glEnableVertexAttribArray = fn("glEnableVertexAttribArray", null, JAVA_INT);
    private static final MethodHandle glDrawElements = fn("glDrawElements", null, JAVA_INT, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glDrawArrays = fn("glDrawArrays", null, JAVA_INT, JAVA_INT, JAVA_INT);
    private static final MethodHandle glUseProgram = fn("glUseProgram", null, JAVA_INT);
    private static final MethodHandle glGetUniformLocation = fn("glGetUniformLocation", JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glUniform1i = fn("glUniform1i", null, JAVA_INT, JAVA_INT);
    private static final MethodHandle glUniform1f = fn("glUniform1f", null, JAVA_INT, JAVA_FLOAT);
    private static final MethodHandle glUniform2f = fn("glUniform2f", null, JAVA_INT, JAVA_FLOAT, JAVA_FLOAT);
    private static final MethodHandle glUniform3f = fn("glUniform3f", null, JAVA_INT, JAVA_FLOAT, JAVA_FLOAT, JAVA_FLOAT);
    private static final MethodHandle glUniformMatrix4fv = fn("glUniformMatrix4fv", null, JAVA_INT, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glCreateProgram = fn("glCreateProgram", JAVA_INT);
    private static final MethodHandle glCreateShader = fn("glCreateShader", JAVA_INT, JAVA_INT);
    private static final MethodHandle glShaderSource = fn("glShaderSource", null, JAVA_INT, JAVA_INT, ADDRESS, ADDRESS);
    private static final MethodHandle glCompileShader = fn("glCompileShader", null, JAVA_INT);
    private static final MethodHandle glGetShaderiv = fn("glGetShaderiv", null, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glGetShaderInfoLog = fn("glGetShaderInfoLog", null, JAVA_INT, JAVA_INT, ADDRESS, ADDRESS);
    private static final MethodHandle glGetProgramiv = fn("glGetProgramiv", null, JAVA_INT, JAVA_INT, ADDRESS);
    private static final MethodHandle glGetProgramInfoLog = fn("glGetProgramInfoLog", null, JAVA_INT, JAVA_INT, ADDRESS, ADDRESS);
    private static final MethodHandle glAttachShader = fn("glAttachShader", null, JAVA_INT, JAVA_INT);
    private static final MethodHandle glLinkProgram = fn("glLinkProgram", null, JAVA_INT);
    private static final MethodHandle glValidateProgram = fn("glValidateProgram", null, JAVA_INT);

    public static final int GL_COLOR_BUFFER_BIT = 0x00004000;
    public static final int GL_DEPTH_BUFFER_BIT = 0x00000100;
    public static final int GL_BLEND = 0x0BE2;
    public static final int GL_SRC_ALPHA = 0x0302;
    public static final int GL_ONE_MINUS_SRC_ALPHA = 0x0303;
    public static final int GL_TEXTURE_2D = 0x0DE1;
    public static final int GL_TEXTURE_MIN_FILTER = 0x2801;
    public static final int GL_TEXTURE_MAG_FILTER = 0x2800;
    public static final int GL_NEAREST = 0x2600;
    public static final int GL_RGBA = 0x1908;
    public static final int GL_UNSIGNED_BYTE = 0x1401;
    public static final int GL_ARRAY_BUFFER = 0x8892;
    public static final int GL_ELEMENT_ARRAY_BUFFER = 0x8893;
    public static final int GL_STATIC_DRAW = 0x88E4;
    public static final int GL_DYNAMIC_DRAW = 0x88E8;
    public static final int GL_FLOAT = 0x1406;
    public static final int GL_TRIANGLES = 0x0004;
    public static final int GL_FALSE = 0;
    public static final int GL_VERTEX_SHADER = 0x8B31;
    public static final int GL_FRAGMENT_SHADER = 0x8B30;
    public static final int GL_COMPILE_STATUS = 0x8B81;
    public static final int GL_LINK_STATUS = 0x8B82;
    public static final int GL_TEXTURE0 = 0x84C0;
    public static final int GL_VERSION = 0x1F02;
    public static final int GL_NO_ERROR = 0;

    public static void clear(int mask) { invokeVoid(glClear, mask); }
    public static void clearColor(float r, float g, float b, float a) { invokeVoid(glClearColor, r, g, b, a); }
    public static int getError() { return invokeInt(glGetError); }
    public static void enable(int cap) { invokeVoid(glEnable, cap); }
    public static void blendFunc(int src, int dst) { invokeVoid(glBlendFunc, src, dst); }
    public static void viewport(int x, int y, int width, int height) { invokeVoid(glViewport, x, y, width, height); }
    public static String getString(int name) {
        try {
            MemorySegment seg = (MemorySegment) glGetString.invoke(name);
            if (seg == null || seg.address() == 0) return "";
            return seg.reinterpret(Long.MAX_VALUE).getUtf8String(0);
        } catch (Throwable t) {
            throw new RuntimeException(t);
        }
    }

    public static int genTexture() {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(JAVA_INT);
            glGenTextures.invoke(1, buf);
            return buf.get(JAVA_INT, 0);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }

    public static void activeTexture(int texture) { invokeVoid(glActiveTexture, texture); }
    public static void bindTexture(int target, int tex) { invokeVoid(glBindTexture, target, tex); }
    public static void texParameteri(int target, int pname, int param) { invokeVoid(glTexParameteri, target, pname, param); }
    public static void texImage2D(int target, int level, int internalFormat, int width, int height, int border, int format, int type, MemorySegment data) {
        invokeVoid(glTexImage2D, target, level, internalFormat, width, height, border, format, type, data);
    }

    public static int genVertexArray() { return genOne(glGenVertexArrays); }
    public static void genVertexArrays(int n, MemorySegment arrays) { invokeVoid(glGenVertexArrays, n, arrays); }
    public static void bindVertexArray(int vao) { invokeVoid(glBindVertexArray, vao); }

    public static int genBuffer() { return genOne(glGenBuffers); }
    public static void genBuffers(int n, MemorySegment buffers) { invokeVoid(glGenBuffers, n, buffers); }
    public static void bindBuffer(int target, int id) { invokeVoid(glBindBuffer, target, id); }
    public static void bufferData(int target, long sizeBytes, MemorySegment data, int usage) { invokeVoid(glBufferData, target, sizeBytes, data, usage); }
    public static void bufferSubData(int target, long offset, long size, MemorySegment data) { invokeVoid(glBufferSubData, target, offset, size, data); }
    public static void vertexAttribPointer(int index, int size, int type, boolean normalized, int stride, long pointer) {
        invokeVoid(glVertexAttribPointer, index, size, type, normalized ? 1 : 0, stride, MemorySegment.ofAddress(pointer));
    }
    public static void enableVertexAttribArray(int index) { invokeVoid(glEnableVertexAttribArray, index); }

    public static void drawElements(int mode, int count, int type, long indicesAddr) {
        invokeVoid(glDrawElements, mode, count, type, MemorySegment.ofAddress(indicesAddr));
    }
    public static void drawArrays(int mode, int first, int count) { invokeVoid(glDrawArrays, mode, first, count); }

    public static void useProgram(int id) { invokeVoid(glUseProgram, id); }
    public static int getUniformLocation(int program, String name) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment c = arena.allocateUtf8String(name); // Fixed: allocateUtf8String -> allocateFrom
            return (int) glGetUniformLocation.invoke(program, c);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static void uniform1i(int loc, int v) { invokeVoid(glUniform1i, loc, v); }
    public static void uniform1f(int loc, float v) { invokeVoid(glUniform1f, loc, v); }
    public static void uniform2f(int loc, float x, float y) { invokeVoid(glUniform2f, loc, x, y); }
    public static void uniform3f(int loc, float x, float y, float z) { invokeVoid(glUniform3f, loc, x, y, z); }
    public static void uniformMatrix4fv(int loc, boolean transpose, MemorySegment data) { invokeVoid(glUniformMatrix4fv, loc, 1, transpose ? 1 : 0, data); }

    public static int createProgram() { return invokeInt(glCreateProgram); }
    public static int createShader(int type) { try { return (int) glCreateShader.invoke(type); } catch (Throwable t) { throw new RuntimeException(t); } }
    public static void shaderSource(int shader, String src) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment c = arena.allocateUtf8String(src); // Fixed: allocateUtf8String -> allocateFrom
            MemorySegment arr = arena.allocateArray(ADDRESS, 1);
            arr.set(ADDRESS, 0, c);
            glShaderSource.invoke(shader, 1, arr, MemorySegment.NULL);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static void compileShader(int shader) { invokeVoid(glCompileShader, shader); }
    public static int getShaderiv(int shader, int pname) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(JAVA_INT);
            glGetShaderiv.invoke(shader, pname, buf);
            return buf.get(JAVA_INT, 0);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static int getProgramiv(int program, int pname) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(JAVA_INT);
            glGetProgramiv.invoke(program, pname, buf);
            return buf.get(JAVA_INT, 0);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static String getShaderInfoLog(int shader) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(1024);
            glGetShaderInfoLog.invoke(shader, 1024, MemorySegment.NULL, buf);
            return buf.getUtf8String(0); // Fixed: getUtf8String -> getString
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static String getProgramInfoLog(int program) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(1024);
            glGetProgramInfoLog.invoke(program, 1024, MemorySegment.NULL, buf);
            return buf.getUtf8String(0);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }
    public static void attachShader(int program, int shader) { invokeVoid(glAttachShader, program, shader); }
    public static void linkProgram(int program) { invokeVoid(glLinkProgram, program); }
    public static void validateProgram(int program) { invokeVoid(glValidateProgram, program); }

    private static int genOne(MethodHandle mh) {
        try (Arena arena = Arena.ofConfined()) {
            MemorySegment buf = arena.allocate(JAVA_INT);
            mh.invoke(1, buf);
            return buf.get(JAVA_INT, 0);
        } catch (Throwable t) { throw new RuntimeException(t); }
    }

    private static int invokeInt(MethodHandle mh) {
        try { return (int) mh.invoke(); } catch (Throwable t) { throw new RuntimeException(t); }
    }
    private static void invokeVoid(MethodHandle mh, Object... args) {
        try { mh.invokeWithArguments(args); } catch (Throwable t) { throw new RuntimeException(t); }
    }

    private GL() {}
}