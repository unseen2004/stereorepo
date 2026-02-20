package game.Shaders;

import game.gl.GL;
import game.math.Matrix4f;
import game.math.Vector3f;
import game.util.ShaderUtils;

import java.lang.foreign.MemorySegment;
import java.util.HashMap;
import java.util.Map;

public final class Shader {

    // Vertex attribute locations
    public static final int VERTEX_ATTRIB = 0;
    public static final int TEXCOORD_ATTRIB = 1;

    // Shader instances
    public static Shader BG;
    public static Shader BIRD;
    public static Shader PIPE;
    public static Shader FADE;
    public static Shader TEXT;
    public static Shader OVERLAY;

    private final int programId;
    private final Map<String, Integer> uniformCache = new HashMap<>();
    private boolean isActive = false;

    private Shader(String vertexPath, String fragmentPath) {
        this.programId = ShaderUtils.load(vertexPath, fragmentPath);
    }

    /**
     * Load all shader programs used by the game.
     */
    public static void loadAll() {
        BG = new Shader("res/Shaders/BG.vert", "res/Shaders/BG.frag");
        BIRD = new Shader("res/Shaders/BIRD.vert", "res/Shaders/BIRD.frag");
        PIPE = new Shader("res/Shaders/PIPE.vert", "res/Shaders/PIPE.frag");
        FADE = new Shader("res/Shaders/fade.vert", "res/Shaders/fade.frag");
        TEXT = new Shader("res/Shaders/TEXT.vert", "res/Shaders/TEXT.frag");
        OVERLAY = new Shader("res/Shaders/OVERLAY.vert", "res/Shaders/OVERLAY.frag");
    }

    /**
     * Get or cache a uniform location.
     */
    public int getUniform(String name) {
        return uniformCache.computeIfAbsent(name, n -> {
            int location = GL.getUniformLocation(programId, n);
            if (location == -1) {
                System.err.println("Could not find uniform variable '" + n + "'!");
            }
            return location;
        });
    }

    public void setUniform1f(String name, float value) {
        ensureEnabled();
        GL.uniform1f(getUniform(name), value);
    }

    public void setUniform2f(String name, float x, float y) {
        ensureEnabled();
        GL.uniform2f(getUniform(name), x, y);
    }

    public void setUniform3f(String name, Vector3f vector) {
        ensureEnabled();
        GL.uniform3f(getUniform(name), vector.x, vector.y, vector.z);
    }

    public void setUniformi(String name, int value) {
        ensureEnabled();
        GL.uniform1i(getUniform(name), value);
    }

    public void setUniformMat4f(String name, Matrix4f matrix) {
        ensureEnabled();
        GL.uniformMatrix4fv(getUniform(name), false, MemorySegment.ofBuffer(matrix.toFloatBuffer()));
    }
    
    private void ensureEnabled() {
        if (!isActive) {
            enable();
        }
    }

    public void enable() {
        GL.useProgram(programId);
        isActive = true;
    }

    public void disable() {
        GL.useProgram(0);
        isActive = false;
    }

}
