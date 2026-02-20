package game.util;

import game.gl.GL;

public class ShaderUtils {

    private ShaderUtils() {

    }

    public static int load(String vertPath, String fragPath) {
        Debug.log("SHADER", "Loading shaders: " + vertPath + ", " + fragPath);
        String vertSource = FileUtils.loadAsString(vertPath);
        String fragSource = FileUtils.loadAsString(fragPath);
        Debug.log("SHADER", "Vertex source length: " + vertSource.length() + ", Fragment source length: " + fragSource.length());
        return create(vertSource, fragSource);
    }

    public static int create(String vert, String frag) {
        int program = GL.createProgram();
        Debug.log("SHADER", "Created program ID: " + program);

        int vertID = GL.createShader(GL.GL_VERTEX_SHADER);
        int fragID = GL.createShader(GL.GL_FRAGMENT_SHADER);
        Debug.log("SHADER", "Created vertex shader ID: " + vertID + ", fragment shader ID: " + fragID);

        GL.shaderSource(vertID, vert);
        GL.shaderSource(fragID, frag);

        GL.compileShader(vertID);
        int vertStatus = GL.getShaderiv(vertID, GL.GL_COMPILE_STATUS);
        Debug.log("SHADER", "Vertex shader compile status: " + vertStatus + " (1=success, 0=fail)");
        if (vertStatus == GL.GL_FALSE) {
            System.err.println("Vertex shader compilation failed! Status: " + vertStatus);
            System.err.println(GL.getShaderInfoLog(vertID));
        }

        GL.compileShader(fragID);
        int fragStatus = GL.getShaderiv(fragID, GL.GL_COMPILE_STATUS);
        Debug.log("SHADER", "Fragment shader compile status: " + fragStatus + " (1=success, 0=fail)");
        if (fragStatus == GL.GL_FALSE) {
            System.err.println("Fragment shader compilation failed! Status: " + fragStatus);
            System.err.println(GL.getShaderInfoLog(fragID));
        }

        GL.attachShader(program, vertID);
        GL.attachShader(program, fragID);
        GL.linkProgram(program);
        GL.validateProgram(program);

        // Check for link errors
        int linkStatus = GL.getProgramiv(program, GL.GL_LINK_STATUS);
        Debug.log("SHADER", "Program link status: " + linkStatus + " (1=success, 0=fail)");
        if (linkStatus == GL.GL_FALSE) {
            System.err.println("Program linking failed! Status: " + linkStatus);
            System.err.println(GL.getProgramInfoLog(program));
            return -1;
        }

        Debug.log("SHADER", "Shader program created successfully with ID: " + program);
        return program;
    }

}
