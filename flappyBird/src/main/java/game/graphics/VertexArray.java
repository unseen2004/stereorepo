package game.graphics;

import game.Shaders.Shader;
import game.gl.GL;
import game.util.BufferUtils;

public class VertexArray {

    private int count;
    private int indicesCount;
    private int vao, vbo, ibo, tbo;

    public VertexArray(int count){
        this.count = count;
        this.indicesCount = 0;
        vao = GL.genVertexArray();
        GL.bindVertexArray(vao);
    }

    public VertexArray(float[] vertices, byte[] indices, float[] textureCoordinates) {
        count = vertices.length / 3;
        indicesCount = indices.length;

        vao = GL.genVertexArray();
        GL.bindVertexArray(vao);

        vbo = GL.genBuffer();
        GL.bindBuffer(GL.GL_ARRAY_BUFFER, vbo);
        GL.bufferData(GL.GL_ARRAY_BUFFER, vertices.length * Float.BYTES, java.lang.foreign.MemorySegment.ofBuffer(BufferUtils.createFloatBuffer(vertices)), GL.GL_STATIC_DRAW);
        GL.vertexAttribPointer(Shader.VERTEX_ATTRIB, 3, GL.GL_FLOAT, false, 0, 0);
        GL.enableVertexAttribArray(Shader.VERTEX_ATTRIB);

        tbo = GL.genBuffer();
        GL.bindBuffer(GL.GL_ARRAY_BUFFER, tbo);
        GL.bufferData(GL.GL_ARRAY_BUFFER, textureCoordinates.length * Float.BYTES, java.lang.foreign.MemorySegment.ofBuffer(BufferUtils.createFloatBuffer(textureCoordinates)), GL.GL_STATIC_DRAW);
        GL.vertexAttribPointer(Shader.TEXCOORD_ATTRIB, 2, GL.GL_FLOAT, false, 0, 0);
        GL.enableVertexAttribArray(Shader.TEXCOORD_ATTRIB);

        ibo = GL.genBuffer();
        GL.bindBuffer(GL.GL_ELEMENT_ARRAY_BUFFER, ibo);
        GL.bufferData(GL.GL_ELEMENT_ARRAY_BUFFER, indices.length, java.lang.foreign.MemorySegment.ofBuffer(BufferUtils.createByteBuffer(indices)), GL.GL_STATIC_DRAW);

        GL.bindBuffer(GL.GL_ARRAY_BUFFER, 0);
        GL.bindVertexArray(0);
        GL.bindBuffer(GL.GL_ELEMENT_ARRAY_BUFFER, 0);
    }

    public void bind() {
        GL.bindVertexArray(vao);
        if (ibo > 0)
            GL.bindBuffer(GL.GL_ELEMENT_ARRAY_BUFFER, ibo);
    }

    public void unbind() {
        if (ibo > 0)
            GL.bindBuffer(GL.GL_ELEMENT_ARRAY_BUFFER, 0);
        GL.bindVertexArray(0);
    }

    public void draw() {
        if (ibo > 0) {
            GL.drawElements(GL.GL_TRIANGLES, indicesCount, GL.GL_UNSIGNED_BYTE, 0);
        } else {
            GL.drawArrays(GL.GL_TRIANGLES, 0, count);
        }
    }

    public void render() {
        bind();
        draw();
    }

}
