package game.graphics;

import game.Shaders.Shader;
import game.gl.GL;
import game.math.Matrix4f;
import game.math.Vector3f;

import java.lang.foreign.Arena;
import java.lang.foreign.MemorySegment;
import java.lang.foreign.ValueLayout;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class TextRenderer {

    private static int vao, vbo;
    private static boolean initialized = false;
    private static Matrix4f projectionMatrix;

    // 5x7 bitmap font patterns
    private static final Map<Character, int[][]> FONT = new HashMap<>();

    static {
        FONT.put('A', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1}});
        FONT.put('B', new int[][]{{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,0}});
        FONT.put('C', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('D', new int[][]{{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,0}});
        FONT.put('E', new int[][]{{1,1,1,1,1},{1,0,0,0,0},{1,0,0,0,0},{1,1,1,1,0},{1,0,0,0,0},{1,0,0,0,0},{1,1,1,1,1}});
        FONT.put('F', new int[][]{{1,1,1,1,1},{1,0,0,0,0},{1,0,0,0,0},{1,1,1,1,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0}});
        FONT.put('G', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,0},{1,0,1,1,1},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('H', new int[][]{{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1}});
        FONT.put('I', new int[][]{{1,1,1,1,1},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{1,1,1,1,1}});
        FONT.put('J', new int[][]{{0,0,0,0,1},{0,0,0,0,1},{0,0,0,0,1},{0,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('K', new int[][]{{1,0,0,0,1},{1,0,0,1,0},{1,0,1,0,0},{1,1,0,0,0},{1,0,1,0,0},{1,0,0,1,0},{1,0,0,0,1}});
        FONT.put('L', new int[][]{{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0},{1,1,1,1,1}});
        FONT.put('M', new int[][]{{1,0,0,0,1},{1,1,0,1,1},{1,0,1,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1}});
        FONT.put('N', new int[][]{{1,0,0,0,1},{1,1,0,0,1},{1,0,1,0,1},{1,0,0,1,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1}});
        FONT.put('O', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('P', new int[][]{{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,0},{1,0,0,0,0},{1,0,0,0,0},{1,0,0,0,0}});
        FONT.put('Q', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,1,0,1},{1,0,0,1,0},{0,1,1,0,1}});
        FONT.put('R', new int[][]{{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{1,1,1,1,0},{1,0,1,0,0},{1,0,0,1,0},{1,0,0,0,1}});
        FONT.put('S', new int[][]{{0,1,1,1,1},{1,0,0,0,0},{1,0,0,0,0},{0,1,1,1,0},{0,0,0,0,1},{0,0,0,0,1},{1,1,1,1,0}});
        FONT.put('T', new int[][]{{1,1,1,1,1},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0}});
        FONT.put('U', new int[][]{{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('V', new int[][]{{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{0,1,0,1,0},{0,1,0,1,0},{0,0,1,0,0}});
        FONT.put('W', new int[][]{{1,0,0,0,1},{1,0,0,0,1},{1,0,0,0,1},{1,0,1,0,1},{1,0,1,0,1},{1,1,0,1,1},{1,0,0,0,1}});
        FONT.put('X', new int[][]{{1,0,0,0,1},{0,1,0,1,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,1,0,1,0},{1,0,0,0,1}});
        FONT.put('Y', new int[][]{{1,0,0,0,1},{1,0,0,0,1},{0,1,0,1,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0}});
        FONT.put('Z', new int[][]{{1,1,1,1,1},{0,0,0,0,1},{0,0,0,1,0},{0,0,1,0,0},{0,1,0,0,0},{1,0,0,0,0},{1,1,1,1,1}});
        FONT.put('0', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,1,1},{1,0,1,0,1},{1,1,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('1', new int[][]{{0,0,1,0,0},{0,1,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,1,1,1,0}});
        FONT.put('2', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{0,0,0,0,1},{0,0,0,1,0},{0,0,1,0,0},{0,1,0,0,0},{1,1,1,1,1}});
        FONT.put('3', new int[][]{{1,1,1,1,0},{0,0,0,0,1},{0,0,0,0,1},{0,1,1,1,0},{0,0,0,0,1},{0,0,0,0,1},{1,1,1,1,0}});
        FONT.put('4', new int[][]{{0,0,0,1,0},{0,0,1,1,0},{0,1,0,1,0},{1,0,0,1,0},{1,1,1,1,1},{0,0,0,1,0},{0,0,0,1,0}});
        FONT.put('5', new int[][]{{1,1,1,1,1},{1,0,0,0,0},{1,1,1,1,0},{0,0,0,0,1},{0,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('6', new int[][]{{0,0,1,1,0},{0,1,0,0,0},{1,0,0,0,0},{1,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('7', new int[][]{{1,1,1,1,1},{0,0,0,0,1},{0,0,0,1,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,1,0,0}});
        FONT.put('8', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,0}});
        FONT.put('9', new int[][]{{0,1,1,1,0},{1,0,0,0,1},{1,0,0,0,1},{0,1,1,1,1},{0,0,0,0,1},{0,0,0,1,0},{0,1,1,0,0}});
        FONT.put(' ', new int[][]{{0,0,0,0,0},{0,0,0,0,0},{0,0,0,0,0},{0,0,0,0,0},{0,0,0,0,0},{0,0,0,0,0},{0,0,0,0,0}});
        FONT.put(':', new int[][]{{0,0,0,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,0,0,0},{0,0,1,0,0},{0,0,1,0,0},{0,0,0,0,0}});
    }

    public static void init(Matrix4f projection) {
        if (initialized) return;
        projectionMatrix = projection;

        // Use simpler GL methods that return single values
        vao = GL.genVertexArray();
        vbo = GL.genBuffer();

        GL.bindVertexArray(vao);
        GL.bindBuffer(GL.GL_ARRAY_BUFFER, vbo);

        try (Arena arena = Arena.ofConfined()) {
            MemorySegment emptyData = arena.allocate(50000 * 6 * 4 * 4);
            GL.bufferData(GL.GL_ARRAY_BUFFER, 50000 * 6 * 4 * 4, emptyData, GL.GL_STATIC_DRAW);
        }

        GL.vertexAttribPointer(0, 4, GL.GL_FLOAT, false, 4 * 4, 0);
        GL.enableVertexAttribArray(0);
        GL.bindBuffer(GL.GL_ARRAY_BUFFER, 0);
        GL.bindVertexArray(0);

        initialized = true;
    }

    public static void renderText(String text, float x, float y, float scale, float r, float g, float b) {
        if (!initialized) return;
        text = text.toUpperCase();

        List<Float> vertices = new ArrayList<>();
        float cursorX = x;
        float pixelSize = scale * 0.15f;
        float charWidth = 5 * pixelSize + pixelSize;

        for (char c : text.toCharArray()) {
            int[][] pattern = FONT.get(c);
            if (pattern == null) {
                cursorX += charWidth;
                continue;
            }
            for (int row = 0; row < 7; row++) {
                for (int col = 0; col < 5; col++) {
                    if (pattern[row][col] == 1) {
                        float px = cursorX + col * pixelSize;
                        float py = y - row * pixelSize;
                        // Triangle 1
                        vertices.add(px); vertices.add(py); vertices.add(0.3f); vertices.add(1.0f);
                        vertices.add(px + pixelSize); vertices.add(py); vertices.add(0.3f); vertices.add(1.0f);
                        vertices.add(px); vertices.add(py - pixelSize); vertices.add(0.3f); vertices.add(1.0f);
                        // Triangle 2
                        vertices.add(px + pixelSize); vertices.add(py); vertices.add(0.3f); vertices.add(1.0f);
                        vertices.add(px + pixelSize); vertices.add(py - pixelSize); vertices.add(0.3f); vertices.add(1.0f);
                        vertices.add(px); vertices.add(py - pixelSize); vertices.add(0.3f); vertices.add(1.0f);
                    }
                }
            }
            cursorX += charWidth;
        }

        if (vertices.isEmpty()) return;

        try (Arena arena = Arena.ofConfined()) {
            MemorySegment data = arena.allocate(vertices.size() * 4L);
            for (int i = 0; i < vertices.size(); i++) {
                data.setAtIndex(ValueLayout.JAVA_FLOAT, i, vertices.get(i));
            }
            GL.bindBuffer(GL.GL_ARRAY_BUFFER, vbo);
            // Re-upload entire buffer data each frame
            GL.bufferData(GL.GL_ARRAY_BUFFER, vertices.size() * 4L, data, GL.GL_STATIC_DRAW);
            GL.bindBuffer(GL.GL_ARRAY_BUFFER, 0);
        }

        Shader.TEXT.enable();
        Shader.TEXT.setUniformMat4f("pr_matrix", projectionMatrix);
        Shader.TEXT.setUniform3f("textColor", new Vector3f(r, g, b));
        GL.bindVertexArray(vao);
        GL.drawArrays(GL.GL_TRIANGLES, 0, vertices.size() / 4);
        GL.bindVertexArray(0);
        Shader.TEXT.disable();
    }

    public static void renderTextCentered(String text, float y, float scale, float r, float g, float b) {
        float pixelSize = scale * 0.15f;
        float charWidth = 5 * pixelSize + pixelSize;
        float totalWidth = text.length() * charWidth - pixelSize;
        float x = -totalWidth / 2.0f;
        renderText(text, x, y, scale, r, g, b);
    }
}
