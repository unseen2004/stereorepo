package game;

import game.Shaders.Shader;
import game.audio.Audio;
import game.gl.GL;
import game.graphics.TextRenderer;
import game.graphics.VertexArray;
import game.input.Input;
import game.level.Level;
import game.level.Menu;
import game.math.Matrix4f;
import game.sdl.Sdl;
import game.util.Debug;

import java.lang.foreign.MemorySegment;

public final class Main implements Runnable {

    private static final int WINDOW_WIDTH = 1280;
    private static final int WINDOW_HEIGHT = 720;
    private static final String WINDOW_TITLE = "Flappy Bird";
    private static final double TARGET_UPS = 60.0;

    private MemorySegment window;
    private MemorySegment glContext;

    private Level level;
    private Menu menu;
    private VertexArray overlayQuad;
    private Matrix4f projectionMatrix;

    private volatile boolean running = false;
    private GameState gameState = GameState.MENU;
    private float gameOverAlpha = 0;
    private boolean firstFrame = true;

    /**
     * Start the game in a new thread.
     */
    public void start() {
        running = true;
        var gameThread = new Thread(this, "GameThread");
        gameThread.start();
    }

    /**
     * Initialize all game systems: SDL, OpenGL, Audio, Shaders, and game objects.
     */
    private void init() {
        Debug.log("INIT", "Debug mode enabled - diagnostic output active");
        Debug.log("INIT", "Initializing SDL...");
        Sdl.initVideo();

        Debug.log("INIT", "Creating window: %dx%d".formatted(WINDOW_WIDTH, WINDOW_HEIGHT));
        window = Sdl.createWindow(WINDOW_TITLE, WINDOW_WIDTH, WINDOW_HEIGHT);

        Debug.log("INIT", "Creating OpenGL 3.3 context...");
        glContext = Sdl.createGlContext(window, 3, 3);

        Debug.log("INIT", "Initializing audio system...");
        Audio.init();

        Debug.log("INIT", "Setting viewport: 0, 0, %d, %d".formatted(WINDOW_WIDTH, WINDOW_HEIGHT));
        GL.viewport(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT);

        Debug.log("INIT", "Enabling blending...");
        GL.enable(GL.GL_BLEND);
        GL.blendFunc(GL.GL_SRC_ALPHA, GL.GL_ONE_MINUS_SRC_ALPHA);
        GL.clearColor(0.0f, 0.5f, 0.5f, 1.0f);

        System.out.println("OpenGL version: " + GL.getString(GL.GL_VERSION));
        Debug.log("INIT", "OpenGL version: " + GL.getString(GL.GL_VERSION));

        initializeShaders();
        initializeGameObjects();

        Debug.log("INIT", "Initialization complete!");
    }

    /**
     * Load and configure all shader programs.
     */
    private void initializeShaders() {
        Debug.log("INIT", "Loading shaders...");
        Shader.loadAll();
        System.out.println("Shaders loaded successfully");

        Debug.log("INIT", "Setting up projection matrix...");
        projectionMatrix = Matrix4f.orthographic(-10.0f, 10.0f, -10.0f * 9.0f / 16.0f, 10.0f * 9.0f / 16.0f, -1.0f, 1.0f);

        // Configure shader uniforms
        configureShaderUniforms(Shader.BG, "BG");
        configureShaderUniforms(Shader.BIRD, "BIRD");
        configureShaderUniforms(Shader.PIPE, "PIPE");

        Debug.log("INIT", "Initializing TextRenderer...");
        TextRenderer.init(projectionMatrix);
        System.out.println("TextRenderer initialized successfully");
    }

    /**
     * Configure common uniforms for a shader.
     */
    private void configureShaderUniforms(Shader shader, String name) {
        Debug.log("INIT", "Setting %s shader uniforms...".formatted(name));
        shader.setUniformMat4f("pr_matrix", projectionMatrix);
        shader.setUniformi("tex", 0);
    }

    /**
     * Initialize game objects: level, menu, overlay.
     */
    private void initializeGameObjects() {

        Debug.log("INIT", "Creating overlay quad...");
        // Create fullscreen overlay quad
        float[] overlayVertices = new float[]{
                -10.0f, -10.0f * 9.0f / 16.0f, 0.5f,
                -10.0f, 10.0f * 9.0f / 16.0f, 0.5f,
                10.0f, 10.0f * 9.0f / 16.0f, 0.5f,
                10.0f, -10.0f * 9.0f / 16.0f, 0.5f,
        };
        byte[] overlayIndices = new byte[]{ 0, 1, 2, 2, 3, 0 };
        float[] overlayTcs = new float[]{ 0, 0, 0, 0, 0, 0, 0, 0 };
        overlayQuad = new VertexArray(overlayVertices, overlayIndices, overlayTcs);

        Debug.log("INIT", "Creating level...");
        level = new Level();
        System.out.println("Level created successfully");

        Debug.log("INIT", "Creating menu...");
        menu = new Menu();
        System.out.println("Menu created successfully");
    }

