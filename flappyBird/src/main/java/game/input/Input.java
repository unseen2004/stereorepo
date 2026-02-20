package game.input;

public final class Input {

    private static final int MAX_KEYS = 512;

    private static final boolean[] keyStates = new boolean[MAX_KEYS];

    private static final boolean[] keyPressed = new boolean[MAX_KEYS];

    public static final int KEY_SPACE = 44;
    public static final int KEY_ENTER = 40;
    public static final int KEY_ESCAPE = 41;
    public static final int KEY_Q = 20;
    public static final int KEY_M = 16;

    private Input() {}

    public static void setKey(int scancode, boolean down) {
        if (isValidScancode(scancode)) {
            if (down && !keyStates[scancode]) {
                keyPressed[scancode] = true;
            }
            keyStates[scancode] = down;
        }
    }

    public static boolean isKeyDown(int scancode) {
        return isValidScancode(scancode) && keyStates[scancode];
    }

    public static boolean isKeyPressed(int scancode) {
        if (isValidScancode(scancode) && keyPressed[scancode]) {
            keyPressed[scancode] = false;
            return true;
        }
        return false;
    }

    private static boolean isValidScancode(int scancode) {
        return scancode >= 0 && scancode < MAX_KEYS;
    }
}
