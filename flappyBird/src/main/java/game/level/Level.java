package game.level;

import game.Shaders.Shader;
import game.graphics.TextRenderer;
import game.graphics.Texture;
import game.graphics.VertexArray;
import game.math.Matrix4f;
import game.math.Vector3f;
import game.util.Debug;

import java.util.Arrays;
import java.util.Random;

public final class Level {

    // Game constants
    private static final int PIPE_COUNT = 10;
    private static final float PIPE_OFFSET = 5.0f;
    private static final float PIPE_SPACING = 3.0f;
    private static final float PIPE_GAP = 11.5f;
    private static final float SCROLL_SPEED = 0.05f;
    private static final float ASPECT_RATIO = 9.0f / 16.0f;
    private static final float HITBOX_MARGIN = 0.15f;

    // Rendering components
    private final VertexArray background;
    private final Texture bgTexture;
    private final Bird bird;
    private final Pipe[] pipes = new Pipe[PIPE_COUNT];
    private final Random random = new Random();

    // Game state
    private int xScroll = 0;
    private int mapOffset = 0;
    private int pipeIndex = 0;
    private boolean isAlive = true;
    private float animationTime = 0;
    private float gameTimer = 0;

    // FPS tracking
    private int frameCount = 0;
    private long lastFpsTime = System.currentTimeMillis();
    private int currentFps = 0;

    public Level() {
        this.background = createBackground();
        this.bgTexture = new Texture("res/bg.jpeg");
        this.bird = new Bird();
        initializePipes();
    }

    /**
     * Create the background mesh.
     */
    private VertexArray createBackground() {
        float h = 10.0f * ASPECT_RATIO;
        float[] vertices = {
            -10.0f, -h, 0.0f,
            -10.0f,  h, 0.0f,
              0.0f,  h, 0.0f,
              0.0f, -h, 0.0f,
        };
        byte[] indices = { 0, 1, 2, 2, 3, 0 };
        float[] texCoords = { 0, 1, 0, 0, 1, 0, 1, 1 };
        return new VertexArray(vertices, indices, texCoords);
    }

    /**
     * Initialize all pipes.
     */
    private void initializePipes() {
        Pipe.create();
        for (int i = 0; i < PIPE_COUNT; i += 2) {
            float x = PIPE_OFFSET + pipeIndex * PIPE_SPACING;
            float y = random.nextFloat() * 4.0f;
            pipes[i] = new Pipe(x, y);
            pipes[i + 1] = new Pipe(x, y - PIPE_GAP);
            pipeIndex += 2;
        }
    }

    /**
     * Spawn new pipes at the end of the pipe array.
     */
    private void spawnNextPipes() {
        int idx = pipeIndex % PIPE_COUNT;
        float x = PIPE_OFFSET + pipeIndex * PIPE_SPACING;
        float y = random.nextFloat() * 4.0f;
        pipes[idx] = new Pipe(x, y);
        pipes[(idx + 1) % PIPE_COUNT] = new Pipe(x, y - PIPE_GAP);
        pipeIndex += 2;
    }

    /**
     * Check collision between bird and any pipe using AABB.
     */
    private boolean checkCollision() {
        float birdX = -xScroll * SCROLL_SPEED;
        float birdY = bird.getY();
        float birdHalfSize = bird.getSize() / 2.0f - HITBOX_MARGIN;

        return Arrays.stream(pipes).anyMatch(pipe -> {
            float bx0 = birdX - birdHalfSize, bx1 = birdX + birdHalfSize;
            float by0 = birdY - birdHalfSize, by1 = birdY + birdHalfSize;
            float px0 = pipe.getX(), px1 = pipe.getX() + Pipe.getWidth();
            float py0 = pipe.getY(), py1 = pipe.getY() + Pipe.getHeight();
            return bx1 > px0 && bx0 < px1 && by1 > py0 && by0 < py1;
        });
    }

    /**
     * Check if bird has gone out of screen bounds.
     */
    private boolean isOutOfBounds() {
        float birdY = bird.getY();
        float screenBound = 10.0f * ASPECT_RATIO + 1.0f;
        return birdY > screenBound || birdY < -screenBound;
    }

    /**
     * Render the entire level: background, pipes, bird, and HUD.
     */
    public void render() {
        renderBackground();
        renderPipes();
        renderBird();
        renderHUD();
    }

    private void renderBackground() {
        bgTexture.bind();
        Shader.BG.enable();
        Shader.BG.setUniform2f("bird", 0, bird.getY());
        background.bind();

        for (int i = mapOffset; i < mapOffset + 4; i++) {
            var translation = new Vector3f(i * 10 + xScroll * 0.03f, 0.0f, 0.0f);
            Shader.BG.setUniformMat4f("vw_matrix", Matrix4f.translate(translation));
            background.draw();
        }

        Shader.BG.disable();
        bgTexture.unbind();
    }

    private void renderBird() {
        Shader.FADE.enable();
        Shader.FADE.setUniform1f("time", animationTime);
        bird.render();
        Shader.FADE.disable();
    }

    private void renderHUD() {
        frameCount++;
        long currentTime = System.currentTimeMillis();
        if (currentTime - lastFpsTime >= 1000) {
            currentFps = frameCount;
            frameCount = 0;
            lastFpsTime = currentTime;
        }

        TextRenderer.renderTextCentered("FLAPPY BIRD", 4.8f, 0.6f, 1.0f, 1.0f, 0.0f);

        int seconds = (int) gameTimer;
        int minutes = seconds / 60;
        seconds = seconds % 60;
        String timerText = String.format("%d:%02d", minutes, seconds);
        TextRenderer.renderText(timerText, 8.0f, 5.3f, 0.4f, 1.0f, 1.0f, 1.0f);

        TextRenderer.renderText("FPS:" + currentFps, -9.8f, 5.3f, 0.25f, 0.7f, 0.7f, 0.7f);
    }

    public void update() {
        if (isAlive) {
            updateScrolling();
            animationTime += 0.01f;
            gameTimer += 1.0f / 60.0f;
        }

        bird.update(isAlive);

        if (isAlive && (checkCollision() || isOutOfBounds())) {
            bird.fall();
            isAlive = false;
        }
    }

    private void updateScrolling() {
        xScroll--;
        if (-xScroll % 335 == 0) {
            mapOffset++;
        }
        if (-xScroll > 250 && -xScroll % 120 == 0) {
            spawnNextPipes();
        }
    }

    private void renderPipes() {
        Shader.PIPE.enable();
        Shader.PIPE.setUniform2f("bird", 0, bird.getY());
        Shader.PIPE.setUniformMat4f("vw_matrix", Matrix4f.translate(new Vector3f(xScroll * SCROLL_SPEED, 0.0f, 0.0f)));

        Pipe.getTexture().bind();
        Pipe.getMesh().bind();

        for (int i = 0; i < PIPE_COUNT; i++) {
            Shader.PIPE.setUniformMat4f("ml_matrix", pipes[i].getMl_matrix());
            Shader.PIPE.setUniformi("top", i % 2 == 0 ? 1 : 0);
            Pipe.getMesh().draw();
        }

        Pipe.getTexture().unbind();
        Pipe.getMesh().unbind();
    }

    public boolean isGameOver() {
        return !isAlive;
    }

    public float getGameTimer() {
        return gameTimer;
    }

    public String getFormattedTime() {
        int seconds = (int) gameTimer;
        int minutes = seconds / 60;
        seconds = seconds % 60;
        return String.format("%d:%02d", minutes, seconds);
    }

}
