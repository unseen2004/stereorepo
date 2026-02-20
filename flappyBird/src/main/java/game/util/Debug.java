package game.util;

import game.gl.GL;

public class Debug {

    public static final boolean ENABLED = "true".equals(System.getProperty("debug"));
    private static int frameCount = 0;
    private static final int MAX_RENDER_LOG_FRAMES = 3;

    public static void log(String message) {
        if (ENABLED) {
            System.out.println("[DEBUG] " + message);
        }
    }

    public static void log(String category, String message) {
        if (ENABLED) {
            if ((category.equals("RENDER") || category.equals("DRAW") || category.equals("TEXTURE")) && frameCount > MAX_RENDER_LOG_FRAMES) {
                return;
            }
            System.out.println("[DEBUG:" + category + "] " + message);
        }
    }

    public static void error(String message) {
        if (ENABLED) {
            System.err.println("[DEBUG:ERROR] " + message);
        }
    }

    public static void checkGLError(String context) {
        if (ENABLED) {
            int error = GL.getError();
            if (error != GL.GL_NO_ERROR) {
                System.err.println("[DEBUG:GL_ERROR] " + context + " - Error code: 0x" + Integer.toHexString(error));
            }
        }
    }

    public static void incrementFrame() {
        frameCount++;
    }
}

