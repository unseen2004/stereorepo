package game.math;

import java.nio.FloatBuffer;

public class Matrix4f {
    // OpenGL uses column-major: element at row r, column c = matrix[r + c * 4]
    public float[] matrix = new float[4 * 4];

    public static Matrix4f identity() {
        Matrix4f result = new Matrix4f();
        // Diagonal: [0+0*4], [1+1*4], [2+2*4], [3+3*4]
        result.matrix[0 + 0 * 4] = 1;
        result.matrix[1 + 1 * 4] = 1;
        result.matrix[2 + 2 * 4] = 1;
        result.matrix[3 + 3 * 4] = 1;
        return result;
    }

    public Matrix4f multiply(Matrix4f other) {
        Matrix4f result = new Matrix4f();
        for (int row = 0; row < 4; row++) {
            for (int col = 0; col < 4; col++) {
                float sum = 0;
                for (int k = 0; k < 4; k++) {
                    sum += this.matrix[row + k * 4] * other.matrix[k + col * 4];
                }
                result.matrix[row + col * 4] = sum;
            }
        }
        return result;
    }

    public static Matrix4f translate(float x, float y, float z) {
        Matrix4f result = identity();
        // Translation goes in column 3: [row + 3*4]
        result.matrix[0 + 3 * 4] = x;
        result.matrix[1 + 3 * 4] = y;
        result.matrix[2 + 3 * 4] = z;
        return result;
    }

    public static Matrix4f translate(Vector3f vector) {
        return translate(vector.x, vector.y, vector.z);
    }

    public static Matrix4f rotateZ(float angle) {
        Matrix4f result = identity();
        float cos = (float) Math.cos(angle);
        float sin = (float) Math.sin(angle);
        result.matrix[0 + 0 * 4] = cos;
        result.matrix[1 + 0 * 4] = sin;
        result.matrix[0 + 1 * 4] = -sin;
        result.matrix[1 + 1 * 4] = cos;
        return result;
    }

    public static Matrix4f orthographic(float left, float right, float bottom, float top, float near, float far) {
        Matrix4f result = identity();
        result.matrix[0 + 0 * 4] = 2f / (right - left);
        result.matrix[1 + 1 * 4] = 2f / (top - bottom);
        result.matrix[2 + 2 * 4] = -2f / (far - near);
        result.matrix[0 + 3 * 4] = -(right + left) / (right - left);
        result.matrix[1 + 3 * 4] = -(top + bottom) / (top - bottom);
        result.matrix[2 + 3 * 4] = -(far + near) / (far - near);
        return result;
    }

    public FloatBuffer toFloatBuffer() {
        return game.util.BufferUtils.createFloatBuffer(matrix);
    }

}
