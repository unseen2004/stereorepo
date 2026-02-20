package game.graphics;

import game.gl.GL;
import game.util.BufferUtils;
import game.util.Debug;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.FileInputStream;
import java.io.IOException;
import java.lang.foreign.MemorySegment;

public class Texture {

    private int width, height;
    private int texture;
    private String path;

    public Texture(String path) {
        this.path = path;
        Debug.log("TEXTURE", "Loading texture: " + path);
        texture = load(path);
    }

    private int load(String path) {
        int[] pixels = null;
        try {
            BufferedImage image = ImageIO.read(new FileInputStream(path));
            width = image.getWidth();
            height = image.getHeight();
            pixels = new int[width * height];
            image.getRGB(0, 0, width, height, pixels, 0, width);
            Debug.log("TEXTURE", "Image loaded: " + width + "x" + height + " pixels");
        } catch (IOException e) {
            Debug.error("Failed to load texture: " + path);
            e.printStackTrace();
        }

        if (pixels == null) {
            Debug.error("Pixels array is null for texture: " + path);
            return 0;
        }

        int[] data = new int[width * height];
        for (int i = 0; i < width * height; i++) {
            int a = (pixels[i] >> 24) & 0xff;
            int r = (pixels[i] >> 16) & 0xff;
            int g = (pixels[i] >> 8) & 0xff;
            int b = (pixels[i]) & 0xff;
            data[i] = a << 24 | b << 16 | g << 8 | r;
        }

        // Sample first pixel for debugging
        if (data.length > 0) {
            Debug.log("TEXTURE", "First pixel ABGR: " + Integer.toHexString(data[0]));
        }

        int tex = GL.genTexture();
        Debug.log("TEXTURE", "Generated texture ID: " + tex);

        GL.bindTexture(GL.GL_TEXTURE_2D, tex);
        GL.texParameteri(GL.GL_TEXTURE_2D, GL.GL_TEXTURE_MIN_FILTER, GL.GL_NEAREST);
        GL.texParameteri(GL.GL_TEXTURE_2D, GL.GL_TEXTURE_MAG_FILTER, GL.GL_NEAREST);
        GL.texImage2D(GL.GL_TEXTURE_2D, 0, GL.GL_RGBA, width, height, 0, GL.GL_RGBA, GL.GL_UNSIGNED_BYTE, MemorySegment.ofBuffer(BufferUtils.createIntBuffer(data)));
        GL.bindTexture(GL.GL_TEXTURE_2D, 0);

        Debug.log("TEXTURE", "Texture uploaded to GPU successfully: " + path);
        return tex;
    }

    public void bind() {
        Debug.log("TEXTURE", "Binding texture ID: " + texture + " (" + path + ")");
         GL.activeTexture(GL.GL_TEXTURE0);
        GL.bindTexture(GL.GL_TEXTURE_2D, texture);
    }

    public void unbind() {
        GL.bindTexture(GL.GL_TEXTURE_2D, 0);
    }

    public int getWidth() {
        return width;
    }
    public int getHeight() {
        return height;
    }
    public int getID() {
        return texture;
    }
}

