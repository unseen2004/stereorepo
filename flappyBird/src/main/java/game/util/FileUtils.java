package game.util;

public class FileUtils {

    private FileUtils() {

    }

    public static String loadAsString(String path) {
        StringBuilder result = new StringBuilder();
        try (java.io.BufferedReader reader = new java.io.BufferedReader(new java.io.FileReader(path))) {
            String line;
            while ((line = reader.readLine()) != null) {
                result.append(line).append('\n');
            }
        } catch (java.io.IOException e) {
            e.printStackTrace();
            System.exit(-1);
        }
        return result.toString();
    }

}
