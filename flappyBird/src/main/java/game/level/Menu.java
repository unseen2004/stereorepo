package game.level;

import game.Shaders.Shader;
import game.graphics.TextRenderer;
import game.graphics.Texture;
import game.graphics.VertexArray;
import game.math.Matrix4f;
import game.math.Vector3f;

public class Menu {

    private VertexArray background;
    private Texture bgTexture;
    private Texture birdTexture;
    private VertexArray birdQuad;
    private float time = 0;

    public Menu() {
        // Full screen background
        float[] vertices = new float[]{
                -10.0f, -10.0f * 9.0f / 16.0f, 0.0f,
                -10.0f, 10.0f * 9.0f / 16.0f, 0.0f,
                10.0f, 10.0f * 9.0f / 16.0f, 0.0f,
                10.0f, -10.0f * 9.0f / 16.0f, 0.0f,
        };

        byte[] indices = new byte[]{
                0, 1, 2,
                2, 3, 0
        };

        float[] tcs = new float[]{
                0.0f, 1.0f,
                0.0f, 0.0f,
                1.0f, 0.0f,
                1.0f, 1.0f,
        };

        background = new VertexArray(vertices, indices, tcs);
        bgTexture = new Texture("res/bg.jpeg");
        birdTexture = new Texture("res/bird.png");

        // Bird quad for logo
        float[] birdVerts = new float[]{
                -1.5f, -1.5f, 0.2f,
                -1.5f, 1.5f, 0.2f,
                1.5f, 1.5f, 0.2f,
                1.5f, -1.5f, 0.2f,
        };
        birdQuad = new VertexArray(birdVerts, indices, tcs);
    }

    public void update() {
        time += 0.05f;
    }

    public void render() {
        bgTexture.bind();
        Shader.BG.enable();
        Shader.BG.setUniform2f("bird", 0, 0);
        Shader.BG.setUniformMat4f("vw_matrix", Matrix4f.translate(new Vector3f(0.0f, 0.0f, 0.0f)));
        background.bind();
        background.draw();
        Shader.BG.disable();
        bgTexture.unbind();

        Shader.BIRD.enable();
        float birdY = (float) Math.sin(time) * 0.5f + 1.0f;
        Shader.BIRD.setUniformMat4f("ml_matrix", Matrix4f.translate(new Vector3f(0.0f, birdY, 0.0f)));

        birdTexture.bind();
        birdQuad.bind();
        birdQuad.draw();
        birdQuad.unbind();
        birdTexture.unbind();
        Shader.BIRD.disable();

        TextRenderer.renderTextCentered("FLAPPY BIRD", 4.0f, 1.0f, 1.0f, 1.0f, 0.0f);

        if ((int)(time * 2) % 2 == 0) {
            TextRenderer.renderTextCentered("PRESS ENTER TO PLAY", -3.0f, 0.5f, 1.0f, 1.0f, 1.0f);
        }

        TextRenderer.renderTextCentered("Q: QUIT", -4.5f, 0.4f, 0.7f, 0.7f, 0.7f);
    }
}

