package game.level;

import game.Shaders.Shader;
import game.audio.Audio;
import game.graphics.Texture;
import game.graphics.VertexArray;
import game.input.Input;
import game.math.Matrix4f;
import game.math.Vector3f;

public final class Bird {

    // Physics constants
    private static final float SIZE = 1.0f;
    private static final float FLAP_STRENGTH = -0.15f;
    private static final float GRAVITY = 0.01f;
    private static final float ROTATION_FACTOR = 90.0f;
    
    // Rendering components
    private final VertexArray mesh;
    private final Texture texture;

    // State
    private final Vector3f position = new Vector3f();
    private float velocity = 0.0f;
    private float rotation = 0.0f;
    private boolean wasSpacePressed = false;

    public Bird() {
        this.mesh = createMesh();
        this.texture = new Texture("res/bird.png");
    }
    
    /**
     * Create the bird's quad mesh.
     */
    private VertexArray createMesh() {
        float halfSize = SIZE / 2.0f;
        float[] vertices = {
            -halfSize, -halfSize, 0.2f,
            -halfSize,  halfSize, 0.2f,
             halfSize,  halfSize, 0.2f,
             halfSize, -halfSize, 0.2f,
        };
        byte[] indices = { 0, 1, 2, 2, 3, 0 };
        float[] texCoords = { 0, 1, 0, 0, 1, 0, 1, 1 };
        
        return new VertexArray(vertices, indices, texCoords);
    }

    public void update(boolean acceptInput) {
        position.y -= velocity;
        
        boolean spacePressed = acceptInput && Input.isKeyDown(Input.KEY_SPACE);
        if (spacePressed) {
            velocity = FLAP_STRENGTH;
            if (!wasSpacePressed) {
                Audio.playJumpSound();
            }
        } else {
            velocity += GRAVITY;
        }
        wasSpacePressed = spacePressed;

        rotation = -velocity * ROTATION_FACTOR;
    }

    public void fall() {
        velocity = FLAP_STRENGTH;
    }

    public void render() {
        Shader.BIRD.enable();
        
        float rotRadians = (float) Math.toRadians(rotation);
        Matrix4f modelMatrix = Matrix4f.translate(position).multiply(Matrix4f.rotateZ(rotRadians));
        Shader.BIRD.setUniformMat4f("ml_matrix", modelMatrix);
        
        texture.bind();
        mesh.render();
        texture.unbind();
        
        Shader.BIRD.disable();
    }

    public float getY() {
        return position.y;
    }

    public float getSize() {
        return SIZE;
    }

}