    /**
     * Main game loop with fixed update rate and variable render rate.
     */
    @Override
    public void run() {
        init();

        long lastTime = System.nanoTime();
        final double nsPerUpdate = 1_000_000_000.0 / TARGET_UPS;
        double delta = 0;
        int updates = 0;
        int frames = 0;
        long timer = System.currentTimeMillis();

        while (running) {
            long now = System.nanoTime();
            delta += (now - lastTime) / nsPerUpdate;
            lastTime = now;

            while (delta >= 1) {
                update();
                updates++;
                delta--;
            }

            render();
            frames++;

            if (System.currentTimeMillis() - timer > 1000) {
                timer += 1000;
                System.out.println("%d ups, %d fps".formatted(updates, frames));
                updates = 0;
                frames = 0;
            }
        }

        cleanup();
    }

    /**
     * Clean up all resources before exit.
     */
    private void cleanup() {
        Audio.cleanup();
        Sdl.deleteGlContext(glContext);
        Sdl.destroyWindow(window);
        Sdl.quit();
    }

    /**
     * Process input and update game state.
     */
    private void update() {
        processEvents();
        updateGameState();
    }

    /**
     * Poll and process SDL events.
     */
    private void processEvents() {
        int[] keyOut = new int[1];
        int type;
        while ((type = Sdl.pollEvent(keyOut)) != 0) {
            switch (type) {
                case Sdl.SDL_QUIT -> running = false;
                case Sdl.SDL_KEYDOWN -> Input.setKey(keyOut[0], true);
                case Sdl.SDL_KEYUP -> Input.setKey(keyOut[0], false);
                default -> { }
            }
        }
    }

    // SDL scancodes
    private static final int KEY_ENTER = 40;
    private static final int KEY_Q = 20;
    private static final int KEY_M = 16;

    /**
     * Update game based on current state using switch expression.
     */
    private void updateGameState() {
        switch (gameState) {
            case MENU -> updateMenuState();
            case PLAYING -> updatePlayingState();
            case GAME_OVER -> updateGameOverState();
        }
    }

    private void updateMenuState() {
        Audio.playMenuMusic();
        menu.update();

        if (Input.isKeyPressed(KEY_ENTER)) {
            level = new Level();
            gameState = GameState.PLAYING;
            gameOverAlpha = 0;
            Audio.playGameMusic();
        }
        if (Input.isKeyPressed(KEY_Q)) {
            running = false;
        }
    }

    private void updatePlayingState() {
        Audio.playGameMusic();
        level.update();

        if (level.isGameOver()) {
            gameState = GameState.GAME_OVER;
            gameOverAlpha = 0;
            Audio.playDeadSound();
            Audio.playDeadMusic();
        }
    }

    private void updateGameOverState() {
        Audio.playDeadMusic();
        level.update();

        if (gameOverAlpha < 0.7f) {
            gameOverAlpha += 0.02f;
        }

        if (Input.isKeyPressed(KEY_ENTER)) {
            level = new Level();
            gameState = GameState.PLAYING;
            gameOverAlpha = 0;
            Audio.playGameMusic();
        }
        if (Input.isKeyPressed(KEY_M)) {
            gameState = GameState.MENU;
            gameOverAlpha = 0;
            Audio.playMenuMusic();
        }
        if (Input.isKeyPressed(KEY_Q)) {
            running = false;
        }
    }

    /**
     * Render the current game state.
     */
    private void render() {
        if (firstFrame) {
            Debug.log("RENDER", "=== First frame render ===");
        }

        GL.clear(GL.GL_COLOR_BUFFER_BIT | GL.GL_DEPTH_BUFFER_BIT);

        switch (gameState) {
            case MENU -> menu.render();
            case PLAYING -> level.render();
            case GAME_OVER -> renderGameOver();
        }

        Sdl.swapWindow(window);

        if (firstFrame) {
            Debug.log("RENDER", "=== First frame complete ===");
            firstFrame = false;
        }
        Debug.incrementFrame();
    }

    /**
     * Render the game over screen with overlay and text.
     */
    private void renderGameOver() {
        level.render();

        Shader.OVERLAY.enable();
        Shader.OVERLAY.setUniformMat4f("pr_matrix", projectionMatrix);
        Shader.OVERLAY.setUniform1f("alpha", gameOverAlpha);
        overlayQuad.bind();
        overlayQuad.draw();
        overlayQuad.unbind();
        Shader.OVERLAY.disable();

        TextRenderer.renderTextCentered("GAME OVER", 3.0f, 1.0f, 1.0f, 0.3f, 0.3f);
        TextRenderer.renderTextCentered("TIME: " + level.getFormattedTime(), 1.5f, 0.6f, 1.0f, 1.0f, 1.0f);
        TextRenderer.renderTextCentered("ENTER: PLAY AGAIN", -0.5f, 0.5f, 0.8f, 1.0f, 0.8f);
        TextRenderer.renderTextCentered("M: MENU", -1.5f, 0.5f, 0.8f, 1.0f, 0.8f);
        TextRenderer.renderTextCentered("Q: QUIT", -2.5f, 0.5f, 0.8f, 1.0f, 0.8f);
    }

    public static void main(String[] args) {
        new Main().start();
    }
}
