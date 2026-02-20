package game;

public enum GameState {
    MENU("Menu"),
    PLAYING("Playing"),
    GAME_OVER("Game Over");

    private final String displayName;

    GameState(String displayName) {
        this.displayName = displayName;
    }

    public String getDisplayName() {
        return displayName;
    }

    public boolean isActive() {
        return this == PLAYING;
    }

    public boolean isGameOver() {
        return this == GAME_OVER;
    }

    public boolean isMenu() {
        return this == MENU;
    }
}

